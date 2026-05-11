## 1. OpenSpec

- [x] 1.1 Add proposal, tasks, and `specs/multi-protocol-idp/spec.md`
- [x] 1.2 Run `openspec validate add-multi-protocol-idp --strict`

## 2. Backend

- [x] 2.1 Add `IdP` SAML/LDAP config structs and defaults; wire listeners/routes when enabled
- [x] 2.2 SAML scaffold endpoint(s); LDAP TLS-capable stub listener
- [x] 2.3 Migrate `oauth2_clients.protocol`; expose in OAuth2 client API

## 3. Frontend + docs

- [x] 3.1 Rename UI to IdP 应用; protocol filter + form; API types
- [x] 3.2 Update `docs/idp-oidc-provider.md`

## 4. Verification

- [x] 4.1 `go build ./...` and `npm run build`
