# Task Management Specification

## ADDED Requirements

### Requirement: Scheduled Task CRUD
系统 SHALL 支持定时任务的创建、读取、更新和删除操作。

#### Scenario: 创建定时任务
- **GIVEN** 用户已登录且有管理权限
- **WHEN** 用户提交定时任务创建请求，包含名称、Cron 表达式、执行内容
- **THEN** 系统创建定时任务并返回任务详情
- **AND** 如果任务启用，系统自动将其加入调度队列

#### Scenario: 更新定时任务
- **GIVEN** 存在一个定时任务
- **WHEN** 用户更新任务的 Cron 表达式或执行内容
- **THEN** 系统更新任务配置
- **AND** 如果任务已启用，系统重新调度该任务

#### Scenario: 删除定时任务
- **GIVEN** 存在一个定时任务
- **WHEN** 用户删除该任务
- **THEN** 系统从调度队列中移除该任务
- **AND** 系统软删除任务记录（保留执行历史）

### Requirement: Scheduled Task Scheduling
系统 SHALL 根据 Cron 表达式自动调度和执行定时任务。

#### Scenario: 按时执行任务
- **GIVEN** 存在一个已启用的定时任务，Cron 表达式为 `0 * * * *`（每小时）
- **WHEN** 到达调度时间点
- **THEN** 系统自动执行该任务
- **AND** 创建执行记录
- **AND** 更新任务的 `last_run_at` 和 `next_run_at`

#### Scenario: 服务重启后恢复调度
- **GIVEN** 服务重启
- **WHEN** 服务初始化完成
- **THEN** 系统加载所有已启用的定时任务到调度队列
- **AND** 根据 Cron 表达式计算下次执行时间

### Requirement: Scheduled Task Enable/Disable
系统 SHALL 支持启用和禁用定时任务。

#### Scenario: 禁用定时任务
- **GIVEN** 存在一个已启用的定时任务
- **WHEN** 用户禁用该任务
- **THEN** 系统将任务从调度队列中移除
- **AND** 任务状态变为禁用
- **AND** 任务不会被自动执行

#### Scenario: 启用定时任务
- **GIVEN** 存在一个已禁用的定时任务
- **WHEN** 用户启用该任务
- **THEN** 系统将任务加入调度队列
- **AND** 计算并设置下次执行时间

### Requirement: Manual Task Execution
系统 SHALL 支持手动立即执行定时任务。

#### Scenario: 立即执行任务
- **GIVEN** 存在一个定时任务（无论启用或禁用）
- **WHEN** 用户点击"立即执行"
- **THEN** 系统立即执行该任务
- **AND** 创建执行记录，标记触发类型为"手动"
- **AND** 不影响任务的正常调度

### Requirement: Task Execution History
系统 SHALL 记录并展示定时任务的执行历史。

#### Scenario: 查看执行历史
- **GIVEN** 定时任务有执行记录
- **WHEN** 用户查看任务详情
- **THEN** 系统展示该任务的所有执行记录
- **AND** 包含执行时间、状态、触发类型、输出日志

#### Scenario: 区分触发类型
- **GIVEN** 任务有多次执行记录
- **WHEN** 用户查看执行历史
- **THEN** 每条记录显示触发类型（手动/定时）

### Requirement: Cron Expression Validation
系统 SHALL 验证用户输入的 Cron 表达式格式。

#### Scenario: 有效的 Cron 表达式
- **GIVEN** 用户创建或更新定时任务
- **WHEN** 输入有效的 Cron 表达式（如 `0 0 * * *`）
- **THEN** 系统接受该表达式
- **AND** 显示下次执行时间预览

#### Scenario: 无效的 Cron 表达式
- **GIVEN** 用户创建或更新定时任务
- **WHEN** 输入无效的 Cron 表达式
- **THEN** 系统拒绝该请求
- **AND** 返回明确的错误信息

### Requirement: Template-based Task Creation
系统 SHALL 支持基于执行模板创建定时任务。

#### Scenario: 从模板创建任务
- **GIVEN** 存在执行模板
- **WHEN** 用户选择模板创建定时任务
- **THEN** 系统自动填充脚本内容和类型
- **AND** 用户可以修改其他配置（Cron 表达式、目标主机）

### Requirement: URL Route Reorganization
系统 SHALL 使用新的 URL 路由结构区分运维执行和任务管理。

#### Scenario: 访问运维执行页面
- **WHEN** 用户访问 `/executions`
- **THEN** 系统显示运维执行管理页面（原 `/tasks`）

#### Scenario: 访问任务管理页面
- **WHEN** 用户访问 `/tasks`
- **THEN** 系统显示定时任务管理页面

#### Scenario: 访问执行模板页面
- **WHEN** 用户访问 `/templates`
- **THEN** 系统显示执行模板管理页面（原 `/task-templates`）
