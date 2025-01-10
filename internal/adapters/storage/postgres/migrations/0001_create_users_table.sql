-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL,
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE users;
