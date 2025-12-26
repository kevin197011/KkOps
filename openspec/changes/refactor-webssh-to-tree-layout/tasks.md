## 1. Header Button Implementation
- [x] 1.1 Add WebSSH button to MainLayout header
  - Add button/icon component in `frontend/src/components/MainLayout.tsx` header section
  - Position button to the left of user avatar dropdown
  - Use appropriate icon (e.g., `TerminalOutlined` or `ConsoleSqlOutlined`)
  - Add click handler to navigate to WebSSH page or open in new window/modal
  - Ensure button is visible on all authenticated pages

## 2. Tree Data Structure
- [x] 2.1 Create tree data transformation function
  - Create utility function to transform flat host list to tree structure
  - Group hosts by project, then by environment
  - Handle hosts without project or environment gracefully
  - Support empty states (no projects, no hosts, etc.)

- [x] 2.2 Add tree node types and interfaces
  - Define TypeScript interfaces for tree nodes (ProjectNode, EnvironmentNode, HostNode)
  - Include metadata (host count, status, etc.) in node data
  - Support custom node properties for rendering

## 3. Tree Component Implementation
- [x] 3.1 Implement tree view component
  - Use Ant Design `Tree` component in `frontend/src/pages/WebSSH.tsx`
  - Configure tree with custom node rendering
  - Implement expand/collapse functionality
  - Add host selection handler

- [x] 3.2 Customize tree node rendering
  - Project nodes: Display project name with host count badge
  - Environment nodes: Display environment name with host count badge and color coding
  - Host nodes: Display hostname, IP address, and status indicator (online/offline icon)
  - Add appropriate icons for each node type

- [x] 3.3 Implement tree search/filter
  - Add search input above tree
  - Filter tree nodes by hostname, IP, project name, environment
  - Highlight matching nodes
  - Support case-insensitive search

## 4. Split-Pane Layout
- [x] 4.1 Implement split-pane layout
  - Use Ant Design `Layout` with `Sider` and `Content` or custom split-pane component
  - Left pane: Tree view (default ~30% width, min 200px)
  - Right pane: Terminal tabs container (default ~70% width)
  - Make divider resizable (use `Resizable` component or custom implementation)

- [x] 4.2 Implement terminal tabs container
  - Right pane contains Ant Design `Tabs` component
  - Each tab contains a Terminal component instance
  - Support tab switching, closing, and reordering
  - Display empty state when no terminals are open

- [x] 4.3 Handle responsive layout
  - On small screens (< 768px): Stack panes vertically or use drawer for tree
  - Ensure terminal tabs remain usable on mobile devices
  - Test layout on various screen sizes

## 5. Host Selection and Terminal Tab Management
- [x] 5.1 Implement host selection logic
  - Track selected host ID in component state
  - When host node clicked: check if terminal tab exists for that host
  - If tab exists: switch to existing tab
  - If tab doesn't exist: create new terminal tab
  - Highlight selected host node in tree

- [x] 5.2 Implement terminal tab management
  - Use Ant Design `Tabs` component for terminal tabs
  - Track terminal tabs in state (array of tab objects)
  - Each tab contains: tab ID, host ID, host name, connection status
  - Support tab switching, closing, and reordering
  - Display host name and connection status in tab title

- [x] 5.3 Integrate terminal component per tab
  - Create Terminal component instance for each tab
  - Pass host ID to Terminal component
  - Handle terminal connection/disconnection per tab
  - Update tab status based on terminal connection state
  - Clean up terminal component when tab is closed

- [x] 5.4 Implement terminal tab context menu
  - Add right-click handler on terminal tabs (use `onContextMenu` prop on Tabs)
  - Display context menu with actions: Clone, Close, Close All, Close Others
  - Use Ant Design `Dropdown` with `Menu` component for context menu
  - Position context menu at mouse click position
  - Implement Clone action: create new tab with same host ID (generate unique tab ID)
  - Implement Close action: close current tab and cleanup terminal connection
  - Implement Close All action: close all tabs and cleanup all terminal connections
  - Implement Close Others action: close all tabs except current one
  - Handle edge cases (only one tab, no tabs, etc.)
  - Add visual feedback (disable actions when not applicable)

## 6. Data Loading and State Management
- [x] 6.1 Load host data for tree
  - Load all hosts (or paginated if needed) on page mount
  - Transform host data to tree structure
  - Handle loading states and errors
  - Refresh tree data when hosts are updated

- [x] 6.2 Manage tree state
  - Track expanded tree keys (which nodes are expanded)
  - Persist expanded state in component state (optional: localStorage)
  - Handle tree node selection state
  - Manage search/filter state

## 7. UI/UX Enhancements
- [x] 7.1 Add tree node tooltips
  - Show host details (IP, status, project, environment) on hover
  - Display full hostname if truncated
  - Show connection status for hosts

- [x] 7.2 Add empty states
  - Show message when no hosts available
  - Show message when no hosts match search filter
  - Provide action buttons (e.g., "Add Host" link)

- [x] 7.3 Add loading indicators
  - Show loading spinner while loading hosts
  - Show loading state for terminal connection
  - Display connection status in tree nodes

## 8. Remove/Update Existing Features
- [ ] 8.1 Remove table-based host list
  - Remove table component and related code from WebSSH page
  - Remove table filters (project, group, tag filters) - tree provides better filtering
  - Keep host loading logic but adapt for tree structure

- [ ] 8.2 Update terminal tab management
  - Remove terminal tabs (single terminal view)
  - Simplify terminal state management
  - Remove tab switching logic

## 9. Testing
- [ ] 9.1 Component tests
  - Test tree rendering with various host configurations
  - Test host selection and terminal connection
  - Test tree search/filter functionality
  - Test split-pane resizing

- [ ] 9.2 Integration tests
  - Test header button navigation
  - Test tree → terminal connection flow
  - Test host switching in terminal
  - Test responsive layout behavior

- [ ] 9.3 Manual testing
  - Test with large host lists (100+ hosts)
  - Test with hosts in multiple projects/environments
  - Test with hosts missing project or environment
  - Test on different screen sizes

## 10. Validation
- [x] 10.1 Code validation
  - Run `npm run build` to verify frontend compiles without errors
  - Fix any linting errors
  - Ensure TypeScript types are correct

- [ ] 10.2 Functional validation
  - Verify header button appears on all pages
  - Verify tree displays correctly with real host data
  - Verify terminal connects when host is selected
  - Verify split-pane resizing works
  - Verify search/filter works correctly

