# Design: Redesign WebSSH Terminal

## Context

The current WebSSH terminal is accessed via sidebar menu and uses a modal-based connection interface. Users want to move it to a header button and redesign it as a split-view page with asset list on the left and terminal on the right.

## Goals

1. Move WebSSH access from menu to header button
2. Redesign WebSSH page as standalone (no MainLayout)
3. Implement split-view layout (assets left, terminal right)
4. Allow direct asset selection without modal
5. Support quick asset switching

## Non-Goals

- Changing WebSocket SSH connection protocol
- Modifying terminal functionality (xterm.js)
- Adding new terminal features beyond layout changes
- Changing backend SSH handling

## Decisions

### Decision: Header Button Placement

**What**: Place WebSSH button in header, between theme toggle and user menu.

**Why**:
- Always visible and accessible
- Doesn't clutter sidebar menu
- Standard location for quick-access features

**Implementation**:
- Use `ConsoleSqlOutlined` or similar terminal icon
- Position in header right section
- Tooltip: "WebSSH 终端"
- Click opens `/ssh/terminal` in new context (standalone)

### Decision: Standalone Page (No MainLayout)

**What**: WebSSH page should not use MainLayout wrapper.

**Why**:
- Full-screen terminal experience
- More space for split-view layout
- Focused interface for terminal operations

**Implementation**:
- Remove `<MainLayout>` wrapper from route
- Add minimal header with back button and title
- Full viewport height layout

### Decision: Split-View Layout

**What**: Left sidebar (assets) + Right panel (terminal).

**Layout**:
- Left panel: 300-400px width, fixed
- Right panel: Remaining width, flexible
- Both panels: Full height (viewport height minus header)

**Why**:
- Provides asset overview while using terminal
- Easy asset switching
- Similar to modern SSH clients (like VS Code Remote)

### Decision: Direct Asset Connection (No Modal)

**What**: Remove connection modal, connect directly by clicking asset in list.

**Why**:
- Faster workflow
- Better UX (no modal interruption)
- More intuitive

**Implementation**:
- Click asset in list → connect immediately
- Show connection status in asset list item
- Disconnect before connecting to new asset (or auto-switch)

### Decision: Asset List Features

**What**: Left panel shows:
- Search/filter input at top
- Scrollable asset list
- Each item shows: name, IP, connection status
- Active connection highlighted

**Implementation**:
- Use Ant Design `List` or `Menu` component
- Show connected asset with checkmark/badge
- Support search by name/IP
- Group by status (optional enhancement)

## Implementation Details

### MainLayout Changes

**File**: `frontend/src/layouts/MainLayout.tsx`

**Remove**:
- WebSSH menu item from `menuItems` array (lines 101-105)

**Add**:
- Import terminal icon (e.g., `ConsoleSqlOutlined`)
- Add WebSSH button in header (between theme toggle and user menu)
- Button onClick: `navigate('/ssh/terminal')`

**Code**:
```tsx
import { ConsoleSqlOutlined } from '@ant-design/icons'

// In header, before theme toggle:
<Button
  type="text"
  icon={<ConsoleSqlOutlined />}
  onClick={() => navigate('/ssh/terminal')}
  title="WebSSH 终端"
>
  SSH
</Button>
```

### WebSSHTerminal Component Changes

**File**: `frontend/src/pages/ssh/WebSSHTerminal.tsx`

**Remove**:
- `MainLayout` import and wrapper
- Connection modal (`connectModalVisible`, `Modal` component)
- Modal-related form and state

**Add**:
- Minimal page header (back button, title)
- Split-view layout (left: asset list, right: terminal)
- Asset list component in left panel
- Direct connection on asset click

**Layout Structure**:
```tsx
<Layout style={{ height: '100vh' }}>
  <Header>Back Button | WebSSH Terminal</Header>
  <Layout style={{ flexDirection: 'row' }}>
    <Sider width={350}>
      {/* Asset List */}
      <Input.Search />
      <List
        dataSource={assets}
        renderItem={(asset) => (
          <List.Item
            onClick={() => handleConnect(asset)}
            className={connectedAsset?.id === asset.id ? 'active' : ''}
          >
            <List.Item.Meta
              title={asset.hostName}
              description={`${asset.ip} - ${asset.ssh_user || 'N/A'}`}
            />
            {connectedAsset?.id === asset.id && <CheckCircleOutlined />}
          </List.Item>
        )}
      />
    </Sider>
    <Content>
      {/* Terminal Area */}
      <div>
        <div>Current: {connectedAsset?.hostName || 'Not connected'}</div>
        <div ref={terminalRef} />
        <Button onClick={handleDisconnect}>Disconnect</Button>
      </div>
    </Content>
  </Layout>
</Layout>
```

### App.tsx Changes

**File**: `frontend/src/App.tsx`

**Modify**:
- Remove `MainLayout` wrapper from `/ssh/terminal` route
- Keep route as standalone

**Code**:
```tsx
<Route
  path="/ssh/terminal"
  element={
    <ProtectedRoute>
      <Suspense fallback={<PageLoading />}>
        <WebSSHTerminal />
      </Suspense>
    </ProtectedRoute>
  }
/>
```

## Risks / Trade-offs

### Risk: Layout Complexity

**Risk**: Split-view layout may be complex to implement and maintain.

**Mitigation**: 
- Use Ant Design Layout components (Sider, Content)
- Keep layout simple and responsive
- Test on different screen sizes

### Risk: Mobile Experience

**Risk**: Split-view may not work well on mobile devices.

**Mitigation**:
- Make left panel collapsible/drawer on mobile
- Consider overlay asset list on small screens
- Test responsive behavior

### Trade-off: Standalone vs MainLayout

**Trade-off**: Standalone page loses navigation consistency but gains space and focus.

**Decision**: Prefer standalone for terminal-focused experience. Users can use back button to return.

## Migration Plan

1. Update MainLayout:
   - Remove menu item
   - Add header button
2. Update WebSSHTerminal component:
   - Remove MainLayout wrapper
   - Add page header
   - Implement split-view layout
   - Add asset list sidebar
   - Remove connection modal
   - Implement direct connection
3. Update App.tsx route
4. Test:
   - Header button opens page correctly
   - Asset list displays and searches
   - Clicking asset connects
   - Terminal works as before
   - Back button returns to previous page

## Open Questions

- Should we support multiple terminal tabs (one per asset)?
  - **Answer**: Out of scope for this change, can be future enhancement

- Should asset list show connection status history?
  - **Answer**: Not in initial implementation, keep simple

- What icon to use for WebSSH button?
  - **Answer**: `ConsoleSqlOutlined` or `TerminalOutlined` (if available)