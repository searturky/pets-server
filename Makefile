.PHONY: build-test build-prod push-test push-prod env-up env-down env-clean k8s-deploy k8s-delete k8s-status migrate-up migrate-down

# 变量
GO := go
DOCKER := docker
DOCKER_COMPOSE := $(if $(shell command -v docker-compose 2>/dev/null),docker-compose,docker compose)
REGISTRY := searturky/pets-server
VERSION := $(shell git describe --tags --always --dirty)


build-test:
	docker build -t $(REGISTRY):$(VERSION) -t $(REGISTRY):latest -f deployments/docker/Dockerfile . --target test

build-prod:
	docker build -t $(REGISTRY):$(VERSION) -t $(REGISTRY):latest -f deployments/docker/Dockerfile . --target prod

# 推送镜像
push-test:
	docker push $(REGISTRY):$(VERSION)
	docker push $(REGISTRY):latest

push-prod:
	docker push $(REGISTRY):$(VERSION)
	docker push $(REGISTRY):latest

# 清理构建产物
clean-images:
	docker rmi $(REGISTRY):$(VERSION)
	docker rmi $(REGISTRY):latest

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

# 下载依赖
install:
	export GOPROXY=https://goproxy.io,https://goproxy.cn,direct && go mod download && go install tool && unset GOPROXY

# 生成wire依赖注入代码
wire:
	go tool wire ./cmd/server/...

# 生成代码（如protobuf等）
generate:
	$(GO) generate ./...

# === 开发环境 ===

# 启动开发环境基础设施
env-up:
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose-basic.yml up -d

# 停止开发环境基础设施
env-down:
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose-basic.yml down

# 停止环境并删除所有卷数据（危险操作，会删除数据库数据）
env-clean:
	@echo -n "⚠️ 确认删除所有卷数据? [y/N] "; \
	read REPLY; \
	if [ "$$REPLY" = "y" ] || [ "$$REPLY" = "Y" ]; then \
		$(DOCKER_COMPOSE) -f deployments/docker/docker-compose-basic.yml down -v; \
		echo "✅ 环境已停止，所有卷数据已删除"; \
	else \
		echo "❌ 操作已取消"; \
	fi

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


