## ADDED Requirements

### Requirement: Tag Search
系统 SHALL 提供标签搜索功能，允许用户通过标签名称快速查找标签。

#### Scenario: Search tags by name
- **WHEN** 用户在搜索框中输入标签名称
- **THEN** 系统 SHALL 实时过滤并显示匹配的标签列表
- **AND** 搜索结果 SHALL 在用户停止输入 300ms 后更新（防抖）

#### Scenario: Clear search
- **WHEN** 用户清空搜索框
- **THEN** 系统 SHALL 显示所有标签

### Requirement: Tag Usage Details
系统 SHALL 允许用户查看使用某个标签的所有主机列表。

#### Scenario: View hosts using a tag
- **WHEN** 用户点击标签的"使用数量"
- **THEN** 系统 SHALL 显示一个弹窗，列出所有使用该标签的主机
- **AND** 弹窗 SHALL 显示主机名称和 IP 地址

#### Scenario: Navigate to host from tag usage
- **WHEN** 用户在标签使用详情弹窗中点击某个主机
- **THEN** 系统 SHALL 跳转到该主机的详情页面或打开主机详情弹窗

### Requirement: Batch Tag Deletion
系统 SHALL 支持批量删除多个标签。

#### Scenario: Select multiple tags for deletion
- **WHEN** 用户在标签列表中选择多个标签
- **THEN** 系统 SHALL 显示批量删除按钮
- **AND** 按钮 SHALL 显示已选择的标签数量

#### Scenario: Confirm batch deletion
- **WHEN** 用户点击批量删除按钮
- **THEN** 系统 SHALL 显示确认弹窗，提示将要删除的标签数量
- **AND** 用户确认后 SHALL 删除所有选中的标签

#### Scenario: Batch deletion with associated hosts
- **WHEN** 用户尝试批量删除包含已关联主机的标签
- **THEN** 系统 SHALL 在确认弹窗中警告用户这些标签正在被使用
- **AND** 删除后 SHALL 自动解除与主机的关联

### Requirement: Color Presets
系统 SHALL 提供常用颜色预设，方便用户快速选择标签颜色。

#### Scenario: Display color presets
- **WHEN** 用户创建或编辑标签时
- **THEN** 系统 SHALL 在颜色选择器旁显示常用颜色预设按钮

#### Scenario: Select preset color
- **WHEN** 用户点击某个预设颜色按钮
- **THEN** 系统 SHALL 自动将该颜色填充到颜色选择器中

### Requirement: Tag Name Uniqueness Validation
系统 SHALL 在创建或编辑标签时验证标签名称的唯一性。

#### Scenario: Create tag with duplicate name
- **WHEN** 用户尝试创建一个已存在名称的标签
- **THEN** 系统 SHALL 显示错误提示"标签名称已存在"
- **AND** 系统 SHALL 阻止表单提交

#### Scenario: Edit tag with duplicate name
- **WHEN** 用户尝试将标签名称修改为已存在的名称
- **THEN** 系统 SHALL 显示错误提示"标签名称已存在"
- **AND** 系统 SHALL 阻止表单提交

#### Scenario: Real-time name validation
- **WHEN** 用户在名称输入框中输入内容
- **THEN** 系统 SHALL 在用户停止输入 500ms 后检查名称是否已存在
- **AND** 如果名称已存在 SHALL 立即显示警告提示
