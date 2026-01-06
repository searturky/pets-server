#!/bin/bash

# 宠物游戏服务器启动脚本
# 自动启动所有依赖服务（PostgreSQL, Redis, NATS）

set -e

echo "🚀 启动宠物游戏服务器..."
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker 未安装，请先安装 Docker${NC}"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo -e "${RED}❌ Docker Compose 未安装，请先安装 Docker Compose${NC}"
    exit 1
fi

# 进入项目根目录
cd "$(dirname "$0")/.."

echo "📦 启动依赖服务（PostgreSQL, Redis, NATS）..."
docker-compose up -d postgres redis nats

echo ""
echo "⏳ 等待服务就绪..."

# 等待 PostgreSQL
echo -n "  - PostgreSQL: "
for i in {1..30}; do
    if docker-compose exec -T postgres pg_isready -U postgres &> /dev/null; then
        echo -e "${GREEN}✓${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}✗ (超时)${NC}"
        exit 1
    fi
    sleep 1
done

# 等待 Redis
echo -n "  - Redis: "
for i in {1..30}; do
    if docker-compose exec -T redis redis-cli ping &> /dev/null; then
        echo -e "${GREEN}✓${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}✗ (超时)${NC}"
        exit 1
    fi
    sleep 1
done

# 等待 NATS
echo -n "  - NATS: "
for i in {1..30}; do
    if curl -s http://localhost:8222/healthz &> /dev/null; then
        echo -e "${GREEN}✓${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${YELLOW}! (可选服务，继续启动)${NC}"
        break
    fi
    sleep 1
done

echo ""
echo "🔧 运行数据库迁移..."
go run cmd/server/main.go migrate 2>/dev/null || echo "  (迁移已执行或不需要)"

echo ""
echo "✨ 启动应用服务器..."
go run cmd/server/main.go

# 脚本退出时清理
trap "echo ''; echo '🛑 停止服务...'; docker-compose down" EXIT

