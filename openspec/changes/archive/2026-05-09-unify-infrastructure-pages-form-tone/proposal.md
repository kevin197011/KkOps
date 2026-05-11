# Change: 基础设施子页面表单与卡片色调统一

## Why

「基础设施」下的子菜单页（项目管理、环境管理、云平台管理、资产管理、标签管理）当前表单与卡片色调不一致：部分页面（如标签管理）使用主题相关的背景、Card 包裹、图标+标题样式；其余页面（项目、环境、云平台、资产）多为裸 div + h2 + Table，无 Card、无统一背景与边框，导致同一分类下视觉不统一。需**检查并统一**所有子页面，尽量保持色调一致。

## What Changes

- **统一页面结构**：所有基础设施列表页采用同一套页面结构——外层容器（可选统一 padding/背景）、主内容区用 Card 包裹（或与布局背景一致），标题区统一为图标 + Typography.Title，表格与筛选项置于 Card 内。
- **统一色调来源**：优先使用 Ant Design 的 theme token（如 colorBgContainer、colorBorder、colorBgElevated）或现有主题变量，避免各页硬编码不同色值；深色/浅色由 ConfigProvider 驱动，保证同一主题下各子页面一致。
- **统一 Modal/表单**：弹窗与表单区域使用同一套背景与边框约定（如 Modal 的 styles.body 与主题一致），不出现「有的页面弹窗有深色背景、有的无」的情况。
- **适用范围**：基础设施下的 5 个子页面——`/projects`（ProjectList）、`/environments`（EnvironmentList）、`/cloud-platforms`（CloudPlatformList）、`/assets`（AssetList）、`/tags`（TagList）；若存在详情页（如 AssetDetail），一并纳入检查与统一。

## Impact

- **受影响规格**：可新增或修改 UI 一致性相关规格（如 infrastructure-pages-ui）
- **受影响代码**：`frontend/src/pages/projects/ProjectList.tsx`、`frontend/src/pages/environments/EnvironmentList.tsx`、`frontend/src/pages/cloudPlatforms/CloudPlatformList.tsx`、`frontend/src/pages/assets/AssetList.tsx`、`frontend/src/pages/tags/TagList.tsx`，以及 `frontend/src/pages/assets/AssetDetail.tsx`（若为弹窗或独立页）。可选：抽公共布局或样式常量（如共享的「基础设施列表页」包装组件或 token 常量）以减少重复。
