# Implementation Tasks

## Backend Changes

- [ ] 在 `handleSSHConnect` 中添加 ZMODEM 状态管理结构体
  - `active` (bool): 是否处于传输状态
  - `direction` (string): "upload" 或 "download"
  - `mode` (string): "zmodem" 或 "trzsz"
  - `bytesTransferred`, `totalBytes` (int64): 进度跟踪
- [ ] 实现 ZMODEM 序列检测函数 `detectZmodemSequence`
  - 检测 `**B` (上传/rz)
  - 检测 `**G` (下载/sz)
  - 检测 `::TRZSZ:` (trzsz 协议)
- [ ] 修改 stdout/stderr 流处理
  - 检测到 ZMODEM 序列时，发送 `zmodem_start` 消息到前端
  - 拦截序列，避免在终端显示
  - 传输期间区分文本和二进制数据
- [ ] 修改 WebSocket 消息处理循环
  - 支持 `websocket.BinaryMessage` 类型（文件数据）
  - 处理 `zmodem_cancel` 消息
  - 处理 `zmodem_data_size` 消息（前端上传时告知文件大小）
  - 在传输模式下，二进制消息直接转发到 SSH stdin

## Frontend Changes

- [ ] 创建 `frontend/src/utils/zmodem.ts` 工具文件
  - `detectZmodemSequence`: 检测 ZMODEM 序列（备用，主要用于调试）
  - `formatBytes`: 格式化字节数为可读字符串
  - `calculateProgress`: 计算传输进度百分比
- [ ] 在 `SSHConnection` 接口中添加 ZMODEM 传输状态
  - `zmodemTransfer`: 传输信息（方向、文件名、进度等）
  - `downloadChunks`: 下载数据缓冲区
- [ ] 修改 WebSocket 消息处理
  - 处理 `zmodem_start` 消息：初始化传输状态
  - 处理 `zmodem_progress` 消息：更新传输进度
  - 处理 `zmodem_end` 消息：完成传输或取消
  - 处理 `zmodem_error` 消息：显示错误
  - 支持接收 `Blob` 类型的二进制消息（下载）
- [ ] 实现文件上传流程
  - 检测到 `zmodem_start` (direction: upload) 时弹出文件选择器
  - 读取文件并分块发送（使用 ArrayBuffer）
  - 发送文件大小信息 (`zmodem_data_size`)
  - 发送文件数据块（WebSocket BinaryMessage）
- [ ] 实现文件下载流程
  - 检测到 `zmodem_start` (direction: download) 时准备接收数据
  - 收集二进制数据块（Blob）
  - 收到 `zmodem_end` 时组装文件并触发浏览器下载
- [ ] 实现传输进度 UI
  - 创建 `ZmodemProgressModal` 组件
  - 显示文件名、进度条、已传输/总大小
  - 提供取消按钮
- [ ] 拦截终端中的 ZMODEM 序列显示
  - 在 `term.onData` 中检测并过滤 ZMODEM 序列

## Testing

- [ ] 测试文件上传（rz）
  - 在远程服务器运行 `rz` 命令
  - 验证文件选择器自动弹出
  - 验证上传进度显示
  - 验证文件成功上传
- [ ] 测试文件下载（sz）
  - 在远程服务器运行 `sz <filename>` 命令
  - 验证文件自动下载
  - 验证下载进度显示
- [ ] 测试传输取消
  - 验证取消按钮功能
  - 验证取消后状态清理
- [ ] 测试错误处理
  - 验证传输失败时的错误提示
  - 验证网络中断处理
