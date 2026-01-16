# Change: Multiple SSH Connections per Host with Tab Switching

## Why

Currently, the WebSSH terminal only supports one connection at a time per asset. Users must disconnect from one host before connecting to another, which creates workflow friction when:

- **Multi-server operations**: Users need to run commands on multiple hosts simultaneously
- **Comparison work**: Users want to compare output from different hosts side-by-side
- **Context switching**: Users need to quickly switch between different hosts without losing connection state
- **Efficiency**: Reconnecting to the same host wastes time and loses session context

Enabling multiple SSH connections with tab-based switching will:
- **Improve workflow efficiency**: Connect to multiple hosts without disconnecting
- **Enable parallel operations**: Run commands on multiple hosts simultaneously
- **Better context management**: Switch between hosts instantly without losing connection state
- **Modern terminal UX**: Match the experience of modern terminal emulators (e.g., iTerm2, VS Code integrated terminal)

## What Changes

- **ADDED**: Support multiple concurrent SSH connections (one per asset)
- **ADDED**: Tab-based UI for switching between active connections
- **ADDED**: Connection list/sidebar showing all active connections
- **MODIFIED**: Terminal area to support tab switching between connections
- **MODIFIED**: Connection management to maintain multiple WebSocket connections
- **ADDED**: Visual indicators for active/active tab and connection status
- **ADDED**: Ability to close individual connections via tab close button
- **MAINTAINED**: All existing SSH connection functionality per connection
- **MAINTAINED**: Existing terminal functionality (xterm.js) for each connection

## Impact

- **UI Changes**:
  - Add tab bar above terminal area showing active connections
  - Connection list in sidebar shows connection status for each asset
  - Each tab represents one active SSH connection
- **User Experience**:
  - Users can have multiple hosts connected simultaneously
  - Quick switching between hosts via tabs
  - Better visibility of active connections
  - Ability to close individual connections
- **Code Changes**:
  - Frontend: Refactor to manage multiple connections and terminal instances
  - Frontend: Add tab UI component for connection switching
  - Frontend: Maintain connection state map (asset ID -> connection metadata)
  - Backend: No changes needed (each WebSocket connection is already independent)

## Layout Structure

### New Tab-Based Terminal Layout
```
┌─────────────────────────────────────────────────────────┐
│  [Back]  WebSSH Terminal                                │
├──────────────┬──────────────────────────────────────────┤
│              │  [Tab: Host1 ✓] [Tab: Host2] [Tab: Host3]│
│  Asset List  │  ┌──────────────────────────────────┐   │
│  ┌────────┐  │  │                                  │   │
│  │ Host1  │  │  │  Terminal output...              │   │
│  │ [2] ✓  │  │  │  (for active tab)                │   │
│  └────────┘  │  │                                  │   │
│  ┌────────┐  │  │                                  │   │
│  │ Host2  │  │  └──────────────────────────────────┘   │
│  │ [1]    │  │                                          │
│  └────────┘  │  [Disconnect Active] [Reconnect]        │
│  ┌────────┐  │                                          │
│  │ Host3  │  │                                          │
│  │ [1]    │  │                                          │
│  └────────┘  │                                          │
│              │                                          │
│  [Search]    │                                          │
└──────────────┴──────────────────────────────────────────┘
```

Notes:
- Numbers in brackets `[2]` indicate number of active connections to that host
- ✓ indicates currently selected/active tab
- Clicking an asset in the list creates a new connection (if not already connected) or switches to existing connection

## Considerations

- **Connection Management**: 
  - Each connection maintains its own WebSocket, Terminal instance, and state
  - Connection state includes: asset, terminal instance, WebSocket reference, connection status
- **Tab Management**:
  - Tabs show asset hostname/IP and connection status
  - Active tab is highlighted
  - Close button (×) on each tab to disconnect
  - Click tab to switch active connection
- **Terminal Instances**:
  - Each terminal instance (xterm.js Terminal) must be properly initialized and disposed
  - Switching tabs should hide/show terminal instances (or dispose and recreate)
  - Terminal state (scrollback, cursor position) should be preserved when switching
- **Resource Management**:
  - Limit maximum concurrent connections (e.g., 10) to prevent resource exhaustion
  - Properly cleanup WebSocket connections and terminal instances when closing
  - Show warning when approaching connection limit
- **Performance**:
  - Lazy initialize terminal instances only when tab is first opened
  - Consider virtual rendering for inactive terminals (if performance issues arise)
- **UI/UX**:
  - Clear visual feedback for active/inactive tabs
  - Connection status indicators (connected, connecting, disconnected, error)
  - Keyboard shortcuts for tab switching (Ctrl+Tab, Ctrl+Shift+Tab, Ctrl+W to close)
