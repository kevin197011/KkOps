// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package websocket

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/service/authorization"
	"github.com/kkops/backend/internal/service/sshkey"
	"github.com/kkops/backend/internal/utils"
)

// zmodemState tracks ZMODEM file transfer state
type zmodemState struct {
	mu            sync.Mutex
	active        bool
	direction     string // "upload" or "download"
	mode          string // "zmodem" or "trzsz" or "sftp"
	filename      string
	totalSize     int64
	bytesSent     int64
	bytesReceived int64
	canceled      bool
	// For SFTP upload
	uploadBuffer  *bytes.Buffer
	uploadChan    chan []byte
	uploadDone    chan bool
}

// detectZmodemSequence detects ZMODEM protocol initiation sequences
// Returns: detected, direction, mode, sequence position
// ZMODEM sequences can be:
//   - **B (two asterisks followed by B) for upload/rz
//   - **G (two asterisks followed by G) for download/sz
//   - **\x18B (may include control characters like CAN 0x18) for upload
//   - **\x18G (may include control characters like CAN 0x18) for download
func detectZmodemSequence(data []byte) (bool, string, string, int) {
	// Check for classic ZMODEM: **B or **G
	// Also handle cases where there may be control characters between ** and B/G
	// Search for pattern: ** followed by optional control chars, then B or G
	// Limit search to reasonable length to avoid false positives
	const maxSequenceLen = 10 // Maximum length of sequence including control chars

	for i := 0; i < len(data)-2 && i < len(data); i++ {
		if data[i] == 0x2A && data[i+1] == 0x2A {
			// Found **, now look for B or G, possibly with control chars in between
			// Check direct match first: **B or **G
			if i+2 < len(data) {
				if data[i+2] == 0x42 { // 'B'
					return true, "upload", "zmodem", i
				}
				if data[i+2] == 0x47 { // 'G'
					return true, "download", "zmodem", i
				}
			}

			// Then check for pattern with control characters: **[\x00-\x1F]*B or **[\x00-\x1F]*G
			// Look ahead up to maxSequenceLen bytes
			j := i + 2
			for j < len(data) && j < i+maxSequenceLen {
				if data[j] == 0x42 { // 'B'
					return true, "upload", "zmodem", i
				}
				if data[j] == 0x47 { // 'G'
					return true, "download", "zmodem", i
				}
				// If we encounter a printable character (not B/G), stop searching
				if data[j] >= 0x20 {
					break
				}
				j++
			}
		}
	}

	// Also check as string for trzsz protocol
	dataStr := string(data)
	if pos := strings.Index(dataStr, "::TRZSZ:"); pos >= 0 {
		if strings.Contains(dataStr[pos:], "TRANSFER:UPLOAD") || strings.Contains(dataStr, "trz") {
			return true, "upload", "trzsz", pos
		}
		if strings.Contains(dataStr[pos:], "TRANSFER:DOWNLOAD") || strings.Contains(dataStr, "tsz") {
			return true, "download", "trzsz", pos
		}
	}
	return false, "", "", -1
}

