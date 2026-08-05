-- Rename tags table to categories and add hierarchical support
-- Disable foreign key constraints temporarily
PRAGMA foreign_keys = OFF;

-- Create new categories table
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(7) DEFAULT '#3B82F6',
    description TEXT,
    parent_id INTEGER,
    level INTEGER DEFAULT 0 CHECK (level >= 0 AND level <= 2),
    sort_order INTEGER DEFAULT 0,
    path VARCHAR(255) DEFAULT '/',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE,
    UNIQUE(name, parent_id)
);

-- Migrate existing tags to categories
INSERT INTO categories (id, name, color, description, level, path, created_at, updated_at)
SELECT id, name, color, description, 0, '/', created_at, updated_at
FROM tags;

-- Create indexes on categories
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_level ON categories(level);
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

-- Rebuild timelogs table with category_id FK pointing to categories
CREATE TABLE timelogs_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER DEFAULT 1,
    start_time DATETIME NOT NULL,
    end_time DATETIME NULL,
    category_id INTEGER NOT NULL,
    task_id INTEGER,
    remark TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

-- Copy data from old timelogs table
INSERT INTO timelogs_new (id, user_id, start_time, end_time, category_id, task_id, remark, created_at, updated_at, deleted_at)
SELECT id, user_id, start_time, end_time, tag_id, task_id, remark, created_at, updated_at, deleted_at
FROM timelogs;

-- Drop old timelogs table
DROP TABLE timelogs;

-- Rename new table to timelogs
ALTER TABLE timelogs_new RENAME TO timelogs;

-- Recreate indexes on timelogs
CREATE INDEX idx_timelogs_category_id ON timelogs(category_id);
CREATE INDEX idx_timelogs_task_id ON timelogs(task_id);
CREATE INDEX idx_timelogs_deleted_at ON timelogs(deleted_at);

-- Drop old tag_id index if exists
DROP INDEX IF EXISTS idx_timelogs_tag_id;

-- Drop old tags table
DROP TABLE IF EXISTS tags;

-- Verify foreign key integrity
-- Note: This is informational. If violations exist, they will cause errors when FK constraints
-- are re-enabled or on subsequent operations. The migration framework doesn't support
-- conditional logic to fail on PRAGMA foreign_key_check results.
PRAGMA foreign_key_check;

-- Re-enable foreign key constraints
-- If there were any FK violations, subsequent operations will fail
PRAGMA foreign_keys = ON;
