# Change: Separate SSH Key Management to Independent Page

## Why
Currently, SSH key management is integrated as a tab within the WebSSH management page. This creates:
- Mixed concerns: SSH key management and WebSSH terminal access are different functionalities
- UI clutter: The WebSSH page becomes complex with multiple tabs (hosts, SSH keys, terminal sessions)
- Navigation confusion: Users who only need to manage SSH keys must navigate through WebSSH page

By separating SSH key management into its own dedicated page, we can:
- Improve separation of concerns (key management is independent from terminal access)
- Simplify the WebSSH page (focus solely on terminal access)
- Provide better user experience (dedicated page for key management)
- Enable independent navigation (keys can be accessed without going through WebSSH)

## What Changes
- **REMOVED**: SSH Key Management tab from WebSSH page
  - Remove "SSH密钥" tab from `frontend/src/pages/WebSSH.tsx`
  - Remove SSH key management UI code from WebSSH component
  - Remove SSH key state management from WebSSH component

- **ADDED**: New SSH Key Management page
  - Create `frontend/src/pages/SSHKeys.tsx` as a dedicated SSH key management page
  - Move all SSH key management functionality from WebSSH page to new page
  - Include full CRUD operations (create, list, edit, delete)
  - Display key metadata (name, type, fingerprint, username, created date)
  - Support key upload and editing

- **ADDED**: SSH Key Management menu item and route
  - Add `/ssh-keys` route to `frontend/src/App.tsx`
  - Add "SSH密钥管理" menu item to `frontend/src/components/MainLayout.tsx`
  - Place menu item logically (e.g., after WebSSH or in a related group)

- **KEPT**: SSH Key functionality and API
  - All backend SSH key APIs remain unchanged (`/api/v1/ssh-keys/*`)
  - SSH key authentication in WebSSH terminal remains unchanged
  - SSH key selection in terminal connection remains unchanged
  - All existing SSH key features continue to work

## Impact
- **Affected specs**: `ssh-key-management` (modified - UI location changed)
- **Affected code**:
  - `frontend/src/pages/WebSSH.tsx` - Remove SSH key management tab and related code
  - `frontend/src/pages/SSHKeys.tsx` - New page (created from extracted code)
  - `frontend/src/App.tsx` - Add `/ssh-keys` route
  - `frontend/src/components/MainLayout.tsx` - Add SSH Keys menu item
- **User experience**: SSH key management becomes a first-class feature with its own dedicated page, improving navigation and separation of concerns