// SSHTerminalHandler handles SSH terminal WebSocket connections
// WS /ws/ssh/connect
// 增加资产访问权限检查：管理员可以连接任意资产，普通用户只能连接已授权的资产
func SSHTerminalHandler(db *gorm.DB, cfg interface{}, sshkeySvc *sshkey.Service, authzSvc *authorization.Service) gin.HandlerFunc {
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
				handleSSHConnect(conn, msg, db, userID.(uint), sshkeySvc, authzSvc)
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

func handleSSHConnect(conn *websocket.Conn, msg map[string]interface{}, db *gorm.DB, userID uint, sshkeySvc *sshkey.Service, authzSvc *authorization.Service) {
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

	// 检查用户对资产的访问权限
	hasAccess, err := authzSvc.HasAssetAccess(userID, assetID)
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Failed to check permission",
		})
		return
	}
	if !hasAccess {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "no permission to access this asset",
		})
		return
	}

	// Get asset
	var asset model.Asset
	if err := db.Preload("SSHKey").First(&asset, assetID).Error; err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "Asset not found",
		})
		return
	}

	// Log asset configuration for debugging
	log.Printf("Asset loaded: id=%d, hostname=%s, ip=%s, ssh_port=%d, ssh_user=%s",
		asset.ID, asset.HostName, asset.IP, asset.SSHPort, asset.SSHUser)

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

	// Establish SSH connection - use asset's configured SSH port
	port := asset.SSHPort
	if port == 0 {
		port = 22 // Default SSH port
		log.Printf("SSH port not configured for asset %d, using default port 22", assetID)
	}

	log.Printf("Connecting to SSH: %s@%s:%d (asset_id=%d, configured_ssh_port=%d)",
		username, asset.IP, port, assetID, asset.SSHPort)

	var sshClient *utils.SSHClient
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

	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshClient.Client())
	if err != nil {
		log.Printf("Failed to create SFTP client: %v", err)
		// Don't fail the connection, just log the error - terminal will still work
		sftpClient = nil
	} else {
		defer sftpClient.Close()
		log.Printf("SFTP client created successfully")
	}

	// Send connected message
	conn.WriteJSON(map[string]interface{}{
		"type": "connected",
		"data": map[string]interface{}{
			"message":    "SSH session established",
			"sftp_ready": sftpClient != nil,
		},
	})

	// Initialize ZMODEM state
	zmodemState := &zmodemState{}

	// Set up bidirectional data streaming
	var wg sync.WaitGroup
	done := make(chan bool, 1)

	// Buffers for detecting ZMODEM sequences (separate for stdout and stderr)
	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}
	const bufferSize = 1024

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
					data := buffer[:n]

					// Check if ZMODEM transfer is active
					zmodemState.mu.Lock()
					isActive := zmodemState.active
					activeDirection := zmodemState.direction
					zmodemState.mu.Unlock()

					if isActive {

						if activeDirection == "download" {
							// During download, send raw binary data
							if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
								log.Printf("Error writing binary data: %v", err)
								return
							}
							zmodemState.mu.Lock()
							zmodemState.bytesReceived += int64(n)
							zmodemState.mu.Unlock()
						} else {
							// During upload, filter out ZMODEM protocol data but allow normal text
							// This allows completion messages and errors to pass through
							dataStr := string(data)

							// Check for upload completion indicators
							// Look for shell prompt patterns that indicate command completion
							hasCompletionIndicator := false

							// Check for common shell prompt patterns
							promptPatterns := []string{
								"root@", // root user prompt
								"# ",    // root prompt ending
								"$ ",    // regular user prompt
								"~#",    // home directory + root prompt
								"~$",    // home directory + user prompt
								":~# ",  // full path prompt
								":~$ ",  // full path prompt
							}

							for _, pattern := range promptPatterns {
								if strings.Contains(dataStr, pattern) {
									// Additional check: make sure it's not part of "waiting" message
									if !strings.Contains(dataStr, "waiting") {
										hasCompletionIndicator = true
										break
									}
								}
							}

							// Check for ZMODEM protocol sequences
							hasZmodemProtocol := false
							if strings.Contains(dataStr, "**") {
								// Check if it's actually a ZMODEM sequence start
								for i := 0; i <= len(data)-3; i++ {
									if data[i] == 0x2A && data[i+1] == 0x2A {
										// Check if followed by B or G (with optional control chars)
										j := i + 2
										for j < len(data) && j < i+10 {
											if data[j] == 0x42 || data[j] == 0x47 { // B or G
												hasZmodemProtocol = true
												break
											}
											if data[j] >= 0x20 {
												break // Printable char that's not B/G, not protocol
											}
											j++
										}
										if hasZmodemProtocol {
											break
										}
									}
								}
							}

							if hasZmodemProtocol {
								// This is ZMODEM protocol data, ignore it
								continue
							}

							// If we see completion indicators, auto-deactivate ZMODEM mode
							if hasCompletionIndicator {
								zmodemState.mu.Lock()
								if zmodemState.active && zmodemState.direction == "upload" {
									log.Printf("Upload completion detected from server output, deactivating ZMODEM mode")
									zmodemState.active = false
									stdoutBuffer.Reset()
									stderrBuffer.Reset()

									// Send completion message
									conn.WriteJSON(map[string]interface{}{
										"type": "zmodem_end",
										"data": map[string]interface{}{
											"success": true,
											"message": "Transfer completed",
										},
									})
								}
								zmodemState.mu.Unlock()
							}

							// Send output (completion messages, errors, normal output, etc.)
							conn.WriteJSON(map[string]interface{}{
								"type": "output",
								"data": dataStr,
							})
							continue
						}
					}

					// Normal output processing (when ZMODEM is not active)
					if !isActive {
						// For normal output: prioritize immediate sending, with minimal buffering for ZMODEM detection
						// Strategy: Send new data immediately, but check if it might be start of ZMODEM sequence
						// Only buffer the tail to check for sequences that span packet boundaries

						var detected bool
						var direction, mode string
						var seqPos int

						// First, check if new data itself contains ZMODEM sequence
						detected, direction, mode, seqPos = detectZmodemSequence(data)
						if detected {
							log.Printf("ZMODEM sequence detected in stdout new data: direction=%s, mode=%s, pos=%d", direction, mode, seqPos)

							// Send any buffered data first
							if stdoutBuffer.Len() > 0 {
								conn.WriteJSON(map[string]interface{}{
									"type": "output",
									"data": stdoutBuffer.String(),
								})
								stdoutBuffer.Reset()
							}

							// Send text before ZMODEM sequence (if any)
							if seqPos > 0 {
								conn.WriteJSON(map[string]interface{}{
									"type": "output",
									"data": string(data[:seqPos]),
								})
							}

							// Activate ZMODEM mode
							zmodemState.mu.Lock()
							zmodemState.active = true
							zmodemState.direction = direction
							zmodemState.mode = mode
							zmodemState.bytesReceived = 0
							zmodemState.bytesSent = 0
							zmodemState.canceled = false
							zmodemState.mu.Unlock()

							conn.WriteJSON(map[string]interface{}{
								"type": "zmodem_start",
								"data": map[string]interface{}{
									"direction": direction,
									"mode":      mode,
								},
							})

							stdoutBuffer.Reset()
							continue
						}

						// Append new data to buffer (for checking sequences that span boundaries)
						stdoutBuffer.Write(data)
						bufBytes := stdoutBuffer.Bytes()

						// Check for ZMODEM sequence in the combined buffer (including previous tail)
						detected, direction, mode, seqPos = detectZmodemSequence(bufBytes)
						if detected {
							log.Printf("ZMODEM sequence detected in stdout buffer: direction=%s, mode=%s, pos=%d", direction, mode, seqPos)

							// Check if text before sequence contains rz/sz waiting messages that should be hidden
							// Find the start of the waiting message (if any)
							hideStart := seqPos
							if seqPos > 0 {
								textBefore := string(bufBytes[:seqPos])
								// Look for common rz/sz waiting patterns
								waitingPatterns := []string{"rz waiting", "sz waiting", "waiting to receive", "waiting to send"}
								for _, pattern := range waitingPatterns {
									if idx := strings.LastIndex(textBefore, pattern); idx >= 0 {
										// Find the start of the line containing this pattern
										lineStart := strings.LastIndex(textBefore[:idx], "\n")
										if lineStart >= 0 {
											hideStart = lineStart + 1
										} else {
											hideStart = idx
										}
										break
									}
								}
							}

							// Send text before the hiding point (if any)
							if hideStart > 0 {
								conn.WriteJSON(map[string]interface{}{
									"type": "output",
									"data": string(bufBytes[:hideStart]),
								})
							}

							// Activate ZMODEM mode
							zmodemState.mu.Lock()
							zmodemState.active = true
							zmodemState.direction = direction
							zmodemState.mode = mode
							zmodemState.bytesReceived = 0
							zmodemState.bytesSent = 0
							zmodemState.canceled = false
							zmodemState.mu.Unlock()

							conn.WriteJSON(map[string]interface{}{
								"type": "zmodem_start",
								"data": map[string]interface{}{
									"direction": direction,
									"mode":      mode,
								},
							})

							stdoutBuffer.Reset()
							continue
						}

						// No ZMODEM sequence detected: send data immediately, keep only minimal tail
						// ZMODEM sequence can be **B, **G, or **\x18B, **\x18G (with control chars)
						// Keep 3-4 bytes to handle sequences that may span packet boundaries
						// We keep 3 bytes to cover **X or **\x18 patterns
						const keepSize = 3
						if len(bufBytes) > keepSize {
							// Check if the tail we want to keep starts with ** (potential ZMODEM sequence start)
							// If not, we can safely send more bytes
							sendLen := len(bufBytes) - keepSize
							tail := bufBytes[sendLen:]

							// If tail doesn't start with **, we can send everything
							// This prevents unnecessary retention of non-ZMODEM data
							if len(tail) < 2 || !(tail[0] == 0x2A && tail[1] == 0x2A) {
								// No ** at start of tail, safe to send all
								conn.WriteJSON(map[string]interface{}{
									"type": "output",
									"data": string(bufBytes),
								})
								stdoutBuffer.Reset()
							} else {
								// Tail starts with **, keep it for next check
								conn.WriteJSON(map[string]interface{}{
									"type": "output",
									"data": string(bufBytes[:sendLen]),
								})
								stdoutBuffer.Reset()
								stdoutBuffer.Write(tail)
							}
						} else if len(bufBytes) > 0 {
							// Buffer is small (<= keepSize): check if it might be start of ZMODEM
							// Only keep if it starts with **, otherwise send immediately
							if len(bufBytes) >= 2 && bufBytes[0] == 0x2A && bufBytes[1] == 0x2A {
								// Might be start of ZMODEM sequence, keep it
								// Don't send yet
							} else {
								// Not a potential ZMODEM sequence, send immediately
								conn.WriteJSON(map[string]interface{}{
									"type": "output",
									"data": string(bufBytes),
								})
								stdoutBuffer.Reset()
							}
						}
					}
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
					data := buffer[:n]

					// Check if ZMODEM transfer is active
					zmodemState.mu.Lock()
					isActive := zmodemState.active
					zmodemState.mu.Unlock()

					if isActive {
						// During ZMODEM transfer, ignore stderr output (it's protocol data)
						// Only stdout is used for download data
						continue
					}

					// Check for ZMODEM sequences in stderr too (use separate buffer)
					// Strategy same as stdout: prioritize immediate sending

					var detected bool
					var direction, mode string
					var seqPos int

					// First, check if new data itself contains ZMODEM sequence
					detected, direction, mode, seqPos = detectZmodemSequence(data)
					if detected {
						log.Printf("ZMODEM sequence detected in stderr new data: direction=%s, mode=%s, pos=%d", direction, mode, seqPos)

						// Send any buffered data first
						if stderrBuffer.Len() > 0 {
							conn.WriteJSON(map[string]interface{}{
								"type": "output",
								"data": stderrBuffer.String(),
							})
							stderrBuffer.Reset()
						}

						// Send text before ZMODEM sequence (if any)
						if seqPos > 0 {
							conn.WriteJSON(map[string]interface{}{
								"type": "output",
								"data": string(data[:seqPos]),
							})
						}

						// Activate ZMODEM mode
						zmodemState.mu.Lock()
						zmodemState.active = true
						zmodemState.direction = direction
						zmodemState.mode = mode
						zmodemState.bytesReceived = 0
						zmodemState.bytesSent = 0
						zmodemState.canceled = false
						zmodemState.mu.Unlock()

						conn.WriteJSON(map[string]interface{}{
							"type": "zmodem_start",
							"data": map[string]interface{}{
								"direction": direction,
								"mode":      mode,
							},
						})

						stderrBuffer.Reset()
						continue
					}

					// Append new data to buffer (for checking sequences that span boundaries)
					stderrBuffer.Write(data)
					bufBytes := stderrBuffer.Bytes()

					// Check for ZMODEM sequence in the combined buffer
					detected, direction, mode, seqPos = detectZmodemSequence(bufBytes)
					if detected {
						log.Printf("ZMODEM sequence detected in stderr buffer: direction=%s, mode=%s, pos=%d", direction, mode, seqPos)

						// Send text before ZMODEM sequence (if any)
						if seqPos > 0 {
							conn.WriteJSON(map[string]interface{}{
								"type": "output",
								"data": string(bufBytes[:seqPos]),
							})
						}

						// Activate ZMODEM mode
						zmodemState.mu.Lock()
						zmodemState.active = true
						zmodemState.direction = direction
						zmodemState.mode = mode
						zmodemState.bytesReceived = 0
						zmodemState.bytesSent = 0
						zmodemState.canceled = false
						zmodemState.mu.Unlock()

						conn.WriteJSON(map[string]interface{}{
							"type": "zmodem_start",
							"data": map[string]interface{}{
								"direction": direction,
								"mode":      mode,
							},
						})

						stderrBuffer.Reset()
						continue
					}

					// No ZMODEM sequence detected: send data immediately, keep only minimal tail
					// ZMODEM sequence can be **B, **G, or **\x18B, **\x18G (with control chars)
					// Keep 3-4 bytes to handle sequences that may span packet boundaries
					// We keep 3 bytes to cover **X or **\x18 patterns
					const keepSize = 3
					if len(bufBytes) > keepSize {
						// Check if the tail we want to keep starts with ** (potential ZMODEM sequence start)
						// If not, we can safely send more bytes
						sendLen := len(bufBytes) - keepSize
						tail := bufBytes[sendLen:]

						// If tail doesn't start with **, we can send everything
						// This prevents unnecessary retention of non-ZMODEM data
						if len(tail) < 2 || !(tail[0] == 0x2A && tail[1] == 0x2A) {
							// No ** at start of tail, safe to send all
							conn.WriteJSON(map[string]interface{}{
								"type": "output",
								"data": string(bufBytes),
							})
							stderrBuffer.Reset()
						} else {
							// Tail starts with **, keep it for next check
							conn.WriteJSON(map[string]interface{}{
								"type": "output",
								"data": string(bufBytes[:sendLen]),
							})
							stderrBuffer.Reset()
							stderrBuffer.Write(tail)
						}
					} else if len(bufBytes) > 0 {
						// Buffer is small (<= keepSize): check if it might be start of ZMODEM
						// Only keep if it starts with **, otherwise send immediately
						if len(bufBytes) >= 2 && bufBytes[0] == 0x2A && bufBytes[1] == 0x2A {
							// Might be start of ZMODEM sequence, keep it
							// Don't send yet
						} else {
							// Not a potential ZMODEM sequence, send immediately
							conn.WriteJSON(map[string]interface{}{
								"type": "output",
								"data": string(bufBytes),
							})
							stderrBuffer.Reset()
						}
					}
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

	// Handle WebSocket messages (both JSON and binary)
	for {
		messageType, messageData, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		if messageType == websocket.BinaryMessage {
			zmodemState.mu.Lock()
			isActive := zmodemState.active
			direction := zmodemState.direction
			mode := zmodemState.mode
			zmodemState.mu.Unlock()

			if isActive {
				if mode == "sftp" && direction == "upload" {
					// SFTP upload: send data to channel
					zmodemState.mu.Lock()
					if zmodemState.uploadChan != nil {
						select {
						case zmodemState.uploadChan <- messageData:
							// Data sent successfully
						default:
							log.Printf("Upload channel full, dropping data")
						}
					}
					zmodemState.bytesSent += int64(len(messageData))
					bytesSent := zmodemState.bytesSent
					totalSize := zmodemState.totalSize
					zmodemState.mu.Unlock()

					// Send progress update
					if totalSize > 0 {
						conn.WriteJSON(map[string]interface{}{
							"type": "sftp_upload_progress",
							"data": map[string]interface{}{
								"bytes_transferred": bytesSent,
								"total_bytes":       totalSize,
							},
						})
					}
				} else if mode == "zmodem" && direction == "upload" {
					// ZMODEM upload: forward to SSH stdin
					if _, err := stdin.Write(messageData); err != nil {
						log.Printf("Error writing to stdin: %v", err)
						break
					}

					// Update progress
					zmodemState.mu.Lock()
					zmodemState.bytesSent += int64(len(messageData))
					bytesSent := zmodemState.bytesSent
					totalSize := zmodemState.totalSize
					zmodemState.mu.Unlock()

					// Send progress update
					if totalSize > 0 {
						conn.WriteJSON(map[string]interface{}{
							"type": "zmodem_progress",
							"data": map[string]interface{}{
								"bytes_transferred": bytesSent,
								"total_bytes":       totalSize,
							},
						})
					}
				}
			}
			continue
		}

		// Parse JSON message
		if messageType != websocket.TextMessage {
			continue // Skip non-text messages when not in ZMODEM mode
		}

		var wsMsg map[string]interface{}
		if err := json.Unmarshal(messageData, &wsMsg); err != nil {
			log.Printf("Error parsing JSON message: %v", err)
			continue
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
					// Check if we should skip input during active transfer
					zmodemState.mu.Lock()
					isActive := zmodemState.active
					zmodemState.mu.Unlock()

					if !isActive {
						stdin.Write([]byte(data))
					}
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
		case "zmodem_cancel":
			// Cancel ZMODEM transfer
			log.Printf("ZMODEM transfer canceled by user")
			zmodemState.mu.Lock()
			zmodemState.canceled = true
			zmodemState.active = false
			zmodemState.mu.Unlock()

			// Clear buffers
			stdoutBuffer.Reset()
			stderrBuffer.Reset()

			conn.WriteJSON(map[string]interface{}{
				"type": "zmodem_end",
				"data": map[string]interface{}{
					"success": false,
					"message": "Transfer canceled",
				},
			})
		case "zmodem_data_size":
			// Frontend sends file size for progress tracking
			if sizeData, ok := wsMsg["data"].(map[string]interface{}); ok {
				if size, ok := sizeData["size"].(float64); ok {
					zmodemState.mu.Lock()
					zmodemState.totalSize = int64(size)
					if filename, ok := sizeData["filename"].(string); ok {
						zmodemState.filename = filename
					}
					zmodemState.mu.Unlock()
				}
			}
		case "zmodem_complete":
			// Transfer completed successfully (frontend has sent all file data)
			log.Printf("ZMODEM transfer completed (file data sent)")

			// Immediately send zmodem_end to frontend for UI update
			// We'll keep ZMODEM mode active briefly to allow server responses through
			conn.WriteJSON(map[string]interface{}{
				"type": "zmodem_end",
				"data": map[string]interface{}{
					"success": true,
					"message": "Transfer completed",
				},
			})

			// Set up timeout to deactivate ZMODEM mode after a short delay
			// This allows server completion responses (prompt, etc.) to pass through
			go func() {
				time.Sleep(500 * time.Millisecond) // Short delay to allow server responses
				zmodemState.mu.Lock()
				if zmodemState.active && zmodemState.direction == "upload" {
					log.Printf("ZMODEM upload cleanup: deactivating mode")
					zmodemState.active = false
					zmodemState.mu.Unlock()

					// Clear buffers
					stdoutBuffer.Reset()
					stderrBuffer.Reset()
				} else {
					zmodemState.mu.Unlock()
				}
			}()
		case "sftp_list":
			// List files in directory
			if sftpClient == nil {
				conn.WriteJSON(map[string]interface{}{
					"type": "sftp_error",
					"data": map[string]interface{}{
						"error": "SFTP client not available",
					},
				})
				break
			}
			if listData, ok := wsMsg["data"].(map[string]interface{}); ok {
				path, _ := listData["path"].(string)
				if path == "" {
					path = "."
				}
				handleSFTPList(conn, sftpClient, path)
			}
		case "sftp_upload":
			// Upload file (file data will be sent as binary message)
			if sftpClient == nil {
				conn.WriteJSON(map[string]interface{}{
					"type": "sftp_error",
					"data": map[string]interface{}{
						"error": "SFTP client not available",
					},
				})
				break
			}
			if uploadData, ok := wsMsg["data"].(map[string]interface{}); ok {
				remotePath, _ := uploadData["remote_path"].(string)
				fileName, _ := uploadData["file_name"].(string)
				size, _ := uploadData["size"].(float64)
				if remotePath == "" || fileName == "" {
					conn.WriteJSON(map[string]interface{}{
						"type": "sftp_error",
						"data": map[string]interface{}{
							"error": "remote_path and file_name are required",
						},
					})
					break
				}
				// Set upload state - file data will come as binary message
				zmodemState.mu.Lock()
				zmodemState.active = true
				zmodemState.direction = "upload"
				zmodemState.mode = "sftp"
				zmodemState.filename = fileName
				zmodemState.totalSize = int64(size)
				zmodemState.bytesSent = 0
				zmodemState.uploadChan = make(chan []byte, 10)
				zmodemState.uploadDone = make(chan bool, 1)
				zmodemState.uploadBuffer = bytes.NewBuffer(nil)
				zmodemState.mu.Unlock()

				// Start SFTP upload in goroutine
				go handleSFTPUpload(conn, sftpClient, remotePath, fileName, zmodemState, done)
			}
		case "sftp_download":
			// Download file
			if sftpClient == nil {
				conn.WriteJSON(map[string]interface{}{
					"type": "sftp_error",
					"data": map[string]interface{}{
						"error": "SFTP client not available",
					},
				})
				break
			}
			if downloadData, ok := wsMsg["data"].(map[string]interface{}); ok {
				remotePath, _ := downloadData["remote_path"].(string)
				if remotePath == "" {
					conn.WriteJSON(map[string]interface{}{
						"type": "sftp_error",
						"data": map[string]interface{}{
							"error": "remote_path is required",
						},
					})
					break
				}
				go handleSFTPDownload(conn, sftpClient, remotePath, zmodemState)
			}
		case "sftp_delete":
			// Delete file or directory
			if sftpClient == nil {
				conn.WriteJSON(map[string]interface{}{
					"type": "sftp_error",
					"data": map[string]interface{}{
						"error": "SFTP client not available",
					},
				})
				break
			}
			if deleteData, ok := wsMsg["data"].(map[string]interface{}); ok {
				path, _ := deleteData["path"].(string)
				if path == "" {
					conn.WriteJSON(map[string]interface{}{
						"type": "sftp_error",
						"data": map[string]interface{}{
							"error": "path is required",
						},
					})
					break
				}
				go handleSFTPDelete(conn, sftpClient, path)
			}
		case "sftp_mkdir":
			// Create directory
			if sftpClient == nil {
				conn.WriteJSON(map[string]interface{}{
					"type": "sftp_error",
					"data": map[string]interface{}{
						"error": "SFTP client not available",
					},
				})
				break
			}
			if mkdirData, ok := wsMsg["data"].(map[string]interface{}); ok {
				path, _ := mkdirData["path"].(string)
				if path == "" {
					conn.WriteJSON(map[string]interface{}{
						"type": "sftp_error",
						"data": map[string]interface{}{
							"error": "path is required",
						},
					})
					break
				}
				go handleSFTPMkdir(conn, sftpClient, path)
			}
		case "sftp_rename":
			// Rename file or directory
			if sftpClient == nil {
				conn.WriteJSON(map[string]interface{}{
					"type": "sftp_error",
					"data": map[string]interface{}{
						"error": "SFTP client not available",
					},
				})
				break
			}
			if renameData, ok := wsMsg["data"].(map[string]interface{}); ok {
				oldPath, _ := renameData["old_path"].(string)
				newPath, _ := renameData["new_path"].(string)
				if oldPath == "" || newPath == "" {
					conn.WriteJSON(map[string]interface{}{
						"type": "sftp_error",
						"data": map[string]interface{}{
							"error": "old_path and new_path are required",
						},
					})
					break
				}
				go handleSFTPRename(conn, sftpClient, oldPath, newPath)
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

// SFTP handler functions

func handleSFTPList(conn *websocket.Conn, sftpClient *sftp.Client, path string) {
	files, err := sftpClient.ReadDir(path)
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "sftp_error",
			"data": map[string]interface{}{
				"error": "Failed to list directory: " + err.Error(),
				"path":  path,
			},
		})
		return
	}

	fileList := make([]map[string]interface{}, 0, len(files))
	for _, file := range files {
		fileInfo := map[string]interface{}{
			"name":    file.Name(),
			"size":    file.Size(),
			"mode":    file.Mode().String(),
			"modTime": file.ModTime().Unix(),
			"isDir":   file.IsDir(),
		}
		fileList = append(fileList, fileInfo)
	}

	conn.WriteJSON(map[string]interface{}{
		"type": "sftp_list",
		"data": map[string]interface{}{
			"path":  path,
			"files": fileList,
		},
	})
}

