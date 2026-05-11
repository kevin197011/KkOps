## ADDED Requirements

### Requirement: Registry APIs

The system SHALL expose read-only Harbor-oriented HTTP endpoints for listing repositories, listing tags for a repository, and fetching vulnerability payload for a tag reference. Endpoints SHALL require `registry:*` and SHALL accept `integration_id` identifying a harbor integration.

#### Scenario: Harbor error

- **WHEN** Harbor returns an error or the integration is not kind harbor
- **THEN** the server SHALL return an appropriate HTTP error and message without panicking.

### Requirement: Registry UI

The web application SHALL provide `/registry` to select a harbor integration, browse repositories, inspect tags, and view raw vulnerability JSON when available.

#### Scenario: Browse tags

- **WHEN** the user selects a repository and requests tags
- **THEN** the UI SHALL list tag names (and related metadata when provided) and allow opening vulnerability details for a selected tag.
