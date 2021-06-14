CREATE EXTENSION citext;

DROP TABLE forums_users CASCADE;
DROP TABLE votes CASCADE;
DROP TABLE posts CASCADE;
DROP TABLE threads CASCADE;
DROP TABLE users CASCADE;
DROP TABLE forums CASCADE;


DROP SEQUENCE post_tree_id;
CREATE SEQUENCE post_tree_id;

CREATE UNLOGGED TABLE users (
    nickname CITEXT UNIQUE NOT NULL COLLATE "POSIX",
    fullname TEXT,
    about TEXT,
    email CITEXT UNIQUE
);

CREATE UNLOGGED TABLE forums (
    id SERIAL PRIMARY KEY,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    title TEXT,
    slug CITEXT UNIQUE NOT NULL, -- человекочетаемый URL
    threads INTEGER DEFAULT 0 NOT NULL,
    posts INTEGER DEFAULT 0 NOT NULL
);

CREATE UNLOGGED TABLE threads (
    id SERIAL PRIMARY KEY,
    title TEXT,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE ,
    message TEXT, -- описание ветки
    votes INTEGER DEFAULT 0 NOT NULL,
    slug CITEXT NOT NULL,
    created TIMESTAMP with time zone
);

CREATE UNLOGGED TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT,
    root_id INTEGER NOT NULL,
    parent INTEGER REFERENCES posts(id) DEFAULT NULL,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    thread INTEGER REFERENCES threads(id) ON DELETE CASCADE NOT NULL,
    created TIMESTAMP with time zone,
    message TEXT,
    is_edited BOOLEAN DEFAULT FALSE,
    tree INTEGER[]
);

CREATE UNLOGGED TABLE votes (
    id SERIAL PRIMARY KEY,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    thread INTEGER REFERENCES threads(id) ON DELETE CASCADE NOT NULL,
    voice INTEGER NOT NULL,
    UNIQUE (user_create, thread)
);

CREATE UNLOGGED TABLE forums_users (
    user_nickname CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL COLLATE "POSIX",
    user_fullname TEXT,
    user_about TEXT,
    user_email CITEXT,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
    UNIQUE (user_nickname, forum)
);


CREATE OR REPLACE FUNCTION add_tree() RETURNS TRIGGER AS
$add_tree$
declare
    parents INTEGER[];
begin
    if (new.parent is null) then
        new.tree := new.tree || new.id;
        new.root_id := new.tree[1];
    else
        select tree from posts where id = new.parent and thread = new.thread
        into parents;

        if (coalesce(array_length(parents, 1), 0) = 0) then
            raise exception 'parent post not exists' USING ERRCODE = '12345';
        end if;

        new.tree := new.tree || parents || new.id;
        new.root_id := new.tree[1];
    end if;
    return new;
end;
$add_tree$ LANGUAGE plpgsql;

create trigger add_path
    before insert on posts for each row
execute procedure add_tree();

-- функция и триггер при создании поста, на увеличение кол-ва постов в forums
CREATE OR REPLACE FUNCTION insert_post() RETURNS TRIGGER AS
$insert_post$
BEGIN
    UPDATE forums SET posts=posts + 1 WHERE forums.slug = NEW.forum;
    RETURN NEW;
END
$insert_post$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post
AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE insert_post();


-- функция и триггер при создании ветки, на увеличение кол-ва веток в forums
CREATE OR REPLACE FUNCTION insert_thread() RETURNS TRIGGER AS
$insert_thread$
BEGIN
    UPDATE forums SET threads=threads + 1 WHERE forums.slug = NEW.forum;
    RETURN NEW;
END
$insert_thread$ LANGUAGE plpgsql;

CREATE TRIGGER insert_thread
AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE insert_thread();


-- функция и триггер при создании ветки и поста, на добавления пользователя в список форума
CREATE OR REPLACE FUNCTION new_forum_user_added() RETURNS TRIGGER AS
$new_forum_user_added$
BEGIN
    DECLARE
        nickAuthor citext;
        fullnameAuthor text;
        emailAuthor citext;
        aboutAuthor text;
    BEGIN
        SELECT nickname, fullname, about, email
        FROM users WHERE nickname = NEW.user_create
        INTO nickAuthor, fullnameAuthor, aboutAuthor, emailAuthor;

        INSERT INTO forums_users(user_fullname, user_about, user_email, user_nickname, forum)
        VALUES (fullnameAuthor, aboutAuthor, emailAuthor, nickAuthor, new.forum)
        ON CONFLICT DO NOTHING;

        RETURN NULL;
    END;
END;
$new_forum_user_added$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS new_forum_user_added ON posts;
CREATE TRIGGER new_forum_user_added
    AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE new_forum_user_added();

DROP TRIGGER IF EXISTS new_forum_user_added ON threads;
CREATE TRIGGER new_forum_user_added
    AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE new_forum_user_added();

-- функция и триггер при создании голоса, на увеличение кол-ва голосов в threads
CREATE OR REPLACE FUNCTION insert_voice() RETURNS TRIGGER AS
$insert_voice$
BEGIN
    UPDATE threads SET votes=(votes + NEW.voice) WHERE id = NEW.thread;
    RETURN NULL;
END
$insert_voice$ LANGUAGE plpgsql;

CREATE TRIGGER insert_vote
AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE insert_voice();


-- функция и триггер при обновлении голоса, на изменение кол-ва голосов в threads
CREATE OR REPLACE FUNCTION update_voice() RETURNS TRIGGER AS
$update_voice$
BEGIN
    UPDATE threads SET votes= votes - OLD.voice + NEW.voice  WHERE id = NEW.thread;

    RETURN NULL;
END
$update_voice$ LANGUAGE plpgsql;

CREATE TRIGGER update_voice
AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE update_voice();


-- index
CREATE INDEX IF NOT EXISTS forum_slug ON forums (slug);

CREATE INDEX IF NOT EXISTS forums_user_user ON forums_users (user_nickname); -- подумать надо ли
CREATE INDEX IF NOT EXISTS forums_user_forum ON forums_users (forum); -- для получения всех юзеров из форума

CREATE INDEX IF NOT EXISTS user_nickname ON users (nickname);

CREATE INDEX IF NOT EXISTS thr_slug ON threads (slug) WHERE slug != '';
CREATE INDEX IF NOT EXISTS thr_forum ON threads (forum); -- для получения всех веток из форума

-- CREATE INDEX IF NOT EXISTS post_thread_id on posts (id, thread);
-- CREATE INDEX IF NOT EXISTS post_thread on posts (thread); -- подумать нужно ли если есть post_thread_id
CREATE INDEX IF NOT EXISTS post_thread_id on posts (thread, id); -- нужно для запросаполучения постов с последующим order by
CREATE INDEX IF NOT EXISTS post_thread_tree on posts (thread, tree); -- для запроса получения постов при сортировки flat
CREATE INDEX IF NOT EXISTS post_thread_root_id on posts (thread, root_id); -- для изменения плана слияния в сортировках tree, tree_parent
CREATE INDEX IF NOT EXISTS post_root_id on posts (root_id); -- не факт что нужно
CREATE INDEX IF NOT EXISTS post_root_id_desc_tree on posts (root_id DESC, tree); -- parent_tree ускоряет на немного
