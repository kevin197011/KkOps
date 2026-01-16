# Tasks: Multiple SSH Connections per Host with Tab Switching

## Frontend Implementation

- [x] **Refactor connection state management**
  - [x] Create connection state structure (connection ID, asset, WebSocket ref, terminal instance, status)
  - [x] Replace single connection state with connection map/array
  - [x] Update connection state management to support multiple connections

- [x] **Implement tab UI component**
  - [x] Add tab bar above terminal area
  - [x] Create Tab component showing asset info and close button
  - [x] Implement tab switching logic
  - [x] Add visual indicators for active/inactive tabs
  - [x] Style tabs to match design system

- [x] **Multiple terminal instances**
  - [x] Refactor terminal initialization to support multiple instances
  - [x] Create terminal instance manager (create, switch, dispose)
  - [x] Ensure each terminal has its own container and xterm.js instance
  - [x] Implement terminal visibility switching (show/hide or dispose/recreate)

- [x] **Connection management**
  - [x] Update `handleConnect` to support creating new connections without disconnecting existing
  - [x] Add connection tracking (which asset has active connection(s))
  - [x] Implement connection close/disconnect for individual connections
  - [x] Add connection status indicators in asset list

- [x] **Asset list updates**
  - [x] Show connection count badge for assets with active connections
  - [x] Update click behavior: create new connection or switch to existing
  - [x] Visual indicators for assets with active connections
  - [ ] Support right-click context menu for connection actions (future enhancement)

- [x] **WebSocket management**
  - [x] Maintain WebSocket references per connection
  - [x] Handle WebSocket messages routing to correct terminal instance
  - [x] Proper cleanup when closing connections
  - [x] Error handling per connection

- [x] **UI enhancements**
  - [x] Add "Close All Connections" button
  - [ ] Add "Close Other Tabs" context menu option (future enhancement)
  - [x] Connection limit warning/notification (MAX_CONNECTIONS = 10)
  - [x] Loading states for new connections

- [ ] **Testing and refinement**
  - [ ] Test multiple connections to same asset
  - [ ] Test switching between connections
  - [ ] Test closing individual connections
  - [ ] Test connection cleanup and resource management
  - [ ] Test edge cases (connection failures, network issues)

## Backend (No Changes Required)

- [x] **Verify backend compatibility**
  - [x] Confirm existing WebSocket handler supports multiple concurrent connections (it should)
  - [x] No backend changes needed as each connection is independent

## Documentation

- [ ] Update user documentation for multi-connection feature
- [ ] Add keyboard shortcuts documentation
- [ ] Update developer documentation if needed
