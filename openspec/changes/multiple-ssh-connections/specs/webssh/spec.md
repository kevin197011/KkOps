# WebSSH Terminal - Multiple Connections

## MODIFIED Requirements

### Multiple Concurrent Connections

#### Requirement: Support Multiple SSH Connections
The WebSSH terminal MUST support multiple concurrent SSH connections to different assets (hosts). Each connection MUST be independent and maintain its own terminal session state.

**Rationale**: Users need to operate on multiple hosts simultaneously without disconnecting from existing connections.

#### Scenario: Multiple Host Connections
```
Given the user has an active SSH connection to Host1
When the user clicks on Host2 in the asset list
Then a new SSH connection to Host2 should be established
And the existing connection to Host1 should remain active
And the user should be able to switch between Host1 and Host2 connections
```

#### Scenario: Multiple Connections to Same Host
```
Given the user has an active SSH connection to Host1
When the user clicks on Host1 again in the asset list
Then a second SSH connection to Host1 should be established
And both connections should be active and independently manageable
```

### Tab-Based Connection Switching

#### Requirement: Tab UI for Connection Management
The terminal area MUST display a tab bar showing all active connections. Each tab MUST represent one active SSH connection and MUST allow the user to switch between connections.

**Rationale**: Tabs provide an intuitive way to manage and switch between multiple connections.

#### Scenario: Tab Switching
```
Given the user has multiple active SSH connections (Host1, Host2, Host3)
When the user clicks on the Host2 tab
Then the Host2 terminal should become visible
And the Host1 and Host3 terminals should be hidden
And the Host2 tab should be visually indicated as active
```

#### Scenario: Close Connection via Tab
```
Given the user has an active SSH connection to Host1 displayed in a tab
When the user clicks the close button (×) on the Host1 tab
Then the SSH connection to Host1 should be terminated
And the Host1 tab should be removed
And the WebSocket connection should be properly closed
And the terminal instance should be disposed
```

### Connection State Management

#### Requirement: Connection State Tracking
The system MUST track the state of each active connection, including: connection ID, associated asset, WebSocket reference, terminal instance, and connection status (connecting, connected, disconnected, error).

**Rationale**: Proper state management is essential for managing multiple connections.

#### Scenario: Connection Status Tracking
```
Given the user has created connections to Host1 and Host2
When Host1 connection is active and Host2 connection is inactive
Then the system should track both connections as active
And switching to Host2 should restore its terminal state
And connection status should be correctly displayed in the asset list
```

### Asset List Connection Indicators

#### Requirement: Visual Connection Indicators
The asset list MUST display visual indicators for assets that have active connections, including the number of active connections per asset.

**Rationale**: Users need to quickly identify which assets have active connections.

#### Scenario: Connection Badge Display
```
Given Host1 has 2 active connections and Host2 has 1 active connection
When viewing the asset list
Then Host1 should display a badge showing "[2]" connections
And Host2 should display a badge showing "[1]" connection
And assets without connections should display no badge
```

#### Scenario: Asset List Click Behavior
```
Given Host1 has an existing active connection
When the user clicks on Host1 in the asset list
Then the system should switch to the existing Host1 connection tab
And if all Host1 tabs are closed, clicking should create a new connection
```

## ADDED Requirements

### Connection Limit

#### Requirement: Maximum Concurrent Connections
The system SHOULD enforce a reasonable maximum number of concurrent connections (e.g., 10) to prevent resource exhaustion. When the limit is reached, the system SHOULD display a warning to the user.

**Rationale**: Limiting connections prevents browser and server resource exhaustion.

#### Scenario: Connection Limit Enforcement
```
Given the user has 10 active connections (the maximum limit)
When the user attempts to create a new connection
Then the system should display a warning message
And the new connection should not be created
And the user should be prompted to close an existing connection first
```

### Connection Cleanup

#### Requirement: Proper Resource Cleanup
When a connection is closed, the system MUST properly cleanup all associated resources, including: WebSocket connection, terminal instance, event listeners, and memory references.

**Rationale**: Proper cleanup prevents memory leaks and resource exhaustion.

#### Scenario: Connection Cleanup
```
Given the user closes a connection via the tab close button
When the connection is terminated
Then the WebSocket should be properly closed
And the xterm.js terminal instance should be disposed
And all event listeners should be removed
And memory references should be cleared
```

### Terminal Instance Management

#### Requirement: Independent Terminal Instances
Each SSH connection MUST have its own independent xterm.js terminal instance. Terminal instances MUST be properly initialized when connections are created and disposed when connections are closed.

**Rationale**: Each connection needs its own terminal to maintain independent state.

#### Scenario: Terminal Instance Lifecycle
```
Given the user creates a new connection to Host1
When the connection is established
Then a new xterm.js Terminal instance should be created
And the terminal should be initialized with proper configuration
And the terminal should be displayed in the terminal area
When the connection is closed
Then the terminal instance should be properly disposed
And all terminal-related resources should be freed
```

## Considerations

- **Performance**: Multiple terminal instances may impact performance. Consider lazy initialization and virtual rendering for inactive terminals.
- **Memory**: Terminal instances and WebSocket connections consume memory. Implement proper cleanup and consider connection limits.
- **User Experience**: Provide clear visual feedback for connection status and active/inactive tabs.
- **Keyboard Shortcuts**: Consider adding keyboard shortcuts for tab navigation (Ctrl+Tab, Ctrl+Shift+Tab, Ctrl+W).
