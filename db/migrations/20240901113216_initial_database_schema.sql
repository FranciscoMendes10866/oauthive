-- +goose Up
-- +goose StatementBegin
CREATE TABLE user (
    id INTEGER PRIMARY KEY NOT NULL,
    name TEXT,
    email TEXT UNIQUE,
    image TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE account (
    id INTEGER PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    provider TEXT NOT NULL,
    provider_account_id TEXT NOT NULL,
    refresh_token TEXT,
    access_token TEXT,
    expires_at INTEGER,
    token_type TEXT,
    id_token TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_account_provider UNIQUE(provider, user_id)
);

CREATE TABLE session (
    id INTEGER PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    account_id INTEGER NOT NULL,
    expires_at INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE
    FOREIGN KEY(account_id) REFERENCES account(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS account;
DROP TABLE IF EXISTS user;
-- +goose StatementEnd
