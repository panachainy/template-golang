-- Create cockroaches table
CREATE TABLE cockroaches (
    id SERIAL PRIMARY KEY,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
