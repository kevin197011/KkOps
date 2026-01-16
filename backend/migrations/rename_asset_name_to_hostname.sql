-- Copyright (c) 2025 kk
--
-- This software is released under the MIT License.
-- https://opensource.org/licenses/MIT

-- Migration: Rename asset.name column to host_name
-- Description: Renames the 'name' column to 'host_name' in the assets table, or creates host_name if name doesn't exist
-- Date: 2025-01-27

-- Check if name column exists and host_name doesn't exist, then rename
DO $$
BEGIN
    -- If name column exists and host_name doesn't, rename it
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'assets' AND column_name = 'name'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'assets' AND column_name = 'host_name'
    ) THEN
        ALTER TABLE assets RENAME COLUMN name TO host_name;
    END IF;

    -- If both name and host_name exist, copy data from name to host_name where host_name is NULL
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'assets' AND column_name = 'name'
    ) AND EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'assets' AND column_name = 'host_name'
    ) THEN
        UPDATE assets SET host_name = name WHERE host_name IS NULL AND name IS NOT NULL;
    END IF;

    -- Ensure host_name is NOT NULL (update any remaining NULL values first)
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'assets' AND column_name = 'host_name'
    ) THEN
        -- Set default value for any NULL host_name
        UPDATE assets SET host_name = 'asset-' || id::text WHERE host_name IS NULL;
        
        -- Now make it NOT NULL if it's not already
        IF EXISTS (
            SELECT 1 FROM information_schema.columns 
            WHERE table_name = 'assets' 
            AND column_name = 'host_name' 
            AND is_nullable = 'YES'
        ) THEN
            ALTER TABLE assets ALTER COLUMN host_name SET NOT NULL;
        END IF;
    END IF;
END $$;

-- Note: Indexes and constraints are automatically preserved by PostgreSQL
-- during column rename operations. No additional migration steps needed.

-- Rollback script (if needed):
-- ALTER TABLE assets RENAME COLUMN host_name TO name;
