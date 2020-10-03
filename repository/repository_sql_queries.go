package repository

const TABLES_DROPPING = `
DROP TABLE IF EXISTS videos_subthemes;
DROP TABLE IF EXISTS videos;
DROP TABLE IF EXISTS avatars;
DROP TABLE IF EXISTS masters_languages
DROP TABLE IF EXISTS masters_subthemes;
DROP TABLE IF EXISTS masters;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS subthemes;
DROP TABLE IF EXISTS themes;
DROP TABLE IF EXISTS languages;
`

const TABLES_CREATION = `
CREATE TABLE languages (
    id SERIAL NOT NULL PRIMARY KEY,
    name text NOT NULL UNIQUE
);

CREATE TABLE themes (
    id SERIAL NOT NULL PRIMARY KEY,
    name text NOT NULL UNIQUE
);

CREATE TABLE subthemes (
    id SERIAL NOT NULL PRIMARY KEY,
    theme_id int NOT NULL REFERENCES themes(id) ON DELETE CASCADE,
    name text NOT NULL UNIQUE
);

CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    email text NOT NULL UNIQUE,
    password text NOT NULL,
    type int NOT NULL CHECK (type = 0 OR type = 1),
    created TIMESTAMPTZ NOT NULL
    CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);

CREATE TABLE masters (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id int NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    username text NOT NULL UNIQUE ,
    fullname text NOT NULL,
    theme int NOT NULL REFERENCES themes(id) ON DELETE SET NULL,
    description text ,
    qualification int NOT NULL CHECK (qualification = 1 OR qualification = 2),
    education_format int NOT NULL CHECK (education_format >=1 AND education_format <= 3),
    avg_price numeric(20, 2) CONSTRAINT non_negative_price CHECK (avg_price >= 0)
);

CREATE TABLE IF NOT EXISTS masters_subthemes (
	master_id int NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    subtheme_id int NOT NULL REFERENCES subthemes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS masters_languages (
	master_id int NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    language_id int NOT NULL REFERENCES languages(id) ON DELETE CASCADE
);

CREATE TABLE avatars (
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename text NOT NULL UNIQUE
);

CREATE TABLE videos (
    id SERIAL NOT NULL PRIMARY KEY,
    master_id int NOT NULL REFERENCES masters(id) ON DELETE SET NULL,
    filename text NOT NULL UNIQUE,
    name text NOT NULL,
    description text ,
    intro boolean  DEFAULT false,
    theme int NOT NULL REFERENCES themes(id) ON DELETE SET NULL,
    uploaded TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS videos_subthemes (
	video_id int NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    subtheme_id int NOT NULL REFERENCES subthemes(id) ON DELETE CASCADE
);
`


const TABLES_FILLING = `
INSERT INTO languages (name) values ('ru'), ('en'), ('be'), ('uk'), ('de'), ('fr'), ('es');

INSERT INTO themes (name) values
    ('media content'), ('music'), ('sports'), ('natural science'),
    ('social science'), ('cooking'), ('painting'), ('craft'), ('languages'), ('photography'), ('design'), ('beauty');

INSERT INTO subthemes (name, theme_id) values
    ('covers design', 1), ('singing', 2), ('instrumental', 2), ('sampling', 2),
    ('cybersport', 3), ('hockey', 3), ('football', 3), ('running', 3),  ('cycling', 3),
    ('programming languages', 4), ('data science', 4), ('theoretical physics', 4), ('math analisys', 4), ('linear algebra', 4),
    ('history', 5), ('philosophy', 5), ('economics', 5),
     ('baking', 6), ('haute cuisine', 6), ('every day meals', 6), ('confectionery making', 6), ('TikTok recipes', 6), ('winemaking', 6), ('brewing', 6), ('cheesemaking', 6),
    ('oil', 7), ('acrylic', 7), ('watercolor', 7), ('gouache', 7),
    ('scrapbooking', 8), ('knitting', 8), ('woodcraft', 8), ('pottery', 8), ('jewellery', 8), ('papercraft', 8),
    ('english', 9),  ('russian', 9), ('german', 9),
    ('nature', 10), ('city', 10),
    ('interior', 11), ('exterior', 11), ('web-design', 11),
    ('make up', 12), ('hairstyling', 12);

`