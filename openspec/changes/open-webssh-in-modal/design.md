# Design: WebSSH Modal Dialog

## Context
The current WebSSH implementation uses routing to navigate to a dedicated page. Users want to access WebSSH terminal functionality without leaving their current page context. This design proposes opening WebSSH in a modal dialog instead.

## Goals
- Open WebSSH in a modal dialog when clicking the header button
- Preserve user's current page context
- Maintain all existing WebSSH functionality
- Support both modal and direct page access (for backward compatibility)

## Non-Goals
- Changing WebSSH functionality or UI layout
- Modifying backend APIs or WebSocket connections
- Removing the direct page route (keep for backward compatibility)
- Implementing modal state persistence across page refreshes

## Decisions

### Decision: Use Ant Design Modal Component
**What**: Use Ant Design `Modal` component with full-screen configuration to display WebSSH interface.

**Why**:
- Consistent with existing UI library (Ant Design)
- Provides built-in close button, backdrop, and keyboard handling
- Supports full-screen mode for terminal interface
- Easy to implement and maintain

**Alternatives considered**:
- Custom modal implementation - **Rejected**: Unnecessary complexity, Ant Design Modal is sufficient
- Drawer component - **Rejected**: Drawer is typically for side panels, not full-screen content
- New window/tab - **Rejected**: More complex, requires window management, less integrated

### Decision: Full-Screen Modal
**What**: Configure modal to use full-screen or near-full-screen dimensions (e.g., 95% width, 90% height).

**Why**:
- Terminal interface needs sufficient space for tree view and terminal tabs
- Full-screen provides better terminal experience
- Maintains visual separation from underlying page

**Alternatives considered**:
- Fixed size modal (e.g., 1200x800) - **Rejected**: May be too small for large terminal windows
- Responsive size based on screen - **Accepted**: Better UX on different screen sizes

### Decision: Component Extraction
**What**: Extract WebSSH content into a reusable component that can be used in both page and modal contexts.

**Why**:
- Avoids code duplication
- Maintains backward compatibility (direct page access still works)
- Easier to maintain and test
- Component can be reused in other contexts if needed

**Alternatives considered**:
- Duplicate WebSSH code for modal - **Rejected**: Code duplication, harder to maintain
- Remove page route entirely - **Rejected**: Breaks backward compatibility, some users may prefer direct page access

### Decision: Modal State Management
**What**: Manage modal open/close state in MainLayout component using React state.

**Why**:
- Simple and straightforward
- State is local to the component that triggers it
- Easy to control from header button

**Alternatives considered**:
- Global state management (Context/Redux) - **Rejected**: Overkill for simple modal state
- URL-based state (query parameter) - **Rejected**: More complex, not necessary for modal state

### Decision: Terminal State Preservation
**What**: When modal is closed, terminal connections are terminated. When reopened, users need to reconnect.

**Why**:
- Simpler implementation
- Avoids resource leaks from keeping connections open
- Clear user expectation (modal closed = connection closed)

**Alternatives considered**:
- Keep connections alive when modal closed - **Rejected**: More complex, potential resource issues, unclear UX
- Persist terminal state in localStorage - **Rejected**: Not requested, adds complexity

## Implementation Details

### Modal Configuration
- **Width**: `95vw` or `1200px` (whichever is smaller)
- **Height**: `90vh` or `800px` (whichever is smaller)
- **Centered**: Yes
- **Mask**: Yes (semi-transparent backdrop)
- **Closable**: Yes (X button in header)
- **Keyboard**: ESC key closes modal
- **Mask click**: Closes modal (optional, can be disabled to prevent accidental closes)

### Component Structure
```
MainLayout
  └─ WebSSHModal (new component)
      └─ WebSSHContent (extracted from WebSSH page)
          ├─ Tree view (left pane)
          └─ Terminal tabs (right pane)
```

### State Management
- `websshModalVisible`: Boolean state in MainLayout to control modal visibility
- Modal opens when header button is clicked
- Modal closes when:
  - User clicks close button (X)
  - User presses ESC key
  - User clicks backdrop (if enabled)

### Backward Compatibility
- Keep `/webssh` route in App.tsx for direct page access
- WebSSH page component can be used directly or wrapped in modal
- Both access methods work independently

## Risks / Trade-offs
- **Risk**: Modal may feel cramped on small screens
  - **Mitigation**: Use responsive sizing, ensure minimum usable dimensions
- **Risk**: Users may accidentally close modal and lose work
  - **Mitigation**: Add confirmation dialog if terminals are active (optional)
- **Risk**: Modal state not persisted across page refreshes
  - **Mitigation**: Acceptable trade-off, users can reopen modal
- **Trade-off**: Terminal connections close when modal closes
  - **Benefit**: Simpler implementation, clear behavior
  - **Cost**: Users need to reconnect when reopening

## Migration Plan
1. Extract WebSSH content into reusable component
2. Create WebSSHModal wrapper component
3. Update MainLayout to use modal instead of navigation
4. Test both modal and direct page access
5. No data migration needed

## Open Questions
- Should we add a confirmation dialog when closing modal with active terminals? (Nice to have, can add later)
- Should modal remember last selected host when reopened? (Nice to have, can add later)

