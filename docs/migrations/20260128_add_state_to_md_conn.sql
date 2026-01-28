-- 为 md_conn 表添加 state 字段和 remark 字段（如果不存在）
-- 执行日期: 2026-01-28

-- 添加 state 字段
ALTER TABLE md_conn ADD COLUMN IF NOT EXISTS state INT NOT NULL DEFAULT 0 COMMENT '连接状态: 0=未检测, 1=有效';

-- 添加 remark 字段（如果之前没有）
ALTER TABLE md_conn ADD COLUMN IF NOT EXISTS remark VARCHAR(512) DEFAULT '' COMMENT '备注';
