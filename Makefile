# Kronos Makefile - Docker 管理命令

.PHONY: help docker-up docker-down docker-build docker-logs docker-restart docker-clean

help: ## 显示帮助信息
	@echo "Kronos Docker 管理命令"
	@echo ""
	@echo "用法: make [target]"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

docker-up: ## 启动所有服务（生产模式）
	docker-compose up -d --build

docker-up-dev: ## 启动所有服务（开发模式，支持热重载）
	docker-compose -f docker-compose.dev.yml up -d --build

docker-down: ## 停止所有服务
	docker-compose down

docker-down-dev: ## 停止所有服务（开发模式）
	docker-compose -f docker-compose.dev.yml down

docker-down-clean: ## 停止服务并删除数据卷
	docker-compose down -v

docker-build: ## 重新构建所有镜像
	docker-compose build --no-cache

docker-logs: ## 查看所有服务日志
	docker-compose logs -f

docker-logs-backend: ## 查看后端日志
	docker-compose logs -f backend

docker-logs-frontend: ## 查看前端日志
	docker-compose logs -f frontend

docker-logs-db: ## 查看数据库日志
	docker-compose logs -f postgres

docker-restart: ## 重启所有服务
	docker-compose restart

docker-restart-backend: ## 重启后端服务
	docker-compose restart backend

docker-ps: ## 查看服务状态
	docker-compose ps

docker-shell-backend: ## 进入后端容器
	docker-compose exec backend sh

docker-shell-db: ## 进入数据库容器
	docker-compose exec postgres psql -U postgres -d kronos

docker-clean: ## 清理所有容器、镜像和数据卷
	docker-compose down -v --rmi all

docker-health: ## 检查服务健康状态
	@echo "检查后端健康状态..."
	@curl -s http://localhost:8080/health || echo "后端服务未响应"
	@echo ""
	@echo "检查数据库连接..."
	@docker-compose exec -T postgres pg_isready -U postgres || echo "数据库未就绪"

