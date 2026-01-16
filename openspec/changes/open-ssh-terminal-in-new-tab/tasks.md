# Tasks: Open SSH Terminal in New Tab

## 1. MainLayout Button Update

- [x] 1.1 Update SSH button onClick handler
  - Change from `navigate('/ssh/terminal')` to `window.open('/ssh/terminal', '_blank')`
  - Update button in `MainLayout.tsx`
  - Verify button still has proper styling and tooltip

- [x] 1.2 Verify navigate import (optional cleanup)
  - Check if `navigate` from `useNavigate` is still used elsewhere in MainLayout
  - If not used, remove unused import (optional cleanup, not required)

## 2. Testing

- [ ] 2.1 Test new tab opening
  - Click SSH button and verify new tab opens
  - Verify terminal page loads correctly in new tab
  - Verify current page remains unchanged

- [ ] 2.2 Test authentication in new tab
  - Verify authentication token works in new tab
  - Verify user can use terminal functionality
  - Test logout behavior in new tab

- [ ] 2.3 Test multiple tabs
  - Open multiple terminal tabs
  - Verify each tab works independently
  - Test closing tabs and returning to original page

## 3. Validation

- [ ] 3.1 Run linter
  - Fix any linting errors
  - Verify no TypeScript errors

- [ ] 3.2 Test in browser
  - Test in Chrome/Edge
  - Test in Firefox
  - Test in Safari (if available)
  - Verify pop-up blockers don't interfere

- [ ] 3.3 Validate OpenSpec proposal
  - `openspec validate open-ssh-terminal-in-new-tab --strict`