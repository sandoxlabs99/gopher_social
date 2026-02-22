CREATE TABLE IF NOT EXISTS followers (
    user_id INTEGER NOT NULL,
    follower_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT prevent_similar_ids CHECK (user_id <> follower_id),
    CONSTRAINT followers_pk PRIMARY KEY(user_id, follower_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_follower FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE
);