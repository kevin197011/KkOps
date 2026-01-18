## ADDED Requirements

### Requirement: 用户资产授权模型

系统 SHALL 提供用户与资产的授权关联机制，用于控制用户对资产的访问权限。

#### Scenario: 创建用户资产授权
- **WHEN** 管理员为用户授权资产访问权限
- **THEN** 系统在 `user_assets` 表创建关联记录
- **AND** 记录包含 user_id、asset_id、created_at、created_by

#### Scenario: 撤销用户资产授权
- **WHEN** 管理员撤销用户的资产访问权限
- **THEN** 系统从 `user_assets` 表删除对应关联记录
- **AND** 用户立即失去对该资产的访问权限

---

### Requirement: 管理员角色识别

系统 SHALL 通过角色的 `is_admin` 字段识别管理员角色，管理员自动拥有所有资产的访问权限。

#### Scenario: 管理员访问所有资产
- **WHEN** 具有管理员角色（is_admin=true）的用户请求资产列表
- **THEN** 系统返回所有资产
- **AND** 不进行 user_assets 表的权限检查

#### Scenario: 管理员连接任意资产 SSH
- **WHEN** 管理员尝试 WebSSH 连接任意资产
- **THEN** 系统允许连接
- **AND** 不进行 user_assets 表的权限检查

---

### Requirement: 普通用户资产访问控制

系统 SHALL 限制普通用户只能访问被授权的资产。

#### Scenario: 普通用户获取资产列表
- **WHEN** 普通用户（非管理员）请求资产列表
- **THEN** 系统只返回在 `user_assets` 表中有授权记录的资产
- **AND** 返回结果不包含未授权的资产

#### Scenario: 普通用户访问已授权资产详情
- **WHEN** 普通用户请求已授权资产的详情
- **THEN** 系统返回资产详情信息

#### Scenario: 普通用户访问未授权资产详情
- **WHEN** 普通用户请求未授权资产的详情
- **THEN** 系统返回 403 Forbidden 错误
- **AND** 响应包含 "insufficient permissions" 错误信息

#### Scenario: 普通用户 SSH 连接已授权资产
- **WHEN** 普通用户尝试 WebSSH 连接已授权的资产
- **THEN** 系统允许建立 SSH 连接

#### Scenario: 普通用户 SSH 连接未授权资产
- **WHEN** 普通用户尝试 WebSSH 连接未授权的资产
- **THEN** 系统拒绝连接
- **AND** 返回 "no permission to access this asset" 错误

---

### Requirement: 任务执行资产权限检查

系统 SHALL 在任务执行时检查用户对目标资产的访问权限。

#### Scenario: 普通用户在已授权资产执行任务
- **WHEN** 普通用户创建任务并选择已授权的目标资产
- **THEN** 系统允许任务创建和执行

#### Scenario: 普通用户在未授权资产执行任务
- **WHEN** 普通用户尝试在未授权的资产上执行任务
- **THEN** 系统拒绝任务创建
- **AND** 返回 "no permission to execute on selected assets" 错误

#### Scenario: 管理员在任意资产执行任务
- **WHEN** 管理员创建任务并选择任意目标资产
- **THEN** 系统允许任务创建和执行

---

### Requirement: 资产授权管理 API

系统 SHALL 提供 RESTful API 用于管理用户的资产授权。

#### Scenario: 获取用户已授权资产列表
- **WHEN** 管理员请求 `GET /api/v1/users/:id/assets`
- **THEN** 系统返回该用户已授权的资产列表
- **AND** 响应包含每个资产的基本信息（id、hostname、ip）

#### Scenario: 批量授权资产给用户
- **WHEN** 管理员请求 `POST /api/v1/users/:id/assets` 并提供 asset_ids 数组
- **THEN** 系统为该用户批量创建资产授权记录
- **AND** 已存在的授权不重复创建
- **AND** 返回授权成功的资产数量

#### Scenario: 撤销单个资产授权
- **WHEN** 管理员请求 `DELETE /api/v1/users/:id/assets/:asset_id`
- **THEN** 系统删除对应的授权记录
- **AND** 返回 204 No Content

#### Scenario: 批量撤销资产授权
- **WHEN** 管理员请求 `DELETE /api/v1/users/:id/assets` 并提供 asset_ids 数组
- **THEN** 系统批量删除对应的授权记录
- **AND** 返回撤销成功的数量

#### Scenario: 普通用户无权管理授权
- **WHEN** 普通用户尝试访问授权管理 API
- **THEN** 系统返回 403 Forbidden 错误

---

### Requirement: 资产已授权用户查询

系统 SHALL 提供查询资产已授权用户列表的功能。

#### Scenario: 查询资产已授权用户
- **WHEN** 管理员请求 `GET /api/v1/assets/:id/users`
- **THEN** 系统返回该资产已授权的用户列表
- **AND** 响应包含每个用户的基本信息（id、username、real_name）

#### Scenario: 批量授权用户访问资产
- **WHEN** 管理员请求 `POST /api/v1/assets/:id/users` 并提供 user_ids 数组
- **THEN** 系统为这些用户批量创建该资产的授权记录
