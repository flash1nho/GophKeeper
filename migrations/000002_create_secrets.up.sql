CREATE TYPE secrets_type AS ENUM ('Card', 'File', 'Text', 'Cred');

CREATE TABLE secrets (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER NOT NULL,
    encrypted_data BYTEA NOT NULL,
    type secrets_type NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX idx_secrets_user_id_and_type ON secrets(user_id, type);
