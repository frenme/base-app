CREATE TABLE
    IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        username VARCHAR(50) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        nationality VARCHAR(300),
        birth_date DATE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );