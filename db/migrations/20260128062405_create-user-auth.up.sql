CREATE TABLE users (
	id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	email_verified BOOLEAN NOT NULL DEFAULT FALSE,
	image TEXT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TYPE enum_role AS
ENUM('admin', 'manager', 'user');

CREATE TABLE accounts (
	id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
	provider_id TEXT NOT NULL,
	password TEXT NOT NULL,
	user_id TEXT NOT NULL,
	user_role enum_role NOT NULL DEFAULT 'user',
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_accounts_user
    	FOREIGN KEY(user_id)
		REFERENCES users(id)
		ON DELETE CASCADE
);

CREATE TABLE sessions (
	id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
	token TEXT NOT NULL,
	expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
	ip_address TEXT,
    user_agent TEXT,
    user_id TEXT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_sessions_user
    	FOREIGN KEY(user_id)
		REFERENCES users(id)
		ON DELETE CASCADE
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_sessions_user_id on sessions(user_id);
CREATE INDEX idx_sessions_expires_at on sessions(expires_at);
CREATE INDEX idx_users_email_verified on users(email_verified);
