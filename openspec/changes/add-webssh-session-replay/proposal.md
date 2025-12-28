# Change: WebSSH Session Replay

## Why

运维人员需要能够回放之前的 WebSSH 终端会话，用于：
- **审计和合规**: 记录所有终端操作，满足安全审计要求
- **问题排查**: 回放历史会话，分析问题发生时的操作过程
- **培训和学习**: 新员工可以观看资深运维人员的操作过程
- **事故分析**: 在安全事件发生后，回放相关会话进行分析

当前系统缺少会话记录和回放功能，无法追溯历史终端操作。

## What Changes

- **ADDED**: 终端会话记录功能
  - 记录所有终端输入（用户输入的命令和按键）
  - 记录所有终端输出（命令执行结果）
  - 记录时间戳和终端尺寸变化
  - 存储会话元数据（用户、主机、开始/结束时间）

- **ADDED**: 会话回放功能
  - 提供回放界面，支持时间轴控制
  - 支持播放、暂停、快进、慢放、跳转
  - 显示回放进度和会话信息
  - 支持搜索会话记录

- **ADDED**: 会话管理功能
  - 会话列表查看（按用户、主机、时间过滤）
  - 会话详情查看
  - 会话删除和归档
  - 会话导出功能（可选）

- **MODIFIED**: WebSSH 终端处理
  - 在 WebSSH 会话中集成记录功能
  - 实时记录输入输出数据
  - 会话结束时保存记录

## Impact

- **Affected specs**: `webssh-management`
- **Affected code**: 
  - `backend/internal/handler/webssh_handler.go` - 添加会话记录逻辑
  - `backend/internal/models/` - 新增会话记录模型
  - `backend/internal/service/` - 新增会话记录服务
  - `backend/internal/repository/` - 新增会话记录仓库
  - `frontend/src/components/Terminal.tsx` - 集成记录功能
  - `frontend/src/pages/WebSSH.tsx` - 添加回放界面
  - `frontend/src/services/webssh.ts` - 添加回放相关 API

- **Database changes**: 
  - 新增 `webssh_sessions` 表（如果不存在）
  - 新增 `webssh_session_records` 表（存储会话记录数据）

- **Performance considerations**:
  - 会话记录会增加存储空间需求
  - 需要定期清理旧记录
  - 回放功能需要高效的数据检索

## Benefits

- ✅ **安全审计**: 完整的操作记录，满足合规要求
- ✅ **问题排查**: 快速定位问题发生时的操作
- ✅ **知识传承**: 通过回放学习最佳实践
- ✅ **责任追溯**: 明确操作责任，提高安全性

