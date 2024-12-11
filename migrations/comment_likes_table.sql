CREATE TABLE IF NOT EXISTS comment_likes (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    comment_id TEXT,
    is_like BOOLEAN,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id),
    UNIQUE (user_id, comment_id)
);