package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
	"github.com/kkops/backend/internal/service"
	"github.com/kkops/backend/internal/utils"
	"golang.org/x/crypto/ssh"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: 在生产环境中应该检查 Origin
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// TerminalMessage WebSocket 消息类型
type TerminalMessage struct {
	Type    string `json:"type"`    // "input", "resize", "ping", "password", "username", "auth_method", "key_id"
	Data    string `json:"data"`    // 输入数据或 resize 数据或密码或用户名
	Rows    int    `json:"rows"`    // 终端行数
	Columns int    `json:"columns"` // 终端列数
	KeyID   uint64 `json:"key_id"`  // SSH密钥ID（用于密钥认证）
}

// WebSSHHandler WebSSH 终端处理器
type WebSSHHandler struct {
	hostRepo     repository.HostRepository
	sshKeyService service.SSHKeyService
}

func NewWebSSHHandler(hostRepo repository.HostRepository, sshKeyService service.SSHKeyService) *WebSSHHandler {
	return &WebSSHHandler{
		hostRepo:     hostRepo,
		sshKeyService: sshKeyService,
	}
}

// HandleTerminal 处理 WebSocket 连接并建立 SSH 终端会话
func (h *WebSSHHandler) HandleTerminal(c *gin.Context) {
	// 从路径获取 host_id
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host id"})
		return
	}

	// 从查询参数或 Header 获取 JWT token（WebSocket 升级前）
	var userID uint64

	// 尝试从查询参数获取 token
	token := c.Query("token")
	if token == "" {
		// 尝试从 Authorization header 获取
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
	}

	if token != "" {
		// 解析 token 获取用户信息
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		userID = claims.UserID
	} else {
		// 尝试从上下文获取（如果中间件已设置）
		if uid, exists := c.Get("user_id"); exists {
			userID = uid.(uint64)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}

	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// 获取主机数据
	host, err := h.hostRepo.GetByID(hostID)
	if err != nil {
		log.Printf("Failed to get host: %v", err)
		sendErrorMessage(conn, "Failed to load host configuration")
		return
	}

	// 建立 SSH 客户端连接
	authMethod, username, err := h.requestAuthentication(conn, host, userID)
	if err != nil {
		log.Printf("Failed to get authentication: %v", err)
		sendErrorMessage(conn, fmt.Sprintf("Authentication failed: %v", err))
		return
	}

	log.Printf("Authentication successful for user %d, username: %s", userID, username)

	sshClient, err := h.establishSSHConnection(host, conn, authMethod, username)
	if err != nil {
		log.Printf("Failed to establish SSH connection: %v", err)
		sendErrorMessage(conn, fmt.Sprintf("SSH connection failed: %v", err))
		return
	}
	defer sshClient.Close()

	log.Printf("SSH connection established to %s:%d", host.IPAddress, host.SSHPort)

	// 创建 SSH 会话
	session, err := sshClient.NewSession()
	if err != nil {
		log.Printf("Failed to create SSH session: %v", err)
		sendErrorMessage(conn, "Failed to create SSH session")
		return
	}
	defer session.Close()

	// 请求 PTY（伪终端）
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // 启用回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	}

	// 默认终端大小
	rows, cols := 40, 160
	if err := session.RequestPty("xterm-256color", rows, cols, modes); err != nil {
		log.Printf("Failed to request PTY: %v", err)
		sendErrorMessage(conn, "Failed to request PTY")
		return
	}

	// 设置 stdin/stdout/stderr
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Printf("Failed to get stdin pipe: %v", err)
		sendErrorMessage(conn, "Failed to setup stdin")
		return
	}
	defer stdin.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Printf("Failed to get stdout pipe: %v", err)
		sendErrorMessage(conn, "Failed to setup stdout")
		return
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Printf("Failed to get stderr pipe: %v", err)
		sendErrorMessage(conn, "Failed to setup stderr")
		return
	}

	// 启动 shell
	if err := session.Shell(); err != nil {
		log.Printf("Failed to start shell: %v", err)
		sendErrorMessage(conn, "Failed to start shell")
		return
	}

	// 发送连接成功消息（包含初始终端大小信息，方便前端确认）
	conn.WriteJSON(TerminalMessage{
		Type:    "connected",
		Data:    "SSH connection established",
		Rows:    rows,
		Columns: cols,
	})

	// 使用 WaitGroup 管理 goroutines
	var wg sync.WaitGroup
	wg.Add(4)

	// 最后活动时间跟踪（用于检测空闲）
	// 包括：用户输入、SSH输出、ping消息
	lastActivity := time.Now()
	var lastActivityMutex sync.Mutex
	updateActivity := func() {
		lastActivityMutex.Lock()
		lastActivity = time.Now()
		lastActivityMutex.Unlock()
	}

	// 用于通知所有goroutine退出
	done := make(chan struct{})

	// 空闲超时检测：如果 30 分钟内没有任何活动（输入或输出），关闭连接
	// 注意：只要有SSH输出（如编译日志），就不会被认为是空闲
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				lastActivityMutex.Lock()
				idleDuration := time.Since(lastActivity)
				lastActivityMutex.Unlock()

				// 如果空闲超过 30 分钟（无任何输入或输出），关闭连接
				if idleDuration > 30*time.Minute {
					log.Printf("SSH connection idle for %v, closing", idleDuration)
					session.Close()
					sshClient.Close()
					conn.Close()
					return
				}
			}
		}
	}()

	// 从 WebSocket 读取输入并发送到 SSH stdin
	go func() {
		defer wg.Done()
		defer stdin.Close()
		defer close(done) // 当WebSocket关闭时，通知其他goroutine退出

		for {
			// 不设置读取超时，让连接保持打开
			// 只要SSH会话还在运行，就保持WebSocket连接
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				return
			}

			// 更新活动时间
			updateActivity()

			var msg TerminalMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				// 如果不是 JSON，直接作为输入发送
				stdin.Write(message)
				continue
			}

			switch msg.Type {
			case "input":
				// 发送输入到 SSH
				stdin.Write([]byte(msg.Data))
			case "resize":
				// 更新终端大小
				if msg.Rows > 0 && msg.Columns > 0 {
					session.WindowChange(msg.Rows, msg.Columns)
					rows, cols = msg.Rows, msg.Columns
				}
			case "ping":
				// 响应 ping（ping 也算活动）
				conn.WriteJSON(TerminalMessage{Type: "pong"})
			}
		}
	}()

	// 从 SSH stdout 读取并发送到 WebSocket
	go func() {
		defer wg.Done()
		buffer := make([]byte, 4096) // 增大缓冲区以处理大量输出
		for {
			select {
			case <-done:
				return
			default:
				n, err := stdout.Read(buffer)
				if n > 0 {
					// 更新活动时间（有输出也算活动，这对于长时间运行的任务很重要）
					updateActivity()
					if err := conn.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
						log.Printf("WebSocket write error (stdout): %v", err)
						return
					}
				}
				if err != nil {
					if err != io.EOF {
						log.Printf("SSH stdout read error: %v", err)
					}
					return
				}
			}
		}
	}()

	// 从 SSH stderr 读取并发送到 WebSocket
	go func() {
		defer wg.Done()
		buffer := make([]byte, 4096) // 增大缓冲区以处理大量输出
		for {
			select {
			case <-done:
				return
			default:
				n, err := stderr.Read(buffer)
				if n > 0 {
					// 更新活动时间（有错误输出也算活动）
					updateActivity()
					if err := conn.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
						log.Printf("WebSocket write error (stderr): %v", err)
						return
					}
				}
				if err != nil {
					if err != io.EOF {
						log.Printf("SSH stderr read error: %v", err)
					}
					return
				}
			}
		}
	}()

	// 等待会话结束
	session.Wait()
	wg.Wait()

	log.Printf("WebSSH session for host %d (user %d) closed", hostID, userID)
}

