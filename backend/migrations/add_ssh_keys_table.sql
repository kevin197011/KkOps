-- add_ssh_keys_table.sql
-- Description: Create ssh_keys table and add ssh_key_id column to hosts table

-- Create ssh_keys table
CREATE TABLE IF NOT EXISTS ssh_keys (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    key_type VARCHAR(20) NOT NULL,
    private_key TEXT NOT NULL,
    public_key TEXT,
    fingerprint VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_ssh_keys_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_ssh_keys_user_id ON ssh_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_ssh_keys_fingerprint ON ssh_keys(fingerprint);
CREATE INDEX IF NOT EXISTS idx_ssh_keys_deleted_at ON ssh_keys(deleted_at);

-- Add ssh_key_id column to hosts table (optional, nullable)
ALTER TABLE hosts ADD COLUMN IF NOT EXISTS ssh_key_id BIGINT;

-- Create index for ssh_key_id
CREATE INDEX IF NOT EXISTS idx_hosts_ssh_key_id ON hosts(ssh_key_id);

-- Add foreign key constraint from hosts.ssh_key_id to ssh_keys.id (optional)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'fk_hosts_ssh_key'
    ) THEN
        ALTER TABLE hosts 
        ADD CONSTRAINT fk_hosts_ssh_key 
        FOREIGN KEY (ssh_key_id) REFERENCES ssh_keys(id) ON DELETE SET NULL;
    END IF;
END $$;