func handleSFTPUpload(conn *websocket.Conn, sftpClient *sftp.Client, remotePath string, fileName string, zmodemState *zmodemState, done chan bool) {
	// Create a timeout for receiving file data
	timeout := time.NewTimer(60 * time.Second)
	defer timeout.Stop()

	// Read data from channel and write to buffer
	zmodemState.mu.Lock()
	uploadChan := zmodemState.uploadChan
	buffer := zmodemState.uploadBuffer
	totalSize := zmodemState.totalSize
	zmodemState.mu.Unlock()

	if uploadChan == nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "sftp_error",
			"data": map[string]interface{}{
				"error": "Upload channel not initialized",
			},
		})
		return
	}

	// Collect all data chunks
	for {
		select {
		case <-done:
			zmodemState.mu.Lock()
			zmodemState.active = false
			close(zmodemState.uploadChan)
			zmodemState.uploadChan = nil
			zmodemState.mu.Unlock()
			return
		case <-timeout.C:
			zmodemState.mu.Lock()
			zmodemState.active = false
			close(zmodemState.uploadChan)
			zmodemState.uploadChan = nil
			zmodemState.mu.Unlock()
			conn.WriteJSON(map[string]interface{}{
				"type": "sftp_error",
				"data": map[string]interface{}{
					"error": "Upload timeout",
				},
			})
			return
		case data, ok := <-uploadChan:
			if !ok {
				// Channel closed, all data received
				goto writeFile
			}
			if buffer != nil {
				buffer.Write(data)
			}
			// Reset timeout on each chunk
			timeout.Reset(60 * time.Second)

			// Check if all data received
			zmodemState.mu.Lock()
			bytesSent := zmodemState.bytesSent
			totalSize := zmodemState.totalSize
			zmodemState.mu.Unlock()

			if totalSize > 0 && bytesSent >= totalSize {
				// All data received
				goto writeFile
			}
		}

		// Check if we have all data
		zmodemState.mu.Lock()
		bytesSent := zmodemState.bytesSent
		totalSize = zmodemState.totalSize
		zmodemState.mu.Unlock()

		if totalSize > 0 && bytesSent >= totalSize {
			// Wait a bit for remaining chunks
			time.Sleep(100 * time.Millisecond)
			// Check if channel is empty
			select {
			case <-timeout.C:
				goto writeFile
			default:
			}
		}
	}

