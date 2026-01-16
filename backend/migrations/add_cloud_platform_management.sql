-- Copyright (c) 2025 kk
--
-- This software is released under the MIT License.
-- https://opensource.org/licenses/MIT

-- Migration: Add Cloud Platform Management
-- Description: Creates cloud_platforms table and migrates existing cloud_platform string values to foreign key relationship
-- Date: 2025-01-27

-- Step 1: Create cloud_platforms table
CREATE TABLE IF NOT EXISTS cloud_platforms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cloud_platforms_name ON cloud_platforms(name) WHERE deleted_at IS NULL;

-- Step 2: Extract unique cloud platform values and create entries
-- Note: This will only insert if the value doesn't already exist
DO $$
DECLARE
    platform_name TEXT;
BEGIN
    FOR platform_name IN 
        SELECT DISTINCT cloud_platform 
        FROM assets 
        WHERE cloud_platform IS NOT NULL 
        AND cloud_platform != ''
        AND cloud_platform NOT IN (SELECT name FROM cloud_platforms WHERE deleted_at IS NULL)
    LOOP
        INSERT INTO cloud_platforms (name, created_at, updated_at)
        VALUES (platform_name, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
    END LOOP;
END $$;

-- Step 3: Add foreign key column to assets table
ALTER TABLE assets ADD COLUMN IF NOT EXISTS cloud_platform_id INTEGER;

-- Step 4: Map existing data from string to foreign key
UPDATE assets
SET cloud_platform_id = (
    SELECT id FROM cloud_platforms 
    WHERE cloud_platforms.name = assets.cloud_platform
    AND cloud_platforms.deleted_at IS NULL
    LIMIT 1
)
WHERE cloud_platform IS NOT NULL 
AND cloud_platform != ''
AND cloud_platform_id IS NULL;

-- Step 5: Add foreign key constraint
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_assets_cloud_platform' 
        AND table_name = 'assets'
    ) THEN
        ALTER TABLE assets 
        ADD CONSTRAINT fk_assets_cloud_platform 
        FOREIGN KEY (cloud_platform_id) REFERENCES cloud_platforms(id);
    END IF;
END $$;

-- Step 6: Drop old cloud_platform column (commented out for safety - uncomment after verification)
-- ALTER TABLE assets DROP COLUMN IF EXISTS cloud_platform;

-- Note: Keep the old column temporarily for rollback safety
-- The column will be dropped in a separate migration after verification

-- Rollback script (if needed):
-- ALTER TABLE assets DROP CONSTRAINT IF EXISTS fk_assets_cloud_platform;
-- ALTER TABLE assets DROP COLUMN IF EXISTS cloud_platform_id;
-- ALTER TABLE assets ADD COLUMN IF NOT EXISTS cloud_platform VARCHAR(50);
-- UPDATE assets SET cloud_platform = (SELECT name FROM cloud_platforms WHERE cloud_platforms.id = assets.cloud_platform_id) WHERE cloud_platform_id IS NOT NULL;
-- DROP TABLE IF EXISTS cloud_platforms;
