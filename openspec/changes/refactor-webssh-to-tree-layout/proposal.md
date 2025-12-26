# Change: Refactor WebSSH to Tree Layout with Header Button

## Why
Currently, WebSSH management is accessed through a dedicated menu item and page. Users need a more convenient and intuitive way to access WebSSH terminal, especially when they need quick terminal access while working on other pages. The current table-based host list doesn't provide a clear hierarchical view of hosts organized by project and environment, which is important for large-scale deployments.

By moving WebSSH access to a header button and implementing a tree-based layout:
- **Quick Access**: Users can access WebSSH from any page via the header button
- **Better Organization**: Tree structure (Project → Environment → Hosts) provides clear hierarchical view
- **Improved UX**: Split-pane layout (tree on left, terminal on right) allows simultaneous browsing and terminal access
- **Space Efficiency**: Tree view is more compact than table view for large host lists

## What Changes
- **ADDED**: WebSSH access button in header (next to user avatar)
  - Add button/icon in `frontend/src/components/MainLayout.tsx` header section
  - Button opens WebSSH in a new independent page/modal
  - Button should be visible on all pages (when authenticated)

- **ADDED**: New WebSSH page with tree layout
  - Create new WebSSH page component with split-pane layout
  - Left pane: Tree view of hosts organized by Project → Environment → Hosts
  - Right pane: Terminal tabs component supporting multiple terminal sessions
  - Tree supports expand/collapse, search/filter, and host selection
  - Support opening multiple terminals for different hosts simultaneously
  - Each terminal tab supports context menu with: Clone, Close, Close All, Close Others

- **MODIFIED**: WebSSH page structure
  - Replace table-based host list with tree-based host list
  - Implement hierarchical grouping: Project (level 1) → Environment (level 2) → Hosts (level 3)
  - Maintain existing terminal functionality (authentication, resize, etc.)

- **ADDED**: Terminal tab management
  - Support multiple terminal tabs (one per host)
  - Each terminal tab displays host name and connection status
  - Terminal tabs support context menu (right-click) with actions:
    - Clone: Create a duplicate terminal connection to the same host
    - Close: Close the current terminal tab
    - Close All: Close all terminal tabs
    - Close Others: Close all terminal tabs except the current one
  - Terminal tabs can be reordered by dragging

- **KEPT**: All existing WebSSH functionality
  - Terminal component remains unchanged (reused for each tab)
  - WebSocket connection logic remains unchanged
  - Authentication flow remains unchanged
  - Backend APIs remain unchanged

## Impact
- **Affected specs**: `webssh-management` (modified - UI layout and access method)
- **Affected code**:
  - `frontend/src/components/MainLayout.tsx` - Add WebSSH button in header
  - `frontend/src/pages/WebSSH.tsx` - Refactor to tree layout with split-pane
  - `frontend/src/App.tsx` - May need route adjustment (if using modal/drawer)
- **User experience**: 
  - WebSSH becomes more accessible via header button
  - Tree view provides better organization for large host lists
  - Split-pane layout improves workflow efficiency