writeFile:
	// All data received, write to SFTP
	fullPath := filepath.Join(remotePath, fileName)
	if remotePath == "." || remotePath == "" {
		fullPath = fileName
	}

	// Create parent directory if it doesn't exist
	parentDir := filepath.Dir(fullPath)
	if parentDir != "." && parentDir != "/" {
		sftpClient.MkdirAll(parentDir)
	}

	dstFile, err := sftpClient.Create(fullPath)
	if err != nil {
		zmodemState.mu.Lock()
		zmodemState.active = false
		if zmodemState.uploadChan != nil {
			close(zmodemState.uploadChan)
			zmodemState.uploadChan = nil
		}
		zmodemState.mu.Unlock()
		conn.WriteJSON(map[string]interface{}{
			"type": "sftp_error",
			"data": map[string]interface{}{
				"error": "Failed to create file: " + err.Error(),
			},
		})
		return
	}

	if buffer != nil {
		_, err = dstFile.Write(buffer.Bytes())
	}
	dstFile.Close()

	if err != nil {
		zmodemState.mu.Lock()
		zmodemState.active = false
		if zmodemState.uploadChan != nil {
			close(zmodemState.uploadChan)
			zmodemState.uploadChan = nil
		}
		zmodemState.mu.Unlock()
		conn.WriteJSON(map[string]interface{}{
			"type": "sftp_error",
			"data": map[string]interface{}{
				"error": "Failed to write file: " + err.Error(),
			},
		})
		return
	}

	zmodemState.mu.Lock()
	bytesSent := zmodemState.bytesSent
	zmodemState.active = false
	if zmodemState.uploadChan != nil {
		close(zmodemState.uploadChan)
		zmodemState.uploadChan = nil
	}
	zmodemState.mu.Unlock()

	conn.WriteJSON(map[string]interface{}{
		"type": "sftp_upload_complete",
		"data": map[string]interface{}{
			"remote_path": fullPath,
			"file_name":   fileName,
			"size":        bytesSent,
		},
	})
}

