# Tasks: Replace SSH Key Test with Detail View

## 1. Frontend Implementation

- [x] 1.1 Remove test-related code from SSHKeyList.tsx
  - Remove `testModalVisible` state
  - Remove `testingKey` state
  - Remove `testForm` form instance
  - Remove `handleTest` function
  - Remove `handleTestSubmit` function
  - Remove `TestSSHKeyRequest` import (if not used elsewhere)
  - Remove test modal JSX

- [x] 1.2 Add detail view state and handler
  - Add `detailModalVisible` state
  - Add `detailKey` state (type: `SSHKey | null`)
  - Create `handleViewDetail` function to open detail modal

- [x] 1.3 Update operation column button
  - Change button text from "测试" to "详情"
  - Change icon from `CheckCircleOutlined` to `EyeOutlined`
  - Update onClick handler from `handleTest` to `handleViewDetail`
  - Update button import if needed

- [x] 1.4 Create detail view modal
  - Add Modal component with appropriate title ("SSH 密钥详情")
  - Import `Descriptions` component from antd
  - Structure detail view with Descriptions component (2 columns, bordered)
  - Display all SSH key fields:
    - ID
    - 名称 (Name)
    - 类型 (Type) - with Tag component
    - SSH 用户 (SSH User)
    - 指纹 (Fingerprint) - span 2 columns
    - 公钥 (Public Key) - span 2 columns, with copy button
    - 描述 (Description) - span 2 columns
    - 最后使用时间 (Last Used At)
    - 创建时间 (Created At)
    - 更新时间 (Updated At)
  - Format dates using `toLocaleString()`

- [x] 1.5 Add public key copy functionality
  - Add copy button next to public key field
  - Implement copy to clipboard using `navigator.clipboard.writeText()`
  - Show success message when copied
  - Handle cases where clipboard API is not available

## 2. UI/UX Improvements

- [x] 2.1 Style public key field
  - Use `wordBreak: 'break-all'` for long public key strings
  - Ensure public key is fully visible and readable

- [x] 2.2 Ensure modal is responsive
  - Test on different screen sizes
  - Ensure modal width is appropriate (90% max, 800px max width)
  - Verify two-column layout works well on all screen sizes

- [x] 2.3 Add appropriate spacing and formatting
  - Ensure consistent spacing between fields
  - Format empty values as '-' consistently
  - Ensure date formatting is user-friendly

## 3. Testing

- [ ] 3.1 Test detail view displays correctly
  - Verify all fields are displayed
  - Verify empty/null fields show '-' appropriately
  - Verify type field shows tag correctly

- [ ] 3.2 Test public key copy functionality
  - Click copy button and verify key is copied to clipboard
  - Verify success message appears
  - Test with different public key formats

- [ ] 3.3 Test date formatting
  - Verify dates are formatted correctly
  - Test with different locales if applicable
  - Verify null dates (last_used_at) show '-' correctly

- [ ] 3.4 Test modal behavior
  - Open detail modal and verify it displays correctly
  - Close modal and verify state is reset
  - Verify modal can be opened for different keys

- [ ] 3.5 Test responsive design
  - Test on mobile devices
  - Test on tablets
  - Test on desktop screens
  - Verify layout adapts appropriately

## 4. Code Cleanup

- [x] 4.1 Remove unused imports
  - Remove `TestSSHKeyRequest` from imports if no longer needed
  - Remove any other unused imports related to test functionality

- [x] 4.2 Verify no test-related code remains
  - Search for any remaining references to test functionality
  - Clean up any unused variables or functions

## 5. Validation

- [ ] 5.1 Run frontend build
  - `npm run build` (or equivalent)
  - Verify no compilation errors

- [ ] 5.2 Run linter
  - Verify no linting errors
  - Fix any warnings if applicable

- [ ] 5.3 Manual testing
  - Navigate to SSH key management page
  - Click "详情" button on a key
  - Verify detail modal displays all information correctly
  - Test copy public key functionality
  - Close modal and verify it works correctly
  - Test with keys that have missing optional fields

- [ ] 5.4 Validate OpenSpec proposal
  - `openspec validate replace-ssh-key-test-with-detail --strict`