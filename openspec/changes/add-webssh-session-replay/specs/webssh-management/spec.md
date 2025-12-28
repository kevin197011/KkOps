## ADDED Requirements

### Requirement: WebSSH Session Recording
The system SHALL record all terminal input and output during WebSSH sessions for later replay and audit purposes.

#### Scenario: Start session recording
- **WHEN** a WebSSH session is established
- **THEN** the system SHALL create a session record with metadata (user, host, start time)
- **AND** the system SHALL begin recording all terminal events
- **AND** the system SHALL assign a unique session ID

#### Scenario: Record terminal input
- **WHEN** user types commands or presses keys in the terminal
- **THEN** the system SHALL record the input event with timestamp
- **AND** the system SHALL store the input data (commands, special keys)
- **AND** the timestamp SHALL be relative to session start time (millisecond precision)

#### Scenario: Record terminal output
- **WHEN** terminal displays output from remote host
- **THEN** the system SHALL record the output event with timestamp
- **AND** the system SHALL store the output data (including ANSI escape sequences)
- **AND** the system SHALL preserve the exact output format

#### Scenario: Record terminal resize
- **WHEN** user resizes the terminal window
- **THEN** the system SHALL record the resize event with timestamp
- **AND** the system SHALL store the new terminal dimensions (rows, columns)

#### Scenario: End session recording
- **WHEN** WebSSH session is closed
- **THEN** the system SHALL finalize the session record
- **AND** the system SHALL calculate session duration
- **AND** the system SHALL compress and store all recorded events
- **AND** the system SHALL update session status to "completed"

### Requirement: Session Record Storage
The system SHALL store session records efficiently with support for querying and retrieval.

#### Scenario: Store session metadata
- **WHEN** session recording is created
- **THEN** the system SHALL store session metadata in `webssh_sessions` table
- **AND** metadata SHALL include: user_id, host_id, hostname, started_at, status
- **AND** the system SHALL create appropriate database indexes

#### Scenario: Store session events
- **WHEN** terminal events are recorded
- **THEN** the system SHALL store events in `webssh_session_records` table
- **AND** events SHALL be stored with session_id, event_type, timestamp_ms, data
- **AND** the system SHALL compress event data to reduce storage size
- **AND** the system SHALL maintain event order by timestamp

#### Scenario: Query session list
- **WHEN** user requests session list
- **THEN** the system SHALL return sessions with pagination
- **AND** the system SHALL support filtering by user, host, date range, status
- **AND** the system SHALL support sorting by start time (descending by default)
- **AND** the system SHALL only return sessions accessible to the user

### Requirement: Session Replay
The system SHALL provide the ability to replay recorded WebSSH sessions with time control.

#### Scenario: Load session for replay
- **WHEN** user selects a session to replay
- **THEN** the system SHALL load session metadata
- **AND** the system SHALL load all session events
- **AND** the system SHALL initialize replay interface

#### Scenario: Play session replay
- **WHEN** user starts replay
- **THEN** the system SHALL display terminal output in real-time
- **AND** the system SHALL respect original timestamps and delays
- **AND** the system SHALL show replay progress (time elapsed / total duration)
- **AND** the system SHALL support playback speed control (0.5x, 1x, 2x, 4x)

#### Scenario: Pause and resume replay
- **WHEN** user pauses replay
- **THEN** the system SHALL stop playback
- **AND** the system SHALL preserve current terminal state
- **WHEN** user resumes replay
- **THEN** the system SHALL continue from paused position

#### Scenario: Jump to time point
- **WHEN** user jumps to a specific time in the session
- **THEN** the system SHALL fast-forward to that time point
- **AND** the system SHALL reconstruct terminal state up to that point
- **AND** the system SHALL continue playback from that point

#### Scenario: Terminal display during replay
- **WHEN** replay is playing
- **THEN** the system SHALL display terminal output using xterm.js
- **AND** the system SHALL support ANSI escape sequences (colors, formatting)
- **AND** the system SHALL handle terminal resize events
- **AND** the system SHALL show input events (commands typed)

### Requirement: Session Management
The system SHALL provide session management capabilities for viewing, searching, and deleting sessions.

#### Scenario: View session list
- **WHEN** user navigates to session list
- **THEN** the system SHALL display sessions in a table
- **AND** the system SHALL show: session ID, hostname, user, start time, duration, status
- **AND** the system SHALL support pagination
- **AND** the system SHALL support search by hostname or user

#### Scenario: View session details
- **WHEN** user clicks on a session
- **THEN** the system SHALL display session details
- **AND** details SHALL include: metadata, statistics (total events, data size)
- **AND** the system SHALL provide "Replay" button

#### Scenario: Delete session
- **WHEN** user deletes a session
- **THEN** the system SHALL validate user has permission
- **AND** the system SHALL delete session and all associated records
- **AND** the system SHALL log deletion in audit log
- **AND** the system SHALL confirm deletion before proceeding

#### Scenario: Filter sessions
- **WHEN** user applies filters
- **THEN** the system SHALL filter sessions by:
  - User (current user or all users for admin)
  - Host (specific host or all hosts)
  - Date range (start date, end date)
  - Status (active, completed, error)
- **AND** the system SHALL update session list accordingly

### Requirement: Replay Interface
The system SHALL provide a user-friendly replay interface with playback controls.

#### Scenario: Replay controls
- **WHEN** replay interface is displayed
- **THEN** the system SHALL provide controls:
  - Play/Pause button
  - Speed selector (0.5x, 1x, 2x, 4x)
  - Progress bar (showing current time / total duration)
  - Time display (current time / total time)
  - Jump to time input
- **AND** controls SHALL be clearly labeled and accessible

#### Scenario: Session information display
- **WHEN** replay is active
- **THEN** the system SHALL display session information:
  - Hostname
  - User who recorded the session
  - Session start time
  - Session duration
  - Total events count

#### Scenario: Terminal display
- **WHEN** replay is playing
- **THEN** the system SHALL display terminal in a resizable container
- **AND** the system SHALL support full-screen mode
- **AND** the system SHALL maintain terminal aspect ratio
- **AND** the system SHALL support scrolling through terminal history

## MODIFIED Requirements

### Requirement: WebSSH Terminal Access
The system SHALL provide web-based SSH terminal access with real-time interaction **and session recording**.

#### Scenario: Terminal Connection (Modified)
- **WHEN** user clicks WebSSH button for a host
- **THEN** system opens modal terminal window
- **AND** establishes WebSocket connection to backend
- **AND** displays connection status
- **AND** automatically starts session recording (if enabled)

