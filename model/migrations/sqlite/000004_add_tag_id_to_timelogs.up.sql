-- 修改timelogs表，添加tag_id外键并保留原tag列用于数据迁移

-- 1. 添加tag_id列
ALTER TABLE timelogs ADD COLUMN tag_id INTEGER;

-- 2. 根据现有tag文本内容设置tag_id
-- 为现有数据匹配对应的tag_id
UPDATE timelogs SET tag_id = (
    SELECT t.id FROM tags t WHERE t.name = '工作'
) WHERE tag LIKE '%工作%' OR tag LIKE '%work%';

UPDATE timelogs SET tag_id = (
    SELECT t.id FROM tags t WHERE t.name = '学习'
) WHERE tag LIKE '%学习%' OR tag LIKE '%study%';

UPDATE timelogs SET tag_id = (
    SELECT t.id FROM tags t WHERE t.name = '会议'
) WHERE tag LIKE '%会议%' OR tag LIKE '%meeting%';

UPDATE timelogs SET tag_id = (
    SELECT t.id FROM tags t WHERE t.name = '开发'
) WHERE tag LIKE '%开发%' OR tag LIKE '%dev%' OR tag LIKE '%EQW%';

-- 其他未匹配的设置为"其他"
UPDATE timelogs SET tag_id = (
    SELECT t.id FROM tags t WHERE t.name = '其他'
) WHERE tag_id IS NULL;

-- 3. 设置外键约束（SQLite中需要重建表）
-- 创建新表
CREATE TABLE timelogs_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER DEFAULT 1,
    start_time DATETIME NOT NULL,
    end_time DATETIME NULL,
    tag_id INTEGER NOT NULL,
    remark TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

-- 复制数据到新表
INSERT INTO timelogs_new (id, user_id, start_time, end_time, tag_id, remark, created_at, updated_at, deleted_at)
SELECT id, user_id, start_time, end_time, tag_id, remark, created_at, updated_at, deleted_at 
FROM timelogs;

-- 删除旧表
DROP TABLE timelogs;

-- 重命名新表
ALTER TABLE timelogs_new RENAME TO timelogs;