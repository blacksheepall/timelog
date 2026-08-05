-- Create temp passwords table
CREATE TABLE temp_passwords (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    password_hash TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE INDEX idx_temp_passwords_expires_at ON temp_passwords(expires_at);
CREATE INDEX idx_temp_passwords_deleted_at ON temp_passwords(deleted_at);
