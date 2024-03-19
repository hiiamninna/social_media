CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INT NOT NULL,
    post VARCHAR(512) NOT NULL,
    tags VARCHAR[] NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INT NOT NULL,
    post_id INT NOT NULL,
    comment VARCHAR(512) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);