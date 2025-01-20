CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE,
    username TEXT UNIQUE,
    password TEXT,
    session_token TEXT,
    provider TEXT
);
