// Package websocket WebSocket 连接管理
package websocket

import (
	"log"
	"sync"
)

// Hub WebSocket 连接中心
// 管理所有客户端连接
type Hub struct {
	// 已注册的客户端
	clients map[int64]*Client

	// 注册请求通道
	register chan *Client

	// 注销请求通道
	unregister chan *Client

	// 广播消息通道
	broadcast chan *Message

	// 互斥锁
	mu sync.RWMutex
}

// Message 消息结构
type Message struct {
	UserID  int64       `json:"userId,omitempty"` // 目标用户ID（0表示广播）
	Type    string      `json:"type"`             // 消息类型
	Payload interface{} `json:"payload"`          // 消息内容
}

// NewHub 创建 Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

// Run 启动 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("Client connected: %d", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: %d", client.UserID)

		case message := <-h.broadcast:
			if message.UserID == 0 {
				// 广播给所有客户端
				h.mu.RLock()
				for _, client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client.UserID)
					}
				}
				h.mu.RUnlock()
			} else {
				// 发送给特定用户
				h.mu.RLock()
				if client, ok := h.clients[message.UserID]; ok {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, message.UserID)
					}
				}
				h.mu.RUnlock()
			}
		}
	}
}

// SendToUser 发送消息给特定用户
func (h *Hub) SendToUser(userID int64, msgType string, payload interface{}) {
	h.broadcast <- &Message{
		UserID:  userID,
		Type:    msgType,
		Payload: payload,
	}
}

// Broadcast 广播消息给所有用户
func (h *Hub) Broadcast(msgType string, payload interface{}) {
	h.broadcast <- &Message{
		UserID:  0,
		Type:    msgType,
		Payload: payload,
	}
}

// IsOnline 检查用户是否在线
func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// OnlineCount 获取在线用户数
func (h *Hub) OnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

