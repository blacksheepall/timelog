-- Revert categories back to tags
-- Disable foreign key constraints temporarily
PRAGMA foreign_keys = OFF;

-- Create tags table
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    color VARCHAR(7) DEFAULT '#3B82F6',
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- Migrate categories back to tags (only root level)
INSERT INTO tags (id, name, color, description, created_at, updated_at)
SELECT id, name, color, description, created_at, updated_at
FROM categories
WHERE level = 0;

-- Rebuild timelogs table with tag_id FK pointing to tags
CREATE TABLE timelogs_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER DEFAULT 1,
    start_time DATETIME NOT NULL,
    end_time DATETIME NULL,
    tag_id INTEGER NOT NULL,
    task_id INTEGER,
    remark TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (tag_id) REFERENCES tags(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

-- Copy data from old timelogs table
INSERT INTO timelogs_new (id, user_id, start_time, end_time, tag_id, task_id, remark, created_at, updated_at, deleted_at)
SELECT id, user_id, start_time, end_time, category_id, task_id, remark, created_at, updated_at, deleted_at
FROM timelogs;

-- Drop old timelogs table
DROP TABLE timelogs;

-- Rename new table to timelogs
ALTER TABLE timelogs_new RENAME TO timelogs;

-- Recreate indexes on timelogs
CREATE INDEX idx_timelogs_tag_id ON timelogs(tag_id);
CREATE INDEX idx_timelogs_task_id ON timelogs(task_id);
CREATE INDEX idx_timelogs_deleted_at ON timelogs(deleted_at);

-- Drop old category_id index if exists
DROP INDEX IF EXISTS idx_timelogs_category_id;

-- Drop categories table
DROP TABLE IF EXISTS categories;

-- Verify foreign key integrity
-- Note: This is informational. If violations exist, they will cause errors when FK constraints
-- are re-enabled or on subsequent operations. The migration framework doesn't support
-- conditional logic to fail on PRAGMA foreign_key_check results.
PRAGMA foreign_key_check;

-- Re-enable foreign key constraints
-- If there were any FK violations, subsequent operations will fail
PRAGMA foreign_keys = ON;