// establishSSHConnection 建立 SSH 客户端连接
func (h *WebSSHHandler) establishSSHConnection(host *models.Host, conn *websocket.Conn, authMethod ssh.AuthMethod, username string) (*ssh.Client, error) {
	// 构建 SSH 客户端配置
	config := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境应该验证主机密钥
		Timeout:         10 * time.Second,
		Auth:            []ssh.AuthMethod{authMethod},
	}

	// 连接到 SSH 服务器
	hostname := host.IPAddress
	if hostname == "" {
		hostname = host.Hostname
	}
	port := host.SSHPort
	if port == 0 {
		port = 22
	}
	address := fmt.Sprintf("%s:%d", hostname, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH server: %w", err)
	}

	// 启动 SSH 连接保活机制
	// 每 30 秒发送一次 keepalive，确保连接在有操作时保持活跃
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// 发送 keepalive 请求，忽略错误（连接关闭时会自动退出）
			_, _, err := client.SendRequest("keepalive@openssh.com", false, nil)
			if err != nil {
				// 连接已关闭，退出keepalive循环
				return
			}
		}
	}()

	return client, nil
}

// requestAuthentication 通过 WebSocket 请求认证方式（密钥或密码）
func (h *WebSSHHandler) requestAuthentication(conn *websocket.Conn, host *models.Host, userID uint64) (ssh.AuthMethod, string, error) {
	// 获取用户的SSH密钥列表
	keys, _, err := h.sshKeyService.ListKeys(userID, 1, 100)
	if err != nil {
		log.Printf("Failed to list SSH keys: %v", err)
		keys = []models.SSHKey{} // 继续，即使获取密钥失败
	}

	// 发送认证方式选择请求
	msg := TerminalMessage{
		Type: "auth_method_request",
		Data: "Please select authentication method (key or password)",
	}
	if err := conn.WriteJSON(msg); err != nil {
		return nil, "", err
	}

	// 等待认证方式响应
	var authMethodType string
	var keyID uint64
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return nil, "", err
		}

		var msg TerminalMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		if msg.Type == "auth_method" {
			authMethodType = strings.ToLower(strings.TrimSpace(msg.Data))
			if authMethodType == "key" || authMethodType == "password" {
				break
			}
		}
	}

	// 根据认证方式处理
	var username string
	if authMethodType == "key" {
		// 密钥认证
		if len(keys) == 0 {
			return nil, "", fmt.Errorf("no SSH keys available")
		}

		// 如果有默认密钥，使用默认密钥；否则请求选择
		if host.SSHKeyID != nil {
			keyID = *host.SSHKeyID
		} else {
			// 发送密钥选择请求
			msg = TerminalMessage{
				Type: "key_selection_request",
				Data: "Please select SSH key",
			}
			if err := conn.WriteJSON(msg); err != nil {
				return nil, "", err
			}

			// 等待密钥选择响应
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					return nil, "", err
				}

				var msg TerminalMessage
				if err := json.Unmarshal(message, &msg); err != nil {
					continue
				}

				if msg.Type == "key_id" {
					keyID = msg.KeyID
					if keyID > 0 {
						break
					}
				}
			}
		}

		// 获取密钥信息（包括用户名）
		key, err := h.sshKeyService.GetKey(keyID, userID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get SSH key: %w", err)
		}

		// 如果密钥有用户名，使用密钥的用户名；否则请求用户名
		if key.Username != "" {
			username = key.Username
		} else {
			// 请求用户名
			username, err = h.requestUsername(conn)
			if err != nil {
				return nil, "", err
			}
		}

		// 获取并解析私钥
		privateKeyContent, err := h.sshKeyService.GetDecryptedPrivateKey(keyID, userID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get private key: %w", err)
		}

		// 解析私钥
		signer, err := ssh.ParsePrivateKey([]byte(privateKeyContent))
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse private key: %w", err)
		}

		return ssh.PublicKeys(signer), username, nil
	} else {
		// 密码认证：先请求用户名
		username, err := h.requestUsername(conn)
		if err != nil {
			return nil, "", err
		}
		// 然后请求密码
		password, err := h.requestPassword(conn)
		if err != nil {
			return nil, "", err
		}
		return ssh.Password(password), username, nil
	}
}

// requestUsername 请求用户名
func (h *WebSSHHandler) requestUsername(conn *websocket.Conn) (string, error) {
	msg := TerminalMessage{
		Type: "username_request",
		Data: "Please enter SSH username:",
	}
	if err := conn.WriteJSON(msg); err != nil {
		return "", err
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return "", err
		}

		var msg TerminalMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		if msg.Type == "username" {
			username := strings.TrimSpace(msg.Data)
			if username != "" {
				return username, nil
			}
		}
	}
}

// requestPassword 请求密码
func (h *WebSSHHandler) requestPassword(conn *websocket.Conn) (string, error) {
	msg := TerminalMessage{
		Type: "password_request",
		Data: "Please enter SSH password:",
	}
	if err := conn.WriteJSON(msg); err != nil {
		return "", err
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return "", err
		}

		var msg TerminalMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		if msg.Type == "password" {
			password := msg.Data
			if password != "" {
				return password, nil
			}
		}
	}
}

// sendErrorMessage 发送错误消息到 WebSocket
func sendErrorMessage(conn *websocket.Conn, message string) {
	msg := TerminalMessage{
		Type: "error",
		Data: message,
	}
	conn.WriteJSON(msg)
	conn.Close()
}

