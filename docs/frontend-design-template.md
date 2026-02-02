# 智能运维管理平台 - 前端设计模板

## 设计理念

采用**Swiss Modernism 2.0 + Minimalism**设计风格，打造专业、极客、高效的运维管理后台。强调数据可视化、信息密度和操作效率。

## 1. 设计系统

### 1.1 色彩系统

基于搜索推荐的SaaS Dashboard配色方案，适配Ant Design 5+主题系统：

#### 亮色主题（Light Mode）
```typescript
const lightTheme: ThemeConfig = {
  token: {
    // 主色调 - 信任蓝
    colorPrimary: '#2563EB',      // 主要操作、链接
    colorSuccess: '#10B981',      // 成功状态
    colorWarning: '#F59E0B',      // 警告状态
    colorError: '#EF4444',        // 错误状态
    colorInfo: '#3B82F6',         // 信息提示
    
    // 背景色
    colorBgBase: '#FFFFFF',       // 主背景
    colorBgContainer: '#F8FAFC',  // 容器背景（浅灰）
    colorBgElevated: '#FFFFFF',   // 浮层背景
    
    // 文字颜色
    colorText: '#1E293B',         // 主文字（slate-800）
    colorTextSecondary: '#64748B', // 次要文字（slate-500）
    colorTextTertiary: '#94A3B8', // 辅助文字（slate-400）
    
    // 边框颜色
    colorBorder: '#E2E8F0',       // 边框（slate-200）
    colorBorderSecondary: '#F1F5F9', // 次要边框
    
    // CTA按钮 - 橙色强调
    colorLink: '#2563EB',
    colorLinkHover: '#1D4ED8',
  },
  components: {
    Button: {
      primaryColor: '#FFFFFF',
      controlHeight: 36,
      borderRadius: 6,
    },
    Card: {
      borderRadius: 8,
      paddingLG: 24,
    },
    Table: {
      borderRadius: 8,
    },
  },
};
```

#### 暗色主题（Dark Mode）
```typescript
const darkTheme: ThemeConfig = {
  algorithm: theme.darkAlgorithm,
  token: {
    // 主色调 - 稍亮的蓝色
    colorPrimary: '#3B82F6',      // 亮色模式稍亮
    colorSuccess: '#10B981',
    colorWarning: '#F59E0B',
    colorError: '#EF4444',
    colorInfo: '#60A5FA',
    
    // 背景色 - OLED友好的深色
    colorBgBase: '#0F172A',       // 主背景（slate-900）
    colorBgContainer: '#1E293B',  // 容器背景（slate-800）
    colorBgElevated: '#334155',   // 浮层背景（slate-700）
    
    // 文字颜色
    colorText: '#F1F5F9',         // 主文字（slate-100）
    colorTextSecondary: '#CBD5E1', // 次要文字（slate-300）
    colorTextTertiary: '#94A3B8', // 辅助文字（slate-400）
    
    // 边框颜色
    colorBorder: '#334155',       // 边框（slate-700）
    colorBorderSecondary: '#475569', // 次要边框（slate-600）
  },
};
```

### 1.2 字体系统

采用**Modern Professional**字体配对：

- **标题字体**：Poppins（几何无衬线，现代专业）
- **正文字体**：Open Sans（人文无衬线，易读性强）

```typescript
// tailwind.config.js 或全局CSS
const fontFamily = {
  heading: ['Poppins', 'sans-serif'],
  body: ['Open Sans', 'sans-serif'],
};

// Google Fonts导入
// @import url('https://fonts.googleapis.com/css2?family=Open+Sans:wght@300;400;500;600;700&family=Poppins:wght@400;500;600;700&display=swap');

// Ant Design字体配置
const themeConfig = {
  token: {
    fontFamily: 'Open Sans, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    fontSizeHeading1: 32,
    fontSizeHeading2: 24,
    fontSizeHeading3: 20,
    fontSizeHeading4: 16,
    fontSize: 14,
    lineHeight: 1.6,
  },
};
```

### 1.3 间距系统

基于8px基础网格系统：

```typescript
const spacing = {
  xs: '4px',   // 0.5
  sm: '8px',   // 1
  md: '16px',  // 2
  lg: '24px',  // 3
  xl: '32px',  // 4
  xxl: '48px', // 6
};
```

### 1.4 圆角系统

