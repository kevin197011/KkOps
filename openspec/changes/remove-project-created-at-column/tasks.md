# Implementation Tasks

## Frontend Changes

- [x] 从 `ProjectList.tsx` 的 `columns` 数组中移除"创建时间"列定义
- [x] 检查是否还需要 moment 库（如果其他列或功能不再使用，可以考虑移除导入）
  - 已移除 moment 导入，因为不再需要

## Validation

- [ ] 在浏览器中访问 `/projects` 页面
- [ ] 确认"创建时间"列已不再显示
- [ ] 确认其他列正常显示
- [ ] 确认表格布局正常
