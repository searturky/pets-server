#!/bin/bash

# NATS JetStream 测试脚本
# 用于验证消息队列是否工作正常

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🧪 NATS JetStream 测试工具${NC}"
echo ""

# 检查 NATS 是否运行
echo -n "检查 NATS 服务状态... "
if ! curl -s http://localhost:8222/healthz > /dev/null 2>&1; then
    echo -e "${RED}✗${NC}"
    echo "NATS 服务未运行，请先启动："
    echo "  docker-compose up -d nats"
    exit 1
fi
echo -e "${GREEN}✓${NC}"

# 显示 NATS 信息
echo ""
echo "📊 NATS 服务器信息："
curl -s http://localhost:8222/varz | jq '{
    "服务器版本": .version,
    "运行时间": .uptime,
    "连接数": .connections,
    "消息数": .in_msgs,
    "字节数": .in_bytes,
    "JetStream": .jetstream
}' 2>/dev/null || echo "  (需要安装 jq: sudo apt install jq)"

# 显示 Stream 信息
echo ""
echo "📦 Stream 列表："
curl -s http://localhost:8222/jsz | jq '.streams[]? | {
    "名称": .name,
    "消息数": .state.messages,
    "字节数": .state.bytes,
    "主题": .config.subjects
}' 2>/dev/null || echo "  (暂无 Stream)"

# 测试菜单
echo ""
echo "选择测试操作："
echo "  1) 发布测试消息"
echo "  2) 订阅消息（监听）"
echo "  3) 查看 Stream 详情"
echo "  4) 清空 Stream"
echo "  5) 退出"
echo ""
read -p "请选择 (1-5): " choice

case $choice in
    1)
        echo ""
        echo "📤 发布测试消息..."
        
        # 检查是否安装了 nats cli
        if command -v nats &> /dev/null; then
            nats pub "game.player.test.level_up" '{"player_id":"test123","old_level":1,"new_level":2}'
            echo -e "${GREEN}✓ 消息已发布${NC}"
        else
            echo -e "${YELLOW}⚠ 未安装 nats cli，使用 curl 方式${NC}"
            echo "建议安装: go install github.com/nats-io/natscli/nats@latest"
        fi
        ;;
        
    2)
        echo ""
        echo "👂 开始监听消息（按 Ctrl+C 停止）..."
        
        if command -v nats &> /dev/null; then
            nats sub "game.>"
        else
            echo -e "${RED}需要安装 nats cli${NC}"
            echo "运行: go install github.com/nats-io/natscli/nats@latest"
        fi
        ;;
        
    3)
        echo ""
        echo "📋 Stream 详细信息："
        
        if command -v nats &> /dev/null; then
            nats stream info game-events
        else
            curl -s http://localhost:8222/jsz | jq '.' 2>/dev/null || echo "需要安装 jq"
        fi
        ;;
        
    4)
        echo ""
        read -p "确定要清空 Stream 吗？(y/N): " confirm
        if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
            if command -v nats &> /dev/null; then
                nats stream purge game-events -f
                echo -e "${GREEN}✓ Stream 已清空${NC}"
            else
                echo -e "${RED}需要安装 nats cli${NC}"
            fi
        fi
        ;;
        
    5)
        echo "退出"
        exit 0
        ;;
        
    *)
        echo -e "${RED}无效选择${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}✅ 完成${NC}"

