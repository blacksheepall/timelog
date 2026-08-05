-- 回滚：恢复到原来的tag文本列

-- 1. 创建带tag文本列的表
CREATE TABLE timelogs_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER DEFAULT 1,
    start_time DATETIME NOT NULL,
    end_time DATETIME NULL,
    tag TEXT,
    remark TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 2. 复制数据并将tag_id转换回tag名称
INSERT INTO timelogs_old (id, user_id, start_time, end_time, tag, remark, created_at, updated_at, deleted_at)
SELECT 
    t.id, 
    t.user_id, 
    t.start_time, 
    t.end_time,
    COALESCE(tags.name, '') as tag,
    t.remark, 
    t.created_at, 
    t.updated_at, 
    t.deleted_at 
FROM timelogs t
LEFT JOIN tags ON t.tag_id = tags.id;

-- 3. 删除当前表
DROP TABLE timelogs;

-- 4. 重命名回原表名
ALTER TABLE timelogs_old RENAME TO timelogs;