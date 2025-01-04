-- +goose Up
CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   username VARCHAR(255) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;
