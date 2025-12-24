# Kronos Frontend

基于 React + TypeScript + Ant Design 的运维中台管理系统前端。

## 技术栈

- React 18
- TypeScript
- Ant Design
- React Router
- Axios

## 快速开始

### 安装依赖

```bash
npm install
```

### 配置环境变量

创建 `.env` 文件：

```env
REACT_APP_API_URL=http://localhost:8080/api/v1
```

### 启动开发服务器

```bash
npm start
```

应用将在 `http://localhost:3000` 启动。

### 构建生产版本

```bash
npm run build
```

## 项目结构

```
frontend/
├── src/
│   ├── components/          # 公共组件
│   │   └── PrivateRoute.tsx
│   ├── contexts/             # Context
│   │   └── AuthContext.tsx
│   ├── pages/                # 页面组件
│   │   ├── Login.tsx
│   │   └── Dashboard.tsx
│   ├── services/             # API服务
│   │   ├── auth.ts
│   │   ├── user.ts
│   │   ├── host.ts
│   │   └── ssh.ts
│   ├── config/               # 配置文件
│   │   └── api.ts
│   └── App.tsx               # 应用入口
└── public/
```

## 功能模块

### 已实现
- ✅ 用户认证（登录）
- ✅ 路由保护
- ✅ API服务封装
- ✅ 基础Dashboard页面

### 待实现
- [ ] 用户管理页面
- [ ] 主机管理页面
- [ ] SSH管理页面
- [ ] 角色和权限管理页面
- [ ] 布局组件（侧边栏、顶部导航）

## API集成

所有API调用通过 `src/services/` 目录下的服务文件进行，使用统一的 `api` 实例（已配置Token自动添加和错误处理）。
