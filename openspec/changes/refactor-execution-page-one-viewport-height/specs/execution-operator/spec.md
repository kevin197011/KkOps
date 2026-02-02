# Execution Operator – Layout Delta (One-Viewport Height, Width Unchanged)

## MODIFIED Requirements

### Requirement: Execution Page Layout
The execution page SHALL use a single-page form layout instead of dual-column workflow list. The page SHALL keep the original content width (single column, centered, maxWidth e.g. 1200) and SHALL fit within one viewport height on typical desktop resolutions. The layout of each functional block SHALL be re-planned so that configuration, execute action, and results are visible within one screen height without excessive page scroll.

#### Scenario: Single-column, width unchanged
- **WHEN** a user navigates to `/executions`
- **THEN** the page SHALL display a single-column layout with:
  - Content area width unchanged (e.g. maxWidth 1200, margin 0 auto)
  - No left/right split; configuration and results SHALL be arranged within the same column width
  - Root container SHALL use viewport height (e.g. minHeight calc(100vh - 120px))

#### Scenario: Re-planned block layout
- **WHEN** a user views the execution page
- **THEN** the functional blocks SHALL be arranged as follows (or equivalent compact arrangement):
  - First row: execution mode selector and execution options side by side (same row)
  - Second row: script editor and host selector side by side (same row), each with maxHeight and internal scroll
  - Third row: execute button (and optional link to task history)
  - Fourth row: execution results area, occupying remaining height with internal scroll
- **AND** on viewport width below a breakpoint (e.g. 992px), script editor and host selector MAY stack vertically (one column each) while keeping width unchanged

#### Scenario: One viewport height
- **WHEN** a user views the execution page on a typical desktop viewport (e.g. 1440×900 or 1920×1080)
- **THEN** the entire page SHALL fit within one viewport height without requiring vertical page scroll to see both configuration and results
- **AND** long content in script area, host list, or results area SHALL use internal scroll (maxHeight + overflow: auto) so that the page itself does not grow beyond one viewport height

#### Scenario: Constrained scroll regions
- **WHEN** the script content or host list or results list is long
- **THEN** the script editor region SHALL have a maximum height and SHALL show internal vertical scroll (overflow: auto)
- **AND** the host selector list SHALL have a maximum height and SHALL show internal vertical scroll
- **AND** the execution results area SHALL use internal scroll so that the overall page layout remains within one viewport height
