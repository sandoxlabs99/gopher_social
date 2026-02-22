CREATE TABLE IF NOT EXISTS user_invitations (
    token BYTEA PRIMARY KEY,
    user_id INTEGER NOT NULL
);

