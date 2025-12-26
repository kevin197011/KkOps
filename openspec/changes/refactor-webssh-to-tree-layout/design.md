# Design: WebSSH Tree Layout with Header Access

## Context
The current WebSSH management page uses a table-based host list with filters. Users navigate to WebSSH via a menu item, and the page displays hosts in a table format. This proposal refactors WebSSH to:
1. Add a header button for quick access
2. Replace table with tree view organized by Project → Environment → Hosts
3. Implement split-pane layout (tree left, terminal right)

## Goals
- Provide quick access to WebSSH from any page via header button
- Organize hosts hierarchically (Project → Environment → Hosts) for better navigation
- Enable simultaneous host browsing and terminal access with split-pane layout
- Maintain all existing WebSSH terminal functionality

## Non-Goals
- Changing WebSSH backend API or WebSocket protocol
- Modifying terminal component functionality
- Changing authentication flow
- Adding new terminal features
- Removing existing WebSSH menu item (can coexist)

## Decisions

### Decision: Header Button Placement
**What**: Add WebSSH button/icon in the header, positioned to the left of the user avatar dropdown.

**Why**: 
- Provides quick access from any page
- Follows common UI patterns (many apps have terminal/console buttons in header)
- Doesn't interfere with existing menu navigation
- Visible and accessible at all times

**Alternatives considered**:
- Floating action button (FAB) - **Rejected**: May overlap content, less discoverable
- Keep only in menu - **Rejected**: Less convenient, requires navigation
- Both menu and header button - **Accepted**: Provides flexibility, users can choose preferred access method

### Decision: Tree Structure Organization
**What**: Organize hosts in a tree with three levels:
- Level 1: Project (grouping by project name)
- Level 2: Environment (grouping by environment within project: dev, test, prod, etc.)
- Level 3: Hosts (individual host nodes)

**Why**:
- Matches common organizational structure (projects contain environments, environments contain hosts)
- Provides clear hierarchy for large host lists
- Enables quick navigation to specific project/environment
- Supports filtering and search at each level

**Alternatives considered**:
- Flat list with filters - **Rejected**: Doesn't scale well, less organized
- Project → Hosts (skip environment) - **Rejected**: Loses important grouping dimension
- Environment → Project → Hosts - **Rejected**: Less intuitive, projects are primary organizational unit

### Decision: Split-Pane Layout
**What**: Use split-pane layout with:
- Left pane: Tree view (resizable, default ~30% width)
- Right pane: Terminal component (default ~70% width)
- Resizable divider between panes

**Why**:
- Allows simultaneous host browsing and terminal access
- Efficient use of screen space
- Common pattern in terminal/IDE applications
- Better workflow (select host → terminal appears immediately)

**Alternatives considered**:
- Modal terminal - **Rejected**: Blocks tree view, less efficient
- Tab-based terminals - **Rejected**: More complex, takes more space
- Full-screen terminal with tree overlay - **Rejected**: Less intuitive, harder to navigate

### Decision: Tree Component
**What**: Use Ant Design `Tree` component with custom node rendering.

**Why**:
- Ant Design Tree provides built-in expand/collapse, selection, search
- Consistent with existing UI library
- Supports custom node rendering (icons, badges, etc.)
- Good performance for large trees

**Alternatives considered**:
- Custom tree component - **Rejected**: Unnecessary complexity, reinventing wheel
- Third-party tree library - **Rejected**: Adds dependency, Ant Design Tree is sufficient

### Decision: Multiple Terminal Tabs Support
**What**: Support opening multiple terminal tabs simultaneously:
- Each host selection opens a new terminal tab (or switches to existing tab if already open for that host)
- Terminal tabs are displayed in the right pane using Ant Design `Tabs` component
- Each tab shows host name and connection status indicator
- Tabs can be closed individually or via context menu

**Why**:
- Users often need to work with multiple hosts simultaneously
- Tabs provide clear organization for multiple terminal sessions
- Common pattern in terminal/IDE applications
- Allows parallel operations across different hosts

**Alternatives considered**:
- Single terminal view - **Rejected**: Too limiting, users need multiple terminals
- Modal terminals - **Rejected**: Blocks tree view, less efficient

### Decision: Terminal Tab Context Menu
**What**: Each terminal tab supports right-click context menu with actions:
- **Clone**: Create a duplicate terminal connection to the same host (new tab, independent connection)
- **Close**: Close the current terminal tab
- **Close All**: Close all terminal tabs
- **Close Others**: Close all terminal tabs except the current one

**Why**:
- Provides quick access to common terminal management operations
- Follows common UI patterns (browsers, IDEs use similar context menus)
- Improves workflow efficiency

**Alternatives considered**:
- Button-based actions only - **Rejected**: Takes more space, less discoverable
- Dropdown menu in tab - **Rejected**: Context menu is more standard for tabs

### Decision: Host Selection Behavior
**What**: When user clicks a host node in tree:
- If no terminal tab exists for that host: Open new terminal tab for selected host
- If terminal tab already exists for that host: Switch to existing tab (don't create duplicate)
- If user wants duplicate: Use "Clone" context menu action on existing tab

**Why**:
- Prevents accidental duplicate tabs for same host
- Allows intentional duplication via Clone action
- Maintains clean tab organization

**Alternatives considered**:
- Always create new tab - **Rejected**: Can lead to many duplicate tabs for same host
- Always switch to existing - **Rejected**: Users may want multiple connections to same host

## Implementation Details

### Tree Node Structure
```
Project: devops
  ├─ Environment: dev
  │   ├─ Host: salt-master (192.168.56.10)
  │   └─ Host: salt-minion (192.168.56.11)
  └─ Environment: prod
      └─ Host: prod-server-01 (10.0.1.10)
```

### Tree Node Rendering
- Project nodes: Show project name, host count badge
- Environment nodes: Show environment name (dev/test/prod), host count badge
- Host nodes: Show hostname, IP address, status indicator (online/offline)

### State Management
- Selected host ID: Track which host is currently selected in tree
- Active terminal tab: Track which terminal tab is currently active
- Terminal tabs: Array of terminal tab objects, each containing:
  - Tab ID (unique identifier)
  - Host ID
  - Host name
  - Connection status
  - Terminal component instance reference
- Tree expanded keys: Track which nodes are expanded
- Terminal connection state: Track terminal connection status per tab

### Search/Filter
- Search input above tree
- Filter tree nodes by hostname, IP, project, environment
- Highlight matching nodes

## Risks / Trade-offs
- **Risk**: Tree view may be slower for very large host lists (1000+ hosts)
  - **Mitigation**: Implement virtual scrolling or lazy loading for tree nodes
- **Risk**: Split-pane may be cramped on small screens
  - **Mitigation**: Make panes responsive, allow full-screen terminal mode
- **Risk**: Users may prefer table view for some use cases
  - **Mitigation**: Keep existing menu item, allow both access methods

## Migration Plan
1. Add header button (non-breaking, additive change)
2. Refactor WebSSH page to tree layout (existing route still works)
3. Keep menu item as alternative access method
4. No data migration needed (uses existing host data structure)

## Open Questions
- Should tree remember expanded state across sessions? (Nice to have, can add later)
- Should we limit the maximum number of terminal tabs? (Consider performance implications)

