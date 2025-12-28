# Design: WebSSH Windows Support

## Context

当前 WebSSH 实现使用标准 SSH 协议，理论上已支持 Windows（因为 Windows OpenSSH Server 遵循标准协议），但可能存在终端兼容性和字符编码问题。需要明确支持 Windows 连接，确保在各种 Windows 环境下都能正常工作。

## Goals / Non-Goals

### Goals
- 支持 Windows 10/11 和 Windows Server 2019+ 连接
- 自动检测操作系统类型
- 适配 Windows 终端特性
- 处理字符编码差异
- 保持与 Linux 连接的兼容性

### Non-Goals
- PowerShell 特定功能（语法高亮、命令补全等）
- Windows 服务管理集成
- RDP 协议支持（仅支持 SSH）

## Decisions

### Decision: 自动检测操作系统类型
**What**: 连接后自动检测远程系统类型
**Why**: 
- 无需用户手动配置
- 自动适配终端类型和编码
- 提升用户体验

**Alternatives considered**:
- 用户手动选择：增加用户负担
- 从主机信息获取：需要预先配置

### Decision: 动态终端类型选择
**What**: 根据操作系统类型动态选择终端类型
**Why**:
- Windows 可能不支持 `xterm-256color`
- 使用兼容的终端类型确保正常显示

**Alternatives considered**:
- 固定使用 `vt100`：兼容性好但功能受限
- 让用户选择：增加复杂度

### Decision: 编码检测和转换
**What**: 自动检测远程系统编码，进行转换
**Why**:
- Windows 可能使用 GBK 等编码
- 确保中文字符正确显示

**Alternatives considered**:
- 强制使用 UTF-8：可能不兼容旧系统
- 用户手动配置：增加用户负担

## Architecture

### OS Detection Flow

```
WebSSH Connection Established
    ↓
Execute Detection Command
    ↓
Parse Output
    ↓
Determine OS Type (Linux/Windows/Unknown)
    ↓
Select Terminal Type
    ↓
Configure Encoding
    ↓
Continue Normal Terminal Session
```

### Detection Methods

#### Method 1: Command-based Detection
```bash
# Linux/Unix
uname -s

# Windows
ver
# 或
systeminfo | findstr /B /C:"OS Name"
```

#### Method 2: Environment Variable
```bash
# Windows PowerShell
$env:OS

# Linux
echo $OSTYPE
```

### Terminal Type Selection

```go
func getTerminalType(osType string) string {
    switch osType {
    case "windows":
        // 尝试 xterm，如果不支持则降级到 vt100
        return "xterm" // 或 "vt100"
    case "linux", "darwin", "unix":
        return "xterm-256color"
    default:
        return "xterm"
    }
}
```

### Encoding Detection

#### Detection Strategy
1. **尝试检测**: 执行命令检测系统编码
2. **默认值**: Linux 使用 UTF-8，Windows 使用系统默认编码
3. **用户配置**: 允许用户覆盖编码设置

#### Windows Encoding Detection
```powershell
# PowerShell
[Console]::OutputEncoding.EncodingName
# 或
chcp
```

### Implementation Details

#### Backend Changes

**1. OS Detection Function**
```go
// detectOSType 检测远程系统类型
func (h *WebSSHHandler) detectOSType(session *ssh.Session) (string, error) {
    // 尝试执行 uname（Linux）
    output, err := session.Output("uname -s 2>/dev/null")
    if err == nil && len(output) > 0 {
        osStr := strings.ToLower(strings.TrimSpace(string(output)))
        if strings.Contains(osStr, "linux") {
            return "linux", nil
        }
        if strings.Contains(osStr, "darwin") {
            return "darwin", nil
        }
    }
    
    // 尝试执行 ver（Windows）
    output, err = session.Output("ver 2>&1")
    if err == nil {
        osStr := strings.ToLower(string(output))
        if strings.Contains(osStr, "windows") || 
           strings.Contains(osStr, "microsoft") {
            return "windows", nil
        }
    }
    
    return "unknown", nil
}
```

**2. Terminal Type Selection**
```go
// getTerminalTypeForOS 根据系统类型获取终端类型
func getTerminalTypeForOS(osType string) string {
    switch osType {
    case "windows":
        return "xterm" // Windows 10/11 通常支持 xterm
    case "linux", "darwin":
        return "xterm-256color"
    default:
        return "xterm"
    }
}
```

**3. Encoding Configuration**
```go
// getEncodingForOS 根据系统类型获取编码
func getEncodingForOS(osType string) string {
    switch osType {
    case "windows":
        return "gbk" // 或从系统检测
    case "linux", "darwin":
        return "utf-8"
    default:
        return "utf-8"
    }
}
```

#### Frontend Changes

**1. Encoding Configuration UI**
- 在连接设置中添加编码选择
- 显示检测到的系统类型和编码
- 允许用户手动调整

**2. Terminal Display**
- 处理不同编码的字符显示
- 确保中文字符正确显示

## Testing Strategy

### Test Scenarios

1. **Windows Server 2019 连接**
   - 测试操作系统检测
   - 测试终端类型适配
   - 测试字符编码

2. **Windows 10/11 连接**
   - 测试不同 Windows 版本
   - 测试 PowerShell 和 CMD

3. **混合环境**
   - 同时连接 Linux 和 Windows
   - 确保互不干扰

4. **编码测试**
   - 测试中文字符显示
   - 测试不同编码配置

## Migration Strategy

### Phase 1: Detection and Basic Support (2-3 days)
- 实现操作系统检测
- 实现终端类型适配
- 基础功能测试

### Phase 2: Encoding Support (2-3 days)
- 实现编码检测
- 实现编码转换
- 编码测试

### Phase 3: Optimization (1-2 days)
- 用户体验优化
- 性能优化
- 文档更新

## Performance Considerations

- **Detection Overhead**: OS 检测增加一次命令执行（约 100-500ms）
- **Encoding Conversion**: 编码转换有轻微性能开销
- **Caching**: 可以缓存 OS 类型到主机信息，避免重复检测

## Security Considerations

- OS 检测命令应该是安全的（只读命令）
- 编码转换不应该引入安全漏洞
- 保持现有的认证和授权机制

