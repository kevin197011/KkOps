-- Create operation_tools table
CREATE TABLE IF NOT EXISTS operation_tools (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    icon VARCHAR(255),
    url VARCHAR(512) NOT NULL,
    order_index INTEGER DEFAULT 0,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for operation_tools
CREATE INDEX IF NOT EXISTS idx_operation_tools_category ON operation_tools(category);
CREATE INDEX IF NOT EXISTS idx_operation_tools_enabled ON operation_tools(enabled);
CREATE INDEX IF NOT EXISTS idx_operation_tools_order ON operation_tools(order_index);
