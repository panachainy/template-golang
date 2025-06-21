-- Create enum type for Provider
CREATE TYPE provider_enum AS ENUM ('local', 'firebase', 'line');

-- Create enum type for Role
CREATE TYPE role_enum AS ENUM ('user', 'admin', 'moderator');

-- Create auths table
CREATE TABLE auths (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    username VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    role role_enum NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true
);

-- Create auth_methods table
CREATE TABLE auth_methods (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    auth_id VARCHAR(36) NOT NULL,
    provider provider_enum NOT NULL,
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
    access_token_secret TEXT,
    CONSTRAINT fk_auth_methods_auth_id FOREIGN KEY (auth_id) REFERENCES auths(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_auths_deleted_at ON auths(deleted_at);
CREATE INDEX idx_auth_methods_deleted_at ON auth_methods(deleted_at);
CREATE INDEX idx_auth_methods_auth_id ON auth_methods(auth_id);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_auths_updated_at BEFORE UPDATE ON auths FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_auth_methods_updated_at BEFORE UPDATE ON auth_methods FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
