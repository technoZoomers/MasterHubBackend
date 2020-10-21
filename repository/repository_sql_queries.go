package repository

const TABLES_DROPPING = `
DROP TABLE IF EXISTS videos_subthemes;
DROP TABLE IF EXISTS videos;
DROP TABLE IF EXISTS avatars;
DROP TABLE IF EXISTS masters_languages;
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
    fullname text NOT NULL DEFAULT '',
    theme int REFERENCES themes(id) ON DELETE SET NULL,
    description text DEFAULT '',
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
    extension text NOT NULL, 
    name text DEFAULT 'noname',
    description text DEFAULT '',
    intro boolean DEFAULT false,
	rating int DEFAULT 0,
    theme int REFERENCES themes(id) ON DELETE SET NULL,
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

INSERT INTO users (email, password, type, created) values ('spiro@mail.ru', '1234', 0, '2020-10-03T13:54:00+00:00');
INSERT INTO masters (user_id, username, fullname, theme, description, qualification, education_format, avg_price) values (1, 'reyamusic', 'Reya Fountain', 2, 'Hi! I''m a flutist', 1, 2, 0);
INSERT INTO masters_subthemes (master_id, subtheme_id) values (1, 3);
INSERT INTO masters_languages (master_id, language_id) values (1, 1), (1, 2);

INSERT INTO users (email, password, type, created) values ('sportsman@mail.ru', '123', 0, '2020-10-03T14:54:00+00:00');
INSERT INTO masters (user_id, username, fullname, theme, description, qualification, education_format, avg_price) values (2, 'alexsportsman', 'Alex Baranoff', 3, 'Hi! I''m a sportsman', 2, 2, 0);
INSERT INTO masters_subthemes (master_id, subtheme_id) values (2, 5), (2, 6);
INSERT INTO masters_languages (master_id, language_id) values (2, 1), (2, 2);

INSERT INTO users (email, password, type, created) values ('interestinguser@mail.ru', '123', 0, '2020-10-13T13:55:00+00:00');
INSERT INTO masters (user_id, username, fullname, theme, description, qualification, education_format, avg_price) values (3, 'interesting', 'Mary Cool', 6, '', 1, 2, 0);
INSERT INTO masters_subthemes (master_id, subtheme_id) values (3, 18), (3, 20);
INSERT INTO masters_languages (master_id, language_id) values (3, 1), (3, 3);

INSERT INTO users (email, password, type, created) values ('roy_aaa@gmail.com', '123', 0, '2020-10-14T11:15:00+00:00');
INSERT INTO masters (user_id, username, fullname, theme, description, qualification, education_format, avg_price) values (4, 'royanderson', 'Roy Anderson', 6, '', 2, 1, 0);
INSERT INTO masters_subthemes (master_id, subtheme_id) values (4, 18), (4, 21);
INSERT INTO masters_languages (master_id, language_id) values (4, 1), (4, 2);


INSERT INTO users (email, password, type, created) values ('cookmaster@gmail.com', '123', 0, '2020-10-15T15:46:00+00:00');
INSERT INTO masters (user_id, username, fullname, theme, description, qualification, education_format, avg_price) values (5, 'cookmaster', 'Jacob Terrier', 6, '', 2, 2, 0);
INSERT INTO masters_subthemes (master_id, subtheme_id) values (5, 21), (5, 22);
INSERT INTO masters_languages (master_id, language_id) values (5, 1), (5, 4);

INSERT INTO users (email, password, type, created) values ('musefan@gmail.com', '123', 0, '2020-10-10T12:30:00+00:00');
INSERT INTO masters (user_id, username, fullname, theme, description, qualification, education_format, avg_price) values (6, 'musefan', 'Ali Torcher', 2, 'I love Muse', 2, 1, 0);
INSERT INTO masters_subthemes (master_id, subtheme_id) values (6, 2), (6, 3);
INSERT INTO masters_languages (master_id, language_id) values (6, 1), (6, 2);

INSERT INTO videos (master_id, filename, extension, intro, uploaded, rating, theme) VALUES (1, 'master_1_video_1', 'webm', false, '2020-10-10T12:30:00+00:00', 112, 1) ,
                                                                             (2, 'master_2_video_2', 'webm', false, '2020-10-10T12:31:00+00:00', 10, 1) ,
                                                                             (1, 'master_1_video_3', 'webm', false, '2020-10-10T12:32:00+00:00', 200, 2) ,
                                                                             (1, 'master_1_video_4', 'webm', false, '2020-10-10T12:33:00+00:00', 4, 10) ,
                                                                             (2, 'master_2_video_5', 'webm', false, '2020-10-10T12:34:00+00:00', 0, 2) ,
                                                                             (2, 'master_2_video_6', 'webm', false, '2020-10-10T12:35:00+00:00', 1, 2) ,
                                                                             (3, 'master_3_video_7', 'webm', false, '2020-10-10T12:36:00+00:00', 143, 3) ,
                                                                             (1, 'master_1_intro', 'webm', true, '2020-10-10T12:37:00+00:00', 10, 3) ,
                                                                             (2, 'master_2_intro', 'webm', true, '2020-10-10T12:38:00+00:00', 1, 6);

`
