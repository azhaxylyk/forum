CREATE TABLE IF NOT EXISTS post_likes (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    post_id TEXT,
    is_like BOOLEAN,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    UNIQUE (user_id, post_id)
);