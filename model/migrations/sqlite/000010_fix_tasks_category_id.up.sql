-- Fix tasks table to use category_id instead of tag_id
PRAGMA foreign_keys = OFF;

-- Create new tasks table with correct column name
CREATE TABLE tasks_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    category_id INTEGER NOT NULL,
    due_date DATETIME NOT NULL,
    estimated_minutes INTEGER NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at DATETIME,
    is_suspended BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Copy data from old table
INSERT INTO tasks_new (
    id, title, description, category_id, due_date, estimated_minutes,
    is_completed, completed_at, is_suspended, created_at, updated_at, deleted_at
)
SELECT 
    id, title, description, tag_id, due_date, estimated_minutes,
    is_completed, completed_at, is_suspended, created_at, updated_at, deleted_at
FROM tasks;

-- Drop old table
DROP TABLE tasks;

-- Rename new table
ALTER TABLE tasks_new RENAME TO tasks;

-- Recreate indexes
CREATE INDEX idx_tasks_category_id ON tasks(category_id);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_tasks_is_completed ON tasks(is_completed);
CREATE INDEX idx_tasks_completed_at ON tasks(completed_at);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);

-- Drop old tag_id index if exists
DROP INDEX IF EXISTS idx_tasks_tag_id;

PRAGMA foreign_keys = ON;
