CREATE TABLE
    IF NOT EXISTS subscriptions (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        artist_id INT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );