## 1. Frontend Component Changes
- [x] 1.1 Remove terminal functionality from Hosts page
  - Remove "打开终端" button from operations column in `frontend/src/pages/Hosts.tsx`
  - Remove terminal tab management state and UI from `frontend/src/pages/Hosts.tsx`
  - Remove `Terminal` component import from `frontend/src/pages/Hosts.tsx`
  - Remove `handleOpenTerminal` and `handleCloseTerminalTab` functions
  - Remove `Tabs` component usage for terminal tabs
  - Simplify `Hosts.tsx` to focus only on host CRUD operations

## 2. Create WebSSH Management Page
- [x] 2.1 Create WebSSH page component
  - Create `frontend/src/pages/WebSSH.tsx`
  - Display host list with columns: hostname, IP address, SSH port, project, groups, tags, status
  - Add "打开终端" button for each host in operations column
  - Implement terminal tab management using Ant Design `Tabs` component
  - Support opening multiple terminal sessions for different hosts
  - Display terminal status indicators (connected, disconnected, connecting)

## 3. Frontend Routing and Menu
- [x] 3.1 Add WebSSH route
  - Add `/webssh` route to `frontend/src/App.tsx`
  - Import and use `WebSSH` component in route configuration
  - Ensure route is protected with `PrivateRoute`
- [x] 3.2 Add WebSSH menu item
  - Add WebSSH menu item to `frontend/src/components/MainLayout.tsx`
  - Use appropriate icon (e.g., `ConsoleSqlOutlined` or `TerminalOutlined`)
  - Place menu item after "主机管理" in menu order
  - Set menu key to `/webssh`

## 4. Testing
- [ ] 4.1 Frontend tests
  - Update Host management page tests to remove terminal-related tests
  - Add WebSSH page component tests
  - Test terminal opening from WebSSH page
  - Test multiple terminal tabs management
  - Test menu navigation to WebSSH page

## 5. Validation
- [x] 5.1 Code validation
  - Run `npm run build` to verify frontend compiles without errors
  - Fix any linting errors
- [x] 5.2 Functional validation
  - Verify "打开终端" button is removed from host management page
  - Verify WebSSH menu item appears in navigation
  - Verify WebSSH page displays host list correctly
  - Verify terminal opens from WebSSH page
  - Verify multiple terminal tabs work correctly
  - Verify terminal functionality works as expected
- [x] 5.3 OpenSpec validation
  - Run `openspec validate separate-webssh-management --strict`
  - Fix any validation errors

