// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package websocket

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/service/sshkey"
	"github.com/kkops/backend/internal/utils"
)

// SSHTerminalHandler handles SSH terminal WebSocket connections
// WS /ws/ssh/connect
func SSHTerminalHandler(db *gorm.DB, cfg interface{}, sshkeySvc *sshkey.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			// For WebSocket, we need to close connection properly
			// Try to upgrade first, then send error and close
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err == nil {
				conn.WriteJSON(map[string]interface{}{
					"type": "error",
					"data": "unauthorized",
				})
				conn.Close()
			} else {
				c.JSON(401, gin.H{"error": "unauthorized"})
			}
			return
		}

		// Upgrade to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		// Send initial connection message
		conn.WriteJSON(map[string]interface{}{
			"type": "ready",
			"data": map[string]string{
				"message": "WebSocket connected, send 'connect' message to establish SSH session",
			},
		})

		// Handle messages - wait for connect message
		for {
			var msg map[string]interface{}
			if err := conn.ReadJSON(&msg); err != nil {
				log.Printf("WebSocket read error: %v", err)
				break
			}

			msgType, ok := msg["type"].(string)
			if !ok {
				conn.WriteJSON(map[string]interface{}{
					"type": "error",
					"data": "Invalid message format",
				})
				continue
			}

			switch msgType {
			case "connect":
				// Handle SSH connection - this will manage the entire session lifecycle
				handleSSHConnect(conn, msg, db, userID.(uint), sshkeySvc)
				return // After SSH session ends, close the WebSocket handler
			default:
				conn.WriteJSON(map[string]interface{}{
					"type": "error",
					"data": "Please send 'connect' message first",
				})
			}
		}
	}
}

