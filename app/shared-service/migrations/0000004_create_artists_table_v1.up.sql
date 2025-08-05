CREATE TABLE
    IF NOT EXISTS artists (
        id SERIAL PRIMARY KEY,
        agency_id INT NOT NULL,
        name VARCHAR(100),
        debut_date DATE,
        image_url VARCHAR(255),
        status VARCHAR(50),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );