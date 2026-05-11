# KkOps 作为 OIDC IdP（身份提供商）

当 KkOps 启用 IdP 功能时，可被第三方系统（如 GitLab、Jenkins、Grafana、Jumpserver 等）配置为 OIDC 身份提供商，实现「一次在 KkOps 登录，多系统统一认证」。

## 前置条件

- 后端配置中 `idp.enabled` 为 `true`（默认开启）。
- 在 KkOps 后台「IdP 应用」中创建应用，获得 `client_id` 与 `client_secret`（创建或重置密钥时仅显示一次，请妥善保存）。应用可登记协议类型（**OIDC** 为当前主要能力；**SAML**、**LDAP** 为协议标识与未来端点扩展预留）。

## IdP 端点（Issuer）

以 Issuer 为 `http://localhost:3000/oidc` 为例（生产环境替换为实际公网地址）：

| 端点 | 方法 | 说明 |
|------|------|------|
| `{issuer}/.well-known/openid-configuration` | GET | OIDC 发现文档 |
| `{issuer}/authorize` | GET | 授权端点（浏览器重定向） |
| `{issuer}/token` | POST | 令牌端点（application/x-www-form-urlencoded） |
| `{issuer}/userinfo` | GET | 用户信息（Bearer access_token） |
| `{issuer}/jwks` | GET | 公钥集（JWKS） |

配置项（如 `idp.issuer`）需与上述 base URL 一致，且为第三方系统可访问的地址。

## 第三方系统侧配置要点

1. **发现文档**  
   填写 Issuer 或 `{issuer}/.well-known/openid-configuration`，由系统自动拉取 `authorization_endpoint`、`token_endpoint`、`userinfo_endpoint`、`jwks_uri` 等。

2. **Client ID / Client Secret**  
   使用在 KkOps「IdP 应用」中创建的应用的 `client_id` 与 `client_secret`。

3. **回调地址 (Redirect URI)**  
   必须与在 KkOps 中为该应用配置的「回调地址」之一完全一致（含协议、域名、路径、端口）。

4. **Scopes**  
   通常请求 `openid`，可选 `profile`、`email`。IdP 支持：`openid`、`profile`、`email`。

5. **授权流程**  
   - 用户访问第三方系统 → 重定向到 KkOps `{issuer}/authorize?response_type=code&client_id=...&redirect_uri=...&scope=openid&state=...`  
   - 未登录时跳转 KkOps IdP 登录页，登录成功后写会话并重定向回 `authorize`，生成授权码并重定向到 `redirect_uri?code=...&state=...`  
   - 第三方系统用 `code` 向 `{issuer}/token` 换取 `id_token` 与 `access_token`（POST，grant_type=authorization_code, code, client_id, client_secret, redirect_uri）。  
   - 使用 `access_token` 调用 `{issuer}/userinfo` 获取用户信息（sub、preferred_username、name、email 等）。

## 示例：Grafana 配置 KkOps 为 OIDC 提供商

- **Root URL**: `http://localhost:3000/oidc`（或你的 KkOps IdP 公网地址）  
- **Auth URL**: 可从发现文档读取，一般为 `{issuer}/authorize`  
- **Token URL**: `{issuer}/token`  
- **User Info URL**: `{issuer}/userinfo`  
- **Client ID / Client Secret**: 在 KkOps「IdP 应用」中创建的应用  
- **Redirect URI**: 在 Grafana 中配置的回调地址，须与 KkOps 中该应用的「回调地址」完全一致  

其他系统（GitLab、Jenkins、Jumpserver 等）在「OIDC / OpenID Connect」或「OAuth2」配置项中填入上述 Issuer 或发现文档 URL，并填写相同的 client_id、client_secret、redirect_uri 即可。

## 安全与运维建议

- 生产环境务必修改 `idp.session_secret`，并确保 `idp.issuer` 为 HTTPS 及正确公网地址。  
- `client_secret` 仅在创建或「重置密钥」时显示一次，请安全保存；泄露后应在 KkOps 中重置该应用的密钥。  
- 回调地址仅允许在「IdP 应用」中预配置的列表，避免开放重定向。

## SAML 2.0 与 LDAP（概览）

生产环境以 **OIDC 授权码流程** 为主。若需逐步引入其他协议，可在配置中分别开启（默认关闭）：

- **`idp.saml.enabled`**：启用后暴露 SAML 元数据与 IdP 侧脚手架路由（实现仍在演进；用于与 SP 对接前的联调占位）。
- **`idp.ldap.enabled`** 与 **`idp.ldap.listen_addr`**：启用可编辑的 LDAP 协议实验监听（当前为接受连接桩，生产请结合 TLS 证书配置）。

具体端点、证书与搜索基 DN 等以发行说明与配置参考为准；登记为 SAML/LDAP 的「IdP 应用」记录与 OAuth2/OIDC 共用同一数据模型，便于统一管理与 RBAC。
