CREATE TABLE IF NOT EXISTS moderation_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL, -- 'delete_post' или 'become_moderator'
    reason TEXT,
    post_id INT, -- Для запросов удаления постов
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
