ALTER TABLE
    IF EXISTS users
ADD
    COLUMN IF NOT EXISTS role_id INT REFERENCES roles(id);

UPDATE users
SET role_id = (
    SELECT id FROM roles
    WHERE name = 'user'
);

ALTER TABLE
    IF EXISTS users
ALTER COLUMN role_id DROP DEFAULT;

ALTER TABLE
    IF EXISTS users
ALTER COLUMN role_id SET NOT NULL;