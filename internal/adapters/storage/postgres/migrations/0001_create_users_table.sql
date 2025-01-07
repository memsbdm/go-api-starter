-- +goose Up
CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   username VARCHAR(255) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;
