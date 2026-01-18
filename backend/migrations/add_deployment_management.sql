-- Create deployment_modules table
CREATE TABLE IF NOT EXISTS deployment_modules (
    id SERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL,
    environment_id BIGINT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    version_source_url VARCHAR(500),
    deploy_script TEXT,
    script_type VARCHAR(50) DEFAULT 'shell',
    timeout INT DEFAULT 600,
    asset_ids TEXT,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE SET NULL
);

-- Create indexes for deployment_modules
CREATE INDEX IF NOT EXISTS idx_deployment_modules_project_id ON deployment_modules (project_id);
CREATE INDEX IF NOT EXISTS idx_deployment_modules_environment_id ON deployment_modules (environment_id);
CREATE INDEX IF NOT EXISTS idx_deployment_modules_deleted_at ON deployment_modules (deleted_at);

-- Create deployments table
CREATE TABLE IF NOT EXISTS deployments (
    id SERIAL PRIMARY KEY,
    module_id BIGINT NOT NULL,
    version VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending',
    asset_ids TEXT,
    output TEXT,
    error TEXT,
    created_by BIGINT,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (module_id) REFERENCES deployment_modules(id) ON DELETE CASCADE
);

-- Create indexes for deployments
CREATE INDEX IF NOT EXISTS idx_deployments_module_id ON deployments (module_id);
CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments (status);
CREATE INDEX IF NOT EXISTS idx_deployments_deleted_at ON deployments (deleted_at);

-- Add environment_id column if it doesn't exist (for existing tables)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'deployment_modules' AND column_name = 'environment_id'
    ) THEN
        ALTER TABLE deployment_modules ADD COLUMN environment_id BIGINT;
        ALTER TABLE deployment_modules ADD CONSTRAINT fk_deployment_modules_environment 
            FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE SET NULL;
        CREATE INDEX idx_deployment_modules_environment_id ON deployment_modules (environment_id);
    END IF;
END $$;

-- Add template_id column if it doesn't exist (for inheriting from execution templates)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'deployment_modules' AND column_name = 'template_id'
    ) THEN
        ALTER TABLE deployment_modules ADD COLUMN template_id BIGINT;
        ALTER TABLE deployment_modules ADD CONSTRAINT fk_deployment_modules_template 
            FOREIGN KEY (template_id) REFERENCES task_templates(id) ON DELETE SET NULL;
        CREATE INDEX idx_deployment_modules_template_id ON deployment_modules (template_id);
    END IF;
END $$;
