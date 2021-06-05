CREATE EXTENSION citext;

DROP TABLE forums_users CASCADE;
DROP TABLE votes CASCADE;
DROP TABLE posts CASCADE;
DROP TABLE threads CASCADE;
DROP TABLE users CASCADE;
DROP TABLE forums CASCADE;


DROP SEQUENCE post_tree_id;
CREATE SEQUENCE post_tree_id;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    nickname CITEXT UNIQUE NOT NULL COLLATE "POSIX",
    fullname TEXT,
    about TEXT,
    email CITEXT UNIQUE
    -- password BYTEA --???
);

CREATE TABLE forums (
    id SERIAL PRIMARY KEY,
    -- user_create TEXT,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    title TEXT,
    slug CITEXT UNIQUE NOT NULL -- человекочетаемый URL
    -- возможная оптимизация в будущем
    -- создать поля с кол-вом сообщений и кол-вом обсуждений
);

CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    title TEXT,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE ,
    message TEXT, -- описание ветки
    -- возможная оптимизация в будущем
    -- создать поля с кол-вом голосов
    slug CITEXT NOT NULL,
    created TIMESTAMP with time zone
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT,
    -- сделать привязку к posts
    parent INTEGER DEFAULT 0 NOT NULL,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    thread INTEGER REFERENCES threads(id) ON DELETE CASCADE NOT NULL,
    created TIMESTAMP with time zone,
    message TEXT,
    is_edited BOOLEAN DEFAULT FALSE,
    tree BIGINT[]
);

CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    thread INTEGER REFERENCES threads(id) ON DELETE CASCADE NOT NULL,
    voice INTEGER NOT NULL,
    UNIQUE (user_create, thread)
);

CREATE TABLE forums_users (
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE NOT NULL
);
