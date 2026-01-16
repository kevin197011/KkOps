# Design: Open SSH Terminal in New Tab

## Context

The SSH button in the header currently uses React Router's `navigate()` to open the terminal page in the same tab. Users want it to open in a new tab instead to preserve their current context.

## Goals

1. Change SSH button to open terminal in new tab
2. Preserve current page context
3. Maintain existing terminal functionality

## Non-Goals

- Changing terminal page functionality
- Adding new terminal features
- Modifying authentication or routing

## Decisions

### Decision: Use window.open() for New Tab

**What**: Replace `navigate('/ssh/terminal')` with `window.open('/ssh/terminal', '_blank')`.

**Why**:
- Standard browser API for opening new tabs
- Preserves current page context
- Works with existing routing and authentication

**Implementation**:
```tsx
onClick={() => window.open('/ssh/terminal', '_blank')}
```

**Considerations**:
- Browsers may block pop-ups if not triggered by direct user interaction (click event is fine)
- `_blank` target opens in new tab
- Same-origin URL, so authentication tokens work

### Decision: Keep Direct Navigation Support

**What**: Allow direct navigation to `/ssh/terminal` (e.g., from bookmarks, direct URLs) to still work.

**Why**:
- Maintains flexibility
- Some users may prefer direct navigation
- No need to break existing functionality

**Implementation**:
- No changes to routes or terminal page
- Only change button behavior in header

## Implementation Details

### MainLayout Changes

**File**: `frontend/src/layouts/MainLayout.tsx`

**Modify**:
- Change SSH button onClick from `navigate('/ssh/terminal')` to `window.open('/ssh/terminal', '_blank')`
- Optionally remove `navigate` import if no longer needed (check if used elsewhere)

**Code**:
```tsx
<Button
  type="text"
  icon={<ConsoleSqlOutlined />}
  onClick={() => window.open('/ssh/terminal', '_blank')}
  title="WebSSH 终端"
  style={{
    color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
  }}
>
  SSH
</Button>
```

## Risks / Trade-offs

### Risk: Pop-up Blockers

**Risk**: Some browsers or extensions may block the new tab if triggered incorrectly.

**Mitigation**: 
- Using `window.open()` with direct user click event is the standard approach
- Modern browsers allow tabs opened from user interactions
- If blocked, users can use browser settings to allow it

### Trade-off: Navigation vs New Tab

**Trade-off**: New tab opens new context but requires tab management.

**Decision**: Prefer new tab for terminal access as it's an auxiliary tool that users typically want alongside their main work.

## Migration Plan

1. Update MainLayout header button onClick handler
2. Test that new tab opens correctly
3. Verify authentication works in new tab
4. Test that terminal functionality works in new tab

## Open Questions

- Should we add any special handling for mobile devices?
  - **Answer**: `window.open()` with `_blank` typically opens new tab on desktop, and new window on mobile (browser-dependent). This is acceptable behavior.

- Should we add any visual indicator that it opens in a new tab?
  - **Answer**: Optional enhancement - could add tooltip text "在新标签页打开" but current tooltip is sufficient.