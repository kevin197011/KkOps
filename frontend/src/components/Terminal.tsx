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

interface TerminalProps {
  hostId: number;
  hostName?: string;
  onClose?: () => void;
}

const Terminal: React.FC<TerminalProps> = ({ hostId, hostName, onClose }) => {
  const terminalRef = useRef<HTMLDivElement>(null);
  const terminalInstanceRef = useRef<XTerm | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const [connecting, setConnecting] = useState(false);
  const [authModalVisible, setAuthModalVisible] = useState(false);
  const [authType, setAuthType] = useState<'username' | 'password' | 'auth_method' | 'key_selection'>('auth_method');
  const [authValue, setAuthValue] = useState('');
  const [authMethod, setAuthMethod] = useState<'key' | 'password'>('password');
  const [selectedKeyId, setSelectedKeyId] = useState<number | undefined>(undefined);
  const [sshKeys, setSshKeys] = useState<SSHKey[]>([]);
  const authInputRef = useRef<InputRef>(null);

  useEffect(() => {
    // 初始化终端
    const terminal = new XTerm({
      theme: {
        background: '#1e1e1e',
        foreground: '#d4d4d4',
        cursor: '#aeafad',
      },
      fontSize: 14,
      fontFamily: 'Consolas, "Courier New", monospace',
      cursorBlink: true,
      cursorStyle: 'block',
      rows: 24,
      cols: 80,
    });

    const fitAddon = new FitAddon();
    terminal.loadAddon(fitAddon);
    terminalInstanceRef.current = terminal;
    fitAddonRef.current = fitAddon;

    if (terminalRef.current) {
      terminal.open(terminalRef.current);
      fitAddon.fit();
    }

    // 处理窗口大小变化
    const handleResize = () => {
      if (fitAddonRef.current) {
        fitAddonRef.current.fit();
        // 发送终端大小到服务器
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN && connected) {
          const dims = terminalInstanceRef.current?.cols && terminalInstanceRef.current?.rows
            ? { rows: terminalInstanceRef.current.rows, cols: terminalInstanceRef.current.cols }
            : { rows: 24, cols: 80 };
          sendTerminalResize(wsRef.current, dims.rows, dims.cols);
        }
      }
    };

    window.addEventListener('resize', handleResize);

    // 清理函数
    return () => {
      window.removeEventListener('resize', handleResize);
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
      sendPassword(wsRef.current, authValue);
      setAuthModalVisible(false);
      setAuthValue('');
    }
  };

  const connect = () => {
    if (connecting || connected) return;

    setConnecting(true);
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
                setAuthType('auth_method');
                setAuthMethod('password'); // 默认选择密码
                setAuthModalVisible(true);
                loadSSHKeys(); // 预加载密钥列表
                return;
              }
              if (msg.type === 'key_selection_request') {
                setAuthType('key_selection');
                setAuthModalVisible(true);
                return;
              }
              if (msg.type === 'username_request') {
                setAuthType('username');
                setAuthValue('');
                setAuthModalVisible(true);
                setTimeout(() => {
                  authInputRef.current?.focus();
                }, 100);
                return;
              }
              if (msg.type === 'password_request') {
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
            } catch (e) {
              // 不是 JSON，是终端输出
            }

            // 显示终端输出（如果收到终端输出，说明连接已建立）
            if (terminal) {
              if (!connected) {
                setConnected(true);
                setConnecting(false);
                setAuthModalVisible(false);
              }
              terminal.write(data);
            }
          },
        (error) => {
          console.error('WebSocket error:', error);
          message.error('连接错误');
          setConnecting(false);
          setConnected(false);
        },
        () => {
          setConnected(false);
          setConnecting(false);
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
          : { rows: 24, cols: 80 };
        sendTerminalResize(ws, dims.rows, dims.cols);
      };

      // 处理终端输入
      terminal.onData((data) => {
        if (ws.readyState === WebSocket.OPEN && connected) {
          sendTerminalInput(ws, data);
        }
      });

      // 处理粘贴
      terminal.onBinary((data) => {
        if (ws.readyState === WebSocket.OPEN && connected) {
          ws.send(data);
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
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <div style={{ padding: '8px', background: '#2d2d2d', borderBottom: '1px solid #3e3e3e' }}>
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
          padding: '8px',
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
            style={{ width: '100%' }}
          >
            <Space direction="vertical">
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

