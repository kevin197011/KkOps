# Execution Operator – Column Order Swap

## MODIFIED Requirements

### Requirement: Execution Page Layout
The execution page SHALL use a single-page form layout with original content width and one-viewport height. Within each row, the **left/right column order** SHALL be set so that the left side shows execution options and host selector respectively, and the right side shows mode selector and script editor respectively, to better match usage habits.

#### Scenario: Re-planned block layout (column order)
- **WHEN** a user views the execution page
- **THEN** the functional blocks SHALL be arranged as follows:
  - First row: **left** = execution options, **right** = execution mode selector (same row)
  - Second row: **left** = host selector, **right** = script editor (same row), each with maxHeight and internal scroll as needed
  - Third row: execute button (and optional link to task history)
  - Fourth row: execution results area, occupying remaining height with internal scroll
- **AND** on viewport width below a breakpoint (e.g. 992px), the two columns in each row MAY stack vertically (left column first, then right column) while keeping width unchanged
