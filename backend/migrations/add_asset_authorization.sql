-- Copyright (c) 2025 kk
--
-- This software is released under the MIT License.
-- https://opensource.org/licenses/MIT

-- 资产授权权限重构迁移脚本
-- 1. 创建用户资产授权关联表
-- 2. 添加角色管理员标识字段

-- PostgreSQL 兼容语法

-- 创建用户资产授权表（如果不存在）
-- 注意：GORM AutoMigrate 会自动创建这个表，这里作为备用迁移
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_tables WHERE tablename = 'user_assets') THEN
        CREATE TABLE user_assets (
            user_id BIGINT NOT NULL,
            asset_id BIGINT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            created_by BIGINT,
            PRIMARY KEY (user_id, asset_id),
            CONSTRAINT fk_user_assets_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
            CONSTRAINT fk_user_assets_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
        );
        
        CREATE INDEX IF NOT EXISTS idx_user_assets_user_id ON user_assets(user_id);
        CREATE INDEX IF NOT EXISTS idx_user_assets_asset_id ON user_assets(asset_id);
    END IF;
END $$;

-- 设置默认 admin 角色为管理员（如果 is_admin 字段存在）
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'roles' AND column_name = 'is_admin') THEN
        UPDATE roles SET is_admin = TRUE WHERE name = 'admin' AND is_admin = FALSE;
    END IF;
END $$;
