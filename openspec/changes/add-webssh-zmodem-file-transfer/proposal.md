# Change: 为 WebSSH 终端添加 rz/sz (ZMODEM) 文件传输支持

## Why

当前 WebSSH 终端仅支持文本交互，无法进行文件传输。运维人员经常需要通过 SSH 传输文件，rz/sz 是最常用的命令行文件传输工具。添加此功能可以：
- 支持文件上传（rz）
- 支持文件下载（sz）
- 提升 WebSSH 的实用性
- 无需额外的文件传输工具

## What Changes

### 后端变更
- 在 `handleSSHConnect` 中检测 ZMODEM 协议序列（`**B` 上传，`**G` 下载）
- 实现二进制数据传输处理（区分文本和 ZMODEM 二进制数据）
- 添加文件传输状态管理（上传/下载进度跟踪）
- 支持 WebSocket 二进制消息类型（BinaryMessage）

### 前端变更
- 检测 ZMODEM 启动序列，自动弹出文件选择器（上传）
- 实现文件分块上传（二进制数据）
- 实现文件下载处理（接收二进制数据并触发浏览器下载）
- 添加传输进度显示 UI（模态框显示上传/下载进度）
- 添加传输取消功能
- 拦截终端中的 ZMODEM 序列显示（避免在终端中显示协议序列）

## Impact

- **受影响的文件**:
  - `backend/internal/handler/websocket/ssh.go`
  - `frontend/src/pages/ssh/WebSSHTerminal.tsx`
  - `frontend/src/utils/zmodem.ts` (新建工具文件)
- **协议支持**: ZMODEM 文件传输协议
- **WebSocket 消息类型**: 新增二进制消息处理

## 非功能性需求

- **性能**: 文件传输使用分块处理，避免内存溢出
- **用户体验**: 提供清晰的进度反馈和错误提示
- **兼容性**: 支持常见的 rz/sz 工具（lrzsz 包）
- **安全性**: 保持现有的 SSH 权限检查机制
