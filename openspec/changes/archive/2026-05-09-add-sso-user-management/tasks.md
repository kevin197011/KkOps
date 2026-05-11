# Tasks: SSO 用户管理

## 1. 配置与依赖

- [x] 1.1 增加 Config.SSO（enabled、OIDC: issuer_url, client_id, client_secret, redirect_url, scopes）
- [x] 1.2 引入 go-oidc/v3、golang.org/x/oauth2，实现 OIDC 登录与 callback

## 2. 数据模型

- [x] 2.1 User 增加 Source、ExternalID、SSOProvider 字段并迁移
- [x] 2.2 本地登录仅对 source=local 校验密码；SSO 用户按 external_id 查找/创建

## 3. 后端 Auth

- [x] 3.1 GET /auth/sso/login：生成 state，重定向到 IdP 授权 URL
- [x] 3.2 GET /auth/sso/callback：换 code、验证 id_token、创建/更新用户、签发 JWT，重定向前端并带 token
- [x] 3.3 GET /auth/sso/config：返回 sso_enabled（及前端所需配置，如登录按钮文案）

## 4. 前端

- [x] 4.1 登录页请求 sso/config，若 enabled 则显示「SSO 登录」并跳转 /api/v1/auth/sso/login
- [x] 4.2 处理 callback 重定向：从 URL 或 cookie 取 token，写入 storage，跳转控制台
- [x] 4.3 用户管理列表/详情展示来源（本地/SSO），SSO 用户不提供修改密码入口

## 5. 文档与配置示例

- [x] 5.1 docs 中增加 SSO 配置说明与示例（Keycloak/通用 OIDC）
