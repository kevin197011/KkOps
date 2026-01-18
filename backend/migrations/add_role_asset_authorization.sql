-- Role Asset Authorization Migration
-- 为角色添加资产授权功能

-- 创建角色资产授权表
CREATE TABLE IF NOT EXISTS role_assets (
    role_id BIGINT NOT NULL,
    asset_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    PRIMARY KEY (role_id, asset_id),
    CONSTRAINT fk_role_assets_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_assets_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
);

-- 创建索引优化查询
CREATE INDEX IF NOT EXISTS idx_role_assets_role_id ON role_assets(role_id);
CREATE INDEX IF NOT EXISTS idx_role_assets_asset_id ON role_assets(asset_id);
