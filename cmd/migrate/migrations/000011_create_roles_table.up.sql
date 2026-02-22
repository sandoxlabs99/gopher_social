CREATE TABLE IF NOT EXISTS roles (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, level, description)
VALUES ('user', 1, 'a user can create posts and comments');

INSERT INTO roles (name, level, description)
VALUES ('moderator', 2, 'a moderator can take action like update/delete other users posts');

INSERT INTO roles (name, level, description)
VALUES ('admin', 3, 'an admin take action like update/delete other users posts and comments');