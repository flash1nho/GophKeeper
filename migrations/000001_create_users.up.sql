CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    secret VARCHAR(64) NOT NULL DEFAULT encode(gen_random_bytes(32), 'hex'),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_users_login ON users(login);
CREATE UNIQUE INDEX ids_users_secret ON users(secret);
