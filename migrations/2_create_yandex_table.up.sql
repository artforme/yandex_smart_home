CREATE TABLE IF NOT EXISTS yandex_tokens (
    user_id INTEGER PRIMARY KEY REFERENCES users(primary_key),
    token TEXT NOT NULL
)