-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id SMALLINT PRIMARY KEY,
    name VARCHAR(10) NOT NULL
);

INSERT INTO roles (id, name) VALUES (0, 'admin');
INSERT INTO roles (id, name) VALUES (1, 'user');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
