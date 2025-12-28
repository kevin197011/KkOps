## MODIFIED Requirements

### Requirement: Result Collection and Display
The system SHALL collect and present batch operation results in a GitHub Actions-style interface.

#### Scenario: Workflow Run List Display
- **WHEN** user views batch operations page
- **THEN** system displays execution history as workflow run list
- **AND** shows status icon, name, ID, command, and timestamp
- **AND** provides visual distinction for different statuses
- **AND** supports click to view detailed logs

#### Scenario: Execution Log Viewer
- **WHEN** user clicks on an execution record
- **THEN** system opens detailed log viewer modal
- **AND** displays terminal-style log output
- **AND** shows per-host execution results
- **AND** supports log copy and download
- **AND** provides fullscreen mode option

#### Scenario: Real-time Log Updates
- **WHEN** execution is in progress
- **THEN** system displays real-time log updates
- **AND** auto-scrolls to latest log entries
- **AND** shows progress bar with completion percentage
- **AND** allows cancellation of running operations
