-- +goose Up
CREATE TABLE refresh_tokens (
    token TEXT UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    user_id   uuid REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    PRIMARY KEY(token)
);

-- +goose Down 
DROP TABLE refresh_tokens;
