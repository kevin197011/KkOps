## ADDED Requirements

### Requirement: AI provider abstraction

The system SHALL expose an internal `Provider` interface with `Chat`, `Embed`, and `Name` operations. Implementations SHALL use HTTP only (no vendor SDK lock-in) for OpenAI-compatible endpoints, Anthropic Messages API, and a generic POST webhook that returns plain text.

#### Scenario: Encrypted integration config

- **WHEN** an administrator stores an integration with `kind` equal to `ai` and a JSON configuration including provider type and credentials
- **THEN** the server SHALL encrypt credentials using the existing integration encryption key and SHALL NOT return secrets from listing APIs.

### Requirement: AI provider APIs

The system SHALL expose `GET /api/v1/ai/providers` listing configured AI integrations (metadata only) and `POST /api/v1/ai/test` sending a minimal ping message through the default or selected provider. Both SHALL require permission `ai:*`.

#### Scenario: Test ping

- **WHEN** an authorized user calls `POST /api/v1/ai/test` with optional `integration_id`
- **THEN** the server SHALL invoke the provider chat with a short ping prompt and SHALL return success or a typed error without logging full prompts at info level.
