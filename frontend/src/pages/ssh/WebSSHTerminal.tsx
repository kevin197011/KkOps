// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import { Layout, Button, message, Space, Tag, Input, Tree, Spin, Typography, theme as antdTheme, Dropdown } from 'antd'
import type { MenuProps } from 'antd'
import { DisconnectOutlined, ArrowLeftOutlined, CheckCircleOutlined, ReloadOutlined, FolderOutlined, FolderOpenOutlined, DatabaseOutlined, CloseOutlined, HomeOutlined, SunOutlined, MoonOutlined, CopyOutlined, CloseCircleOutlined, UserOutlined, LogoutOutlined } from '@ant-design/icons'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'
import { assetApi, Asset } from '@/api/asset'
import { projectApi, Project } from '@/api/project'
import { environmentApi, Environment } from '@/api/environment'
import { useThemeStore } from '@/stores/theme'

const { Header, Sider, Content } = Layout
const { Search } = Input
const { Text } = Typography

// Connection state interface
interface SSHConnection {
  id: string // Unique connection ID (e.g., `${asset.id}-${timestamp}`)
  asset: Asset
  terminal: Terminal | null
  fitAddon: FitAddon | null
  terminalRef: HTMLDivElement | null
  ws: WebSocket | null
  status: 'connecting' | 'connected' | 'disconnected' | 'error'
  createdAt: number
  containerRef: React.RefObject<HTMLDivElement>
}

const MAX_CONNECTIONS = 10

