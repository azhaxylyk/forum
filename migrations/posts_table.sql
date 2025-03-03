CREATE TABLE IF NOT EXISTS posts (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    content TEXT,
    image_path TEXT, -- Новое поле для хранения пути к изображению
    created_at DATETIME,
    likes INTEGER DEFAULT 0,
    dislikes INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
