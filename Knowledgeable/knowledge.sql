PRAGMA foreign_keys = ON;


-- USERS TABLE
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- PAGES TABLE
DROP TABLE IF EXISTS pages;

CREATE TABLE pages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    language TEXT NOT NULL CHECK(language IN ('en', 'da')) DEFAULT 'en',
    content TEXT NOT NULL,
    last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster search
CREATE INDEX idx_pages_language ON pages(language);
CREATE INDEX idx_pages_title ON pages(title);
