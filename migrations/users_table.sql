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

INSERT INTO users (id, email, username, password, provider, role) 
VALUES ('00000000-0000-0000-0000-000000000000', 'admin@gmail.com', 'admin', '$2a$10$KIX/Y6wVKTKuB1XJh8RUyO7EN0H3aY5h6F9q/NUz5JavjH1KT0D6C', 'local', 'admin')
ON CONFLICT (email) DO NOTHING;