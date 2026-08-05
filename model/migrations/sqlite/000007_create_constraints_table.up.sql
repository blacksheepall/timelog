-- Create constraints table
CREATE TABLE constraints (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    end_reason TEXT,
    punishment_quote TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- Create index for active constraints
CREATE INDEX idx_constraints_is_active ON constraints(is_active);
CREATE INDEX idx_constraints_start_date ON constraints(start_date);
CREATE INDEX idx_constraints_end_date ON constraints(end_date);