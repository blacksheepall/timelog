-- 回滚迁移：恢复 end_time 为 NOT NULL

-- 1. 创建旧表结构
CREATE TABLE timelogs_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER DEFAULT 1,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    tag TEXT,
    remark TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 2. 复制数据（只复制有end_time的记录）
INSERT INTO timelogs_old (id, user_id, start_time, end_time, tag, remark, created_at, updated_at, deleted_at)
SELECT id, user_id, start_time, 
       COALESCE(end_time, start_time) as end_time,  -- 用start_time替换NULL的end_time
       tag, remark, created_at, updated_at, deleted_at 
FROM timelogs
WHERE end_time IS NOT NULL OR end_time IS NULL;

-- 3. 删除当前表
DROP TABLE timelogs;

-- 4. 重命名回原表名
ALTER TABLE timelogs_old RENAME TO timelogs;