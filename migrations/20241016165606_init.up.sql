BEGIN;

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) DEFAULT '',
    nickname VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(200) NOT NULL,
    auth_method INTEGER DEFAULT 1,
    confirmed BOOLEAN DEFAULT FALSE,
    first_language VARCHAR(3) NOT NULL,
    created_at BIGINT NOT NULL
);

COMMIT;