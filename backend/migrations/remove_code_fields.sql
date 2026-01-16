-- Copyright (c) 2025 kk
--
-- This software is released under the MIT License.
-- https://opensource.org/licenses/MIT

-- Migration: Remove code fields from all tables
-- Description: Removes the unused 'code' columns and their unique indexes from assets, projects, asset_categories, environments, roles, and permissions tables
-- Date: 2025-01-27

-- Migration 1: Assets
DROP INDEX IF EXISTS idx_assets_code;
ALTER TABLE assets DROP COLUMN IF EXISTS code;

-- Migration 2: Projects
DROP INDEX IF EXISTS idx_projects_code;
ALTER TABLE projects DROP COLUMN IF EXISTS code;
-- Add unique constraint on name if needed
CREATE UNIQUE INDEX IF NOT EXISTS idx_projects_name ON projects(name) WHERE deleted_at IS NULL;

-- Migration 3: Asset Categories
DROP INDEX IF EXISTS idx_asset_categories_code;
ALTER TABLE asset_categories DROP COLUMN IF EXISTS code;

-- Migration 4: Environments
DROP INDEX IF EXISTS idx_environments_code;
ALTER TABLE environments DROP COLUMN IF EXISTS code;
-- Add unique constraint on name if needed
CREATE UNIQUE INDEX IF NOT EXISTS idx_environments_name ON environments(name) WHERE deleted_at IS NULL;

-- Migration 5: Roles
DROP INDEX IF EXISTS idx_roles_code;
ALTER TABLE roles DROP COLUMN IF EXISTS code;
-- Add unique constraint on name if needed
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_name ON roles(name) WHERE deleted_at IS NULL;

-- Migration 6: Permissions
DROP INDEX IF EXISTS idx_permissions_code;
ALTER TABLE permissions DROP COLUMN IF EXISTS code;
-- Add unique constraint on resource + action if needed
CREATE UNIQUE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action) WHERE deleted_at IS NULL;

-- Note: Indexes and constraints are automatically preserved by PostgreSQL
-- during column drop operations. No additional migration steps needed.

-- Rollback script (if needed):
-- ALTER TABLE assets ADD COLUMN code VARCHAR(50);
-- ALTER TABLE projects ADD COLUMN code VARCHAR(50);
-- ALTER TABLE asset_categories ADD COLUMN code VARCHAR(50);
-- ALTER TABLE environments ADD COLUMN code VARCHAR(50);
-- ALTER TABLE roles ADD COLUMN code VARCHAR(50);
-- ALTER TABLE permissions ADD COLUMN code VARCHAR(100);
-- CREATE UNIQUE INDEX idx_assets_code ON assets(code);
-- CREATE UNIQUE INDEX idx_projects_code ON projects(code);
-- CREATE UNIQUE INDEX idx_asset_categories_code ON asset_categories(code);
-- CREATE UNIQUE INDEX idx_environments_code ON environments(code);
-- CREATE UNIQUE INDEX idx_roles_code ON roles(code);
-- CREATE UNIQUE INDEX idx_permissions_code ON permissions(code);
-- DROP INDEX IF EXISTS idx_projects_name;
-- DROP INDEX IF EXISTS idx_environments_name;
-- DROP INDEX IF EXISTS idx_roles_name;
-- DROP INDEX IF EXISTS idx_permissions_resource_action;
