## Context
The current implementation integrates WebSSH terminal access directly into the host management page. Users can click "打开终端" button in the host list operations column to open terminals. This creates UI clutter and mixes concerns - host management (CRUD operations) and terminal access are intertwined.

## Goals
- Separate WebSSH terminal access into its own dedicated management page
- Remove terminal functionality from host management page
- Provide a focused interface for WebSSH operations
- Maintain all existing WebSSH functionality (terminal I/O, authentication, resize, etc.)

## Non-Goals
- Changing WebSSH backend API or WebSocket protocol
- Modifying terminal component functionality
- Changing authentication flow
- Adding new WebSSH features

## Decisions

### Decision: Create Dedicated WebSSH Page
**What**: Create a new `/webssh` route with a dedicated page component that displays hosts and manages terminal sessions.

**Why**: 
- Separation of concerns: host management focuses on CRUD, WebSSH focuses on terminal access
- Better UX: users can have a dedicated workspace for terminal operations
- Cleaner UI: host management page is less cluttered

**Alternatives considered**:
- Keep terminal in host management but move to a separate tab - **Rejected**: Still mixes concerns
- Modal-based terminal - **Rejected**: Less convenient for multiple terminals

### Decision: Remove Terminal Button from Host Management
**What**: Remove "打开终端" button and all terminal-related UI from host management page.

**Why**: 
- Host management should focus on host CRUD operations
- Terminal access is available in dedicated WebSSH page
- Reduces UI complexity

**Alternatives considered**:
- Keep button but redirect to WebSSH page - **Rejected**: Creates unnecessary navigation, better to have direct access in WebSSH page

### Decision: Menu Item Placement
**What**: Place WebSSH menu item after "主机管理" in the navigation menu.

**Why**: 
- Logical grouping: WebSSH is related to host management
- Users who manage hosts will likely use WebSSH
- Maintains menu order consistency

## Risks / Trade-offs

### Risk: Users may not find WebSSH page
**Mitigation**: 
- Clear menu item label ("WebSSH管理" or "终端管理")
- Appropriate icon for visibility
- Documentation updates

### Risk: Breaking existing workflows
**Mitigation**: 
- This is a UI reorganization, not a functional change
- All terminal functionality remains the same
- Users can still access terminals, just from a different page

## Migration Plan
1. Create new WebSSH page component
2. Add route and menu item
3. Remove terminal functionality from host management page
4. Test thoroughly
5. Deploy

No data migration needed - this is purely a UI reorganization.

## Open Questions
None - this is a straightforward UI separation.

