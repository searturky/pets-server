# NATS JetStream 快速启动指南

本指南帮助你快速从 RabbitMQ 迁移到 NATS JetStream。

## ✅ 已完成的改动

### 1. 代码更新
- ✅ `publisher.go` - 替换为 NATS JetStream 实现
- ✅ `main.go` - 更新消息队列初始化逻辑
- ✅ `config.go` - 更新配置结构
- ✅ `config.yaml` - 更新配置文件

### 2. 三种实现方案
1. **NATSPublisher** - NATS JetStream（主要方案，推荐）
2. **RedisStreamPublisher** - Redis Stream（轻量级备用方案）
3. **NoopPublisher** - 空实现（开发/测试）

### 3. 自动降级策略
```
NATS 连接失败 → 尝试 Redis Stream → 使用 Noop Publisher
```

## 🚀 快速开始

### 步骤 1: 启动 NATS Server（推荐 Docker）

```bash
# 启动 NATS with JetStream
docker run -d --name nats-jetstream \
  -p 4222:4222 \
  -p 8222:8222 \
  nats:latest \
  -js \
  -m 8222

# 查看日志
docker logs -f nats-jetstream

# 查看 Web 监控页面
open http://localhost:8222
```

### 步骤 2: 验证 NATS 连接

```bash
# 安装 NATS CLI（可选）
go install github.com/nats-io/natscli/nats@latest

# 检查 NATS 状态
nats server check

# 查看 Stream 列表
nats stream list
```

### 步骤 3: 更新配置文件

编辑 `configs/config.yaml`（已自动更新）：

```yaml
# 消息队列配置 (NATS JetStream)
mq:
  nats_url: "nats://localhost:4222"
  stream_name: "game-events"
```

### 步骤 4: 运行项目

```bash
# 安装依赖（如果还没安装）
go mod tidy

# 启动服务
go run cmd/server/main.go
```

成功启动后，你应该看到：
```
NATS JetStream connected successfully (stream: game-events)
Server starting on 0.0.0.0:8080
```

## 📊 验证事件发布

### 方式 1: 通过 NATS CLI 监听

```bash
# 监听所有游戏事件
nats sub "game.>"

# 监听特定类型事件
nats sub "game.player.*.level_up"
```

### 方式 2: 查看 Stream 信息

```bash
# 查看 Stream 详情
nats stream info game-events

# 查看最新消息
nats stream view game-events
```

### 方式 3: 通过 Web 监控

打开浏览器访问：http://localhost:8222/

## 🔄 备用方案：使用 Redis Stream

如果不想使用 NATS，可以用 Redis Stream：

### 1. 清空 NATS 配置

```yaml
# configs/config.yaml
mq:
  nats_url: ""  # 清空或注释掉
```

### 2. 确保 Redis 配置正确

```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

### 3. 重启服务

```bash
go run cmd/server/main.go
```

应该看到：
```
Warning: Failed to connect NATS, trying Redis Stream: ...
Redis Stream connected successfully (stream: game:events)
```

### 4. 验证 Redis Stream

```bash
# 查看 Stream 长度
redis-cli XLEN game:events

# 查看最新消息
redis-cli XREVRANGE game:events + - COUNT 5
```

## 🧪 开发/测试：使用 Noop Publisher

如果暂时不需要消息队列功能：

```yaml
# configs/config.yaml
mq:
  nats_url: ""  # 清空

# 同时停止 Redis（可选）
```

应该看到：
```
MQ not configured, using noop publisher
Using NoopPublisher (events will not be actually published)
```

事件会被打印到日志，但不会真正发布。

## 📝 代码示例

### 发布领域事件

```go
// 定义事件
type PlayerLevelUpEvent struct {
    PlayerID  string
    OldLevel  int
    NewLevel  int
    Timestamp time.Time
}

func (e *PlayerLevelUpEvent) EventName() string {
    return fmt.Sprintf("game.player.%s.level_up", e.PlayerID)
}

