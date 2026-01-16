# Change: Replace SSH Key Test with Detail View

## Why

Currently, the SSH key management list has a "测试" (Test) button that opens a modal to test SSH connection with the key. However, this test functionality requires additional input (host, port, username) which may not be intuitive for users who simply want to view the key details.

Replacing the test button with a "详情" (Detail) button will:
- Provide a more intuitive way to view complete SSH key information
- Align with common UI patterns (similar to asset detail view)
- Remove the need for additional inputs when users just want to see key information
- Improve user experience by showing all key properties in a structured format

The test functionality is still useful but may be better suited as a separate action or feature elsewhere in the system.

## What Changes

- **MODIFIED**: Replace "测试" button with "详情" button in SSH key list operations
- **MODIFIED**: Change button icon from `CheckCircleOutlined` to `EyeOutlined` (or similar detail icon)
- **MODIFIED**: Replace test connection modal with detail view modal
- **ADDED**: Detail view modal displaying all SSH key information in a structured format (using `Descriptions` component)
- **REMOVED**: Test connection modal and related form
- **REMOVED**: `handleTest` and `handleTestSubmit` functions
- **MAINTAINED**: Backend test API endpoint remains available (for potential future use)

## Impact

- **UI Changes**: Operation column in SSH key list will show "详情" instead of "测试"
- **User Experience**: Users can easily view complete SSH key information without additional input
- **Code Simplification**: Removes test form and related state management
- **Consistency**: Aligns with asset detail view pattern

## Detail View Content

The detail modal should display:
- ID
- 名称 (Name)
- 类型 (Type) - with tag styling
- SSH 用户 (SSH User)
- 指纹 (Fingerprint)
- 公钥 (Public Key) - full content with copy option
- 描述 (Description)
- 最后使用时间 (Last Used At) - if available
- 创建时间 (Created At)
- 更新时间 (Updated At)

## Considerations

- Use `Descriptions` component (similar to asset detail) for consistent UI
- Public key should be displayed in full (possibly with copy button)
- Modal should be read-only (no editing in detail view)
- Consider adding a "复制公钥" (Copy Public Key) button for convenience
- The test functionality can be re-added later if needed through a different UI pattern