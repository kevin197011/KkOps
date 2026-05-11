# Change: Multi-protocol IdP (SAML and LDAP scaffolds)

## Why

Some enterprise systems cannot use OIDC; KkOps should offer SAML 2.0 and LDAP (read-only bind) entry points behind feature flags while keeping a single application registry.

## What Changes

- OpenSpec delta for `multi-protocol-idp`.
- Backend: extend `IdPConfig` with `saml` and `ldap` flags; SAML endpoint scaffold (signed response stub where applicable); LDAP listener scaffold; `oauth2_clients.protocol` column (`oidc` | `saml` | `ldap`).
- Frontend: rename "OAuth2 应用" to "IdP 应用", protocol filter and form field; document high-level in `docs/idp-oidc-provider.md`.

## Impact

- Affected specs: `multi-protocol-idp` (new delta).
- Affected code: `backend/internal/config`, `backend/internal/model/oauth2_client.go`, `backend/internal/idp/`, `backend/cmd/server/main.go`, `backend/internal/handler/oauth2client`, `frontend`, `docs/idp-oidc-provider.md`.