const WebSSHTerminal = () => {
  const navigate = useNavigate()
  const terminalContainerRef = useRef<HTMLDivElement>(null)
  const [connections, setConnections] = useState<Map<string, SSHConnection>>(new Map())
  const [connectionOrder, setConnectionOrder] = useState<string[]>([]) // Track tab order for drag & drop
  const [activeConnectionId, setActiveConnectionId] = useState<string | null>(null)
  const [draggedTabId, setDraggedTabId] = useState<string | null>(null) // Track currently dragged tab
  const [assets, setAssets] = useState<Asset[]>([])
  const [projects, setProjects] = useState<Project[]>([])
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [filteredAssets, setFilteredAssets] = useState<Asset[]>([])
  const [searchQuery, setSearchQuery] = useState('')
  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([])
  const [autoExpandParent, setAutoExpandParent] = useState(true)
  const { mode, toggleMode } = useThemeStore()
  const { user, clearAuth } = useAuthStore()
  const { token } = antdTheme.useToken()
  const treeContainerRef = useRef<HTMLDivElement>(null)
  // Track which terminals are being initialized to prevent duplicate initialization
  const initializingTerminalsRef = useRef<Set<string>>(new Set())
  
  // Get active connection
  const activeConnection = activeConnectionId ? connections.get(activeConnectionId) : null

  useEffect(() => {
    fetchAssets()
    fetchProjects()
    fetchEnvironments()

    return () => {
      // Cleanup all connections on unmount
      connections.forEach((conn) => {
        if (conn.ws) {
          try {
            conn.ws.close(1000, 'Component unmounting')
          } catch (e) {
            // Ignore errors during cleanup
          }
        }
        if (conn.terminal) {
          try {
            conn.terminal.dispose()
          } catch (e) {
            // Ignore errors during cleanup
          }
        }
      })
    }
  }, [])

  useEffect(() => {
    // Filter assets based on search query
    if (!searchQuery.trim()) {
      setFilteredAssets(assets)
    } else {
      const query = searchQuery.toLowerCase()
      setFilteredAssets(
        assets.filter(
          (asset) =>
            asset.hostName.toLowerCase().includes(query) ||
            asset.ip.toLowerCase().includes(query) ||
            (asset.ssh_user && asset.ssh_user.toLowerCase().includes(query))
        )
      )
    }
  }, [assets, searchQuery])

  // 强制设置 Tree 缩进单元宽度为 16px（两个空格）
  useEffect(() => {
    const updateTreeIndent = () => {
      if (treeContainerRef.current) {
        const indentUnits = treeContainerRef.current.querySelectorAll('.ant-tree-indent-unit')
        indentUnits.forEach((unit) => {
          const element = unit as HTMLElement
          element.style.width = '16px'
          element.style.minWidth = '16px'
          element.style.maxWidth = '16px'
          element.style.flex = '0 0 16px'
        })
      }
    }

    // 初始设置
    updateTreeIndent()

    // 使用 MutationObserver 监听 DOM 变化，当 Tree 更新时重新设置
    const observer = new MutationObserver(updateTreeIndent)
    if (treeContainerRef.current) {
      observer.observe(treeContainerRef.current, {
        childList: true,
        subtree: true,
      })
    }

    // 延迟执行，确保 Tree 已渲染
    const timer = setTimeout(updateTreeIndent, 100)

    return () => {
      observer.disconnect()
      clearTimeout(timer)
    }
  }, [filteredAssets, expandedKeys])

  // Handle window resize for active terminal
  useEffect(() => {
    const handleResize = () => {
      const conn = activeConnectionId ? connections.get(activeConnectionId) : null
      if (conn && conn.terminal && conn.fitAddon) {
        conn.fitAddon.fit()
        // Send resize message to server
        if (conn.ws && conn.ws.readyState === WebSocket.OPEN && conn.terminal) {
          const cols = conn.terminal.cols
          const rows = conn.terminal.rows
          conn.ws.send(JSON.stringify({
            type: 'resize',
            data: { cols, rows },
          }))
        }
      }
    }
    window.addEventListener('resize', handleResize)
    // Initial fit - use requestAnimationFrame for better performance
    if (activeConnectionId) {
      requestAnimationFrame(handleResize)
    }
    return () => window.removeEventListener('resize', handleResize)
  }, [activeConnectionId, connections])

  // Initialize terminal immediately when container is ready (don't wait for 'connected' status)
  // This provides instant feedback to users
  useEffect(() => {
    const connectionsArray = Array.from(connections.entries())
    connectionsArray.forEach(([connId, conn]) => {
      // Initialize terminal as soon as container is ready, even if status is 'connecting'
      // This eliminates delay and provides immediate visual feedback
      // IMPORTANT: Check if terminal already exists to avoid duplicate initialization
      if ((conn.status === 'connecting' || conn.status === 'connected') && !conn.terminal && conn.containerRef.current) {
        // Use ref to track initialization in progress - prevents race conditions
        if (initializingTerminalsRef.current.has(connId)) {
          console.log(`[Terminal Init ${connId}] Already initializing, skipping`)
          return
        }
        const containerElement = conn.containerRef.current
        // Double-check: verify container doesn't already have a terminal instance
        if (containerElement && containerElement.querySelector('.xterm')) {
          console.log(`[Terminal Init ${connId}] Terminal already exists in container, skipping`)
          return
        }
        // Mark as initializing BEFORE any async operations
        initializingTerminalsRef.current.add(connId)
        console.log(`[Terminal Init ${connId}] Initializing terminal for asset ${conn.asset.hostName} (${conn.asset.ip}), status: ${conn.status}`)
        // Capture connId explicitly for this terminal instance
        const initConnId = connId
        if (containerElement) {
          // Initialize terminal immediately
          const term = new Terminal({
            cursorBlink: true,
            fontSize: 13,
            lineHeight: 1.4,
            letterSpacing: 0.2,
            fontFamily: '"JetBrains Mono", "Fira Code", "Cascadia Code", "Consolas", "Monaco", "Courier New", monospace',
            // Rich color theme for code syntax highlighting and ANSI colors
            // Based on popular terminal color schemes with enhanced visibility
            theme: {
              background: '#0D1117',      // Dark background (GitHub Dark style)
              foreground: '#E6EDF3',       // Light gray text for readability
              cursor: '#58A6FF',           // Blue cursor
              cursorAccent: '#0D1117',     // Cursor text color
              selection: '#264F78',        // Selection background
              selectionBackground: 'rgba(56, 139, 253, 0.4)',
              // Standard ANSI colors (0-7)
              black: '#484F58',            // Dim black (visible on dark bg)
              red: '#FF7B72',              // Coral red
              green: '#7EE787',            // Bright green
              yellow: '#FFA657',           // Orange-yellow
              blue: '#79C0FF',             // Sky blue
              magenta: '#D2A8FF',          // Purple/magenta
              cyan: '#56D4DD',             // Cyan/teal
              white: '#E6EDF3',            // White/light gray
              // Bright ANSI colors (8-15)
              brightBlack: '#6E7681',      // Brighter gray
              brightRed: '#FFA198',        // Light coral
              brightGreen: '#A5F3C0',      // Light green
              brightYellow: '#FFD580',     // Light orange
              brightBlue: '#A5D6FF',       // Light blue
              brightMagenta: '#E2C5FF',    // Light purple
              brightCyan: '#7EEAEA',       // Light cyan
              brightWhite: '#FFFFFF',      // Pure white
            },
            allowTransparency: true,
            disableStdin: false,
          })

          const fit = new FitAddon()
          term.loadAddon(fit)
          term.open(containerElement)

          // Update connection with terminal immediately
          // Use functional update to ensure we have the latest state and avoid duplicates
          setConnections((prev) => {
            const newConnections = new Map(prev)
            const updatedConn = newConnections.get(initConnId)
            // Double-check: don't overwrite if terminal already exists
            if (updatedConn && !updatedConn.terminal) {
              newConnections.set(initConnId, {
                ...updatedConn,
                terminal: term,
                fitAddon: fit,
                terminalRef: containerElement,
              })
              console.log(`[Terminal Init ${initConnId}] Terminal instance stored in connection state`)
            } else if (updatedConn?.terminal) {
              console.log(`[Terminal Init ${initConnId}] Terminal already exists, skipping storage`)
              // Dispose the duplicate terminal we just created
              try {
                term.dispose()
              } catch (e) {
                console.warn(`[Terminal Init ${initConnId}] Failed to dispose duplicate terminal:`, e)
              }
              // Clear initializing flag
              initializingTerminalsRef.current.delete(initConnId)
              return prev // Don't update state if terminal already exists
            }
            return newConnections
          })

          // Setup terminal with proper layout timing
          // Wait for container to be visible and layout to be calculated
          requestAnimationFrame(() => {
            // Double RAF to ensure layout is complete (especially after display: none -> block)
            requestAnimationFrame(() => {
              // Get the latest connection state to ensure we have the WebSocket reference
              setConnections((prev) => {
                const latestConn = prev.get(initConnId)
                if (fit && term && containerElement && latestConn) {
                  // Ensure container is visible before fitting
                  const containerStyle = window.getComputedStyle(containerElement)
                  if (containerStyle.display !== 'none' && containerElement.offsetWidth > 0 && containerElement.offsetHeight > 0) {
                    try {
                      fit.fit()
                    } catch (error) {
                      console.warn(`[Terminal Init ${initConnId}] Fit failed, will retry:`, error)
                      // Retry after a short delay
                      setTimeout(() => {
                        if (fit && term && containerElement) {
                          try {
                            fit.fit()
                          } catch (e) {
                            console.error(`[Terminal Init ${initConnId}] Fit retry failed:`, e)
                          }
                        }
                      }, 100)
                    }
                  } else {
                    // If container not ready, wait a bit more
                    setTimeout(() => {
                      if (fit && term && containerElement && containerElement.offsetWidth > 0) {
                        try {
                          fit.fit()
                        } catch (e) {
                          console.error(`[Terminal Init ${initConnId}] Fit failed in timeout:`, e)
                        }
                      }
                    }, 100)
                  }
                  
                  // Reset terminal completely to ensure clean state
                  // This clears all content and resets cursor to top-left
                  term.reset()
                  // Force scroll to top using DOM manipulation
                  setTimeout(() => {
                    const viewport = term.element?.querySelector('.xterm-viewport')
                    if (viewport) {
                      viewport.scrollTop = 0
                    }
                    term.scrollToTop()
                  }, 0)
                  
                  // Setup terminal input handler with latest WebSocket reference
                  // Use the captured initConnId to ensure input goes to correct connection
                  // IMPORTANT: Check if handler already set to prevent duplicate handlers
                  const terminalConnId = initConnId
                  // @ts-ignore - Add custom property to track if handler is set
                  if (term._inputHandlerSet) {
                    console.log(`[Terminal Init ${terminalConnId}] Input handler already set, skipping`)
                    return prev
                  }
                  // @ts-ignore - Mark that handler has been set
                  term._inputHandlerSet = true
                  console.log(`[Terminal Init ${terminalConnId}] Setting up input handler`)
                  term.onData((data) => {
                    // Get the latest connection state when input is received
                    setConnections((prev2) => {
                      const currentConn = prev2.get(terminalConnId)
                      if (!currentConn) {
                        console.warn(`[Terminal ${terminalConnId}] Connection not found when sending input`)
                        return prev2
                      }
                      if (currentConn?.ws && currentConn.ws.readyState === WebSocket.OPEN) {
                        try {
                          console.log(`[Terminal ${terminalConnId}] Sending input to WebSocket`)
                          currentConn.ws.send(JSON.stringify({
                            type: 'input',
                            data: { data },
                          }))
                        } catch (error) {
                          console.error(`[Terminal ${terminalConnId}] Failed to send input:`, error)
                        }
                      } else {
                        console.warn(`[Terminal ${terminalConnId}] WebSocket not ready (state: ${currentConn.ws?.readyState})`)
                      }
                      return prev2
                    })
                  })
                  
                  // Only focus if this is the active connection
                  if (activeConnectionId === initConnId) {
                    term.focus()
                  }
                  // Refresh terminal only if rows/cols are properly initialized
                  // fit.fit() already triggers a refresh, so we only need to refresh if needed
                  setTimeout(() => {
                    if (term && term.rows && term.cols && term.rows > 0 && term.cols > 0) {
                      try {
                        // Only refresh if terminal dimensions are valid integers
                        const startRow = Math.floor(0)
                        const endRow = Math.floor(term.rows - 1)
                        if (startRow >= 0 && endRow >= startRow && Number.isInteger(endRow)) {
                          term.refresh(startRow, endRow)
                        }
                      } catch (error) {
                        console.warn(`[Terminal Init ${terminalConnId}] Refresh failed:`, error)
                      }
                    }
                  }, 50)
                  console.log(`[Terminal Init ${terminalConnId}] Terminal setup complete`)
                  // Clear initializing flag after setup is complete
                  initializingTerminalsRef.current.delete(terminalConnId)
                }
                return prev
              })
            })
          })
        }
      }
    })
  }, [connections, mode])

  // Focus active terminal when switching tabs and auto-fit to container
  useEffect(() => {
    if (activeConnectionId) {
      const conn = connections.get(activeConnectionId)
      if (conn?.terminal && conn?.fitAddon) {
        // Use requestAnimationFrame for smoother fitting after tab switch
        requestAnimationFrame(() => {
          setTimeout(() => {
            try {
              // Fit terminal to container
              conn.fitAddon?.fit()
              
              // Send resize message to server with updated dimensions
              if (conn.ws && conn.ws.readyState === WebSocket.OPEN && conn.terminal) {
                const cols = conn.terminal.cols
                const rows = conn.terminal.rows
                console.log(`[Tab Switch ${activeConnectionId}] Sending resize: ${cols}x${rows}`)
                conn.ws.send(JSON.stringify({
                  type: 'resize',
                  data: { cols, rows },
                }))
              }
              
              // Focus terminal
              conn.terminal?.focus()
            } catch (error) {
              console.warn(`[Tab Switch ${activeConnectionId}] Failed to fit terminal:`, error)
            }
          }, 50)
        })
      }
    }
  }, [activeConnectionId])

  // Terminal always uses dark theme - no need to update on mode change
  // This provides better readability for terminal content regardless of app theme

  const fetchAssets = async () => {
    try {
      const response = await assetApi.list({ page: 1, page_size: 1000, status: 'active' })
      const assetsList = response.data.data || []
      setAssets(assetsList)
      setFilteredAssets(assetsList) // Initialize filtered assets
    } catch (error: any) {
      message.error('获取资产列表失败')
    }
  }

  const fetchProjects = async () => {
    try {
      const response = await projectApi.list()
      setProjects(Array.isArray(response.data) ? response.data : [])
    } catch (error: any) {
      console.error('获取项目列表失败:', error)
    }
  }

  const fetchEnvironments = async () => {
    try {
      const response = await environmentApi.list()
      setEnvironments(Array.isArray(response.data) ? response.data : [])
    } catch (error: any) {
      console.error('获取环境列表失败:', error)
    }
  }


  const handleConnect = async (asset: Asset, username?: string) => {
    // Check connection limit
    if (connections.size >= MAX_CONNECTIONS) {
      message.warning(`最多支持 ${MAX_CONNECTIONS} 个并发连接，请先关闭其他连接`)
      return
    }

    // Allow multiple connections to the same host
    // Create new connection with unique ID (asset.id + timestamp + random)
    const connId = `${asset.id}-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
    const containerRef = React.createRef<HTMLDivElement>()
    
    const newConnection: SSHConnection = {
      id: connId,
      asset,
      terminal: null,
      fitAddon: null,
      terminalRef: null,
      ws: null,
      status: 'connecting',
      createdAt: Date.now(),
      containerRef,
    }

    setConnections((prev) => new Map(prev).set(connId, newConnection))
    setConnectionOrder((prev) => [...prev, connId]) // Add to order array
    setActiveConnectionId(connId)

    const token = localStorage.getItem('token')
    if (!token) {
      message.error('未登录')
      navigate('/login')
      setConnections((prev) => {
        const newConnections = new Map(prev)
        newConnections.delete(connId)
        return newConnections
      })
      return
    }

    // Build WebSocket URL
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws/ssh/connect?token=${encodeURIComponent(token)}`

    try {
      const ws = new WebSocket(wsUrl)

      // Update connection with WebSocket
      setConnections((prev) => {
        const newConnections = new Map(prev)
        const conn = newConnections.get(connId)
        if (conn) {
          newConnections.set(connId, {
            ...conn,
            ws,
          })
        }
        return newConnections
      })

      ws.onopen = () => {
        console.log('WebSocket opened, sending connect message')
        try {
          // Get initial terminal size from container or use defaults
          // This ensures server creates PTY with correct dimensions from the start
          let cols = 80
          let rows = 24
          
          // Try to get size from existing terminal if already initialized
          setConnections((prev) => {
            const conn = prev.get(connId)
            if (conn?.terminal) {
              cols = conn.terminal.cols || 80
              rows = conn.terminal.rows || 24
            } else if (conn?.containerRef?.current) {
              // Estimate size based on container dimensions
              // Using approximate character dimensions for 13px font
              const container = conn.containerRef.current
              const charWidth = 8.4  // Approximate width of monospace character
              const charHeight = 18.2 // Approximate height with line-height 1.4
              cols = Math.max(80, Math.floor(container.clientWidth / charWidth))
              // Subtract extra space for padding to ensure last line is visible
              const availableHeight = container.clientHeight - 12 // Account for container padding
              rows = Math.max(24, Math.floor(availableHeight / charHeight))
            }
            return prev
          })
          
          console.log(`[Connect ${connId}] Sending initial size: ${cols}x${rows}`)
          ws.send(JSON.stringify({
            type: 'connect',
            data: {
              asset_id: asset.id,
              username: username || asset.ssh_user || '',
              auth_type: 'key',
              size: { cols, rows },
            },
          }))
        } catch (error) {
          console.error('Failed to send connect message:', error)
          message.error('发送连接请求失败')
          ws.close()
          setConnections((prev) => {
            const newConnections = new Map(prev)
            const conn = newConnections.get(connId)
            if (conn) {
              newConnections.set(connId, { ...conn, status: 'error' })
            }
            return newConnections
          })
        }
      }

      // Create message handler with explicitly captured connId
      const messageConnId = connId
      ws.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data)
          console.log(`[WebSocket ${messageConnId}] Message received:`, msg.type, msg)

          if (msg.type === 'ready') {
            console.log('WebSocket ready')
          } else if (msg.type === 'connected') {
            console.log('SSH connected')
            
            // Update connection status - don't clear terminal, preserve server output
            setConnections((prev) => {
              const newConnections = new Map(prev)
              const conn = newConnections.get(messageConnId)
              if (conn) {
                const prevStatus = conn.status
                // Update status
                const updatedConn = {
                  ...conn,
                  status: 'connected' as const,
                }
                newConnections.set(messageConnId, updatedConn)
                console.log(`[WebSocket ${messageConnId}] Connection status updated to connected`)
                
                // If terminal is initialized and status changed from connecting to connected
                // Fit terminal and scroll to top after initial output
                if (updatedConn.terminal && updatedConn.fitAddon && prevStatus === 'connecting') {
                  // Use setTimeout to ensure DOM is ready
                  setTimeout(() => {
                    try {
                      // Fit terminal to container
                      updatedConn.fitAddon?.fit()
                      console.log(`[WebSocket ${messageConnId}] Terminal fitted after connection`)
                      
                      // Send resize message to server with new dimensions
                      if (updatedConn.ws && updatedConn.ws.readyState === WebSocket.OPEN && updatedConn.terminal) {
                        const cols = updatedConn.terminal.cols
                        const rows = updatedConn.terminal.rows
                        console.log(`[WebSocket ${messageConnId}] Sending resize: ${cols}x${rows}`)
                        updatedConn.ws.send(JSON.stringify({
                          type: 'resize',
                          data: { cols, rows },
                        }))
                      }
                      
                      // Scroll to top after initial output has time to arrive
                      // Use longer delay (1.5s) to ensure all welcome messages are received
                      setTimeout(() => {
                        if (updatedConn.terminal) {
                          // First, scroll to very top
                          updatedConn.terminal.scrollToTop()
                          
                          // Also force viewport scroll
                          const viewport = updatedConn.terminal.element?.querySelector('.xterm-viewport')
                          if (viewport) {
                            ;(viewport as HTMLElement).scrollTop = 0
                          }
                          
                          // Use scrollLines to ensure we're at the beginning
                          updatedConn.terminal.scrollLines(-10000)
                          
                          console.log(`[WebSocket ${messageConnId}] Scrolled to top after connection`)
                        }
                      }, 1500)
                      
                      // Focus terminal if it's the active one
                      setActiveConnectionId((currentActiveId) => {
                        if (currentActiveId === messageConnId && updatedConn.terminal) {
                          updatedConn.terminal.focus()
                        }
                        return currentActiveId
                      })
                    } catch (error) {
                      console.warn(`[WebSocket ${messageConnId}] Failed to fit terminal:`, error)
                    }
                  }, 100)
                }
              } else {
                console.error(`[WebSocket ${messageConnId}] Connection not found when updating status`)
              }
              return newConnections
            })
          } else if (msg.type === 'output' || msg.type === 'data') {
            // IMPORTANT: Use the captured connId to ensure messages go to correct terminal
            // The connId is captured in the closure when ws.onmessage is set
            const capturedConnId = messageConnId // Use the captured connId from message handler
            setConnections((prev) => {
              const conn = prev.get(capturedConnId)
              if (!conn) {
                console.warn(`[WebSocket ${capturedConnId}] Connection not found in state`)
                return prev
              }
              if (!conn.terminal) {
                console.warn(`[WebSocket ${capturedConnId}] Terminal not initialized yet`)
                return prev
              }
              
              // Auto-update status to 'connected' when receiving output (if still connecting)
              // This handles cases where server sends output before 'connected' message
              const newConnections = new Map(prev)
              let updatedConn = conn
              let shouldUpdateStatus = false
              
              if (conn.status === 'connecting') {
                // If we're receiving output, connection is likely established
                // Update status to connected automatically
                updatedConn = {
                  ...conn,
                  status: 'connected' as const,
                }
                newConnections.set(capturedConnId, updatedConn)
                shouldUpdateStatus = true
                console.log(`[WebSocket ${capturedConnId}] Auto-updated status to connected (received output)`)
                
                // Auto-fit terminal when status changes to connected
                if (updatedConn.terminal && updatedConn.fitAddon) {
                  setTimeout(() => {
                    try {
                      updatedConn.fitAddon?.fit()
                      // Send resize to server
                      if (updatedConn.ws && updatedConn.ws.readyState === WebSocket.OPEN && updatedConn.terminal) {
                        const cols = updatedConn.terminal.cols
                        const rows = updatedConn.terminal.rows
                        console.log(`[WebSocket ${capturedConnId}] Sending resize after auto-connect: ${cols}x${rows}`)
                        updatedConn.ws.send(JSON.stringify({
                          type: 'resize',
                          data: { cols, rows },
                        }))
                      }
                    } catch (error) {
                      console.warn(`[WebSocket ${capturedConnId}] Failed to fit after auto-connect:`, error)
                    }
                  }, 100)
                }
              }
              
              const data = msg.data || msg.content || ''
              const dataStr = typeof data === 'string' ? data : String(data)
              // Write to this specific terminal instance directly
              if (dataStr) {
                updatedConn.terminal.write(dataStr)
                // Let xterm.js handle scrolling naturally
                // It will auto-scroll to bottom only if viewport was already at bottom
                // Focus terminal if status was just updated
                if (shouldUpdateStatus) {
                  setTimeout(() => {
                    setActiveConnectionId((currentActiveId) => {
                      if (currentActiveId === capturedConnId && updatedConn.terminal) {
                        updatedConn.terminal.focus()
                      }
                      return currentActiveId
                    })
                  }, 0)
                }
              }
              return shouldUpdateStatus ? newConnections : prev
            })
          } else if (msg.type === 'error') {
            const errorMsg = msg.data || msg.message || '连接错误'
            console.error('SSH error:', errorMsg)
            message.error(errorMsg)
            
            setConnections((prev) => {
              const newConnections = new Map(prev)
              const conn = newConnections.get(messageConnId)
              if (conn) {
                newConnections.set(messageConnId, { ...conn, status: 'error' })
                if (conn.terminal) {
                  conn.terminal.writeln(`\r\n错误: ${errorMsg}`)
                }
              }
              return newConnections
            })
            
            ws.close()
          } else if (msg.type === 'disconnected') {
            console.log('SSH disconnected')
            setConnections((prev) => {
              const newConnections = new Map(prev)
              const conn = newConnections.get(messageConnId)
              if (conn) {
                newConnections.set(messageConnId, { ...conn, status: 'disconnected' })
                if (conn.terminal) {
                  const reason = msg.data?.reason || msg.message || ''
                  conn.terminal.writeln(`\r\n\x1b[33m⚠ SSH 会话已结束${reason ? `: ${reason}` : ''}\x1b[0m`)
                }
              }
              return newConnections
            })
          }
        } catch (error) {
          // If parsing fails, treat as plain text output
          // Only write if connection is connected to avoid confusion
          console.log('Non-JSON message, treating as text output')
          setConnections((prev) => {
            const conn = prev.get(messageConnId)
            if (conn?.terminal && conn.status === 'connected') {
              // Write plain text directly to terminal
              const dataStr = typeof event.data === 'string' ? event.data : String(event.data)
              conn.terminal.write(dataStr)
            }
            return prev
          })
        }
      }

      ws.onerror = (error) => {
        console.error(`[WebSocket ${messageConnId}] Error:`, error)
        message.error('WebSocket 连接错误，请检查网络连接')
        setConnections((prev) => {
          const newConnections = new Map(prev)
          const conn = newConnections.get(messageConnId)
          if (conn) {
            newConnections.set(messageConnId, { ...conn, status: 'error' })
          }
          return newConnections
        })
      }

      ws.onclose = (event) => {
        console.log(`[WebSocket ${messageConnId}] Closed:`, event.code, event.reason)
        
        setConnections((prev) => {
          const newConnections = new Map(prev)
          const conn = newConnections.get(messageConnId)
          if (conn) {
            newConnections.set(messageConnId, { ...conn, status: 'disconnected', ws: null })
          }
          return newConnections
        })
        
        // Show message for unexpected disconnections (not user-initiated)
        if (event.code !== 1000 && event.code !== 1001) {
          const conn = connections.get(connId)
          if (conn?.terminal) {
            conn.terminal.writeln('\r\n\x1b[31m⚠ 连接已断开\x1b[0m')
          }
          if (!event.wasClean) {
            message.info('连接已关闭')
          }
        }
      }
      
      // Set connection timeout
      const connectionTimeout = setTimeout(() => {
        if (ws.readyState === WebSocket.CONNECTING) {
          console.warn('Connection timeout')
          message.error('连接超时，请重试')
          ws.close()
          setConnections((prev) => {
            const newConnections = new Map(prev)
            newConnections.delete(connId)
            return newConnections
          })
          if (activeConnectionId === connId) {
            setActiveConnectionId(null)
          }
        }
      }, 10000)

      ws.addEventListener('open', () => {
        clearTimeout(connectionTimeout)
      })
      
    } catch (error) {
      console.error(`[Connection ${connId}] Connection error:`, error)
      message.error(`连接失败: ${error instanceof Error ? error.message : '未知错误'}`)
      setConnections((prev) => {
        const newConnections = new Map(prev)
        newConnections.delete(connId)
        return newConnections
      })
      setActiveConnectionId((currentActiveId) => {
        if (currentActiveId === connId) {
          return null
        }
        return currentActiveId
      })
    }
  }


  const handleAssetClick = (asset: Asset) => {
    // Always create a new connection when clicking an asset
    // This allows multiple connections to the same host
    handleConnect(asset)
  }

  // Get connection count for an asset
  const getConnectionCount = (assetId: number): number => {
    return Array.from(connections.values()).filter(
      (conn) => conn.asset.id === assetId && conn.status === 'connected'
    ).length
  }

  // Check if asset has active connections
  const hasActiveConnection = (assetId: number): boolean => {
    return getConnectionCount(assetId) > 0
  }

  const handleDisconnect = (connId?: string) => {
    const targetConnId = connId || activeConnectionId
    if (!targetConnId) return

    const connection = connections.get(targetConnId)
    if (!connection) return

    // Close WebSocket
    if (connection.ws) {
      try {
        if (connection.ws.readyState === WebSocket.OPEN) {
          connection.ws.send(JSON.stringify({ type: 'disconnect' }))
        }
      } catch (error) {
        console.error('Failed to send disconnect message:', error)
      }
      connection.ws.close(1000, 'User disconnected')
    }

    // Dispose terminal
    if (connection.terminal) {
      try {
        connection.terminal.onData(() => {}) // Remove input handler
        connection.terminal.dispose()
      } catch (error) {
        console.error('Failed to dispose terminal:', error)
      }
    }

    // Remove connection from both connections map and order array
    setConnections((prev) => {
      const newConnections = new Map(prev)
      newConnections.delete(targetConnId)
      return newConnections
    })
    setConnectionOrder((prev) => prev.filter((id) => id !== targetConnId))

    // If this was the active connection, switch to another or clear
    if (activeConnectionId === targetConnId) {
      const remainingConnections = connectionOrder.filter((id) => id !== targetConnId)
      if (remainingConnections.length > 0) {
        setActiveConnectionId(remainingConnections[0])
      } else {
        setActiveConnectionId(null)
      }
    }
  }

  const handleCloseAllConnections = () => {
    connections.forEach((conn, connId) => {
      handleDisconnect(connId)
    })
  }

  // Handle user logout
  const handleLogout = async () => {
    try {
      await authApi.logout()
    } catch (error) {
      // Ignore API errors, proceed with client-side logout
    }
    clearAuth()
    message.success('已成功登出')
    navigate('/login')
  }

  // User menu items for dropdown
  const userMenuItems: MenuProps['items'] = [
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '登出',
      onClick: handleLogout,
    },
  ]

  // Close all connections except the specified one
  const handleCloseOthers = (connIdToKeep: string) => {
    connections.forEach((conn, connId) => {
      if (connId !== connIdToKeep) {
        handleDisconnect(connId)
      }
    })
  }

  // Clone connection - create a new connection to the same asset
  const handleCloneConnection = (connId: string) => {
    const conn = connections.get(connId)
    if (!conn) return
    
    // Check connection limit
    if (connections.size >= MAX_CONNECTIONS) {
      message.warning(`最多支持 ${MAX_CONNECTIONS} 个并发连接，请先关闭其他连接`)
      return
    }
    
    // Create new connection to the same asset
    handleConnect(conn.asset)
    message.success(`正在克隆连接到 ${conn.asset.hostName}`)
  }

  // Reconnect - disconnect and reconnect to the same asset
  const handleReconnect = (connId: string) => {
    const conn = connections.get(connId)
    if (!conn) return
    
    const asset = conn.asset
    
    message.loading({ content: `正在重连到 ${asset.hostName}...`, key: `reconnect-${connId}`, duration: 0 })
    
    // Disconnect current connection
    handleDisconnect(connId)
    
    // Reconnect to the same asset after a short delay to ensure cleanup completes
    // handleConnect will automatically set the new connection as active
    setTimeout(() => {
      handleConnect(asset, asset.ssh_user || undefined)
      message.success({ content: `已重连到 ${asset.hostName}`, key: `reconnect-${connId}` })
    }, 300)
  }

  // Get display name for a tab with auto-numbering for same host
  const getTabDisplayName = (connId: string): string => {
    const conn = connections.get(connId)
    if (!conn) return ''
    
    // Find all connections to the same host
    const sameHostConnections = connectionOrder.filter(id => {
      const c = connections.get(id)
      return c && c.asset.id === conn.asset.id
    })
    
    // If only one connection to this host, just show hostname
    if (sameHostConnections.length <= 1) {
      return conn.asset.hostName
    }
    
    // Multiple connections - add index
    const index = sameHostConnections.indexOf(connId) + 1
    return `${conn.asset.hostName} (${index})`
  }

  // Handle drag start
  const handleDragStart = (e: React.DragEvent<HTMLDivElement>, connId: string) => {
    setDraggedTabId(connId)
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', connId)
    // Add drag visual
    if (e.currentTarget) {
      e.currentTarget.style.opacity = '0.5'
    }
  }

  // Handle drag end
  const handleDragEnd = (e: React.DragEvent<HTMLDivElement>) => {
    setDraggedTabId(null)
    if (e.currentTarget) {
      e.currentTarget.style.opacity = '1'
    }
  }

  // Handle drag over
  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
  }

  // Handle drop - reorder tabs
  const handleDrop = (e: React.DragEvent<HTMLDivElement>, targetConnId: string) => {
    e.preventDefault()
    const draggedId = draggedTabId
    if (!draggedId || draggedId === targetConnId) return

    setConnectionOrder(prevOrder => {
      const newOrder = [...prevOrder]
      const draggedIndex = newOrder.indexOf(draggedId)
      const targetIndex = newOrder.indexOf(targetConnId)
      
      if (draggedIndex === -1 || targetIndex === -1) return prevOrder
      
      // Remove dragged item and insert at target position
      newOrder.splice(draggedIndex, 1)
      newOrder.splice(targetIndex, 0, draggedId)
      
      return newOrder
    })
  }

  // Get right-click menu items for a tab
  const getTabContextMenuItems = (connId: string): MenuProps['items'] => {
    const otherConnectionsCount = connections.size - 1
    const canClone = connections.size < MAX_CONNECTIONS
    return [
      {
        key: 'clone',
        label: '克隆',
        icon: <CopyOutlined />,
        disabled: !canClone,
        onClick: () => handleCloneConnection(connId),
      },
      {
        key: 'reconnect',
        label: '重连',
        icon: <ReloadOutlined />,
        onClick: () => handleReconnect(connId),
      },
      {
        type: 'divider',
      },
      {
        key: 'close',
        label: '关闭',
        icon: <CloseOutlined />,
        onClick: () => handleDisconnect(connId),
      },
      {
        key: 'closeOthers',
        label: '关闭其他',
        icon: <CloseCircleOutlined />,
        disabled: otherConnectionsCount === 0,
        onClick: () => handleCloseOthers(connId),
      },
      {
        key: 'closeAll',
        label: '关闭所有',
        icon: <DisconnectOutlined />,
        onClick: () => handleCloseAllConnections(),
      },
    ]
  }

  // Build tree data structure: Project -> Environment -> Asset
  const buildTreeData = () => {
    if (!searchQuery.trim()) {
      // Group assets by project and environment
      const treeData: any[] = []
      
      // Get project map and environment map
      const projectMap = new Map<number, Project>()
      projects.forEach(p => projectMap.set(p.id, p))
      
      const envMap = new Map<number, Environment>()
      environments.forEach(e => envMap.set(e.id, e))
      
      // Group assets by project_id
      const assetsByProject = new Map<number, Asset[]>()
      const assetsWithoutProject: Asset[] = []
      
      filteredAssets.forEach(asset => {
        if (asset.project_id && projectMap.has(asset.project_id)) {
          if (!assetsByProject.has(asset.project_id)) {
            assetsByProject.set(asset.project_id, [])
          }
          assetsByProject.get(asset.project_id)!.push(asset)
        } else {
          assetsWithoutProject.push(asset)
        }
      })
      
      // Build project nodes
      assetsByProject.forEach((projectAssets, projectId) => {
        const project = projectMap.get(projectId)!
        
        // Group assets by environment within this project
        const assetsByEnv = new Map<number, Asset[]>()
        const assetsWithoutEnv: Asset[] = []
        
        projectAssets.forEach(asset => {
          if (asset.environment_id && envMap.has(asset.environment_id)) {
            if (!assetsByEnv.has(asset.environment_id)) {
              assetsByEnv.set(asset.environment_id, [])
            }
            assetsByEnv.get(asset.environment_id)!.push(asset)
          } else {
            assetsWithoutEnv.push(asset)
          }
        })
        
        const envNodes: any[] = []
        
        // Build environment nodes
        assetsByEnv.forEach((envAssets, envId) => {
          const env = envMap.get(envId)!
          envNodes.push({
            title: env.name,
            key: `env-${envId}`,
            icon: <FolderOutlined style={{ 
              color: mode === 'dark' ? '#94A3B8' : '#64748B',
              fontSize: 12,
            }} />,
            children: envAssets.map(asset => ({
              title: asset.hostName,
              key: `asset-${asset.id}`,
            icon: <DatabaseOutlined style={{ 
              color: mode === 'dark' ? '#94A3B8' : '#64748B',
              fontSize: 12,
            }} />,
              isLeaf: true,
              asset: asset,
              count: envAssets.length,
            })),
            count: envAssets.length,
          })
        })
        
        // Add assets without environment
        if (assetsWithoutEnv.length > 0) {
          envNodes.push({
            title: '未分类环境',
            key: `project-${projectId}-no-env`,
            icon: <FolderOutlined style={{ 
              color: mode === 'dark' ? '#94A3B8' : '#64748B',
              fontSize: 12,
            }} />,
            children: assetsWithoutEnv.map(asset => ({
              title: asset.hostName,
              key: `asset-${asset.id}`,
            icon: <DatabaseOutlined style={{ 
              color: mode === 'dark' ? '#94A3B8' : '#64748B',
              fontSize: 12,
            }} />,
              isLeaf: true,
              asset: asset,
            })),
            count: assetsWithoutEnv.length,
          })
        }
        
        const totalAssets = projectAssets.length
        treeData.push({
          title: project.name,
          key: `project-${projectId}`,
            icon: <FolderOpenOutlined style={{ 
              color: mode === 'dark' ? '#94A3B8' : '#64748B',
              fontSize: 13,
            }} />,
          children: envNodes,
          count: totalAssets,
        })
      })
      
      // Add assets without project
      if (assetsWithoutProject.length > 0) {
        // Group unassigned assets by environment
        const unassignedByEnv = new Map<number, Asset[]>()
        const unassignedWithoutEnv: Asset[] = []
        
        assetsWithoutProject.forEach(asset => {
          if (asset.environment_id && envMap.has(asset.environment_id)) {
            if (!unassignedByEnv.has(asset.environment_id)) {
              unassignedByEnv.set(asset.environment_id, [])
            }
            unassignedByEnv.get(asset.environment_id)!.push(asset)
          } else {
            unassignedWithoutEnv.push(asset)
          }
        })
        
        const unassignedEnvNodes: any[] = []
        
        unassignedByEnv.forEach((envAssets, envId) => {
          const env = envMap.get(envId)!
          unassignedEnvNodes.push({
            title: env.name,
            key: `no-project-env-${envId}`,
            icon: <FolderOutlined style={{ 
              color: mode === 'dark' ? '#94A3B8' : '#64748B',
              fontSize: 12,
            }} />,
            children: envAssets.map(asset => ({
              title: asset.hostName,
              key: `asset-${asset.id}`,
            icon: <DatabaseOutlined style={{ 
              color: mode === 'dark' ? '#94A3B8' : '#64748B',
              fontSize: 12,
            }} />,
              isLeaf: true,
              asset: asset,
            })),
            count: envAssets.length,
          })
        })
        
        if (unassignedWithoutEnv.length > 0) {
          unassignedEnvNodes.push({
            title: '未分类环境',
            key: 'no-project-no-env',
            icon: <FolderOutlined style={{ 
              color: mode === 'dark' ? '#94A3B8' : '#64748B',
              fontSize: 12,
            }} />,
            children: unassignedWithoutEnv.map(asset => ({
              title: asset.hostName,
              key: `asset-${asset.id}`,
            icon: <DatabaseOutlined style={{ 
              color: mode === 'dark' ? '#94A3B8' : '#64748B',
              fontSize: 12,
            }} />,
              isLeaf: true,
              asset: asset,
            })),
            count: unassignedWithoutEnv.length,
          })
        }
        
        if (unassignedEnvNodes.length > 0) {
          const totalUnassigned = assetsWithoutProject.length
          treeData.push({
            title: '未分类项目',
            key: 'no-project',
            icon: <FolderOpenOutlined style={{ 
              color: mode === 'dark' ? '#94A3B8' : '#64748B',
              fontSize: 13,
            }} />,
            children: unassignedEnvNodes,
            count: totalUnassigned,
          })
        }
      }
      
      // Auto expand all on first load (including root)
      // Only set if we haven't initialized expandedKeys yet
      if (expandedKeys.length === 0 && treeData.length > 0) {
        const allKeys: React.Key[] = ['root'] // Root is always expanded by default
        treeData.forEach((project: any) => {
          allKeys.push(project.key)
          if (project.children) {
            project.children.forEach((env: any) => {
              allKeys.push(env.key)
            })
          }
        })
        // Use setTimeout to avoid calling setState during render
        setTimeout(() => {
          setExpandedKeys(allKeys)
          setAutoExpandParent(true)
        }, 0)
      }
      
      // Wrap all tree data under a root node
      const totalCount = filteredAssets.length
      return [{
        title: '根目录',
        key: 'root',
        icon: <HomeOutlined style={{ 
          color: mode === 'dark' ? '#94A3B8' : '#64748B',
          fontSize: 14,
        }} />,
        children: treeData,
        count: totalCount,
      }]
    } else {
      // If searching, show flat list of matching assets
      return filteredAssets.map(asset => ({
        title: `${asset.hostName} (${asset.ip})`,
        key: `asset-${asset.id}`,
        icon: <DatabaseOutlined style={{ color: mode === 'dark' ? '#94A3B8' : '#64748B' }} />,
        isLeaf: true,
        asset: asset,
      }))
    }
  }

  const onExpand = (expandedKeysValue: React.Key[]) => {
    // Check if root node expansion changed (only when not searching)
    if (!searchQuery.trim()) {
      const wasRootExpanded = expandedKeys.includes('root')
      const isRootExpanded = expandedKeysValue.includes('root')
      
      // If root node was collapsed, collapse all children
      if (wasRootExpanded && !isRootExpanded) {
        setExpandedKeys([])
        setAutoExpandParent(false)
        return
      }
      
      // If root node was expanded, expand all projects and environments
      if (!wasRootExpanded && isRootExpanded) {
        const allKeys: React.Key[] = ['root']
        // Build keys from current data without calling buildTreeData recursively
        const projectMap = new Map<number, Project>()
        projects.forEach(p => projectMap.set(p.id, p))
        const envMap = new Map<number, Environment>()
        environments.forEach(e => envMap.set(e.id, e))
        
        // Group assets by project
        const assetsByProject = new Map<number, Asset[]>()
        filteredAssets.forEach(asset => {
          if (asset.project_id && projectMap.has(asset.project_id)) {
            if (!assetsByProject.has(asset.project_id)) {
              assetsByProject.set(asset.project_id, [])
            }
            assetsByProject.get(asset.project_id)!.push(asset)
          }
        })
        
        // Add project keys and environment keys
        assetsByProject.forEach((projectAssets, projectId) => {
          allKeys.push(`project-${projectId}`)
          const assetsByEnv = new Map<number, Asset[]>()
          projectAssets.forEach(asset => {
            if (asset.environment_id && envMap.has(asset.environment_id)) {
              if (!assetsByEnv.has(asset.environment_id)) {
                assetsByEnv.set(asset.environment_id, [])
              }
              assetsByEnv.get(asset.environment_id)!.push(asset)
            }
          })
          assetsByEnv.forEach((_, envId) => {
            allKeys.push(`env-${envId}`)
          })
          // Add no-env keys
          const hasNoEnv = projectAssets.some(asset => !asset.environment_id || !envMap.has(asset.environment_id))
          if (hasNoEnv) {
            allKeys.push(`project-${projectId}-no-env`)
          }
        })
        
        // Add unassigned project keys
        const hasUnassigned = filteredAssets.some(asset => !asset.project_id || !projectMap.has(asset.project_id))
        if (hasUnassigned) {
          allKeys.push('no-project')
          // Add unassigned environment keys
          const unassignedByEnv = new Map<number, Asset[]>()
          filteredAssets.forEach(asset => {
            if ((!asset.project_id || !projectMap.has(asset.project_id)) && asset.environment_id && envMap.has(asset.environment_id)) {
              if (!unassignedByEnv.has(asset.environment_id)) {
                unassignedByEnv.set(asset.environment_id, [])
              }
              unassignedByEnv.get(asset.environment_id)!.push(asset)
            }
          })
          unassignedByEnv.forEach((_, envId) => {
            allKeys.push(`no-project-env-${envId}`)
          })
          const hasUnassignedNoEnv = filteredAssets.some(asset => 
            (!asset.project_id || !projectMap.has(asset.project_id)) && (!asset.environment_id || !envMap.has(asset.environment_id))
          )
          if (hasUnassignedNoEnv) {
            allKeys.push('no-project-no-env')
          }
        }
        
        setExpandedKeys(allKeys)
        setAutoExpandParent(false)
        return
      }
    }
    
    // Normal expand/collapse behavior
    setExpandedKeys(expandedKeysValue)
    setAutoExpandParent(false)
  }

  const onSelect = (selectedKeys: React.Key[], info: any) => {
    if (info.node?.asset) {
      handleAssetClick(info.node.asset)
    }
  }

  const headerBg = mode === 'dark' ? '#0A0E27' : '#FFFFFF'
  const sidebarBg = mode === 'dark' ? '#121212' : '#FAFAFA'
  const terminalBg = '#000000' // Pure black for maximum terminal contrast

  return (
    <Layout style={{ height: '100vh', background: mode === 'dark' ? '#0A0E27' : '#F5F5F5' }}>
      <Header
        style={{
          padding: '0 24px',
          background: headerBg,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          borderBottom: mode === 'dark' 
            ? '1px solid rgba(255, 255, 255, 0.08)' 
            : '1px solid rgba(0, 0, 0, 0.06)',
          boxShadow: mode === 'dark'
            ? '0 2px 8px rgba(0, 0, 0, 0.3)'
            : '0 1px 2px rgba(0, 0, 0, 0.04)',
          position: 'relative',
          zIndex: 10,
        }}
      >
        {/* 左侧：Logo 和返回按钮 */}
        <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
          <img
            src={mode === 'dark' ? '/logo-dark.svg' : '/logo-light.svg'}
            alt="KkOps"
            style={{
              height: 28,
              width: 'auto',
            }}
          />
          <div style={{ 
            width: 1, 
            height: 24, 
            background: mode === 'dark' ? 'rgba(255, 255, 255, 0.12)' : 'rgba(0, 0, 0, 0.12)',
          }} />
          <Button
            type="text"
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate('/dashboard')}
            style={{ 
              color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
              transition: 'all 0.2s ease',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = mode === 'dark' ? 'rgba(255, 255, 255, 0.08)' : 'rgba(0, 0, 0, 0.04)'
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = 'transparent'
            }}
          >
            返回
          </Button>
        </div>

        {/* 右侧：连接数、主题切换和用户 */}
        <Space size="middle">
          {connections.size > 0 && (
            <>
              <Text style={{ 
                color: mode === 'dark' ? '#94A3B8' : '#64748B',
                fontSize: 12,
              }}>
                连接数: {Array.from(connections.values()).filter(c => c.status === 'connected').length}/{connections.size}
              </Text>
              {connections.size > 1 && (
                <Button
                  type="text"
                  size="small"
                  danger
                  onClick={handleCloseAllConnections}
                  style={{
                    borderRadius: 6,
                    height: 28,
                    padding: '0 12px',
                    fontSize: 12,
                  }}
                >
                  关闭全部
                </Button>
              )}
            </>
          )}
          <div
            role="button"
            tabIndex={0}
            aria-label={mode === 'dark' ? '切换到亮色模式' : '切换到暗色模式'}
            title={mode === 'dark' ? '切换到亮色模式' : '切换到暗色模式'}
            style={{
              fontSize: 18,
              cursor: 'pointer',
              color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: 32,
              height: 32,
              borderRadius: 6,
              transition: 'all 0.2s ease',
            }}
            onClick={toggleMode}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault()
                toggleMode()
              }
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = mode === 'dark' ? 'rgba(255, 255, 255, 0.08)' : 'rgba(0, 0, 0, 0.04)'
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = 'transparent'
            }}
          >
            {mode === 'dark' ? <SunOutlined /> : <MoonOutlined />}
          </div>
          {user && (
            <Dropdown
              menu={{ items: userMenuItems }}
              placement="bottomRight"
              trigger={['click']}
            >
              <div
                role="button"
                tabIndex={0}
                aria-label="用户菜单"
                aria-haspopup="menu"
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 8,
                  cursor: 'pointer',
                  padding: '4px 12px',
                  borderRadius: 6,
                  color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
                  fontSize: 14,
                  transition: 'all 0.2s ease',
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault()
                  }
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = mode === 'dark' ? 'rgba(255, 255, 255, 0.08)' : 'rgba(0, 0, 0, 0.04)'
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = 'transparent'
                }}
              >
                <UserOutlined />
                <span>{user.real_name || user.username}</span>
              </div>
            </Dropdown>
          )}
        </Space>
      </Header>
      <Layout style={{ flexDirection: 'row', background: 'transparent', position: 'relative' }}>
        <Sider
          width={360}
          style={{
            background: sidebarBg,
            borderRight: mode === 'dark' 
              ? '1px solid rgba(255, 255, 255, 0.08)' 
              : '1px solid rgba(0, 0, 0, 0.06)',
            overflow: 'hidden',
            boxShadow: mode === 'dark'
              ? '4px 0 12px rgba(0, 0, 0, 0.4)'
              : '2px 0 8px rgba(0, 0, 0, 0.04)',
            display: 'flex',
            flexDirection: 'column',
          }}
        >
          <div style={{ 
            padding: 20,
            borderBottom: mode === 'dark' 
              ? '1px solid rgba(255, 255, 255, 0.08)' 
              : '1px solid rgba(0, 0, 0, 0.06)',
          }}>
            <div style={{ 
              marginBottom: 16, 
              display: 'flex', 
              justifyContent: 'space-between', 
              alignItems: 'center' 
            }}>
              <Text strong style={{ 
                color: mode === 'dark' ? '#F1F5F9' : '#1E293B', 
                fontSize: 12,
                fontWeight: 600,
                fontFamily: '"IBM Plex Sans", sans-serif',
                letterSpacing: '-0.01em',
              }}>
                资产列表
              </Text>
              <Button
                type="text"
                size="small"
                icon={<ReloadOutlined />}
                onClick={() => {
                  fetchAssets()
                  fetchProjects()
                  fetchEnvironments()
                }}
                style={{ 
                  color: mode === 'dark' ? '#94A3B8' : '#64748B',
                  transition: 'all 0.2s ease',
                  borderRadius: 6,
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = mode === 'dark' ? 'rgba(255, 255, 255, 0.08)' : 'rgba(0, 0, 0, 0.04)'
                  e.currentTarget.style.color = mode === 'dark' ? '#F1F5F9' : '#1E293B'
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = 'transparent'
                  e.currentTarget.style.color = mode === 'dark' ? '#94A3B8' : '#64748B'
                }}
              >
                刷新
              </Button>
            </div>
            <Search
              placeholder="搜索主机名、IP 或用户名"
              allowClear
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              style={{ 
                marginBottom: 0,
              }}
            />
          </div>
          <div 
            ref={treeContainerRef}
            style={{ 
              flex: 1, 
              overflowY: 'auto', 
              overflowX: 'hidden',
              padding: '8px 12px',
            }}>
            {assets.length === 0 ? (
              <div style={{ 
                textAlign: 'center', 
                padding: '60px 20px', 
                color: mode === 'dark' ? '#64748B' : '#94A3B8',
              }}>
                <Spin />
                <div style={{ marginTop: 16 }}>
                  <Text style={{ 
                    fontSize: 13,
                    color: mode === 'dark' ? '#64748B' : '#94A3B8',
                  }}>
                    加载中...
                  </Text>
                </div>
              </div>
            ) : buildTreeData().length === 0 ? (
              <div style={{ 
                textAlign: 'center', 
                padding: '60px 20px', 
                color: mode === 'dark' ? '#64748B' : '#94A3B8',
              }}>
                <Text style={{ 
                  fontSize: 13,
                  color: mode === 'dark' ? '#64748B' : '#94A3B8',
                }}>
                  {searchQuery ? '未找到匹配的资产' : '暂无资产'}
                </Text>
              </div>
            ) : (
              <Tree
                showIcon
                defaultExpandAll={false}
                expandedKeys={expandedKeys}
                autoExpandParent={autoExpandParent}
                onExpand={onExpand}
                onSelect={onSelect}
                selectedKeys={[]} // Keep empty to allow clicking same asset multiple times
                treeData={buildTreeData()}
                blockNode
                indent={16}
                style={{
                  background: 'transparent',
                  color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
                  fontFamily: '"IBM Plex Sans", sans-serif',
                }}
                titleRender={(node: any) => {
                  if (node.asset) {
                    // Asset node - Enhanced styling
                    const asset = node.asset as Asset
                    const connectionCount = getConnectionCount(asset.id)
                    const isActive = activeConnection?.asset.id === asset.id
                    const isConnecting = Array.from(connections.values()).some(
                      (conn) => conn.asset.id === asset.id && conn.status === 'connecting'
                    )
                    
                    return (
                      <div
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: 6,
                          padding: '2px 4px',
                          margin: '1px 0',
                          cursor: 'pointer',
                          height: '24px',
                          lineHeight: '24px',
                          transition: 'background-color 0.15s ease',
                        }}
                        onClick={(e) => {
                          e.stopPropagation()
                          handleAssetClick(asset)
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.backgroundColor = mode === 'dark' 
                            ? 'rgba(255, 255, 255, 0.04)' 
                            : 'rgba(0, 0, 0, 0.02)'
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.backgroundColor = 'transparent'
                        }}
                      >
                        <Text
                          style={{
                            color: isActive
                              ? (mode === 'dark' ? '#60A5FA' : '#2563EB')
                              : (mode === 'dark' ? '#E2E8F0' : '#1E293B'),
                            fontSize: 12,
                            fontWeight: isActive ? 500 : 400,
                            fontFamily: '"IBM Plex Sans", sans-serif',
                            flex: 1,
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                          }}
                        >
                          {asset.hostName}
                        </Text>
                        {connectionCount > 0 && (
                          <Tag
                            style={{
                              margin: 0,
                              padding: '0 6px',
                              fontSize: 10,
                              height: 16,
                              lineHeight: '16px',
                              borderRadius: 8,
                              background: mode === 'dark' 
                                ? 'rgba(59, 130, 246, 0.2)' 
                                : 'rgba(59, 130, 246, 0.1)',
                              border: 'none',
                              color: mode === 'dark' ? '#60A5FA' : '#2563EB',
                              fontFamily: '"IBM Plex Sans", sans-serif',
                              flexShrink: 0,
                            }}
                          >
                            {connectionCount}
                          </Tag>
                        )}
                        {isConnecting && (
                          <Spin size="small" style={{ flexShrink: 0 }} />
                        )}
                        <Text
                          style={{
                            color: mode === 'dark' ? '#64748B' : '#94A3B8',
                            fontSize: 10,
                            fontFamily: '"IBM Plex Mono", monospace',
                            marginLeft: 4,
                            flexShrink: 0,
                          }}
                        >
                          {asset.ip}
                        </Text>
                      </div>
                    )
                  } else if (node.key === 'root') {
                    // Root node - Special styling for collapse/expand all
                    const isExpanded = expandedKeys.includes('root')
                    return (
                      <div
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: 6,
                          padding: '2px 4px',
                          margin: '1px 0',
                          cursor: 'pointer',
                          height: '24px',
                          lineHeight: '24px',
                          transition: 'background-color 0.15s ease',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.backgroundColor = mode === 'dark' 
                            ? 'rgba(255, 255, 255, 0.06)' 
                            : 'rgba(0, 0, 0, 0.04)'
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.backgroundColor = 'transparent'
                        }}
                      >
                        <Text
                          style={{
                            color: mode === 'dark' ? '#E2E8F0' : '#1E293B',
                            fontSize: 13,
                            fontWeight: 600,
                            fontFamily: '"IBM Plex Sans", sans-serif',
                            flex: 1,
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                          }}
                        >
                          {node.title}
                        </Text>
                        {node.count !== undefined && (
                          <Text
                            style={{
                              color: mode === 'dark' ? '#94A3B8' : '#64748B',
                              fontSize: 11,
                              fontFamily: '"IBM Plex Sans", sans-serif',
                              flexShrink: 0,
                            }}
                          >
                            ({node.count})
                          </Text>
                        )}
                      </div>
                    )
                  } else {
                    // Project or Environment node - Enhanced styling
                    const nodeLevel = node.key.startsWith('project-') ? 'project' : 'environment'
                    const isProject = nodeLevel === 'project'
                    
                    return (
                      <div
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: 6,
                          padding: '2px 4px',
                          margin: '1px 0',
                          cursor: 'pointer',
                          height: '24px',
                          lineHeight: '24px',
                          transition: 'background-color 0.15s ease',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.backgroundColor = mode === 'dark' 
                            ? 'rgba(255, 255, 255, 0.04)' 
                            : 'rgba(0, 0, 0, 0.02)'
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.backgroundColor = 'transparent'
                        }}
                      >
                        <Text
                          style={{
                            color: mode === 'dark' ? '#E2E8F0' : '#1E293B',
                            fontSize: isProject ? 13 : 12,
                            fontWeight: isProject ? 600 : 500,
                            fontFamily: '"IBM Plex Sans", sans-serif',
                            flex: 1,
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                          }}
                        >
                          {node.title}
                        </Text>
                        {node.count !== undefined && (
                          <Text
                            style={{
                              color: mode === 'dark' ? '#94A3B8' : '#64748B',
                              fontSize: 11,
                              fontFamily: '"IBM Plex Sans", sans-serif',
                              flexShrink: 0,
                            }}
                          >
                            ({node.count})
                          </Text>
                        )}
                      </div>
                    )
                  }
                }}
              />
            )}
          </div>
        </Sider>
        <Content
          style={{
            display: 'flex',
            flexDirection: 'column',
            background: 'transparent',
            overflow: 'hidden',
          }}
        >
          {/* Tab Bar */}
          {connections.size > 0 && (
            <div
              style={{
                padding: '8px 16px',
                borderBottom: mode === 'dark' 
                  ? '1px solid rgba(255, 255, 255, 0.08)' 
                  : '1px solid rgba(0, 0, 0, 0.06)',
                background: mode === 'dark' ? '#121212' : '#FFFFFF',
                display: 'flex',
                gap: 4,
                overflowX: 'auto',
                overflowY: 'hidden',
              }}
            >
              {connectionOrder.map((connId) => {
                const conn = connections.get(connId)
                if (!conn) return null
                const isActive = activeConnectionId === conn.id
                const isDragging = draggedTabId === conn.id
                return (
                  <Dropdown
                    key={conn.id}
                    menu={{ items: getTabContextMenuItems(conn.id) }}
                    trigger={['contextMenu']}
                  >
                    <div
                      draggable
                      onDragStart={(e) => handleDragStart(e, conn.id)}
                      onDragEnd={handleDragEnd}
                      onDragOver={handleDragOver}
                      onDrop={(e) => handleDrop(e, conn.id)}
                      onClick={() => setActiveConnectionId(conn.id)}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 8,
                        padding: '6px 12px',
                        borderRadius: '6px 6px 0 0',
                        cursor: isDragging ? 'grabbing' : 'grab',
                        background: isActive
                          ? (mode === 'dark' ? '#0A0E27' : '#F5F5F5')
                          : 'transparent',
                        borderBottom: isActive
                          ? `2px solid ${mode === 'dark' ? '#60A5FA' : '#2563EB'}`
                          : '2px solid transparent',
                        color: isActive
                          ? (mode === 'dark' ? '#F1F5F9' : '#1E293B')
                          : (mode === 'dark' ? '#94A3B8' : '#64748B'),
                        transition: 'all 0.2s ease',
                        whiteSpace: 'nowrap',
                        minWidth: 120,
                        maxWidth: 200,
                        opacity: isDragging ? 0.5 : 1,
                      }}
                      onMouseEnter={(e) => {
                        if (!isActive && !isDragging) {
                          e.currentTarget.style.background = mode === 'dark' 
                            ? 'rgba(255, 255, 255, 0.04)' 
                            : 'rgba(0, 0, 0, 0.02)'
                        }
                      }}
                      onMouseLeave={(e) => {
                        if (!isActive) {
                          e.currentTarget.style.background = 'transparent'
                        }
                      }}
                    >
                      <Text
                        style={{
                          fontSize: 12,
                          fontWeight: isActive ? 500 : 400,
                          overflow: 'hidden',
                          textOverflow: 'ellipsis',
                          flex: 1,
                        }}
                      >
                        {getTabDisplayName(conn.id)}
                      </Text>
                      {conn.status === 'connecting' && (
                        <Spin size="small" />
                      )}
                      {conn.status === 'connected' && (
                        <div
                          style={{
                            width: 8,
                            height: 8,
                            borderRadius: '50%',
                            background: '#10B981',
                          }}
                        />
                      )}
                      {conn.status === 'error' && (
                        <div
                          style={{
                            width: 8,
                            height: 8,
                            borderRadius: '50%',
                            background: '#EF4444',
                          }}
                        />
                      )}
                      <Button
                        type="text"
                        size="small"
                        icon={<CloseOutlined />}
                        onClick={(e) => {
                          e.stopPropagation()
                          handleDisconnect(conn.id)
                        }}
                        style={{
                          padding: 0,
                          width: 16,
                          height: 16,
                          minWidth: 16,
                          color: mode === 'dark' ? '#94A3B8' : '#64748B',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.color = '#EF4444'
                          e.currentTarget.style.background = mode === 'dark' 
                            ? 'rgba(239, 68, 68, 0.1)' 
                            : 'rgba(239, 68, 68, 0.05)'
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.color = mode === 'dark' ? '#94A3B8' : '#64748B'
                          e.currentTarget.style.background = 'transparent'
                        }}
                      />
                    </div>
                  </Dropdown>
                )
              })}
            </div>
          )}

          {/* Terminal Area */}
          <div
            style={{
              flex: 1,
              padding: '24px 24px 12px 24px', // Reduced bottom padding to give more space for terminal
              minHeight: 0,
              display: 'flex',
              flexDirection: 'column',
              overflow: 'hidden',
              background: mode === 'dark' ? '#0A0E27' : '#F5F5F5',
            }}
          >
            <div
              style={{
                flex: 1,
                width: '100%',
                minHeight: 0,
                display: 'flex',
                flexDirection: 'column',
                borderRadius: 8,
                boxShadow: mode === 'dark'
                  ? '0 4px 24px rgba(0, 0, 0, 0.5)'
                  : '0 2px 12px rgba(0, 0, 0, 0.08)',
                border: mode === 'dark'
                  ? '1px solid rgba(255, 255, 255, 0.08)'
                  : '1px solid rgba(0, 0, 0, 0.06)',
                overflow: 'hidden',
                position: 'relative',
                padding: 0,
                boxSizing: 'border-box',
                background: terminalBg,
              }}
            >
              {!activeConnection && (
                <div style={{
                  position: 'absolute',
                  top: '50%',
                  left: '50%',
                  transform: 'translate(-50%, -50%)',
                  textAlign: 'center',
                  color: mode === 'dark' ? '#64748B' : '#94A3B8',
                  zIndex: 10,
                  pointerEvents: 'none',
                  padding: '0 20px',
                }}>
                  <Text style={{ 
                    fontSize: 12,
                    fontFamily: '"IBM Plex Sans", sans-serif',
                  }}>
                    请从左侧列表选择资产进行连接
                  </Text>
                </div>
              )}
              {/* Render terminal containers for each connection */}
              {/* All containers use absolute positioning, only active one is visible */}
              {Array.from(connections.values()).map((conn) => {
                const isActive = activeConnectionId === conn.id
                return (
                  <div
                    key={conn.id}
                    ref={conn.containerRef}
                    style={{
                      display: isActive ? 'block' : 'none',
                      width: '100%',
                      height: '100%',
                      padding: '12px 16px 8px 16px', // Minimal padding: top 12px, bottom 8px for spacing
                      margin: 0,
                      boxSizing: 'border-box',
                      position: 'absolute',
                      top: 0,
                      left: 0,
                      right: 0,
                      bottom: 0,
                      background: 'transparent',
                      outline: 'none',
                      overflow: 'hidden',
                      zIndex: isActive ? 1 : 0,
                    }}
                    className="xterm-container"
                    tabIndex={isActive ? 0 : -1}
                    onClick={() => {
                      // Focus terminal when clicking on container
                      if (conn.terminal && isActive) {
                        conn.terminal.focus()
                      }
                    }}
                  />
                )
              })}
            </div>
          </div>
        </Content>
      </Layout>
    </Layout>
  )
}

export default WebSSHTerminal
