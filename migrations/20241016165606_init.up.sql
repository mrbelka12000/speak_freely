BEGIN;

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) DEFAULT '',
    nickname VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(200) NOT NULL,
    auth_method INTEGER DEFAULT 1
);

COMMIT;