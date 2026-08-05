-- 修改 end_time 字段，允许为空
-- SQLite不支持直接修改列，需要重建表

-- 1. 创建新表
CREATE TABLE timelogs_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER DEFAULT 1,
    start_time DATETIME NOT NULL,
    end_time DATETIME NULL,  -- 允许为空
    tag TEXT,
    remark TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 2. 复制数据
INSERT INTO timelogs_new (id, user_id, start_time, end_time, tag, remark, created_at, updated_at, deleted_at)
SELECT id, user_id, start_time, end_time, tag, remark, created_at, updated_at, deleted_at 
FROM timelogs;

-- 3. 删除旧表
DROP TABLE timelogs;

-- 4. 重命名新表
ALTER TABLE timelogs_new RENAME TO timelogs;