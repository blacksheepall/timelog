-- Drop indexes for tasks table
DROP INDEX IF EXISTS idx_tasks_deleted_at;
DROP INDEX IF EXISTS idx_tasks_completed_at;
DROP INDEX IF EXISTS idx_tasks_is_completed;
DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_tasks_tag_id;

-- Drop tasks table
DROP TABLE IF EXISTS tasks;