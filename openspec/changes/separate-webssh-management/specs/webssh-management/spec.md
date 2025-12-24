## MODIFIED Requirements

### Requirement: Terminal Integration with Host Management
The system SHALL provide a dedicated WebSSH management page for terminal access, separate from host CRUD operations.

#### Scenario: Terminal access from WebSSH management page
- **WHEN** a user navigates to the WebSSH management page
- **THEN** the system SHALL display a list of all hosts
- **AND** the system SHALL display an "打开终端" button for each host in the operations column
- **AND** the system SHALL allow opening a terminal for any host in the list
- **AND** the system SHALL support opening multiple terminals for different hosts

#### Scenario: Terminal tab management
- **WHEN** a user opens multiple terminals from the WebSSH management page
- **THEN** the system SHALL display each terminal in a separate tab
- **AND** the system SHALL allow switching between terminal tabs
- **AND** the system SHALL allow closing individual terminal tabs
- **AND** the system SHALL maintain terminal sessions when switching tabs

#### Scenario: Terminal status display
- **WHEN** a terminal is open in the WebSSH management page
- **THEN** the system SHALL display the host name in the terminal tab title
- **AND** the system SHALL display connection status (connected, disconnected, connecting)
- **AND** the system SHALL display connection status indicator in the terminal header

#### Scenario: WebSSH management page access
- **WHEN** a user wants to access WebSSH terminals
- **THEN** the system SHALL provide a "WebSSH管理" menu item in the main navigation
- **AND** the system SHALL navigate to the WebSSH management page when the menu item is clicked
- **AND** the system SHALL NOT display "打开终端" button in the host management page operations column

