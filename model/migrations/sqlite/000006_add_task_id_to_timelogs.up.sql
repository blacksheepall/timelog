-- Add task_id column to timelogs table
ALTER TABLE timelogs ADD COLUMN task_id INTEGER;

-- Note: SQLite doesn't support adding foreign key constraints to existing tables
-- The foreign key constraint FOREIGN KEY (task_id) REFERENCES tasks(id) 
-- is handled in the application code

-- Create index for task_id column
CREATE INDEX idx_timelogs_task_id ON timelogs(task_id);