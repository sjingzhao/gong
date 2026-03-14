CREATE TABLE authors
(
    id   INT PRIMARY KEY NOT NULL,
    name CHAR(50)        NOT NULL,
    bio  text
);

CREATE TABLE users
(
    id         INT PRIMARY KEY NOT NULL,
    name       CHAR(50)        NOT NULL,
    email      text,
    created_at INT
);