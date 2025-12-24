-- Migration: Remove SSH Management Tables
-- Description: Remove SSH-related tables (ssh_connections, ssh_keys, ssh_sessions)
--              and SSH-related columns from hosts table (ssh_username, ssh_key_id)
-- Date: 2025-01-XX

BEGIN;

-- Drop SSH-related tables
DROP TABLE IF EXISTS ssh_sessions CASCADE;
DROP TABLE IF EXISTS ssh_connections CASCADE;
DROP TABLE IF EXISTS ssh_keys CASCADE;

-- Remove SSH-related columns from hosts table (if they exist)
ALTER TABLE hosts DROP COLUMN IF EXISTS ssh_username;
ALTER TABLE hosts DROP COLUMN IF EXISTS ssh_key_id;

-- Remove SSH-related foreign key constraints (if they exist)
-- Note: CASCADE in DROP TABLE should handle this, but we'll be explicit
DO $$
BEGIN
    -- Drop foreign key constraint if it exists
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'hosts_ssh_key_id_fkey'
    ) THEN
        ALTER TABLE hosts DROP CONSTRAINT hosts_ssh_key_id_fkey;
    END IF;
END $$;

COMMIT;