func handleSFTPDownload(conn *websocket.Conn, sftpClient *sftp.Client, remotePath string, zmodemState *zmodemState) {
	srcFile, err := sftpClient.Open(remotePath)
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "sftp_error",
			"data": map[string]interface{}{
				"error":      "Failed to open file: " + err.Error(),
				"remote_path": remotePath,
			},
		})
		return
	}
	defer srcFile.Close()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "sftp_error",
			"data": map[string]interface{}{
				"error":      "Failed to get file info: " + err.Error(),
				"remote_path": remotePath,
			},
		})
		return
	}

	// Send download start message
	conn.WriteJSON(map[string]interface{}{
		"type": "sftp_download_start",
		"data": map[string]interface{}{
			"remote_path": remotePath,
			"file_name":   filepath.Base(remotePath),
			"size":        fileInfo.Size(),
		},
	})

	// Read and send file in chunks
	buffer := make([]byte, 8192) // 8KB chunks
	totalRead := int64(0)
	fileName := filepath.Base(remotePath)

	for {
		n, err := srcFile.Read(buffer)
		if n > 0 {
			totalRead += int64(n)
			// Send binary data
			if err := conn.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
				log.Printf("Error sending file chunk: %v", err)
				return
			}

			// Send progress update
			conn.WriteJSON(map[string]interface{}{
				"type": "sftp_download_progress",
				"data": map[string]interface{}{
					"remote_path":      remotePath,
					"file_name":        fileName,
					"bytes_transferred": totalRead,
					"total_bytes":       fileInfo.Size(),
				},
			})
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			conn.WriteJSON(map[string]interface{}{
				"type": "sftp_error",
				"data": map[string]interface{}{
					"error":      "Failed to read file: " + err.Error(),
					"remote_path": remotePath,
				},
			})
			return
		}
	}

	// Send download complete message
	conn.WriteJSON(map[string]interface{}{
		"type": "sftp_download_complete",
		"data": map[string]interface{}{
			"remote_path": remotePath,
			"file_name":   fileName,
			"size":        totalRead,
		},
	})
}

