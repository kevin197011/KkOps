## MODIFIED Requirements

### Requirement: Formula Deployment Execution
The system SHALL execute Formula deployments and display results in a GitHub Actions-style interface.

#### Scenario: Deployment Run List Display
- **WHEN** user views Formula deployment page
- **THEN** system displays deployment history as workflow run list
- **AND** shows status icon, Formula name, ID, and timestamp
- **AND** provides visual distinction for different statuses
- **AND** supports click to view detailed logs

#### Scenario: Deployment Log Viewer
- **WHEN** user clicks on a deployment record
- **THEN** system opens detailed log viewer modal
- **AND** displays terminal-style deployment output
- **AND** shows per-host state apply results
- **AND** supports log copy and download
- **AND** provides fullscreen mode option

#### Scenario: Real-time Deployment Updates
- **WHEN** deployment is in progress
- **THEN** system displays real-time deployment updates
- **AND** auto-scrolls to latest log entries
- **AND** shows progress bar with completion percentage
- **AND** allows cancellation of running deployments