```typescript
const borderRadius = {
  none: '0px',
  sm: '4px',   // 小元素（标签、徽章）
  base: '6px', // 按钮、输入框
  md: '8px',   // 卡片、表格
  lg: '12px',  // 大卡片、模态框
  full: '9999px', // 圆形
};
```

### 1.5 阴影系统

```typescript
const shadows = {
  sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
  base: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)',
  md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
  lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
  xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
};
```

## 2. 布局结构

### 2.1 整体布局

采用经典的**侧边栏导航 + 顶部栏 + 内容区**布局：

```
┌─────────────────────────────────────────┐
│  TopBar (固定)                           │
├──────────┬──────────────────────────────┤
│          │  Breadcrumb (可选)            │
│ Sidebar  ├──────────────────────────────┤
│ (可折叠) │  Content Area                │
│          │  (可滚动)                     │
│          │                               │
│          │                               │
└──────────┴──────────────────────────────┘
```

### 2.2 侧边栏（Sidebar）

- **宽度**：展开 240px，折叠 64px
- **背景**：亮色模式 `#FFFFFF`，暗色模式 `#1E293B`
- **高度**：100vh（固定）
- **导航项高度**：48px
- **激活状态**：左侧边框高亮（4px，主色） + 背景色变化
- **图标**：24x24px，与文字间距 12px

### 2.3 顶部栏（TopBar）

- **高度**：64px
- **背景**：亮色模式 `#FFFFFF`，暗色模式 `#1E293B`
- **内容**：
  - 左侧：Logo + 项目名称
  - 中间：全局搜索（可选）
  - 右侧：通知、主题切换、用户菜单
- **阴影**：底部阴影（sm级别）
- **固定定位**：`position: fixed, top: 0`

### 2.4 内容区

- **内边距**：24px（lg间距）
- **最大宽度**：无限制（全宽布局）
- **背景**：`colorBgContainer`
- **顶部间距**：考虑TopBar高度，`padding-top: 64px`

## 3. 组件设计规范

### 3.1 卡片（Card）

```typescript
// 标准卡片
<Card 
  className="hover:shadow-md transition-shadow duration-200 cursor-pointer"
  style={{ borderRadius: 8 }}
>
  <Card.Meta
    title="卡片标题"
    description="卡片描述"
  />
</Card>

// 数据统计卡片
<Card>
  <Statistic
    title="总资产数"
    value={1234}
    prefix={<Icon />}
    suffix="台"
    valueStyle={{ color: '#2563EB' }}
  />
</Card>
```

**规范**：
- 圆角：8px
- 内边距：24px
- 悬停效果：阴影增强（sm → md）
- 可点击卡片：添加 `cursor-pointer`

### 3.2 表格（Table）

```typescript
<Table
  columns={columns}
  dataSource={data}
  pagination={{
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total) => `共 ${total} 条`,
    pageSizeOptions: ['10', '20', '50', '100'],
  }}
  rowSelection={rowSelection}
  scroll={{ x: 'max-content' }}
/>
```

**规范**：
- 斑马纹：启用（提高可读性）
- 行高：48px（中等密度）
- 悬停：行背景色变化
- 固定列：重要列（操作列、ID列）固定
- 响应式：小屏幕启用横向滚动

### 3.3 表单（Form）

```typescript
<Form
  layout="vertical"
  labelCol={{ span: 24 }}
  wrapperCol={{ span: 24 }}
>
  <Form.Item
    label="字段名称"
    name="field"
    rules={[{ required: true, message: '请输入' }]}
  >
    <Input placeholder="请输入" />
  </Form.Item>
</Form>
```

**规范**：
- 布局：垂直布局（`layout="vertical"`）
- 标签宽度：24（占满一行）
- 标签样式：`font-weight: 500`，`color: colorText`
- 输入框高度：36px（`controlHeight: 36`）

### 3.4 按钮（Button）

```typescript
// 主按钮
<Button type="primary" size="middle">
  确认
</Button>

// 次要按钮
<Button>取消</Button>

// 危险操作
<Button type="primary" danger>删除</Button>

// 图标按钮
<Button icon={<Icon />} />
```

**规范**：
- 高度：36px（middle）
- 圆角：6px
- 间距：按钮之间 8px
- 主要操作：`type="primary"`
- 危险操作：`danger`属性

### 3.5 标签（Tag）

