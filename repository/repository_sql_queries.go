package repository

const TABLES_DROPPING = `
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS videos_subthemes;
DROP TABLE IF EXISTS videos;
DROP TABLE IF EXISTS avatars;
DROP TABLE IF EXISTS masters_languages;
DROP TABLE IF EXISTS masters_subthemes;
DROP TABLE IF EXISTS masters;
DROP TABLE IF EXISTS students;
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
    type int NOT NULL CHECK (type = 1 OR type = 2),
    created TIMESTAMPTZ NOT NULL
    CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);

CREATE TABLE students (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id int NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    username text NOT NULL UNIQUE ,
    fullname text NOT NULL DEFAULT ''
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

CREATE TABLE masters_subthemes (
	master_id int NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    subtheme_id int NOT NULL REFERENCES subthemes(id) ON DELETE CASCADE
);

CREATE TABLE masters_languages (
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

CREATE TABLE videos_subthemes (
	video_id int NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    subtheme_id int NOT NULL REFERENCES subthemes(id) ON DELETE CASCADE
);

CREATE TABLE chats (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id_master int NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    user_id_student int NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    type int NOT NULL,
    created TIMESTAMPTZ NOT NULL
);


CREATE TABLE messages (
    id SERIAL NOT NULL PRIMARY KEY,
    info boolean DEFAULT false,
    user_id int REFERENCES users(id) ON DELETE SET NULL,
    chat_id int NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    text text NOT NULL,
    created TIMESTAMPTZ NOT NULL
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

INSERT INTO users (email, password, type, created) values
                                                          ('spiro@mail.ru', '1234', 1, '2020-10-03T13:54:00+00:00'),
                                                          ('sportsman@mail.ru', '123', 1, '2020-10-03T14:54:00+00:00'),
                                                          ('interestinguser@mail.ru', '123', 1, '2020-10-13T13:55:00+00:00'),
                                                          ('roy_aaa@gmail.com', '123', 1, '2020-10-14T11:15:00+00:00'),
                                                          ('cookmaster@gmail.com', '123', 1, '2020-10-15T15:46:00+00:00'),
                                                          ('musefan@gmail.com', '123', 1, '2020-10-10T12:30:00+00:00');


INSERT INTO masters (user_id, username, fullname, theme, description, qualification, education_format, avg_price) values
                                        (1, 'reyamusic', 'Reya Fountain', 2, 'Hi! I''m a flutist', 1, 2, 0),
                                        (2, 'alexsportsman', 'Alex Baranoff', 3, 'Hi! I''m a sportsman', 2, 2, 0),
                                        (3, 'interesting', 'Mary Cool', 6, '', 1, 2, 0),
                                        (4, 'royanderson', 'Roy Anderson', 6, '', 2, 1, 0),
                                        (5, 'cookmaster', 'Jacob Terrier', 6, '', 2, 2, 0),
                                        (6, 'musefan', 'Ali Torcher', 2, 'I love Muse', 2, 1, 0);
                                                                                                                         ;
;
INSERT INTO masters_subthemes (master_id, subtheme_id) values (1, 3);
INSERT INTO masters_languages (master_id, language_id) values (1, 1), (1, 2);

INSERT INTO masters_subthemes (master_id, subtheme_id) values (2, 5), (2, 6);
INSERT INTO masters_languages (master_id, language_id) values (2, 1), (2, 2);

INSERT INTO masters_subthemes (master_id, subtheme_id) values (3, 18), (3, 20);
INSERT INTO masters_languages (master_id, language_id) values (3, 1), (3, 3);

INSERT INTO masters_subthemes (master_id, subtheme_id) values (4, 18), (4, 21);
INSERT INTO masters_languages (master_id, language_id) values (4, 1), (4, 2);

INSERT INTO masters_subthemes (master_id, subtheme_id) values (5, 21), (5, 22);
INSERT INTO masters_languages (master_id, language_id) values (5, 1), (5, 4);

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


INSERT INTO users (email, password, type, created) values
                                                          ('hwllo1@mail.ru', '1234', 2, '2020-10-23T13:54:00+00:00'),
                                                          ('usertest1@mail.ru', '123', 2, '2020-10-13T14:54:00+00:00'),
                                                          ('studentbest@mail.ru', '123', 2, '2020-10-23T13:55:00+00:00'),
                                                          ('musiclove777@gmail.com', '123', 2, '2020-10-06T11:15:00+00:00'),
                                                          ('suova@gmail.com', '123', 2, '2020-10-12T15:46:00+00:00'),
                                                          ('whoami@gmail.com', '123', 2, '2020-10-11T12:30:00+00:00');


INSERT INTO students (user_id, username, fullname) values
                                        (7, 'camillaharris', 'Camilla Harris'),
                                        (8, 'rebeccaaaa', 'Rebecca Cox'),
                                        (9, 'lovetostudy', 'Max Levinson'),
                                        (10, 'musiclover', 'Alexandra Spiridonova'),
                                        (11, 'suovaMail', 'Anastasia Kuznetsova'),
                                        (12, 'siberiacalling', 'Anita Smirnova');

INSERT INTO chats (user_id_master, user_id_student, type, created) values (1, 7, 1, '2020-10-24T12:30:00+00:00'),
                                                                (2, 8, 1, '2020-10-24T12:31:00+00:00'),
                                                                (1, 9, 1, '2020-10-24T12:32:00+00:00'),
                                                                (2, 10, 1, '2020-10-24T12:33:00+00:00'),
                                                                (1, 11, 1, '2020-10-24T12:34:00+00:00'),
                                                                (2, 7, 1, '2020-10-24T12:35:00+00:00'),
                                                                (3, 8, 1, '2020-10-24T12:36:00+00:00'),
                                                                (3, 12, 1, '2020-10-24T12:37:00+00:00');

INSERT INTO messages (user_id, chat_id, text, created) values (7, 1, 'random text 1', '2020-10-24T12:40:00+00:00'),
                                                              (7, 1, 'random text 2', '2020-10-24T12:41:00+00:00'),
                                                              (7, 1, 'random text 3', '2020-10-24T12:42:00+00:00'),
                                                              (8, 2, 'random text 4', '2020-10-24T12:43:00+00:00'),
                                                              (8, 2, 'random text 5', '2020-10-24T12:44:00+00:00'),
                                                              (9, 3, 'random text 6', '2020-10-24T12:45:00+00:00'),
                                                              (9, 3, 'random text 7', '2020-10-24T12:46:00+00:00'),
                                                              (10, 4, 'random text 8', '2020-10-24T12:47:00+00:00'),
                                                              (11, 5, 'random text 9', '2020-10-24T12:48:00+00:00');


INSERT INTO messages (info, chat_id, text, created) values (true, 1, 'videocall 1', '2020-10-23T12:47:00+00:00'),
                                                          (true, 1, 'videocall 2', '2020-10-24T12:40:01+00:00'),
                                                          (true, 3,  'videocall 3', '2020-10-24T12:48:00+00:00'),
                                                          (true, 4, 'videocall 4', '2020-10-24T12:49:00+00:00');
`
