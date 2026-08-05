-- Revert tasks table back to tag_id
PRAGMA foreign_keys = OFF;

CREATE TABLE tasks_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    tag_id INTEGER NOT NULL,
    due_date DATETIME NOT NULL,
    estimated_minutes INTEGER NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at DATETIME,
    is_suspended BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

INSERT INTO tasks_new (
    id, title, description, tag_id, due_date, estimated_minutes,
    is_completed, completed_at, is_suspended, created_at, updated_at, deleted_at
)
SELECT 
    id, title, description, category_id, due_date, estimated_minutes,
    is_completed, completed_at, is_suspended, created_at, updated_at, deleted_at
FROM tasks;

DROP TABLE tasks;
ALTER TABLE tasks_new RENAME TO tasks;

CREATE INDEX idx_tasks_tag_id ON tasks(tag_id);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_tasks_is_completed ON tasks(is_completed);
CREATE INDEX idx_tasks_completed_at ON tasks(completed_at);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);

DROP INDEX IF EXISTS idx_tasks_category_id;

PRAGMA foreign_keys = ON;
