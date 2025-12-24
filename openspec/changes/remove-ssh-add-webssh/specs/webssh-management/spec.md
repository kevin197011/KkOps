## ADDED Requirements

### Requirement: WebSSH Terminal Access
The system SHALL provide interactive SSH terminal access directly from host management, using host data for connection configuration.

#### Scenario: Open terminal for a host
- **WHEN** a user clicks "打开终端" button for a host in the host management page
- **THEN** the system SHALL open a WebSSH terminal in a new tab
- **AND** the system SHALL use the host's data for SSH connection:
  - hostname: from `host.hostname` or `host.ip_address`
  - port: from `host.ssh_port` (default: 22 if not set)
- **AND** the system SHALL prompt the user for username and password via the terminal interface
- **AND** the system SHALL establish the SSH connection using the provided credentials
- **AND** the system SHALL display the interactive terminal interface

#### Scenario: Terminal connection with host data
- **WHEN** a WebSSH terminal connection is established
- **THEN** the system SHALL use the host's IP address or hostname for connection
- **AND** the system SHALL use the host's SSH port (default: 22)
- **AND** the system SHALL NOT require pre-configured SSH connection data
- **AND** the system SHALL establish the connection on-demand

#### Scenario: Multiple terminal sessions
- **WHEN** a user opens multiple terminals for different hosts
- **THEN** the system SHALL support multiple concurrent terminal sessions
- **AND** each terminal session SHALL be displayed in a separate tab
- **AND** each terminal session SHALL maintain its own WebSocket connection
- **AND** terminal sessions SHALL be independent (input/output isolated)

#### Scenario: Terminal authentication
- **WHEN** a user opens a WebSSH terminal
- **THEN** the system SHALL authenticate the user via JWT token in WebSocket handshake
- **AND** the system SHALL prompt for SSH username and password via the terminal interface
- **AND** the system SHALL establish SSH connection using the provided credentials
- **AND** the system SHALL NOT store SSH credentials (password entered each time)

### Requirement: Terminal Controls
The system SHALL provide controls for managing terminal sessions, including connection, disconnection, and terminal manipulation.

#### Scenario: Connect terminal
- **WHEN** a user opens a terminal tab
- **THEN** the system SHALL automatically attempt to establish the WebSocket connection
- **AND** the system SHALL display connection status (connecting, connected, disconnected)
- **AND** the system SHALL prompt for credentials when connection is established

#### Scenario: Disconnect terminal
- **WHEN** a user clicks "断开" button or closes the terminal tab
- **THEN** the system SHALL close the WebSocket connection
- **AND** the system SHALL close the SSH connection
- **AND** the system SHALL update the terminal display to show disconnected state

#### Scenario: Clear terminal
- **WHEN** a user clicks "清空" button
- **THEN** the system SHALL clear the terminal display
- **AND** the system SHALL NOT disconnect the SSH connection
- **AND** the system SHALL maintain the terminal session

#### Scenario: Resize terminal
- **WHEN** a user resizes the browser window or terminal container
- **THEN** the system SHALL automatically resize the terminal to fit the container
- **AND** the system SHALL send terminal size (rows/columns) to the SSH server
- **AND** the system SHALL update the terminal display accordingly

### Requirement: Terminal I/O
The system SHALL provide real-time bidirectional communication between the browser terminal and the SSH server.

#### Scenario: Terminal input
- **WHEN** a user types in the terminal
- **THEN** the system SHALL send the input to the SSH server via WebSocket
- **AND** the system SHALL display the input in the terminal (with echo)
- **AND** the system SHALL support special keys (Enter, Backspace, Arrow keys, etc.)

#### Scenario: Terminal output
- **WHEN** the SSH server sends output (stdout)
- **THEN** the system SHALL receive the output via WebSocket
- **AND** the system SHALL display the output in the terminal in real-time
- **AND** the system SHALL support ANSI color codes and terminal formatting

#### Scenario: Terminal error output
- **WHEN** the SSH server sends error output (stderr)
- **THEN** the system SHALL receive the error output via WebSocket
- **AND** the system SHALL display the error output in the terminal
- **AND** the system SHALL handle error output the same way as standard output

#### Scenario: Terminal paste
- **WHEN** a user pastes text into the terminal
- **THEN** the system SHALL send the pasted text to the SSH server
- **AND** the system SHALL display the pasted text in the terminal
- **AND** the system SHALL support multi-line paste

### Requirement: Terminal Integration with Host Management
The system SHALL integrate WebSSH terminal access into the host management page.

#### Scenario: Terminal access from host list
- **WHEN** a user views the host management page
- **THEN** the system SHALL display an "打开终端" button for each host
- **AND** the system SHALL allow opening a terminal for any host in the list
- **AND** the system SHALL support opening multiple terminals for different hosts

#### Scenario: Terminal tab management
- **WHEN** a user opens multiple terminals
- **THEN** the system SHALL display each terminal in a separate tab
- **AND** the system SHALL allow switching between terminal tabs
- **AND** the system SHALL allow closing individual terminal tabs
- **AND** the system SHALL maintain terminal sessions when switching tabs

#### Scenario: Terminal status display
- **WHEN** a terminal is open
- **THEN** the system SHALL display the host name in the terminal tab title
- **AND** the system SHALL display connection status (connected, disconnected, connecting)
- **AND** the system SHALL display connection status indicator in the terminal header

