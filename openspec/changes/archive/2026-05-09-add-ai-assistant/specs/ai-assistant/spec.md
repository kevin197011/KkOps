## ADDED Requirements

### Requirement: Streaming chat API

The system SHALL expose `POST /api/v1/ai/chat` accepting optional `session_id`, message history, and optional `context`. The response SHALL use `Content-Type: text/event-stream` and stream assistant tokens as SSE events until completion.

#### Scenario: Authorized chat

- **WHEN** a user with `ai:*` sends a chat request
- **THEN** the server SHALL persist messages to the session when `session_id` is provided and SHALL stream model output to the client.

### Requirement: Tool bridge

The assistant SHALL resolve internal read-only tools including `list_alerts`, `get_incident`, `query_metric`, `search_logs`, `list_pods`, `pipeline_status`, and `list_integrations` via a marker protocol (`<<TOOL: name args>>`) embedded in model output, re-injecting JSON results into the conversation.

#### Scenario: Tool round-trip

- **WHEN** the model output contains a valid tool marker
- **THEN** the server SHALL execute the tool server-side and SHALL append the result before continuing generation.

### Requirement: Chat sessions UI

The web application SHALL provide `/ai/chat` with a session list, streaming transcript, and markdown rendering for assistant messages.

#### Scenario: Keyboard shortcut

- **WHEN** the user presses Ctrl+L or Cmd+L while the shell is focused
- **THEN** the application SHALL navigate to `/ai/chat`.
