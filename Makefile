.PHONY: env-up env-down up down k8s-deploy k8s-delete k8s-status migrate-up migrate-down

# 变量
GO := go
DOCKER := docker
DOCKER_COMPOSE := $(if $(shell command -v docker-compose 2>/dev/null),docker-compose,docker compose)
REGISTRY := searturky/pets-server
VERSION := $(shell git describe --tags --always --dirty)

# 服务列表
SERVICES := gateway user-service feishu-service logistics-service catering-service booking-service

# 构建所有服务
build:
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		$(GO) build -o bin/$$service ./services/$$service/cmd/...; \
	done

# 清理构建产物
clean:
	rm -rf bin/
	rm -rf vendor/

# 运行测试
test:
	$(GO) test -v ./...

# 运行测试（带覆盖率）
test-coverage:
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# 代码检查
lint:
	golangci-lint run ./...

# 格式化代码
fmt:
	$(GO) fmt ./...

# 下载依赖
deps:
	$(GO) mod download
	$(GO) mod tidy

# 生成代码（如protobuf等）
generate:
	$(GO) generate ./...

# === Docker 相关 ===

# 构建所有Docker镜像
docker-build:
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..."; \
		$(DOCKER) build \
			--build-arg SERVICE_NAME=$$service \
			--build-arg SERVICE_PATH=./services/$$service/cmd \
			-t $(REGISTRY)/$$service:$(VERSION) \
			-t $(REGISTRY)/$$service:latest \
			-f deployments/docker/Dockerfile.base .; \
	done

# 推送所有Docker镜像
docker-push:
	@for service in $(SERVICES); do \
		echo "Pushing Docker image for $$service..."; \
		$(DOCKER) push $(REGISTRY)/$$service:$(VERSION); \
		$(DOCKER) push $(REGISTRY)/$$service:latest; \
	done

# === 开发环境 ===

# 启动开发环境基础设施
env-up:
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.dev.yml up -d

# 停止开发环境基础设施
env-down:
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.dev.yml down

# 启动完整环境
up:
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yml up -d

# 停止完整环境
down:
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yml down

# === Kubernetes 相关 ===

# 部署到Kubernetes
k8s-deploy:
	kubectl apply -f deployments/k8s/namespace.yaml
	kubectl apply -f deployments/k8s/configmap.yaml
	kubectl apply -f deployments/k8s/secrets.yaml
	kubectl apply -f deployments/k8s/

# 删除Kubernetes部署
k8s-delete:
	kubectl delete -f deployments/k8s/

# 查看Kubernetes状态
k8s-status:
	kubectl get all -n enterprise-platform

# === 数据库迁移 ===

# 运行数据库迁移
migrate-up:
	@echo "Running database migrations..."
	# TODO: 添加迁移命令

# 回滚数据库迁移
migrate-down:
	@echo "Rolling back database migrations..."
	# TODO: 添加回滚命令


