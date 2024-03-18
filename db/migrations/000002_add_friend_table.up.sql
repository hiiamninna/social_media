CREATE TABLE IF NOT EXISTS friends (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INT,
    added_by INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);