# Change: SSO 用户管理 — 对接常用运维系统统一授权登录

## Why

- 企业内通常已有统一身份（如 Keycloak、Azure AD、钉钉/飞书、蓝鲸、JumpServer 等），运维系统需要与现有 IdP 对接，避免多套账号密码。
- 支持 OIDC 即可覆盖大部分常用运维/企业 IdP，实现一次登录、多系统统一授权。

## What Changes

### 数据模型
- **ADDED**: User 表增加 `source`（local | sso）、`external_id`、`sso_provider`，用于区分本地用户与 SSO 用户并关联 IdP 身份。
- **MODIFIED**: 本地登录仅允许 `source=local` 且具备密码的用户；SSO 回调时按 `external_id` 查找或自动创建用户并签发本系统 JWT。

### 后端
- **ADDED**: 配置项 `sso`（enabled、oidc: issuer_url、client_id、client_secret、redirect_url、scopes）。
- **ADDED**: OIDC 流程：`GET /auth/sso/login` 跳转 IdP，`GET /auth/sso/callback` 处理 code 换 token、验证 id_token、创建/更新用户、签发 JWT 并重定向到前端（带 token 或 cookie）。
- **MODIFIED**: 登录逻辑：禁止对 `source=sso` 用户使用密码登录（可选策略：允许若已设置密码则仍可本地登录，按产品决定）；用户管理列表展示“来源”（本地/SSO）。

### 前端
- **ADDED**: 登录页在 SSO 开启时展示“SSO 登录”入口，点击跳转后端 `/auth/sso/login`。
- **ADDED**: 回调页（或根路径）解析 URL 中的 token（或从 cookie 读取）后写入 storage 并跳转控制台。
- **MODIFIED**: 用户管理列表/详情展示“来源”（本地/SSO），SSO 用户不展示“修改密码”或置灰。

## Impact

- **Affected code**: `backend/internal/config`, `backend/internal/model/user.go`, `backend/internal/service/auth`, `backend/internal/handler/auth`, `backend/cmd/server/main.go`, 前端登录页与路由、用户管理列表/详情。
- **Database**: User 表新增字段（可空、带默认值），兼容现有用户。

## Security

- state 参数防 CSRF；redirect_uri 严格校验；client_secret 仅服务端使用；id_token 必须校验签名与 issuer。