func handleSSHConnect(conn *websocket.Conn, msg map[string]interface{}, db *gorm.DB, userID uint, sshkeySvc *sshkey.Service) {
	data, ok := msg["data"].(map[string]interface{})
	if !ok {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Invalid connect message format",
		})
		return
	}

	// Parse connection parameters
	assetIDFloat, ok := data["asset_id"].(float64)
	if !ok {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "asset_id is required",
		})
		return
	}
	assetID := uint(assetIDFloat)

	// Get asset
	var asset model.Asset
	if err := db.Preload("SSHKey").First(&asset, assetID).Error; err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Asset not found",
		})
		return
	}

	// Check asset status
	if asset.Status != "active" {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Asset is not active",
		})
		return
	}

	authType, _ := data["auth_type"].(string)
	if authType == "" {
		authType = "key" // Default to key authentication
	}

	username, _ := data["username"].(string)
	if username == "" {
		username = asset.SSHUser
		if username == "" && asset.SSHKey != nil {
			username = asset.SSHKey.SSHUser
		}
	}

	if username == "" {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Username is required",
		})
		return
	}

	// Establish SSH connection
	port := asset.SSHPort
	if port == 0 {
		port = 22
	}

	var sshClient *utils.SSHClient
	var err error
	timeout := 30 * time.Second

	if asset.SSHKeyID == nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "SSH key is required for asset",
		})
		return
	}

	// Get SSH key and decrypt private key
	var sshKey model.SSHKey
	if err := db.First(&sshKey, *asset.SSHKeyID).Error; err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "SSH key not found",
		})
		return
	}

	// Use the SSH key's owner user ID for decryption
	privateKeyBytes, err := sshkeySvc.GetDecryptedPrivateKey(sshKey.UserID, *asset.SSHKeyID)
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Failed to get private key: " + err.Error(),
		})
		return
	}

	passphraseBytes, _ := sshkeySvc.GetDecryptedPassphrase(sshKey.UserID, *asset.SSHKeyID)
	if len(passphraseBytes) > 0 {
		sshClient, err = utils.NewSSHClientWithPassphrase(asset.IP, port, username, privateKeyBytes, passphraseBytes, timeout)
	} else {
		sshClient, err = utils.NewSSHClient(asset.IP, port, username, privateKeyBytes, timeout)
	}

	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "SSH connection failed: " + err.Error(),
		})
		return
	}
	defer sshClient.Close()

	// Create SSH session with PTY
	session, err := sshClient.Client().NewSession()
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Failed to create SSH session: " + err.Error(),
		})
		return
	}
	defer session.Close()

	// Request PTY with default size
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Default terminal size
	cols := 80
	rows := 24

	if sizeData, ok := data["size"].(map[string]interface{}); ok {
		if c, ok := sizeData["cols"].(float64); ok {
			cols = int(c)
		}
		if r, ok := sizeData["rows"].(float64); ok {
			rows = int(r)
		}
	}

	if err := session.RequestPty("xterm-256color", rows, cols, modes); err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Failed to request PTY: " + err.Error(),
		})
		return
	}

	// Set up stdin/stdout/stderr
	stdin, err := session.StdinPipe()
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Failed to get stdin pipe: " + err.Error(),
		})
		return
	}
	defer stdin.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Failed to get stdout pipe: " + err.Error(),
		})
		return
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Failed to get stderr pipe: " + err.Error(),
		})
		return
	}

	// Start shell
	if err := session.Shell(); err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Failed to start shell: " + err.Error(),
		})
		return
	}

	// Send connected message
	conn.WriteJSON(map[string]interface{}{
		"type": "connected",
		"data": map[string]interface{}{
			"message": "SSH session established",
		},
	})

	// Set up bidirectional data streaming
	var wg sync.WaitGroup
	done := make(chan bool, 1)

	// Stream stdout to WebSocket
	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		for {
			select {
			case <-done:
				return
			default:
				n, err := stdout.Read(buffer)
				if n > 0 {
					conn.WriteJSON(map[string]interface{}{
						"type": "output",
						"data": string(buffer[:n]),
					})
				}
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading stdout: %v", err)
					}
					return
				}
			}
		}
	}()

	// Stream stderr to WebSocket
	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		for {
			select {
			case <-done:
				return
			default:
				n, err := stderr.Read(buffer)
				if n > 0 {
					conn.WriteJSON(map[string]interface{}{
						"type": "output",
						"data": string(buffer[:n]),
					})
				}
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading stderr: %v", err)
					}
					return
				}
			}
		}
	}()

	// Handle WebSocket messages
	for {
		var wsMsg map[string]interface{}
		if err := conn.ReadJSON(&wsMsg); err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		msgType, ok := wsMsg["type"].(string)
		if !ok {
			continue
		}

		switch msgType {
		case "input":
			// Send input to SSH session
			if inputData, ok := wsMsg["data"].(map[string]interface{}); ok {
				if data, ok := inputData["data"].(string); ok {
					stdin.Write([]byte(data))
				}
			}
		case "resize":
			// Handle terminal resize
			if resizeData, ok := wsMsg["data"].(map[string]interface{}); ok {
				var newCols, newRows int
				if c, ok := resizeData["cols"].(float64); ok {
					newCols = int(c)
				}
				if r, ok := resizeData["rows"].(float64); ok {
					newRows = int(r)
				}
				if newCols > 0 && newRows > 0 {
					session.WindowChange(newRows, newCols)
					cols = newCols
					rows = newRows
				}
			}
		case "disconnect":
			// User requested disconnect
			done <- true
			session.Close()
			conn.WriteJSON(map[string]interface{}{
				"type": "disconnected",
				"data": map[string]interface{}{
					"reason": "User requested disconnect",
				},
			})
			return
		default:
			log.Printf("Unknown message type: %s", msgType)
		}
	}

	// Cleanup
	done <- true
	session.Close()

	// Wait for streams to finish
	wg.Wait()

	conn.WriteJSON(map[string]interface{}{
		"type": "disconnected",
		"data": map[string]interface{}{
			"reason": "Connection closed",
		},
	})
}
