## ADDED Requirements

### Requirement: IdP application protocol

The system SHALL persist a protocol discriminator on each IdP application record with allowed values `oidc`, `saml`, and `ldap`, defaulting to `oidc` for existing rows.

#### Scenario: Protocol returned in API

- **WHEN** a client reads or lists IdP applications
- **THEN** each record SHALL include the `protocol` field.

### Requirement: SAML IdP scaffold

The system SHALL expose optional SAML-related HTTP endpoints when `idp.saml.enabled` is true, sufficient for validating configuration and future SSO flows.

#### Scenario: SAML disabled by default

- **WHEN** `idp.saml.enabled` is false
- **THEN** SAML-specific routes SHALL NOT be registered or SHALL return disabled responses consistent with existing IdP enablement patterns.

### Requirement: LDAP IdP scaffold

The system SHALL support an optional LDAP listener when `idp.ldap.enabled` is true, using configuration for bind-oriented read-only verification.

#### Scenario: LDAP disabled by default

- **WHEN** `idp.ldap.enabled` is false
- **THEN** the application SHALL not listen for LDAP clients.

### Requirement: Documentation

The system documentation SHALL describe OIDC, SAML, and LDAP IdP modes at a high level for operators.

#### Scenario: Docs mention protocols

- **WHEN** an operator reads `docs/idp-oidc-provider.md`
- **THEN** the document SHALL reference SAML and LDAP alongside OIDC.
