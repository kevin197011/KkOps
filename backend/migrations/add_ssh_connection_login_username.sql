-- Add login_username to ssh_connection_records (操作用户: KkOps login user; username remains 连线用户)
ALTER TABLE ssh_connection_records
ADD COLUMN IF NOT EXISTS login_username VARCHAR(100) DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_ssh_connection_records_login_username ON ssh_connection_records (login_username);
