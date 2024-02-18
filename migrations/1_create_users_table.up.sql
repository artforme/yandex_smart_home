CREATE TABLE IF NOT EXISTS users (
    primary_key SERIAL PRIMARY KEY,
    user_id TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_user_id ON users(user_id)