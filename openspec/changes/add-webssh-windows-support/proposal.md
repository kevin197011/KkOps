# Change: WebSSH Windows Support

## Why

当前 WebSSH 实现主要针对 Linux/Unix 系统优化，但实际运维环境中经常需要管理 Windows 服务器。Windows 10/11 和 Windows Server 2019+ 已内置 OpenSSH Server，支持标准 SSH 协议，理论上应该可以连接，但可能存在以下问题：

- **终端兼容性**: Windows 终端特性与 Linux 不同，可能导致显示异常
- **字符编码**: Windows 可能使用不同的字符编码（如 GBK），导致中文显示乱码
- **用户体验**: 缺少 Windows 特定的优化和提示

需要明确支持 Windows 连接，确保在各种 Windows 环境下都能正常工作。

## What Changes

- **ADDED**: 操作系统类型检测
  - 连接后自动检测远程系统类型（Linux/Windows）
  - 在主机信息中标识操作系统类型
  - 根据系统类型调整连接参数

- **ADDED**: 终端类型适配
  - Linux/Unix: 使用 `xterm-256color`
  - Windows: 使用 `vt100` 或 `vt220`（或检测支持的终端类型）
  - 自动根据系统类型选择终端类型

- **ADDED**: 字符编码支持
  - 检测远程系统的字符编码
  - 支持编码配置（UTF-8, GBK, CP936 等）
  - 自动进行编码转换

- **MODIFIED**: WebSSH 连接逻辑
  - 添加操作系统检测步骤
  - 根据系统类型调整终端类型
  - 处理编码转换

- **ADDED**: Windows 特定优化（可选）
  - 系统类型标识显示
  - Windows 路径处理提示
  - PowerShell/CMD 识别

## Impact

- **Affected specs**: `webssh-management`
- **Affected code**:
  - `backend/internal/handler/webssh_handler.go` - 添加 OS 检测和终端类型适配
  - `backend/internal/models/host.go` - 添加操作系统类型字段（可选）
  - `frontend/src/components/Terminal.tsx` - 添加编码配置和显示优化
  - `frontend/src/services/webssh.ts` - 添加编码配置参数

- **Database changes**: 
  - 可选：在 `hosts` 表添加 `os_type` 字段（用于缓存系统类型）

- **Backward compatibility**: 
  - ✅ 完全向后兼容，不影响现有 Linux 连接
  - ✅ 现有功能不受影响

## Benefits

- ✅ **扩展支持范围**: 支持 Windows 服务器管理
- ✅ **统一管理界面**: 在一个平台管理 Linux 和 Windows
- ✅ **提升用户体验**: Windows 用户可以使用熟悉的 WebSSH 界面
- ✅ **降低复杂度**: 无需为 Windows 单独部署管理工具

## Complexity Assessment

**复杂度**: 🟡 中等

**工作量估算**: 5-7 天

**风险**: 🟢 低（标准协议，成熟技术）

详见 `ANALYSIS.md` 详细分析文档。

