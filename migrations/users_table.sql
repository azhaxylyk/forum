CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
	password TEXT,
    session_token VARCHAR(255),
    provider VARCHAR(50) NOT NULL, -- Для Google/GitHub
    role VARCHAR(20) NOT NULL DEFAULT 'user', -- 'user', 'moderator', 'admin'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
