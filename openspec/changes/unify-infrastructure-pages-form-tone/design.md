# Design: 基础设施子页面表单与卡片色调统一

## Context

用户反馈：基础设施下子菜单页面的表单色调没有保持一致，希望把所有子页面检查一遍并尽量保持色调一致。当前差异：TagList 使用自定义背景色（#0F172A/#F5F5F5）、Card 包裹、图标+标题；ProjectList、EnvironmentList、CloudPlatformList、AssetList 多为裸 div + h2 + Table，无 Card、无统一背景。

## Goals / Non-Goals

### Goals
- 五类基础设施列表页（项目、环境、云平台、资产、标签）视觉上统一：同一主题下背景、卡片、标题、表格、弹窗表单的色调与风格一致。
- 使用主题 token 驱动颜色，避免各页硬编码不同 hex 值；深色/浅色由全局主题决定。

### Non-Goals
- 不改变各页面的业务逻辑、表格列、筛选项功能。
- 不强制改动其他分类（任务管理、安全管理、系统管理）的页面风格；本变更仅针对「基础设施」下子页面。

## Decisions

### 1. 参考基准

- **以「标签管理」页为视觉参考**：当前 TagList 已有 Card、主题相关背景、图标+标题、Modal body 背景，统一时其余四页向该风格靠拢；同时将 TagList 中硬编码色值（如 #0F172A、#1E293B、#F5F5F5）改为 theme token 或 CSS 变量（若有），以便全局主题变更时一起生效。

### 2. 统一结构

- **页面容器**：根节点使用同一约定——例如 `padding: 24`、`background` 使用 `token.colorBgContainer` 或继承 Layout Content 背景，不强制 minHeight 除非布局需要。
- **主内容区**：列表+筛选+操作按钮置于一张 **Card** 内；Card 的 `styles` 或 `style` 使用 token（如 `colorBgContainer`/`colorBgElevated`、`colorBorderSecondary`），保证圆角、边框、阴影与现有设计系统一致。
- **标题区**：各页统一为「左侧图标 + Typography.Title level 3」，图标可用各模块已有图标（FolderOutlined、GlobalOutlined、CloudOutlined、DatabaseOutlined、TagsOutlined），标题文案为各页名称（项目管理、环境管理、云平台管理、资产管理、标签管理）；右侧为「新增」等主操作按钮。

### 3. 色调与 Token

- **背景**：页面背景与 Card 背景均通过 `theme.useToken()` 的 `colorBgContainer` 或 `colorBgLayout`、`colorBgElevated` 等获取，不写死 hex。
- **边框与阴影**：Card 的 border、boxShadow 使用 token 或 Ant Design 默认，保证深色/浅色一致。
- **Modal**：表单弹窗的 `styles.body` 使用与页面一致的 token（如 `background: token.colorBgElevated`），确保深色模式下弹窗内容区与主内容区色调一致。

### 4. 实施顺序

- 先统一 ProjectList、EnvironmentList、CloudPlatformList、AssetList 四页：增加 Card 包裹、标题区改为图标+Title、页面/Modal 使用 token。
- 再调整 TagList：保留现有结构与 Card，仅将硬编码色值替换为 theme token，使与其他四页同一 token 体系。
- 若有 AssetDetail（详情页或弹窗）：检查背景与边框是否与列表页、Modal 一致，必要时改为 token。

### 5. 可选：公共组件或常量

- 若五页重复较多，可抽「基础设施列表页」公共布局组件（接收 title、icon、children、extra），或集中定义一份「基础设施页 Card/标题」的 style 常量，减少重复并保证后续新增子页面时风格一致；否则可在各页内联相同结构，优先保证统一再考虑抽取。

## Risks / Trade-offs

| 风险 | 缓解 |
|------|------|
| 与 MainLayout Content 背景重复 | 使用与 Content 相同的 token，或仅 Card 内使用 colorBgElevated，外层不重复设背景 |
| 深色模式下部分页面未接入主题 | 所有涉及背景/边框的地方均使用 token，避免手写 #xxx |

## Migration Plan

1. 在 ProjectList、EnvironmentList、CloudPlatformList、AssetList 中增加 Card 包裹、统一标题区（图标+Title）、页面容器与 Modal 使用 theme token。
2. 在 TagList 中将硬编码背景/边框色改为 theme token，保持现有布局与图标+标题。
3. 检查 AssetDetail 等详情页，统一色调。
4. 手动在浅色/深色主题下各打开五页，确认色调一致。
