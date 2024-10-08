-- +goose Up
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    artist VARCHAR(255) NOT NULL,
    s_group VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    lyrics TEXT NULL NULL,
    release_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

-- +goose Down
DROP TABLE songs;