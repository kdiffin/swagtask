CREATE TABLE
    sessions (
        id TEXT PRIMARY KEY,
        user_id INTEGER REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT now ()
    );