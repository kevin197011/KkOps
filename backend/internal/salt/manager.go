package salt

import (
	"log"
	"sync"

	"github.com/kkops/backend/internal/config"
)

// Manager Salt客户端管理器，用于管理Salt客户端的生命周期和热重载
type Manager struct {
	client *Client
	mu     sync.RWMutex
}

var (
	saltManager *Manager
	once        sync.Once
)

// GetManager 获取Salt客户端管理器单例
func GetManager() *Manager {
	once.Do(func() {
		saltManager = &Manager{}
	})
	return saltManager
}

// GetClient 获取当前的Salt客户端（线程安全）
func (m *Manager) GetClient() *Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client
}

// SetClient 设置Salt客户端（线程安全）
func (m *Manager) SetClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.client = client
}

// InitializeClient 初始化Salt客户端
func (m *Manager) InitializeClient(cfg config.SaltConfig) *Client {
	client := NewClient(cfg)
	m.SetClient(client)
	log.Println("Salt API client initialized via Manager")
	return client
}

// ReloadClient 重新加载Salt客户端（用于配置更新后的热重载）
// 注意：此方法需要在调用前确保配置已从数据库重新加载
func (m *Manager) ReloadClient() error {
	cfg := config.AppConfig.Salt
	if cfg.APIURL == "" || cfg.Username == "" {
		// 配置不完整，清空客户端
		m.SetClient(nil)
		log.Println("Salt API configuration incomplete, client cleared")
		return nil
	}

	// 创建新客户端
	client := NewClient(cfg)
	m.SetClient(client)
	log.Println("Salt API client reloaded successfully")
	return nil
}

