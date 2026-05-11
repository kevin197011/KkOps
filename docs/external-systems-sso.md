# SSO 应用门户与外部系统跳转

**SSO（单点登录）**：用户登录一次，即可访问多个系统。KkOps 支持两种用法：

1. **KkOps 登录用 SSO**：在登录页使用「SSO 登录」通过 Keycloak 等 IdP 登录（见 [sso-config.md](./sso-config.md)）。
2. **SSO 应用门户**：用户登录 KkOps 后，在「外部系统跳转」一键打开 GitLab、Jenkins、Grafana 等；若这些系统与 KkOps **共用同一 IdP**，打开即已登录（同 SSO 应用）；也可配置为「Token 跳转」向目标系统附带 KkOps 身份与权限。

## 两种打开方式

| 方式 | 说明 | 适用场景 |
|------|------|----------|
| **同 SSO 应用** | 只配置应用 URL，点击打开直接访问；目标系统与 KkOps 共用同一 IdP（如 Keycloak），浏览器已有 IdP 会话，目标系统用同一 IdP 校验即已登录 | GitLab、Jenkins、Grafana、K8s Dashboard 等已接入同一 Keycloak 的场景 |
| **Token 跳转** | KkOps 签发短期 JWT，带用户与权限，跳转到目标系统的指定路径；目标系统需提供接收并校验 token 的接口 | 自研或支持「JWT 断言」登录的运维系统 |

## 同 SSO 应用流程（推荐）

```
用户登录 KkOps（通过 Keycloak OIDC）
  → 在 KkOps 点击「打开 GitLab」
  → 新窗口打开 https://gitlab.example.com
  → GitLab 使用同一 Keycloak；浏览器已有会话，Keycloak 直接放行
  → 用户无需再次输入账号密码
```

## Token 跳转流程

1. 用户在 KkOps 进入「外部系统跳转」，点击某系统的「打开」。
2. 后端根据当前用户生成短期 JWT（5 分钟有效），携带用户 ID、用户名、邮箱、姓名、角色、权限及映射后的角色，使用与目标系统约定的**共享密钥**签名。
3. 前端在新窗口打开：`{目标系统 base_url}{login_path}?token={JWT}`。
4. 目标系统在 `login_path` 接口中解析 `token` 参数，用同一共享密钥校验 JWT，根据 payload 创建或更新本地用户会话，并按 `mapped_roles` 或 `permissions` 做权限控制。

## JWT Payload（HS256）

目标系统解码并验证签名后，可读取以下字段（均为标准或自定义 claim）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `sub` | number | KkOps 用户 ID |
| `username` | string | 用户名 |
| `email` | string | 邮箱 |
| `real_name` | string | 姓名 |
| `roles` | string[] | KkOps 角色名列表 |
| `permissions` | string[] | KkOps 权限列表，格式为 `resource:action`，如 `assets:*`、`users:read` |
| `mapped_roles` | string[] | 经「角色映射」后的目标系统角色名，可用于目标系统 RBAC |
| `exp` | number | 过期时间（Unix 秒） |
| `iat` | number | 签发时间（Unix 秒） |

## 目标系统接入要点

1. **共享密钥**：在 KkOps 添加外部系统时填写「共享密钥」，与目标系统约定一致；目标系统用该密钥做 HS256 验签。
2. **登录路径**：目标系统需提供 GET（或 POST）接口，例如 `/api/sso/consume` 或 `/sso/login`，从 query 或 body 读取 `token`，校验后创建会话并重定向到业务页。
3. **校验步骤**：
   - 解析 `token` 为 JWT，使用共享密钥验证签名（算法 HS256）。
   - 检查 `exp` 未过期。
   - 使用 `sub`/`username` 等创建或关联本地用户，并应用 `mapped_roles` 或 `permissions` 做权限控制。
4. **角色映射**：在 KkOps 配置的「角色映射」为 JSON，如 `{"admin":"administrator","ops":"operator"}`。后端会将当前用户的 KkOps 角色按此映射为 `mapped_roles` 数组传给目标系统，目标系统可直接按 `mapped_roles` 赋权。

## KkOps 侧配置

- **系统管理** → **外部系统跳转**（或运维导航旁菜单）→ **添加外部系统**。
- 填写：名称、基础 URL、登录路径、共享密钥、可选角色映射（JSON）、排序、启用。
- 用户点击「打开」即会带 JWT 跳转到目标系统。

## 安全说明

- 共享密钥仅保存在 KkOps 与目标系统，不可暴露给前端或第三方。
- JWT 有效期较短（默认 5 分钟），仅用于一次跳转登录。
- 目标系统必须校验签名与 `exp`，避免重放与伪造。
