# Tasks: Redesign WebSSH Terminal

## 1. MainLayout Updates

- [x] 1.1 Remove WebSSH menu item
  - Remove WebSSH entry from `menuItems` array in `MainLayout.tsx`
  - Remove `/ssh/terminal` key from menu items

- [x] 1.2 Add WebSSH button to header
  - Import terminal icon (`ConsoleSqlOutlined` or similar)
  - Add button in header, between theme toggle and user menu
  - Set onClick to navigate to `/ssh/terminal`
  - Add tooltip: "WebSSH 终端"
  - Style button to match other header buttons

## 2. WebSSHTerminal Component Redesign

- [x] 2.1 Remove MainLayout wrapper
  - Remove `MainLayout` import
  - Remove `<MainLayout>` wrapper from component return

- [x] 2.2 Add standalone page header
  - Create minimal header with back button and title
  - Back button navigates to previous page (or `/assets` if no history)
  - Title: "WebSSH 终端"
  - Use Ant Design Layout Header component

- [x] 2.3 Implement split-view layout
  - Use Ant Design Layout with Sider and Content
  - Left Sider: 350px width, fixed
  - Right Content: Flexible width, remaining space
  - Full height layout (100vh minus header)

- [x] 2.4 Create asset list sidebar
  - Add search input at top of sidebar
  - Implement scrollable asset list
  - Display asset name, IP, and SSH user for each item
  - Show connection status indicator (checkmark/badge for connected)
  - Highlight currently connected asset
  - Use Ant Design List or Menu component

- [x] 2.5 Remove connection modal
  - Remove `connectModalVisible` state
  - Remove `Modal` component for connection
  - Remove connection form from modal
  - Remove modal-related handlers

- [x] 2.6 Implement direct asset connection
  - Add click handler on asset list items
  - Connect immediately on click (no modal)
  - Handle disconnect before connecting to new asset
  - Show loading state during connection
  - Update asset list to reflect connection status

- [x] 2.7 Update terminal area layout
  - Move terminal to right panel
  - Add current asset info display above terminal
  - Add connection controls (disconnect button)
  - Keep terminal resize functionality

## 3. App.tsx Route Update

- [x] 3.1 Remove MainLayout from route
  - Update `/ssh/terminal` route to not use MainLayout wrapper
  - Keep ProtectedRoute wrapper
  - Ensure route works correctly

## 4. State Management Updates

- [x] 4.1 Refactor connection state
  - Track currently connected asset ID
  - Update state when asset clicked
  - Handle disconnect/switch logic
  - Update asset list highlighting based on connection status

- [x] 4.2 Add asset search/filter state
  - Add search query state
  - Filter assets based on search input
  - Update asset list display

## 5. UI/UX Improvements

- [x] 5.1 Style asset list
  - Add hover effects on asset items
  - Style active/connected asset differently
  - Add spacing and padding
  - Ensure list is scrollable

- [x] 5.2 Improve terminal area
  - Add header showing current asset info
  - Style disconnect button
  - Ensure terminal fills available space
  - Handle terminal resize on window resize

- [x] 5.3 Add loading and error states
  - Show loading indicator when connecting
  - Display error messages if connection fails
  - Handle connection state transitions

## 6. Responsive Design

- [ ] 6.1 Make layout responsive
  - On mobile: Make left panel collapsible/drawer
  - Ensure terminal remains usable on small screens
  - Test on different screen sizes
  - Adjust breakpoints as needed

## 7. Testing

- [ ] 7.1 Test header button
  - Verify button appears in header
  - Click button and verify navigation to terminal page
  - Verify button styling matches other header buttons

- [ ] 7.2 Test asset list
  - Verify assets are displayed correctly
  - Test search/filter functionality
  - Verify connection status indicators
  - Test clicking assets to connect

- [ ] 7.3 Test terminal connection
  - Verify connection works from asset list click
  - Test switching between assets
  - Verify disconnect functionality
  - Test terminal input/output

- [ ] 7.4 Test navigation
  - Verify back button returns to previous page
  - Test browser back button
  - Verify page works as standalone (no MainLayout)

- [ ] 7.5 Test responsive layout
  - Test on desktop (wide screen)
  - Test on tablet
  - Test on mobile
  - Verify layout adapts correctly

## 8. Code Cleanup

- [x] 8.1 Remove unused code
  - Remove modal-related imports
  - Remove unused state variables
  - Clean up unused functions

- [x] 8.2 Verify imports
  - Ensure all new imports are correct
  - Remove unused imports
  - Verify component dependencies

## 9. Validation

- [ ] 9.1 Run linter
  - Fix any linting errors
  - Fix any TypeScript errors

- [ ] 9.2 Build and test
  - Run frontend build
  - Verify no compilation errors
  - Test in browser

- [ ] 9.3 Validate OpenSpec proposal
  - `openspec validate redesign-webssh-terminal --strict`