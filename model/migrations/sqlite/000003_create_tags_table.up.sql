-- 创建tags表
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    color VARCHAR(7) DEFAULT '#3B82F6',  -- 默认蓝色
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 插入一些默认标签
INSERT INTO tags (name, color, description) VALUES 
('工作', '#EF4444', '工作相关的时间记录'),
('学习', '#10B981', '学习和培训时间'),
('会议', '#F59E0B', '各种会议时间'),
('开发', '#8B5CF6', '软件开发和编程'),
('休息', '#6B7280', '休息和放松时间'),
('运动', '#F97316', '体育锻炼和健身'),
('其他', '#6366F1', '其他未分类活动');