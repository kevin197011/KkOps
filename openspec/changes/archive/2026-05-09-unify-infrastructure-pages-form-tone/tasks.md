# Tasks: 基础设施子页面表单与卡片色调统一

## 1. 检查与约定

- [x] 1.1 约定基础设施列表页统一结构：页面容器（padding + 背景 token）、Card 包裹主内容、标题区（图标 + Typography.Title）、Modal 使用 token 背景
- [x] 1.2 确认 theme token 用法：使用 `theme.useToken()` 的 colorBgContainer / colorBgElevated / colorBorder 等，不在页面内硬编码 hex

## 2. 项目管理、环境管理、云平台管理

- [x] 2.1 ProjectList：增加 Card 包裹 Table；标题区改为图标（FolderOutlined）+ Typography.Title；页面容器与 Card 使用 theme token；Modal body 使用 token 背景
- [x] 2.2 EnvironmentList：同上，图标 GlobalOutlined，标题「环境管理」
- [x] 2.3 CloudPlatformList：同上，图标 CloudOutlined，标题「云平台管理」

## 3. 资产管理、标签管理

- [x] 3.1 AssetList：增加 Card 包裹筛选区+Table；标题区改为图标（DatabaseOutlined）+ Typography.Title；页面容器与 Card、Modal 使用 theme token
- [x] 3.2 TagList：保留 Card 与图标+标题；将硬编码背景/边框色（#0F172A、#1E293B、#F5F5F5 等）改为 theme token，与其余四页统一

## 4. 详情页与验收

- [x] 4.1 AssetDetail：检查背景与边框，若为独立页或弹窗，使用与列表页一致的 token
- [x] 4.2 在浅色与深色主题下分别打开五类子页面，确认表单与卡片色调一致
