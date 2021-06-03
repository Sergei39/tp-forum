CREATE EXTENSION citext;

DROP TABLE votes CASCADE;
DROP TABLE posts CASCADE;
DROP TABLE threads CASCADE;
DROP TABLE users CASCADE;
DROP TABLE forums CASCADE;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    nickname CITEXT UNIQUE NOT NULL,
    fullname TEXT,
    about TEXT,
    email CITEXT UNIQUE
    -- password BYTEA --???
);

CREATE TABLE forums (
    id SERIAL PRIMARY KEY,
    -- user_create TEXT,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE,
    title TEXT,
    slug CITEXT UNIQUE -- человекочетаемый URL
    -- возможная оптимизация в будущем
    -- создать поля с кол-вом сообщений и кол-вом обсуждений
);

CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    title TEXT,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE,
    message TEXT, -- описание ветки
    -- возможная оптимизация в будущем
    -- создать поля с кол-вом голосов
    slug TEXT,
    created TIMESTAMP
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT,
    parent INTEGER DEFAULT 0 REFERENCES posts(id) ON DELETE CASCADE,
    forum INTEGER REFERENCES forums(id) ON DELETE CASCADE,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE,
    thread INTEGER REFERENCES threads(id) ON DELETE CASCADE,
    created TIMESTAMP,
    message TEXT,
    is_edited BOOLEAN DEFAULT FALSE
);

CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE,
    thread INTEGER REFERENCES threads(id) ON DELETE CASCADE,
    voice INTEGER
);
