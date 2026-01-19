# Project Management - Remove Created At Column

## REMOVED Requirements

### Requirement: 项目列表创建时间列显示

项目列表页面 **SHALL NOT** 显示"创建时间"列。

#### Scenario: 查看项目列表时不显示创建时间
1. 用户访问 `/projects` 页面
2. 系统显示项目列表表格
3. **期望**: 表格中不包含"创建时间"列
4. **期望**: 其他列（ID、项目名称、描述、操作）正常显示
