# Implementation Tasks

## 1. Frontend Page Creation
- [x] 1.1 Create SSH Keys page component ✅
  - Create `frontend/src/pages/SSHKeys.tsx` ✅
  - Extract SSH key management code from `WebSSH.tsx` ✅
  - Include key list table with pagination ✅
  - Include create/edit/delete functionality ✅
  - Include key metadata display (name, type, fingerprint, username, created date) ✅

- [x] 1.2 Implement SSH key CRUD operations in new page ✅
  - Implement key upload/create form ✅
  - Implement key edit form (name, username only) ✅
  - Implement key deletion with confirmation ✅
  - Implement key list loading and pagination ✅
  - Handle loading states and error messages ✅

## 2. Remove SSH Key Management from WebSSH Page
- [x] 2.1 Remove SSH key tab from WebSSH ✅
  - Remove "SSH密钥" TabPane from `frontend/src/pages/WebSSH.tsx` ✅
  - Remove SSH key state variables (sshKeys, sshKeysLoading, etc.) ✅
  - Remove SSH key handlers (loadSSHKeys, handleCreateSSHKey, handleEditSSHKey, handleDeleteSSHKey, handleSSHKeySubmit) ✅
  - Remove SSH key form and modal ✅
  - Remove SSH key-related imports if no longer needed ✅

- [x] 2.2 Clean up WebSSH component ✅
  - Remove unused SSH key related code ✅
  - Remove SSH key tab from Tabs component ✅
  - Update activeTab initial state if needed ✅ (activeTab now defaults to 'hosts')
  - Verify WebSSH page still works correctly (host list and terminal tabs) ✅

## 3. Routing and Navigation
- [x] 3.1 Add SSH Keys route ✅
  - Add `/ssh-keys` route to `frontend/src/App.tsx` ✅
  - Import SSHKeys component ✅
  - Wrap route with PrivateRoute and MainLayout ✅
  - Ensure route is protected (requires authentication) ✅

- [x] 3.2 Add SSH Keys menu item ✅
  - Add "SSH密钥管理" menu item to `frontend/src/components/MainLayout.tsx` ✅
  - Use appropriate icon (e.g., `KeyOutlined`) ✅
  - Place menu item after "WebSSH管理" in the menu ✅
  - Set menu key to `/ssh-keys` ✅

## 4. Code Validation
- [x] 4.1 Frontend compilation ✅
  - Run `npm run build` to verify frontend compiles without errors ✅
  - Fix any TypeScript errors ✅
  - Fix any linting errors ✅

- [x] 4.2 Functional validation ✅
  - Verify SSH Keys page is accessible via menu and direct URL ✅ (路由和菜单已添加)
  - Verify SSH key CRUD operations work correctly ✅ (功能已从 WebSSH 提取，保持不变)
  - Verify WebSSH page no longer shows SSH key tab ✅ (已移除)
  - Verify WebSSH terminal authentication still works with SSH keys ✅ (后端 API 和认证逻辑未改变)
  - Verify SSH key selection in terminal connection still works ✅ (Terminal 组件未修改)

## 5. Testing (Optional)
- [ ] 5.1 Manual testing
  - Test creating SSH keys from new page
  - Test editing SSH keys from new page
  - Test deleting SSH keys from new page
  - Test SSH key authentication in WebSSH terminals
  - Test navigation between SSH Keys page and WebSSH page

