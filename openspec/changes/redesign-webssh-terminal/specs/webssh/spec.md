## MODIFIED Requirements

### Requirement: WebSSH Terminal Access
The system SHALL provide access to WebSSH terminal functionality through a header button and standalone page interface.

#### Scenario: Access WebSSH from header
- **WHEN** a user is viewing any page in the application
- **THEN** a WebSSH button is visible in the header (top right area)
- **AND** clicking the button navigates to the WebSSH terminal page
- **AND** the button is always visible regardless of current page

#### Scenario: WebSSH as standalone page
- **WHEN** a user navigates to the WebSSH terminal page
- **THEN** the page is displayed without the standard MainLayout wrapper
- **AND** the page has a minimal header with back button and title
- **AND** the page uses full viewport height for terminal interface

### Requirement: WebSSH Terminal Connection
The system SHALL provide browser-based SSH terminal access to remote hosts with an improved split-view interface.

#### Scenario: Split-view layout
- **WHEN** a user opens the WebSSH terminal page
- **THEN** the page displays in a split-view layout:
  - **LEFT PANEL**: Asset list sidebar showing all available assets
  - **RIGHT PANEL**: SSH terminal connection area
- **AND** the left panel is fixed width (350px)
- **AND** the right panel takes remaining width
- **AND** both panels use full height

#### Scenario: Asset list in sidebar
- **WHEN** a user views the WebSSH terminal page
- **THEN** the left sidebar displays a list of all available assets
- **AND** each asset item shows: name, IP address, and SSH user (if available)
- **AND** a search/filter input is available at the top of the asset list
- **AND** the currently connected asset is highlighted/indicated
- **AND** assets can be searched by name or IP address

#### Scenario: Direct asset connection
- **WHEN** a user clicks on an asset in the left sidebar
- **THEN** the system immediately attempts to connect to that asset
- **AND** no modal dialog is shown
- **AND** connection status is displayed in the asset list item
- **AND** if already connected to another asset, the current connection is disconnected before connecting to the new asset

#### Scenario: Terminal display
- **WHEN** a user is connected to an asset via WebSSH
- **THEN** the terminal is displayed in the right panel
- **AND** the current asset information is shown above the terminal
- **AND** connection controls (disconnect button) are available
- **AND** terminal input/output works bidirectionally
- **AND** the terminal resizes correctly when window is resized

#### Scenario: Multiple terminal sessions
- **WHEN** a user is connected to an asset
- **AND** the user clicks on a different asset
- **THEN** the system disconnects from the current asset
- **AND** connects to the newly selected asset
- **AND** the terminal displays the new connection

#### Scenario: Terminal input/output
- **WHEN** a user types commands in the terminal
- **THEN** input is sent to the remote host via SSH
- **AND** output from the remote host is displayed in the terminal
- **AND** the communication is bidirectional and real-time

#### Scenario: Terminal resize
- **WHEN** a user resizes the browser window
- **THEN** the terminal PTY size is updated on the remote host
- **AND** applications on the remote host adjust to the new size

#### Scenario: Disconnect SSH session
- **WHEN** a user disconnects an SSH session
- **THEN** the SSH connection to the remote host is closed
- **AND** the WebSocket connection is closed
- **AND** the asset list is updated to show no active connection
- **AND** the terminal is cleared and ready for a new connection

## REMOVED Requirements

### Requirement: WebSSH Menu Item
**Reason**: WebSSH access has been moved to header button for better discoverability and reduced menu clutter.

**Migration**: Users can now access WebSSH via the header button (top right) instead of the sidebar menu.