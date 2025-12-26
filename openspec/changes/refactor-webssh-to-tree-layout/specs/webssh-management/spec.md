## MODIFIED Requirements

### Requirement: WebSSH Management Page Access
The system SHALL provide multiple access methods for WebSSH management:
- A dedicated menu item in the main navigation menu
- A header button positioned to the left of the user avatar dropdown
- Both access methods SHALL open the same WebSSH management page

#### Scenario: Access via header button
- **WHEN** user is authenticated and viewing any page
- **THEN** a WebSSH button/icon is visible in the header
- **AND** clicking the button opens the WebSSH management page

#### Scenario: Access via menu item
- **WHEN** user navigates to the WebSSH menu item
- **THEN** the WebSSH management page is displayed
- **AND** the page shows the tree-based host list and terminal interface

## MODIFIED Requirements

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

#### Scenario: Tree structure display
- **WHEN** user opens WebSSH management page
- **THEN** hosts are displayed in a tree structure
- **AND** projects are shown as top-level nodes
- **AND** environments are shown as second-level nodes under their project
- **AND** hosts are shown as third-level nodes under their environment

#### Scenario: Tree node expansion
- **WHEN** user clicks on a project or environment node
- **THEN** the node expands/collapses to show/hide child nodes
- **AND** the expanded state is maintained during the session

#### Scenario: Host selection from tree
- **WHEN** user clicks on a host node in the tree
- **THEN** the host is selected and highlighted
- **AND** the terminal component connects to the selected host
- **AND** the terminal is displayed in the right pane

## MODIFIED Requirements

### Requirement: WebSSH Page Layout
The WebSSH management page SHALL use a split-pane layout:
- Left pane: Tree view of hosts (resizable, default ~30% width)
- Right pane: Terminal component for selected host (default ~70% width)
- Resizable divider between panes to adjust width distribution

The layout SHALL be responsive:
- On large screens: Side-by-side split-pane layout
- On small screens: Stacked layout or drawer-based tree view

#### Scenario: Split-pane layout
- **WHEN** user opens WebSSH management page
- **THEN** the page displays with tree on left and terminal on right
- **AND** the divider between panes is resizable
- **AND** users can adjust the width of each pane

#### Scenario: Responsive layout
- **WHEN** user views WebSSH page on a small screen (< 768px)
- **THEN** the layout adapts to stacked or drawer-based view
- **AND** the terminal remains usable

## MODIFIED Requirements

### Requirement: Multiple Terminal Tabs Management
The system SHALL support multiple terminal tabs simultaneously:
- Each terminal tab represents an independent terminal connection to a host
- Users can open multiple terminal tabs for different hosts
- Each tab displays the host name and connection status
- Tabs can be switched, closed, and reordered

#### Scenario: Opening multiple terminals
- **WHEN** user selects a host from the tree
- **THEN** a new terminal tab is created for that host (if not already open)
- **AND** the terminal connection is established
- **AND** the tab displays the host name and connection status
- **WHEN** user selects another host
- **THEN** a new terminal tab is created for the new host
- **AND** both terminal tabs remain open and accessible

#### Scenario: Switching to existing terminal tab
- **WHEN** user selects a host that already has an open terminal tab
- **THEN** the system switches to the existing terminal tab
- **AND** no duplicate tab is created

#### Scenario: Terminal tab context menu
- **WHEN** user right-clicks on a terminal tab
- **THEN** a context menu is displayed with options: Clone, Close, Close All, Close Others
- **WHEN** user selects "Clone"
- **THEN** a new terminal tab is created for the same host (duplicate connection)
- **WHEN** user selects "Close"
- **THEN** the current terminal tab is closed and connection is terminated
- **WHEN** user selects "Close All"
- **THEN** all terminal tabs are closed and all connections are terminated
- **WHEN** user selects "Close Others"
- **THEN** all terminal tabs except the current one are closed

#### Scenario: Terminal tab reordering
- **WHEN** user drags a terminal tab
- **THEN** the tab can be reordered among other tabs
- **AND** the tab order is maintained during the session

#### Scenario: Terminal connection state per tab
- **WHEN** terminal is connecting to a host
- **THEN** loading indicator is shown in the tab
- **AND** the tab title shows "连接中..." status
- **WHEN** terminal is connected
- **THEN** the tab title shows host name and "已连接" status
- **WHEN** terminal is disconnected
- **THEN** the tab title shows "已断开" status

