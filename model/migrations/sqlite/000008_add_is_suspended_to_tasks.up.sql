-- Add is_suspended field to tasks table
ALTER TABLE tasks ADD COLUMN is_suspended BOOLEAN DEFAULT FALSE;