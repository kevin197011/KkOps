# Change: Open SSH Terminal in New Tab

## Why

Currently, when users click the SSH button in the header, it navigates to the WebSSH terminal page in the current tab/window. This replaces the current page and users lose their current context.

Opening the SSH terminal in a new tab will:
- **Preserve context**: Users can keep their current work open while accessing the terminal
- **Better workflow**: Users can reference information from other pages while using the terminal
- **Parallel work**: Users can work on multiple things simultaneously
- **Standard behavior**: Opening external tools in new tabs is a common pattern for tools like terminals, editors, etc.

## What Changes

- **MODIFIED**: Change SSH header button to open WebSSH terminal in a new browser tab/window
- **ADDED**: Use `window.open()` to open new tab with `/ssh/terminal` URL
- **MAINTAINED**: Existing WebSSH terminal page functionality remains unchanged
- **MAINTAINED**: Direct navigation to `/ssh/terminal` still works (e.g., from bookmarks, direct URLs)

## Impact

- **User Experience**: 
  - SSH terminal opens in new tab, preserving current page
  - Users can work with multiple tabs simultaneously
- **Code Changes**:
  - Update MainLayout header button onClick handler
  - Change from `navigate()` to `window.open()`

## Considerations

- **Browser behavior**: Some browsers may block pop-ups if not triggered by user interaction
- **URL handling**: New tab opens with same domain, so authentication tokens should work
- **Multiple tabs**: Users can open multiple terminal tabs if needed
- **Closing behavior**: Users can close terminal tab and return to previous page easily