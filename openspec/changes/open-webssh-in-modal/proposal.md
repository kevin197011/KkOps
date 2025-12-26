# Change: Open WebSSH in Modal Dialog

## Why
Currently, clicking the WebSSH button in the header navigates to the WebSSH page via routing, which changes the current page context. Users need to access WebSSH terminal functionality without leaving their current work context. Opening WebSSH in a modal dialog provides:
- **Context Preservation**: Users can access WebSSH without losing their current page state
- **Quick Access**: Modal opens instantly without full page navigation
- **Better UX**: Users can quickly check terminal output and return to their work
- **Non-blocking**: Modal can be minimized or kept open while working on other tasks

## What Changes
- **MODIFIED**: WebSSH access method
  - Change header button from navigation (`Link`/`navigate`) to modal trigger
  - Open WebSSH interface in a full-screen modal dialog instead of routing
  - Modal should be dismissible and preserve state when reopened

- **ADDED**: WebSSH Modal Component
  - Create a modal wrapper component for WebSSH functionality
  - Use Ant Design `Modal` component with full-screen configuration
  - Support modal open/close state management
  - Preserve terminal connections when modal is closed and reopened

- **MODIFIED**: WebSSH Page Component
  - Extract WebSSH content into a reusable component that can be used in both page and modal contexts
  - Ensure component works correctly in modal context (proper sizing, scroll handling)
  - Handle modal-specific behaviors (close button, escape key, backdrop click)

- **KEPT**: All existing WebSSH functionality
  - Tree-based host list (Project → Environment → Hosts)
  - Multiple terminal tabs support
  - Terminal tab context menu (Clone, Close, Close All, Close Others)
  - Terminal connection and authentication flow
  - All backend APIs and WebSocket connections

## Impact
- **Affected specs**: `webssh-management` (modified - access method changed from routing to modal)
- **Affected code**:
  - `frontend/src/components/MainLayout.tsx` - Change button from Link to modal trigger
  - `frontend/src/pages/WebSSH.tsx` - Extract content component, add modal wrapper
  - `frontend/src/App.tsx` - Keep route for direct access (optional, for backward compatibility)
- **User experience**: 
  - WebSSH opens in modal without page navigation
  - Users can access terminal while staying on current page
  - Modal can be closed and reopened without losing terminal state (if implemented)