```typescript
// 状态标签
<Tag color="success">运行中</Tag>
<Tag color="warning">警告</Tag>
<Tag color="error">错误</Tag>

// 自定义颜色标签
<Tag color="#2563EB">生产环境</Tag>
```

**规范**：
- 圆角：4px（sm）
- 内边距：4px 8px
- 字体大小：12px
- 颜色语义化：success/warning/error/info

### 3.6 徽章（Badge）

```typescript
<Badge count={5} offset={[-8, 8]}>
  <Avatar />
</Badge>
```

**规范**：
- 数字徽章：红色背景（`#EF4444`）
- 圆点徽章：小圆点（8px）

## 4. 数据可视化

### 4.1 图表库

推荐使用 **Recharts** 或 **Ant Design Charts**（基于 G2Plot）：

- **趋势图**：Line Chart / Area Chart
- **对比图**：Bar Chart / Column Chart
- **分布图**：Pie Chart（少量分类）
- **实时监控**：Time Series Line Chart

### 4.2 图表配色

```typescript
const chartColors = {
  primary: '#3B82F6',      // 主数据系列
  secondary: '#10B981',    // 次要数据系列
  accent: '#F59E0B',       // 强调数据
  danger: '#EF4444',       // 异常数据
  
  // 多系列配色
  series: [
    '#3B82F6', '#10B981', '#F59E0B', 
    '#EF4444', '#8B5CF6', '#EC4899'
  ],
};
```

### 4.3 图表规范

- **填充透明度**：Area Chart使用20%透明度
- **交互**：启用hover提示 + 缩放（可选）
- **响应式**：图表容器使用百分比宽度
- **无障碍**：为色盲用户添加图案叠加（可选）

## 5. 交互设计

### 5.1 动画

```typescript
// 标准过渡
transition: all 0.2s ease-in-out;

// 颜色过渡
transition-colors duration-200

// 悬停效果
.hover-effect:hover {
  transform: translateY(-2px);
  box-shadow: shadow-md;
}
```

**原则**：
- 持续时间：150-300ms（快速响应）
- 缓动函数：`ease-in-out`
- 避免：过度动画、过慢动画（>500ms）

### 5.2 反馈

- **点击反馈**：按钮按下状态（active）
- **加载状态**：Loading Spinner + 文字提示
- **成功反馈**：Message.success（2秒自动消失）
- **错误反馈**：Message.error（3秒自动消失）
- **确认操作**：Modal.confirm（危险操作）

### 5.3 悬停状态

所有可交互元素必须有悬停反馈：

- **按钮**：背景色变深 / 阴影增强
- **卡片**：阴影增强（sm → md）
- **链接**：下划线 + 颜色变化
- **表格行**：背景色变化

```typescript
// 示例
className="hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200 cursor-pointer"
```

## 6. 响应式设计

### 6.1 断点

```typescript
const breakpoints = {
  xs: '< 576px',   // 手机
  sm: '≥ 576px',   // 平板（竖屏）
  md: '≥ 768px',   // 平板（横屏）
  lg: '≥ 992px',   // 桌面
  xl: '≥ 1200px',  // 大桌面
  xxl: '≥ 1600px', // 超大桌面
};
```

### 6.2 布局适配

- **手机（< 768px）**：
  - 侧边栏：抽屉式（Drawer）
  - 表格：横向滚动
  - 卡片：单列布局
  - TopBar：简化（隐藏搜索）

- **平板（768px - 992px）**：
  - 侧边栏：可折叠
  - 表格：启用横向滚动
  - 卡片：2列布局

- **桌面（≥ 992px）**：
  - 侧边栏：正常展开
  - 表格：完整显示
  - 卡片：3-4列布局

## 7. 无障碍设计（A11y）

### 7.1 键盘导航

- 所有可交互元素：`tabIndex={0}`
- 焦点可见：清晰的高亮边框
- 操作快捷键：支持常见快捷键（Ctrl+S保存等）

### 7.2 颜色对比度

- **正文文字**：≥ 4.5:1（WCAG AA）
- **大文字（18px+）**：≥ 3:1
- **交互元素**：≥ 3:1

### 7.3 屏幕阅读器

- 图标按钮：添加 `aria-label`
- 表单字段：关联 `<label>`
- 状态提示：使用 `aria-live` 区域

## 8. 页面模板示例

### 8.1 列表页模板

