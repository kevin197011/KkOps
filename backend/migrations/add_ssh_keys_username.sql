-- add_ssh_keys_username.sql
-- Description: Add username column to ssh_keys table

-- Add username column to ssh_keys table (optional, nullable)
ALTER TABLE ssh_keys ADD COLUMN IF NOT EXISTS username VARCHAR(100);