func handleSFTPDelete(conn *websocket.Conn, sftpClient *sftp.Client, path string) {
	// Check if it's a directory or file
	fileInfo, err := sftpClient.Stat(path)
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "sftp_error",
			"data": map[string]interface{}{
				"error": "Failed to stat path: " + err.Error(),
				"path":  path,
			},
		})
		return
	}

	if fileInfo.IsDir() {
		err = sftpClient.RemoveDirectory(path)
	} else {
		err = sftpClient.Remove(path)
	}

	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "sftp_error",
			"data": map[string]interface{}{
				"error": "Failed to delete: " + err.Error(),
				"path":  path,
			},
		})
		return
	}

	conn.WriteJSON(map[string]interface{}{
		"type": "sftp_delete_complete",
		"data": map[string]interface{}{
			"path": path,
		},
	})
}

func handleSFTPMkdir(conn *websocket.Conn, sftpClient *sftp.Client, path string) {
	err := sftpClient.MkdirAll(path)
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "sftp_error",
			"data": map[string]interface{}{
				"error": "Failed to create directory: " + err.Error(),
				"path":  path,
			},
		})
		return
	}

	conn.WriteJSON(map[string]interface{}{
		"type": "sftp_mkdir_complete",
		"data": map[string]interface{}{
			"path": path,
		},
	})
}

func handleSFTPRename(conn *websocket.Conn, sftpClient *sftp.Client, oldPath string, newPath string) {
	err := sftpClient.Rename(oldPath, newPath)
	if err != nil {
		conn.WriteJSON(map[string]interface{}{
			"type": "sftp_error",
			"data": map[string]interface{}{
				"error":    "Failed to rename: " + err.Error(),
				"old_path": oldPath,
				"new_path": newPath,
			},
		})
		return
	}

	conn.WriteJSON(map[string]interface{}{
		"type": "sftp_rename_complete",
		"data": map[string]interface{}{
			"old_path": oldPath,
			"new_path": newPath,
		},
	})
}
