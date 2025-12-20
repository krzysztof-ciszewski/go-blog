CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT,
    provider VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    provider_user_id VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    access_token TEXT NOT NULL,
    access_token_secret TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    id_token TEXT NOT NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_provider_user_id ON users(provider, provider_user_id);
