# Database Migrations

数据库迁移脚本目录。

## Migration Tool

本项目使用 GORM 的 AutoMigrate 功能进行数据库迁移，迁移脚本将在应用启动时自动执行。

## Usage

迁移将在应用启动时自动执行，无需手动运行迁移命令。

## Migration Scripts

迁移脚本将按照以下顺序创建：
1. Phase 2: 用户、角色、权限、API Token 相关表
2. Phase 3: 项目、资产、标签相关表
3. Phase 4: 执行主机、任务相关表
4. Phase 5: SSH 密钥、SSH 会话相关表