// 在业务代码中发布
func (s *PetService) LevelUp(ctx context.Context, petID string) error {
    // ... 业务逻辑 ...
    
    // 发布事件
    event := &PlayerLevelUpEvent{
        PlayerID:  pet.OwnerID,
        OldLevel:  pet.Level - 1,
        NewLevel:  pet.Level,
        Timestamp: time.Now(),
    }
    
    if err := s.eventPublisher.Publish(ctx, event); err != nil {
        log.Printf("Failed to publish event: %v", err)
        // 注意：事件发布失败不影响主流程
    }
    
    return nil
}
```

## 🔧 故障排查

### 问题 1: NATS 连接失败

```
Error: failed to connect to NATS: dial tcp [::1]:4222: connect: connection refused
```

**解决方案**：
1. 确认 NATS Server 已启动：`docker ps | grep nats`
2. 检查端口占用：`lsof -i :4222`
3. 检查配置文件中的 URL

### 问题 2: Stream 已存在错误

```
Error: stream name already in use
```

**解决方案**：这是正常的，代码已经处理了这个错误，不影响使用。

### 问题 3: 找不到 nats.go 包

```
Error: could not import github.com/nats-io/nats.go
```

**解决方案**：
```bash
go get github.com/nats-io/nats.go
go mod tidy
```

## 📈 性能对比

| 指标 | RabbitMQ | NATS JetStream | Redis Stream |
|------|----------|----------------|--------------|
| **吞吐量** | ~2万/s | ~100万/s | ~10万/s |
| **延迟** | 5-10ms | < 1ms | 1-5ms |
| **内存占用** | 高 | 低 | 极低 |
| **部署复杂度** | 高 | 低 | 极低 |
| **Go 支持** | 一般 | 优秀 | 优秀 |

## 🎯 推荐配置

### 开发环境
```yaml
mq:
  nats_url: ""  # 使用 Noop Publisher，无需额外部署
```

### 测试环境
```yaml
mq:
  nats_url: "nats://localhost:4222"  # Docker 本地部署
```

### 生产环境
```yaml
mq:
  nats_url: "nats://nats-cluster:4222"  # 集群部署
  stream_name: "game-events"
```

## 📚 进阶主题

### 1. 订阅事件（Consumer）

```go
// 创建消费者（示例）
js, _ := nc.JetStream()

// 订阅特定事件
sub, err := js.Subscribe("game.player.*.level_up", func(msg *nats.Msg) {
    var event PlayerLevelUpEvent
    json.Unmarshal(msg.Data, &event)
    
    // 处理事件
    log.Printf("Player %s leveled up to %d", event.PlayerID, event.NewLevel)
    
    // 确认消息
    msg.Ack()
})
```

### 2. 事件回放

```bash
# 从头开始读取所有事件
nats consumer add game-events replay \
  --filter "game.player.>" \
  --deliver all \
  --replay instant

# 查看消费进度
nats consumer info game-events replay
```

### 3. 监控告警

```bash
# 查看 Stream 统计
nats stream report

# 查看 Consumer 统计
nats consumer report game-events
```

## 🔗 相关资源

- [NATS 官方文档](https://docs.nats.io/)
- [NATS JetStream](https://docs.nats.io/nats-concepts/jetstream)
- [Go 客户端文档](https://pkg.go.dev/github.com/nats-io/nats.go)
- [Redis Stream](https://redis.io/docs/data-types/streams/)

## ❓ 常见问题

**Q: 需要卸载 RabbitMQ 相关依赖吗？**  
A: 可以，运行 `go mod tidy` 会自动清理未使用的依赖。

**Q: 可以同时使用 NATS 和 Redis Stream 吗？**  
A: 当前实现是自动降级模式，只会使用其中一个。如需同时使用，需要修改代码。

**Q: 消息会丢失吗？**  
A: NATS JetStream 使用文件存储，默认持久化，不会丢失。

**Q: 如何从 RabbitMQ 数据迁移？**  
A: 领域事件是实时产生的，无需迁移历史数据。新系统上线后会自动使用新的消息队列。

---

## 🎉 完成！

现在你的项目已经成功从 RabbitMQ 迁移到 NATS JetStream！

有任何问题请查看 `internal/infrastructure/messaging/README.md` 获取更详细的文档。

