# Implementation Tasks

## Frontend Changes

- [x] 在 `ProjectList.tsx` 中导入 `moment` 库
- [x] 修改创建时间列的 `render` 函数，使用 `moment` 格式化时间
  - 格式：`YYYY-MM-DD HH:mm:ss`
  - 处理空值情况，显示 `-`
- [x] 测试验证
  - [x] 验证创建时间正确显示
  - [x] 验证时间格式符合预期
  - [x] 验证空值/无效时间处理正确

## Validation

- [x] 在浏览器中访问 `/projects` 页面
- [x] 确认"创建时间"列正确显示格式化后的时间
- [x] 确认时间格式与其他页面（如任务管理）保持一致
