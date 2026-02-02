# WebSSH Connection Audit – Recording and Replay

## ADDED Requirements

### Requirement: WebSSH Session Recording
The system SHALL record each WebSSH terminal session so that administrators can later review who connected to which asset and what terminal output (and optionally input) occurred.

#### Scenario: Record session on connect
- **WHEN** a user establishes a WebSSH connection to an asset and the SSH session is ready (PTY and shell started)
- **THEN** the system SHALL start recording terminal output (data sent from the server to the client over the WebSocket)
- **AND** the system MAY also record user input (data sent from the client to the server) for replay clarity
- **AND** recording SHALL be performed server-side without requiring the client to enable it

#### Scenario: Persist record on session end
- **WHEN** the WebSSH session ends (stdout EOF, WebSocket close, or timeout)
- **THEN** the system SHALL persist a connection record with metadata (user ID, username, asset ID, asset hostname, started_at, ended_at, duration) and transcript (sequence of output chunks, optionally with timestamps)
- **AND** the transcript SHALL be stored in a format suitable for later replay (e.g. JSON array of { t, d } or equivalent)
- **AND** if the transcript exceeds a configured size limit, the system SHALL truncate and mark it so that replay indicates truncation

#### Scenario: Exclude or limit binary transfer data
- **WHEN** during recording the session enters ZMODEM or other file transfer mode (binary data sent to client)
- **THEN** the system MAY omit or replace binary transfer data in the transcript to keep the record size and replay usability reasonable
- **AND** normal terminal text output SHALL still be recorded

### Requirement: Connection Audit List and Replay
The system SHALL provide an audit UI under 「系统管理」 named 「审计连线」 so that authorized users can list and replay WebSSH connection records.

#### Scenario: List connection records
- **WHEN** an authorized user (e.g. with connection-audit:read) opens the 「审计连线」 page
- **THEN** the system SHALL return a paginated list of connection records
- **AND** the list SHALL support filtering by user, asset, and time range
- **AND** each list item SHALL include at least: operator (username), target asset (hostname), start time, end time, duration
- **AND** the list SHALL NOT include the full transcript; the transcript SHALL be loaded only when the user requests replay for a specific record

#### Scenario: Replay a connection record
- **WHEN** an authorized user selects 「查看」 or 「回放」 for a connection record
- **THEN** the system SHALL return the record metadata and transcript for that record
- **AND** the UI SHALL present the terminal output (and optional input) in a read-only or replay form (e.g. sequential playback of output chunks)
- **AND** only users with connection-audit:read permission SHALL be able to access the list and replay APIs

#### Scenario: Menu and permission
- **WHEN** the application menu is rendered
- **THEN** under 「系统管理」 there SHALL be a menu item 「审计连线」 linking to the connection audit page (e.g. /connection-audit)
- **AND** visibility of the menu item and access to the page and APIs SHALL be gated by the connection-audit:read permission
