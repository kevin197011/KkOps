# Change: Separate WebSSH Management from Host Management

## Why
Currently, WebSSH terminal access is integrated directly into the host management page with a "打开终端" button in the operations column. This creates UI clutter and mixes host management concerns with terminal access. Users need a dedicated WebSSH management page where they can view all hosts and open terminals in a focused interface, separate from host CRUD operations.

## What Changes
- **REMOVED**: "打开终端" button from host management page operations column
- **REMOVED**: Terminal tab management from host management page
- **ADDED**: New WebSSH management page at `/webssh` route
- **ADDED**: WebSSH management menu item in main navigation
- **MODIFIED**: WebSSH terminal access now available only from the dedicated WebSSH management page
- **MODIFIED**: Host management page focuses solely on host CRUD operations

## Impact
- **Affected specs**: `webssh-management` capability
- **Affected code**:
  - `frontend/src/pages/Hosts.tsx` - Remove terminal button and tab management
  - `frontend/src/pages/WebSSH.tsx` - New page for WebSSH management
  - `frontend/src/App.tsx` - Add `/webssh` route
  - `frontend/src/components/MainLayout.tsx` - Add WebSSH menu item
- **User experience**: WebSSH becomes a first-class feature with its own dedicated page, improving separation of concerns and user workflow

