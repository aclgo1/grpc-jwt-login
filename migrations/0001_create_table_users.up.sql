CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL CHECK (last_name <> ''),
    email VARCHAR(64) UNIQUE NOT NULL CHECK (email <> ''),
    password VARCHAR(250) NOT NULL CHECK (octet_length(password) <> 0),
    role VARCHAR(20) DEFAULT 'client',
    verified VARCHAR(3) DEFAULT 'no',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;