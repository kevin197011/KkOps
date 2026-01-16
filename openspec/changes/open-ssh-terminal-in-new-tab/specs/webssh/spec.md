## MODIFIED Requirements

### Requirement: WebSSH Terminal Access
The system SHALL provide access to WebSSH terminal functionality through a header button that opens in a new browser tab.

#### Scenario: Open WebSSH from header button
- **WHEN** a user clicks the SSH button in the header
- **THEN** a new browser tab opens with the WebSSH terminal page
- **AND** the current page remains open and unchanged
- **AND** the user can work with both pages simultaneously

#### Scenario: Direct navigation to WebSSH
- **WHEN** a user navigates directly to `/ssh/terminal` (e.g., from bookmark, direct URL)
- **THEN** the terminal page loads in the current tab
- **AND** all terminal functionality works as expected
- **Note**: Direct navigation still works for flexibility, only header button opens new tab