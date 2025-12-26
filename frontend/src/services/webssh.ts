/**
 * WebSSH Service
 * Handles WebSocket communication for WebSSH terminal connections
 */

// WebSSH 终端消息类型
export interface TerminalMessage {
  type: 'input' | 'resize' | 'ping' | 'pong' | 'username' | 'password' | 'username_request' | 'password_request' | 'error' | 'auth_method' | 'auth_method_request' | 'key_id' | 'key_selection_request' | 'connected';
  data?: string;
  rows?: number;
  columns?: number;
  key_id?: number;
}

/**
 * 创建 WebSSH WebSocket 连接
 * @param hostId 主机ID
 * @param onMessage 消息处理回调
 * @param onError 错误处理回调
 * @param onClose 关闭处理回调
 * @returns WebSocket 连接
 */
export function createWebSSHConnection(
  hostId: number,
  onMessage: (data: string) => void,
  onError: (error: Event) => void,
  onClose: () => void
): WebSocket {
  // 获取 JWT token
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('未登录，请先登录');
  }

  // 构建 WebSocket URL
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host;
  
  // 修复 WebSocket URL 构建逻辑
  let wsUrl: string;
  if (window.location.port === '3000') {
    // 开发环境：前端在 3000 端口，后端在 8080 端口
    wsUrl = `${protocol}//${window.location.hostname}:8080/api/v1/webssh/terminal/${hostId}`;
  } else {
    // 生产环境：使用相同的 host
    wsUrl = `${protocol}//${host}/api/v1/webssh/terminal/${hostId}`;
  }
  
  // 添加认证 token 到查询参数
  const wsUrlWithAuth = `${wsUrl}?token=${encodeURIComponent(token)}`;
  
  console.log('WebSSH WebSocket URL:', wsUrlWithAuth);

  // 创建 WebSocket 连接
  const ws = new WebSocket(wsUrlWithAuth);

  ws.onopen = () => {
    console.log('WebSSH WebSocket connected');
  };

  ws.onmessage = (event) => {
    console.log('WebSSH received message:', event.data);
    
    // 检查是否是 JSON 消息（用于认证请求等）
    try {
      const msg: TerminalMessage = JSON.parse(event.data);
      
      // 处理特殊消息类型
      if (msg.type === 'username_request' || msg.type === 'password_request' || 
          msg.type === 'auth_method_request' || msg.type === 'key_selection_request') {
        // 这些消息会在 Terminal 组件中处理
        console.log('Received auth request:', msg.type);
        onMessage(event.data);
        return;
      }
      
      if (msg.type === 'error') {
        console.error('WebSSH server error:', msg.data);
        onError(new ErrorEvent('error', { message: msg.data || 'Unknown server error' }));
        return;
      }
      
      // 处理 connected 消息
      if (msg.type === 'connected') {
        console.log('SSH connection established');
        onMessage(event.data);
        return;
      }
      
      // 其他 JSON 消息（如 pong）可以忽略或处理
      if (msg.type !== 'pong') {
        onMessage(event.data);
      }
    } catch (e) {
      // 不是 JSON，是终端输出，直接传递
      onMessage(event.data);
    }
  };

  ws.onerror = (error) => {
    console.error('WebSSH WebSocket error:', error);
    onError(error);
  };

  ws.onclose = () => {
    console.log('WebSSH WebSocket closed');
    onClose();
  };

  return ws;
}

/**
 * 发送终端输入
 * @param ws WebSocket 连接
 * @param data 输入数据
 */
export function sendTerminalInput(ws: WebSocket, data: string): void {
  if (ws.readyState === WebSocket.OPEN) {
    const msg: TerminalMessage = {
      type: 'input',
      data: data,
    };
    ws.send(JSON.stringify(msg));
  }
}

/**
 * 发送用户名
 * @param ws WebSocket 连接
 * @param username 用户名
 */
export function sendUsername(ws: WebSocket, username: string): void {
  if (ws.readyState === WebSocket.OPEN) {
    const msg: TerminalMessage = {
      type: 'username',
      data: username,
    };
    ws.send(JSON.stringify(msg));
  }
}

/**
 * 发送密码
 * @param ws WebSocket 连接
 * @param password 密码
 */
export function sendPassword(ws: WebSocket, password: string): void {
  if (ws.readyState === WebSocket.OPEN) {
    const msg: TerminalMessage = {
      type: 'password',
      data: password,
    };
    ws.send(JSON.stringify(msg));
  }
}

/**
 * 发送终端大小变更
 * @param ws WebSocket 连接
 * @param rows 行数
 * @param columns 列数
 */
export function sendTerminalResize(ws: WebSocket, rows: number, columns: number): void {
  if (ws.readyState === WebSocket.OPEN) {
    const msg: TerminalMessage = {
      type: 'resize',
      rows: rows,
      columns: columns,
    };
    ws.send(JSON.stringify(msg));
  }
}

/**
 * 发送认证方式
 * @param ws WebSocket 连接
 * @param method 认证方式 ('key' 或 'password')
 */
export function sendAuthMethod(ws: WebSocket, method: 'key' | 'password'): void {
  if (ws.readyState === WebSocket.OPEN) {
    const msg: TerminalMessage = {
      type: 'auth_method',
      data: method,
    };
    ws.send(JSON.stringify(msg));
  }
}

/**
 * 发送SSH密钥ID
 * @param ws WebSocket 连接
 * @param keyId SSH密钥ID
 */
export function sendKeyId(ws: WebSocket, keyId: number): void {
  if (ws.readyState === WebSocket.OPEN) {
    const msg: TerminalMessage = {
      type: 'key_id',
      key_id: keyId,
    };
    ws.send(JSON.stringify(msg));
  }
}

/**
 * 发送 ping（保持连接）
 * @param ws WebSocket 连接
 */
export function sendPing(ws: WebSocket): void {
  if (ws.readyState === WebSocket.OPEN) {
    const msg: TerminalMessage = {
      type: 'ping',
    };
    ws.send(JSON.stringify(msg));
  }
}

