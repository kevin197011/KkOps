# Execution Operator – Layout Delta (One-Screen)

## MODIFIED Requirements

### Requirement: Execution Page Layout
The execution page SHALL use a single-page form layout instead of dual-column workflow list, and SHALL fit the main configuration and execution results within one viewport on typical desktop resolutions so that operations are convenient without excessive scrolling.

#### Scenario: Single-page execution form
- **WHEN** a user navigates to `/executions`
- **THEN** the page SHALL display a single-page form layout with:
  - Execution mode selector at the top
  - Script editor in the middle
  - Host selector below script editor
  - Execution options below host selector
  - Execute button at the bottom
  - Results area (hidden until execution starts or shown as placeholder)
- **AND** SHALL NOT display a workflow list on the left
- **AND** SHALL provide a link to view saved tasks in `/tasks` page

#### Scenario: One-screen layout on typical viewport
- **WHEN** a user views the execution page on a typical desktop viewport (e.g. 1440×900 or 1920×1080)
- **THEN** the main configuration area (mode, script, hosts, options, execute button) and the execution results area SHALL be visible within the same viewport without requiring the user to scroll the page to see both
- **AND** the layout SHALL use a split (e.g. left configuration, right results) or a constrained vertical layout so that after execution starts, results appear alongside or below the configuration without pushing it off-screen

#### Scenario: Constrained scroll regions
- **WHEN** the script content or host list is long
- **THEN** the script editor region SHALL have a maximum height and SHALL show internal vertical scroll (overflow) so that the page itself does not grow unbounded
- **AND** the host selector list SHALL have a maximum height and SHALL show internal vertical scroll so that the page itself does not grow unbounded
- **AND** the execution results area SHALL use internal scroll for its content so that the overall page layout remains within one viewport

#### Scenario: Responsive one-screen behavior
- **WHEN** the viewport width is below a defined breakpoint (e.g. 992px)
- **THEN** the layout MAY switch to a vertical stack (configuration above, results below)
- **AND** the configuration and results areas SHALL still use height constraints and internal scroll so that the key operations and results remain visible within one screen with minimal page scroll
