import React, { useEffect, useRef, useState } from 'react';
import { Terminal as XTerm } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import 'xterm/css/xterm.css';
import { Button, Space, message, Input, Modal, InputRef, Radio, Select } from 'antd';
import { DisconnectOutlined, ReloadOutlined } from '@ant-design/icons';
import {
  createWebSSHConnection,
  sendTerminalInput,
  sendUsername,
  sendPassword,
  sendTerminalResize,
  sendAuthMethod,
  sendKeyId,
  TerminalMessage,
} from '../services/webssh';
import { sshKeyService, SSHKey } from '../services/sshKey';

// 认证信息类型
export interface AuthInfo {
  authMethod: 'key' | 'password';
  username: string;
  password?: string;
  keyId?: number;
}

interface TerminalProps {
  hostId: number;
  hostName?: string;
  authInfo?: AuthInfo; // 克隆时传入的认证信息
  onClose?: () => void;
  onStatusChange?: (status: 'connecting' | 'connected' | 'disconnected') => void;
  onAuthInfoUpdate?: (authInfo: AuthInfo) => void; // 认证成功后回调
}

const Terminal: React.FC<TerminalProps> = ({ hostId, hostName, authInfo, onClose, onStatusChange, onAuthInfoUpdate }) => {
  const terminalRef = useRef<HTMLDivElement>(null);
  const terminalInstanceRef = useRef<XTerm | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const connectedRef = useRef(false); // 用于在闭包中跟踪连接状态
  const [connected, setConnected] = useState(false);
  const [connecting, setConnecting] = useState(false);
  const [authModalVisible, setAuthModalVisible] = useState(false);
  const [authType, setAuthType] = useState<'username' | 'password' | 'auth_method' | 'key_selection'>('auth_method');
  const [authValue, setAuthValue] = useState('');
  const [authMethod, setAuthMethod] = useState<'key' | 'password'>(authInfo?.authMethod || 'password');
  const [selectedKeyId, setSelectedKeyId] = useState<number | undefined>(authInfo?.keyId);
  const [sshKeys, setSshKeys] = useState<SSHKey[]>([]);
  const authInputRef = useRef<InputRef>(null);
  
  // 存储当前认证信息，用于回调
  const currentAuthRef = useRef<{ username: string; password?: string; keyId?: number }>({
    username: authInfo?.username || '',
    password: authInfo?.password,
    keyId: authInfo?.keyId,
  });

  // 如果有预设的认证信息，初始化相关状态
  useEffect(() => {
    if (authInfo) {
      setAuthMethod(authInfo.authMethod);
      if (authInfo.keyId) {
        setSelectedKeyId(authInfo.keyId);
      }
      currentAuthRef.current = {
        username: authInfo.username,
        password: authInfo.password,
        keyId: authInfo.keyId,
      };
    }
  }, [authInfo]);

  useEffect(() => {
    // 初始化终端
    const terminal = new XTerm({
      theme: {
        background: '#1e1e1e',
        foreground: '#d4d4d4',
        cursor: '#aeafad',
        selectionBackground: '#264f78', // 选择高亮颜色
      },
      fontSize: 14,
      fontFamily: 'Consolas, "Courier New", monospace',
      cursorBlink: true,
      cursorStyle: 'block',
      rows: 40,
      cols: 160,
    });

    const fitAddon = new FitAddon();
    terminal.loadAddon(fitAddon);
    terminalInstanceRef.current = terminal;
    fitAddonRef.current = fitAddon;

    if (terminalRef.current) {
      terminal.open(terminalRef.current);
      // 延迟执行 fit，确保 DOM 完全渲染
      setTimeout(() => {
        if (fitAddonRef.current) {
          fitAddonRef.current.fit();
        }
      }, 100);
    }

    // 选择即复制功能：鼠标选择文本后自动复制到剪贴板
    let selectionTimeout: NodeJS.Timeout | null = null;
    const handleSelectionChange = () => {
      // 使用防抖，避免选择过程中频繁触发
      if (selectionTimeout) {
        clearTimeout(selectionTimeout);
      }
      selectionTimeout = setTimeout(() => {
        const selection = terminal.getSelection();
        if (selection && selection.length > 0) {
          navigator.clipboard.writeText(selection).then(() => {
            // 复制成功，显示简短提示
            message.success({ content: '已复制到剪贴板', duration: 1 });
          }).catch((err) => {
            console.error('复制失败:', err);
          });
        }
      }, 200); // 200ms 防抖，等待用户完成选择
    };
    terminal.onSelectionChange(handleSelectionChange);

    // 防抖定时器
    let resizeTimeout: NodeJS.Timeout | null = null;

    // 处理终端大小变化 - 自动适配容器
    const handleResize = () => {
      // 使用防抖避免频繁触发
      if (resizeTimeout) {
        clearTimeout(resizeTimeout);
      }
      resizeTimeout = setTimeout(() => {
        if (fitAddonRef.current && terminalInstanceRef.current) {
          try {
            fitAddonRef.current.fit();
            // 只有在 SSH 连接建立后才发送终端大小到服务器
            // 使用 connectedRef 避免闭包问题
            if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN && connectedRef.current) {
              const dims = terminalInstanceRef.current.cols && terminalInstanceRef.current.rows
                ? { rows: terminalInstanceRef.current.rows, cols: terminalInstanceRef.current.cols }
                : { rows: 40, cols: 160 };
              sendTerminalResize(wsRef.current, dims.rows, dims.cols);
            }
          } catch (e) {
            // 忽略 fit 错误
          }
        }
      }, 50);
    };

    // 监听窗口大小变化
    window.addEventListener('resize', handleResize);

    // 使用 ResizeObserver 监听容器大小变化（更准确）
    let resizeObserver: ResizeObserver | null = null;
    const currentTerminalRef = terminalRef.current;
    if (currentTerminalRef) {
      resizeObserver = new ResizeObserver(() => {
        handleResize();
      });
      resizeObserver.observe(currentTerminalRef);
    }

    // 清理函数
    return () => {
      window.removeEventListener('resize', handleResize);
      if (resizeObserver && currentTerminalRef) {
        resizeObserver.unobserve(currentTerminalRef);
        resizeObserver.disconnect();
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
      terminal.dispose();
    };
  }, []);

  const loadSSHKeys = async () => {
    try {
      const response = await sshKeyService.list(1, 100);
      setSshKeys(response.ssh_keys);
    } catch (error: any) {
      console.error('Failed to load SSH keys:', error);
    }
  };

  const handleAuthSubmit = () => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      return;
    }

    if (authType === 'auth_method') {
      // 提交认证方式选择
      sendAuthMethod(wsRef.current, authMethod);
      if (authMethod === 'key') {
        // 如果选择密钥，需要先加载密钥列表
        loadSSHKeys();
        setAuthType('key_selection');
        setAuthModalVisible(true);
      } else {
        // 如果选择密码，请求用户名
        setAuthType('username');
        setAuthValue('');
        setAuthModalVisible(true);
        setTimeout(() => {
          authInputRef.current?.focus();
        }, 100);
      }
      return;
    }

    if (authType === 'key_selection') {
      // 提交密钥选择
      if (!selectedKeyId) {
        message.warning('请选择SSH密钥');
        return;
      }
      // 保存密钥ID
      currentAuthRef.current.keyId = selectedKeyId;
      sendKeyId(wsRef.current, selectedKeyId);
      setAuthType('username');
      setAuthValue('');
      setAuthModalVisible(true);
      setTimeout(() => {
        authInputRef.current?.focus();
      }, 100);
      return;
    }

    if (authType === 'username') {
      if (!authValue.trim()) {
        message.warning('请输入用户名');
        return;
      }
      // 保存用户名
      currentAuthRef.current.username = authValue.trim();
      sendUsername(wsRef.current, authValue.trim());
      if (authMethod === 'password') {
        setAuthType('password');
        setAuthValue('');
        setAuthModalVisible(true);
        setTimeout(() => {
          authInputRef.current?.focus();
        }, 100);
      } else {
        // 密钥认证不需要密码
        setAuthModalVisible(false);
        setAuthValue('');
      }
      return;
    }

    if (authType === 'password') {
      if (!authValue.trim()) {
        message.warning('请输入密码');
        return;
      }
      // 保存密码
      currentAuthRef.current.password = authValue;
      sendPassword(wsRef.current, authValue);
      setAuthModalVisible(false);
      setAuthValue('');
    }
  };

  const connect = () => {
    if (connecting || connected) return;

    setConnecting(true);
    if (onStatusChange) onStatusChange('connecting');
    const terminal = terminalInstanceRef.current;
    if (!terminal) return;

    try {
      // 创建 WebSocket 连接
      const ws = createWebSSHConnection(
        hostId,
          (data) => {
            // 检查是否是认证请求消息
            try {
              const msg: TerminalMessage = JSON.parse(data);
              if (msg.type === 'auth_method_request') {
                // 如果有预设的认证信息（克隆），自动发送
                if (authInfo) {
                  setAuthMethod(authInfo.authMethod);
                  sendAuthMethod(ws, authInfo.authMethod);
                  return;
                }
                setAuthType('auth_method');
                setAuthMethod('password'); // 默认选择密码
                setAuthModalVisible(true);
                loadSSHKeys(); // 预加载密钥列表
                return;
              }
              if (msg.type === 'key_selection_request') {
                // 如果有预设的密钥ID（克隆），自动发送
                if (authInfo?.keyId) {
                  setSelectedKeyId(authInfo.keyId);
                  sendKeyId(ws, authInfo.keyId);
                  return;
                }
                setAuthType('key_selection');
                setAuthModalVisible(true);
                return;
              }
              if (msg.type === 'username_request') {
                // 如果有预设的用户名（克隆），自动发送
                if (authInfo?.username) {
                  currentAuthRef.current.username = authInfo.username;
                  sendUsername(ws, authInfo.username);
                  return;
                }
                setAuthType('username');
                setAuthValue('');
                setAuthModalVisible(true);
                setTimeout(() => {
                  authInputRef.current?.focus();
                }, 100);
                return;
              }
              if (msg.type === 'password_request') {
                // 如果有预设的密码（克隆），自动发送
                if (authInfo?.password) {
                  currentAuthRef.current.password = authInfo.password;
                  sendPassword(ws, authInfo.password);
                  return;
                }
                setAuthType('password');
                setAuthValue('');
                setAuthModalVisible(true);
                setTimeout(() => {
                  authInputRef.current?.focus();
                }, 100);
                return;
              }
              if (msg.type === 'error') {
                message.error(msg.data || '连接错误');
                setConnecting(false);
                setConnected(false);
                setAuthModalVisible(false);
                return;
              }
              // 处理连接成功消息
              if (msg.type === 'connected') {
                setConnected(true);
                connectedRef.current = true;
                setConnecting(false);
                setAuthModalVisible(false);
                if (onStatusChange) onStatusChange('connected');
                // 保存认证信息供克隆使用
                if (onAuthInfoUpdate && currentAuthRef.current.username) {
                  onAuthInfoUpdate({
                    authMethod: authMethod,
                    username: currentAuthRef.current.username,
                    password: currentAuthRef.current.password,
                    keyId: currentAuthRef.current.keyId,
                  });
                }
                // 清除之前的"正在连接"消息，显示连接成功
                terminal.clear();
                terminal.writeln('\x1b[32m✓ SSH 连接成功\x1b[0m\r\n');
                // 连接成功后，多次 fit 确保终端填满容器
                const doFit = () => {
                  if (fitAddonRef.current && terminalInstanceRef.current && wsRef.current) {
                    fitAddonRef.current.fit();
                    const rows = terminalInstanceRef.current.rows;
                    const cols = terminalInstanceRef.current.cols;
                    sendTerminalResize(wsRef.current, rows, cols);
                  }
                };
                // 延迟多次执行确保布局稳定
                setTimeout(doFit, 50);
                setTimeout(doFit, 200);
                setTimeout(doFit, 500);
                // 不显示这个 JSON 消息到终端
                return;
              }
            } catch (e) {
              // 不是 JSON，是终端输出
            }

            // 显示终端输出（如果收到终端输出，说明连接已建立）
            if (terminal) {
              if (!connectedRef.current) {
                // 首次收到终端输出，清除"正在连接"消息
                terminal.clear();
                terminal.writeln('\x1b[32m✓ SSH 连接成功\x1b[0m\r\n');
                setConnected(true);
                connectedRef.current = true;
                setConnecting(false);
                setAuthModalVisible(false);
                if (onStatusChange) onStatusChange('connected');
              }
              terminal.write(data);
            }
          },
        (error) => {
          console.error('WebSocket error:', error);
          message.error('连接错误');
          setConnecting(false);
          setConnected(false);
          connectedRef.current = false; // 同步更新 ref
          if (onStatusChange) onStatusChange('disconnected');
        },
        () => {
          setConnected(false);
          connectedRef.current = false; // 同步更新 ref
          setConnecting(false);
          if (onStatusChange) onStatusChange('disconnected');
          if (terminal) {
            terminal.writeln('\r\n\x1b[31m✗ 连接已断开\x1b[0m\r\n');
          }
        }
      );

      wsRef.current = ws;

      ws.onopen = () => {
        setConnecting(false);
        // 等待服务器请求用户名和密码
        terminal.writeln('\r\n\x1b[33m正在连接，请等待认证...\x1b[0m\r\n');
        
        // 发送初始终端大小
        const dims = terminal.cols && terminal.rows
          ? { rows: terminal.rows, cols: terminal.cols }
          : { rows: 40, cols: 160 };
        sendTerminalResize(ws, dims.rows, dims.cols);
      };

      // 处理终端输入 - 使用 wsRef 避免闭包问题
      terminal.onData((data) => {
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
          sendTerminalInput(wsRef.current, data);
        }
      });

      // 处理粘贴 - 使用 wsRef 避免闭包问题
      terminal.onBinary((data) => {
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
          sendTerminalInput(wsRef.current, data);
        }
      });

    } catch (error: any) {
      console.error('Failed to connect:', error);
      message.error(error.message || '连接失败');
      setConnecting(false);
    }
  };

  const disconnect = () => {
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setConnected(false);
    connectedRef.current = false; // 同步更新 ref
    if (onStatusChange) onStatusChange('disconnected');
  };

  const clear = () => {
    if (terminalInstanceRef.current) {
      terminalInstanceRef.current.clear();
    }
  };

  useEffect(() => {
    // 自动连接
    connect();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [hostId]);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%', flex: 1 }}>
      <div style={{ padding: '8px', background: '#2d2d2d', borderBottom: '1px solid #3e3e3e', flexShrink: 0 }}>
        <Space>
          <span style={{ color: '#d4d4d4' }}>
            {hostName || `主机 #${hostId}`}
          </span>
          {connected ? (
            <>
              <span style={{ color: '#4ec9b0' }}>●</span>
              <span style={{ color: '#d4d4d4', fontSize: '12px' }}>已连接</span>
              <Button
                size="small"
                icon={<DisconnectOutlined />}
                onClick={disconnect}
                danger
              >
                断开
              </Button>
            </>
          ) : connecting || authModalVisible ? (
            <>
              <span style={{ color: '#ffa500' }}>●</span>
              <span style={{ color: '#d4d4d4', fontSize: '12px' }}>连接中...</span>
            </>
          ) : (
            <>
              <span style={{ color: '#f48771' }}>●</span>
              <span style={{ color: '#d4d4d4', fontSize: '12px' }}>未连接</span>
              <Button size="small" onClick={connect}>
                连接
              </Button>
            </>
          )}
          <Button size="small" icon={<ReloadOutlined />} onClick={clear}>
            清空
          </Button>
          {onClose && (
            <Button size="small" onClick={onClose}>
              关闭
            </Button>
          )}
        </Space>
      </div>
      <div
        ref={terminalRef}
        style={{
          flex: 1,
          background: '#1e1e1e',
          overflow: 'hidden',
        }}
      />
      {/* 认证模态框 */}
      <Modal
        title={
          authType === 'auth_method' ? '选择认证方式' :
          authType === 'key_selection' ? '选择SSH密钥' :
          authType === 'username' ? '输入SSH用户名' :
          '输入SSH密码'
        }
        open={authModalVisible}
        onOk={handleAuthSubmit}
        onCancel={() => {
          setAuthModalVisible(false);
          setAuthValue('');
          setSelectedKeyId(undefined);
          disconnect();
        }}
        okText="确定"
        cancelText="取消"
        destroyOnClose
      >
        {authType === 'auth_method' && (
          <Radio.Group
            value={authMethod}
            onChange={(e) => setAuthMethod(e.target.value)}
            style={{ width: '100%', marginTop: 16 }}
          >
            <Space direction="vertical" size="middle">
              <Radio value="password">密码认证</Radio>
              <Radio value="key">SSH密钥认证</Radio>
            </Space>
          </Radio.Group>
        )}
        {authType === 'key_selection' && (
          <Select
            style={{ width: '100%' }}
            placeholder="请选择SSH密钥"
            value={selectedKeyId}
            onChange={(value) => setSelectedKeyId(value)}
            options={sshKeys.map((key) => ({
              label: `${key.name} (${key.key_type.toUpperCase()})`,
              value: key.id,
            }))}
          />
        )}
        {(authType === 'username' || authType === 'password') && (
          <Input
            ref={authInputRef}
            type={authType === 'password' ? 'password' : 'text'}
            placeholder={authType === 'username' ? '请输入SSH用户名' : '请输入SSH密码'}
            value={authValue}
            onChange={(e) => setAuthValue(e.target.value)}
            onPressEnter={handleAuthSubmit}
            autoFocus
          />
        )}
      </Modal>
    </div>
  );
};

export default Terminal;