```
┌─────────────────────────────────────────┐
│  Breadcrumb                             │
├─────────────────────────────────────────┤
│  Page Title                    [新增]    │
├─────────────────────────────────────────┤
│  [搜索框] [筛选器] [排序]                │
├─────────────────────────────────────────┤
│  Table (带分页)                         │
└─────────────────────────────────────────┘
```

### 8.2 详情页模板

```
┌─────────────────────────────────────────┐
│  Breadcrumb                             │
├─────────────────────────────────────────┤
│  Page Title              [编辑] [删除]   │
├─────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────┐              │
│  │ Card 1  │  │ Card 2  │              │
│  └─────────┘  └─────────┘              │
├─────────────────────────────────────────┤
│  ┌───────────────────────────────────┐ │
│  │  Main Content Card                │ │
│  │                                    │ │
│  └───────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

### 8.3 仪表盘模板

```
┌─────────────────────────────────────────┐
│  Breadcrumb                             │
├─────────────────────────────────────────┤
│  Page Title              [刷新] [导出]   │
├─────────────────────────────────────────┤
│  ┌───┐ ┌───┐ ┌───┐ ┌───┐              │
│  │Stat│ │Stat│ │Stat│ │Stat│              │
│  └───┘ └───┘ └───┘ └───┘              │
├─────────────────────────────────────────┤
│  ┌───────────────────────────────────┐ │
│  │  Chart Area                       │ │
│  │                                    │ │
│  └───────────────────────────────────┘ │
├─────────────────────────────────────────┤
│  ┌──────────────┐ ┌──────────────┐    │
│  │  Table 1     │ │  Table 2     │    │
│  └──────────────┘ └──────────────┘    │
└─────────────────────────────────────────┘
```

## 9. 实现建议

### 9.1 Ant Design主题配置

```typescript
// src/theme/index.ts
import { ThemeConfig } from 'antd';
import { theme } from 'antd';

const { darkAlgorithm, defaultAlgorithm } = theme;

export const lightTheme: ThemeConfig = {
  algorithm: defaultAlgorithm,
  token: {
    // ... 颜色配置
  },
  components: {
    // ... 组件配置
  },
};

export const darkTheme: ThemeConfig = {
  algorithm: darkAlgorithm,
  token: {
    // ... 暗色配置
  },
};
```

### 9.2 主题切换

```typescript
// 使用Context或状态管理
const [isDark, setIsDark] = useState(false);

<ConfigProvider theme={isDark ? darkTheme : lightTheme}>
  <App />
</ConfigProvider>
```

### 9.3 全局样式

```css
/* 字体导入 */
@import url('https://fonts.googleapis.com/css2?family=Open+Sans:wght@300;400;500;600;700&family=Poppins:wght@400;500;600;700&display=swap');

/* 全局样式重置 */
* {
  box-sizing: border-box;
}

body {
  font-family: 'Open Sans', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

h1, h2, h3, h4, h5, h6 {
  font-family: 'Poppins', sans-serif;
  font-weight: 600;
}
```

## 10. 设计检查清单

### 视觉质量
- [ ] 使用SVG图标（Heroicons/Lucide），不使用Emoji
- [ ] 所有图标尺寸一致（24x24px）
- [ ] 悬停状态不引起布局偏移
- [ ] 品牌Logo使用正确的SVG

### 交互
- [ ] 所有可点击元素有 `cursor-pointer`
- [ ] 悬停状态提供清晰的视觉反馈
- [ ] 过渡动画流畅（150-300ms）
- [ ] 键盘导航焦点可见

### 主题
- [ ] 亮色模式文字对比度足够（≥4.5:1）
- [ ] 暗色模式文字对比度足够
- [ ] 边框在两种模式下都可见
- [ ] 两种模式都经过测试

### 布局
- [ ] 响应式设计（320px, 768px, 1024px, 1440px）
- [ ] 移动端无横向滚动
- [ ] 固定导航不遮挡内容
- [ ] 内容区域有适当的内边距

### 无障碍
- [ ] 图片有alt文本
- [ ] 表单字段有label
- [ ] 颜色不是唯一的指示器
- [ ] 支持 `prefers-reduced-motion`

---

**文档版本**：v1.0  
**创建日期**：2025-01-27  
**基于**：ui-ux-pro-max设计系统 + Ant Design 5+
