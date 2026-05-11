# SSO 配置说明（OIDC）

支持通过 OIDC 与常用运维/企业 IdP（如 Keycloak、Azure AD、GitLab、蓝鲸、JumpServer 等）对接，实现统一授权登录。

## 配置项

在 `config.yaml` 或环境变量中配置：

```yaml
sso:
  enabled: true
  oidc:
    issuer_url: "https://your-idp.example.com/realms/kkops"   # IdP 发现地址（不含 .well-known/openid-configuration）
    client_id: "kkops"
    client_secret: "your-client-secret"
    redirect_url: "https://kkops-api.example.com/api/v1/auth/sso/callback"  # 回调地址（必须在 IdP 中配置）
    frontend_base_url: "https://kkops.example.com"   # 可选；前后端分离时前端地址，登录成功后重定向到此
    scopes: ["openid", "profile", "email"]          # 可选，默认即此
```

环境变量示例（若使用 viper 自动绑定）：

- `SSO_ENABLED=true`
- `SSO_OIDC_ISSUER_URL=https://...`
- `SSO_OIDC_CLIENT_ID=kkops`
- `SSO_OIDC_CLIENT_SECRET=...`
- `SSO_OIDC_REDIRECT_URL=https://...`
- `SSO_OIDC_FRONTEND_BASE_URL=https://...`

## IdP 端配置

1. 在 IdP 中创建 OIDC 客户端（confidential），配置：
   - **Redirect URI**: 与上面 `redirect_url` 一致（如 `https://kkops-api.example.com/api/v1/auth/sso/callback`）
   - **Scopes**: 至少包含 `openid`，建议 `profile`、`email`
2. 获取 **Client ID**、**Client Secret**、**Issuer URL**（如 Keycloak 为 `https://host/realms/your-realm`）

## 流程说明

1. 用户点击登录页「SSO 登录」→ 跳转后端 `GET /api/v1/auth/sso/login` → 重定向到 IdP 授权页。
2. 用户在 IdP 登录并授权 → IdP 重定向到 `redirect_url?code=...&state=...`。
3. 后端 `GET /api/v1/auth/sso/callback` 用 code 换取 id_token，校验后按 `sub` 查找或创建用户（source=sso），签发本系统 JWT，重定向到前端 `{frontend_base_url}/auth/callback?token=...`（未配置 frontend_base_url 时使用相对路径 `/auth/callback?token=...`）。
4. 前端 `/auth/callback` 将 token 写入 storage，拉取用户信息与权限后跳转控制台。

## 用户来源

- **本地用户**（source=local）：使用用户名+密码登录，可修改密码。
- **SSO 用户**（source=sso）：仅能通过 SSO 登录，用户管理列表中显示「SSO」来源，编辑时无「重置密码」选项。首次 SSO 登录时按 IdP 的 `sub`、`preferred_username`/`name`、`email` 自动创建用户；若用户名已存在则返回错误，需在 IdP 中调整或由管理员预先创建对应用户。
