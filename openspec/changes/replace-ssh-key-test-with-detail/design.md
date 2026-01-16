# Design: Replace SSH Key Test with Detail View

## Context

The current SSH key management interface has a "测试" (Test) button that opens a form to test SSH connections. Users want to replace this with a "详情" (Detail) button that shows complete key information in a read-only view.

## Goals

1. Replace test button with detail button in SSH key list
2. Display complete SSH key information in a structured, readable format
3. Maintain consistency with existing detail view patterns (e.g., asset detail)
4. Provide easy access to key information without additional user input

## Non-Goals

- Removing the backend test API endpoint (kept for potential future use)
- Adding editing capabilities in the detail view
- Adding new key information that doesn't already exist

## Decisions

### Decision: Use Descriptions Component for Detail View

**What**: Display SSH key details using Ant Design's `Descriptions` component in a modal.

**Why**: 
- Consistent with asset detail view pattern
- Provides structured, readable layout
- Easy to implement and maintain
- Supports multi-column layout for better information density

**Implementation**:
- Use `Modal` component with `Descriptions` inside
- Two-column layout for most fields
- Full-width for long content (public key, description)

### Decision: Remove Test Functionality from List

**What**: Remove the test button, test modal, and related handlers from the SSH key list.

**Why**:
- Simplifies the UI and code
- Test functionality can be re-added later if needed
- Focuses the list on key management operations (view, edit, delete)

**Alternatives considered**:
- Keep both test and detail buttons: Rejected - too many buttons in operation column
- Move test to detail view: Rejected - detail view should be read-only

### Decision: Icon Selection for Detail Button

**What**: Use `EyeOutlined` icon for the detail button.

**Why**:
- Clearly indicates "view" action
- Standard icon for detail/view operations
- Different from edit (`EditOutlined`) and delete icons

## Implementation Details

### Frontend Changes

#### Component: `frontend/src/pages/ssh/SSHKeyList.tsx`

**Remove**:
- `testModalVisible` state
- `testingKey` state
- `testForm` form instance
- `handleTest` function
- `handleTestSubmit` function
- Test modal JSX
- `TestSSHKeyRequest` import (if not used elsewhere)

**Add**:
- `detailModalVisible` state
- `detailKey` state (to store the key being viewed)
- `handleViewDetail` function (opens detail modal)
- Detail modal JSX using `Descriptions` component
- Import `EyeOutlined` icon (or `InfoCircleOutlined`)
- Import `Descriptions` component from antd

**Modify**:
- Operation column button: Change "测试" to "详情"
- Change icon from `CheckCircleOutlined` to `EyeOutlined`
- Change button onClick from `handleTest` to `handleViewDetail`

**Detail Modal Structure**:
```tsx
<Modal
  title="SSH 密钥详情"
  open={detailModalVisible}
  onCancel={() => {
    setDetailModalVisible(false)
    setDetailKey(null)
  }}
  footer={[
    <Button key="close" onClick={() => {
      setDetailModalVisible(false)
      setDetailKey(null)
    }}>
      关闭
    </Button>,
  ]}
  width="90%"
  style={{ maxWidth: 800 }}
>
  {detailKey ? (
    <Descriptions column={2} bordered>
      <Descriptions.Item label="ID">{detailKey.id}</Descriptions.Item>
      <Descriptions.Item label="名称">{detailKey.name}</Descriptions.Item>
      <Descriptions.Item label="类型">
        <Tag>{detailKey.type || 'unknown'}</Tag>
      </Descriptions.Item>
      <Descriptions.Item label="SSH 用户">{detailKey.ssh_user || '-'}</Descriptions.Item>
      <Descriptions.Item label="指纹" span={2}>
        {detailKey.fingerprint || '-'}
      </Descriptions.Item>
      <Descriptions.Item label="公钥" span={2}>
        <div style={{ wordBreak: 'break-all' }}>
          {detailKey.public_key || '-'}
        </div>
        {detailKey.public_key && (
          <Button 
            size="small" 
            onClick={() => {
              navigator.clipboard.writeText(detailKey.public_key)
              message.success('公钥已复制到剪贴板')
            }}
          >
            复制
          </Button>
        )}
      </Descriptions.Item>
      <Descriptions.Item label="描述" span={2}>
        {detailKey.description || '-'}
      </Descriptions.Item>
      <Descriptions.Item label="最后使用时间">
        {detailKey.last_used_at ? new Date(detailKey.last_used_at).toLocaleString() : '-'}
      </Descriptions.Item>
      <Descriptions.Item label="创建时间">
        {new Date(detailKey.created_at).toLocaleString()}
      </Descriptions.Item>
      <Descriptions.Item label="更新时间">
        {new Date(detailKey.updated_at).toLocaleString()}
      </Descriptions.Item>
    </Descriptions>
  ) : null}
</Modal>
```

## Risks / Trade-offs

### Risk: Loss of Test Functionality

**Risk**: Users may need to test SSH keys but no longer have easy access to this feature.

**Mitigation**: 
- Test functionality is still available via API (can be re-added later)
- Consider adding test functionality in asset connection or another context where it's more relevant

### Trade-off: Simplicity vs Functionality

**Trade-off**: Removing test reduces functionality but simplifies the UI.

**Decision**: Prioritize simplicity and consistency. Detail view is more commonly needed than test functionality.

## Migration Plan

1. Update `SSHKeyList.tsx` component:
   - Remove test-related code
   - Add detail view modal
   - Update button in operation column
2. Test detail view displays all key information correctly
3. Verify public key copy functionality works
4. Ensure modal is responsive and looks good on different screen sizes

## Open Questions

- Should we add a "复制公钥" button in the detail view?
  - **Answer**: Yes, it's a common use case and improves UX

- Should we format dates in a specific locale?
  - **Answer**: Use `toLocaleString()` for user-friendly date formatting