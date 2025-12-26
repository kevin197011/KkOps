## MODIFIED Requirements

### Requirement: WebSSH Access Method
The system SHALL provide WebSSH terminal access through a modal dialog when triggered from the header button, while maintaining backward compatibility for direct page access via routing.

#### Scenario: Access via header button (modal)
- **WHEN** user clicks the WebSSH button in the header
- **THEN** a full-screen modal dialog opens displaying the WebSSH interface
- **AND** the current page remains visible in the background (with semi-transparent overlay)
- **AND** the modal contains the complete WebSSH interface (tree view and terminal tabs)
- **WHEN** user closes the modal (via close button, ESC key, or backdrop click)
- **THEN** the modal closes and user returns to the previous page context
- **AND** any active terminal connections are terminated

#### Scenario: Access via direct route (backward compatibility)
- **WHEN** user navigates directly to `/webssh` route
- **THEN** the WebSSH interface is displayed as a full page (existing behavior)
- **AND** all WebSSH functionality works as before
- **AND** the interface is not wrapped in a modal

#### Scenario: Modal behavior
- **WHEN** the WebSSH modal is open
- **THEN** the modal is centered on screen
- **AND** the modal uses responsive sizing (95% viewport width, 90% viewport height, with maximum dimensions)
- **AND** the modal has a close button (X) in the header
- **AND** pressing ESC key closes the modal
- **AND** clicking the backdrop (mask) may close the modal (configurable)
- **WHEN** user interacts with terminal in the modal
- **THEN** all terminal functionality works identically to the page version
- **AND** terminal connections, authentication, and tab management function correctly

#### Scenario: Modal state management
- **WHEN** user opens the WebSSH modal
- **THEN** the modal state is managed locally in the MainLayout component
- **WHEN** user closes the modal
- **THEN** all terminal connections are terminated
- **WHEN** user reopens the modal
- **THEN** the modal opens fresh (no preserved terminal state)
- **AND** users need to reconnect to terminals

## KEPT Requirements

### Requirement: WebSSH Host List Display
The system SHALL display hosts in a hierarchical tree structure organized by:
- Level 1: Project (grouping hosts by project)
- Level 2: Environment (grouping hosts by environment within each project: dev, test, prod, etc.)
- Level 3: Hosts (individual host nodes)

The tree SHALL support:
- Expand/collapse of project and environment nodes
- Visual indicators for host status (online/offline)
- Host count badges on project and environment nodes
- Search and filter functionality

### Requirement: Multiple Terminal Tabs Management
The system SHALL support multiple terminal tabs simultaneously:
- Each terminal tab represents an independent terminal connection to a host
- Users can open multiple terminal tabs for different hosts
- Each tab displays the host name and connection status
- Tabs can be switched, closed, and reordered

### Requirement: Terminal Tab Context Menu
Each terminal tab SHALL support right-click context menu with actions:
- **Clone**: Create a duplicate terminal connection to the same host (new tab, independent connection)
- **Close**: Close the current terminal tab
- **Close All**: Close all terminal tabs
- **Close Others**: Close all terminal tabs except the current one

### Requirement: Terminal Connection Management
The system SHALL support terminal connections with:
- WebSocket-based connection to backend
- SSH authentication (password or SSH key)
- Terminal resize handling
- Connection status indicators
- Error handling and reconnection logic

