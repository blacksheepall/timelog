PRAGMA foreign_keys = OFF;

CREATE TABLE categories_new (
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

INSERT INTO categories_new (
    id, name, color, description, parent_id, level, sort_order, path,
    created_at, updated_at, deleted_at
)
SELECT
    id, name, color, description, parent_id, level - 1, sort_order, path,
    created_at, updated_at, deleted_at
FROM categories;

DROP TABLE categories;

ALTER TABLE categories_new RENAME TO categories;

CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_level ON categories(level);
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

PRAGMA foreign_keys = ON;
