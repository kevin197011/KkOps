# Change: RDP Connection Support

## Why

当前系统仅支持 SSH 协议连接，但 Windows 服务器管理通常需要图形界面访问。虽然可以通过 SSH 进行命令行管理，但某些场景下需要图形界面操作：

- **GUI 应用管理**: 某些 Windows 应用只有图形界面
- **桌面操作**: 需要直接操作 Windows 桌面
- **用户培训**: 新用户更熟悉图形界面
- **故障排查**: 图形界面更直观地查看问题

需要添加 RDP（Remote Desktop Protocol）连接支持，提供 Windows 图形界面远程访问能力。

## What Changes

- **ADDED**: RDP 连接管理功能
  - RDP 连接配置（主机、端口、用户名、密码、域）
  - 连接参数配置（分辨率、颜色深度、安全模式）
  - 连接创建、编辑、删除
  - 连接列表和搜索

- **ADDED**: RDP 会话管理
  - 会话创建和建立
  - 会话状态跟踪
  - 会话历史记录
  - 会话断开和清理

- **ADDED**: RDP 客户端集成
  - 集成 Apache Guacamole 或类似 RDP Web 客户端
  - 浏览器中的图形界面显示
  - 鼠标和键盘输入支持
  - 剪贴板同步（可选）

- **ADDED**: 数据库表结构
  - `rdp_connections` 表（连接配置）
  - `rdp_sessions` 表（会话记录）

- **MODIFIED**: 主机管理
  - 添加 RDP 连接选项
  - 支持同时配置 SSH 和 RDP
  - 主机详情显示连接方式

## Impact

- **Affected specs**: `webssh-management` (扩展为远程访问管理)
- **Affected code**:
  - `backend/internal/models/` - 新增 RDP 模型
  - `backend/internal/repository/` - 新增 RDP 仓库
  - `backend/internal/service/` - 新增 RDP 服务
  - `backend/internal/handler/` - 新增 RDP 处理器
  - `frontend/src/pages/` - 新增或修改 RDP 管理页面
  - `frontend/src/components/` - 新增 RDP 客户端组件

- **Database changes**: 
  - 新增 `rdp_connections` 表
  - 新增 `rdp_sessions` 表

- **External dependencies**:
  - Apache Guacamole 服务器（推荐）
  - 或 noVNC + xrdp（替代方案）

- **Infrastructure changes**:
  - 需要部署 Guacamole 服务（Java 应用）
  - 需要配置数据库（Guacamole 使用）
  - 需要网络配置（RDP 端口 3389）

## Benefits

- ✅ **图形界面支持**: 支持 Windows 图形界面远程访问
- ✅ **统一管理平台**: 在一个平台管理 SSH 和 RDP 连接
- ✅ **提升用户体验**: 图形界面更直观易用
- ✅ **满足 Windows 管理需求**: 支持 Windows 服务器的完整管理

## Complexity Assessment

**复杂度**: 🔴 **高**

**工作量估算**: 12-18 天

**技术风险**: 🟡 中-高

**主要原因**:
- RDP 协议复杂度高（二进制协议，多种子协议）
- 需要集成第三方服务（Guacamole）
- 需要处理图形界面、鼠标键盘、剪贴板等
- 性能要求高（实时图形传输）

详见 `ANALYSIS.md` 详细分析文档。

## Alternatives Considered

1. **仅使用 SSH**: 功能受限，无法访问图形界面
2. **使用第三方 RDP 客户端**: 用户需要安装额外软件，体验不统一
3. **使用 VNC**: 需要 Windows 安装 VNC 服务器，配置复杂

## Recommendation

**建议**: 如果 Windows 管理需求不紧急，建议：
1. **先实施 SSH Windows 支持**（5-7 天，复杂度中等）
2. **根据实际需求再考虑 RDP 支持**（12-18 天，复杂度高）

如果确实需要图形界面，推荐使用 **Apache Guacamole** 作为 RDP 网关。

