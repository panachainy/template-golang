-- Create auths table
CREATE TABLE auths (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    username VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    role VARCHAR(50) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true
);

-- Create index on deleted_at for soft deletes
CREATE INDEX idx_auths_deleted_at ON auths(deleted_at);

-- Create auth_methods table
CREATE TABLE auth_methods (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    auth_id VARCHAR(36) REFERENCES auths(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    user_id VARCHAR(255),
    name VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    nick_name VARCHAR(255),
    description TEXT,
    avatar_url VARCHAR(500),
    location VARCHAR(255),
    access_token TEXT,
    refresh_token TEXT,
    id_token TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    access_token_secret TEXT
);

-- Create indexes for auth_methods
CREATE INDEX idx_auth_methods_deleted_at ON auth_methods(deleted_at);
CREATE INDEX idx_auth_methods_auth_id ON auth_methods(auth_id);
CREATE INDEX idx_auth_methods_provider ON auth_methods(provider);
