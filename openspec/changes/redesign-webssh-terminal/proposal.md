# Change: Redesign WebSSH Terminal - Move to Header Button and Split View

## Why

Currently, the WebSSH terminal is accessible via a menu item and uses a modal-based connection interface. This design has limitations:

- **Menu clutter**: The WebSSH menu item takes up space in the sidebar menu
- **Poor discoverability**: Users may not easily find the terminal feature in the menu
- **Limited asset visibility**: Users can only see one asset at a time in the modal
- **Inefficient workflow**: Switching between assets requires disconnecting and reconnecting
- **Modal-based UX**: The connection modal doesn't provide a good overview of available assets

Moving WebSSH to a header button and redesigning it as a split-view page will:
- **Reduce menu clutter**: Free up menu space by moving to header
- **Improve discoverability**: Header button is always visible and easy to find
- **Better asset overview**: Left sidebar shows all assets for quick selection
- **Faster asset switching**: Users can switch between assets without modal dialogs
- **Improved UX**: Split-view provides a better terminal experience similar to modern SSH clients

## What Changes

- **MODIFIED**: Remove WebSSH menu item from sidebar menu in MainLayout
- **ADDED**: Add WebSSH button in header (top right, next to theme toggle and user menu)
- **MODIFIED**: Redesign WebSSHTerminal page as standalone page (without MainLayout wrapper)
- **MODIFIED**: Change page layout to split-view:
  - **Left panel**: Asset list sidebar showing all available assets
  - **Right panel**: SSH terminal connection area
- **MODIFIED**: Remove connection modal, replace with direct connection from asset list
- **ADDED**: Click asset in left panel to connect/disconnect
- **MAINTAINED**: Keep existing WebSocket SSH connection functionality
- **MAINTAINED**: Keep existing terminal functionality (xterm.js)

## Impact

- **UI Changes**: 
  - Menu has one less item
  - Header has a new WebSSH button (terminal icon)
  - WebSSH page has new split-view layout
- **User Experience**: 
  - Faster access to terminal
  - Better overview of available assets
  - Easier asset switching
- **Code Changes**:
  - MainLayout: Remove menu item, add header button
  - WebSSHTerminal: Redesign layout, remove modal, add asset list sidebar
  - App.tsx: Keep route but remove MainLayout wrapper

## Layout Structure

### New Page Layout (Standalone, no MainLayout)
```
┌─────────────────────────────────────────────────────────┐
│  [Back]  WebSSH Terminal                                │
├──────────────┬──────────────────────────────────────────┤
│              │  [Terminal Area]                         │
│  Asset List  │                                          │
│  ┌────────┐  │  ┌──────────────────────────────────┐   │
│  │ Asset1 │  │  │                                  │   │
│  │ ✓      │  │  │  Terminal output...              │   │
│  └────────┘  │  │                                  │   │
│  ┌────────┐  │  │                                  │   │
│  │ Asset2 │  │  │                                  │   │
│  └────────┘  │  └──────────────────────────────────┘   │
│  ┌────────┐  │                                          │
│  │ Asset3 │  │  [Disconnect] [Reconnect]              │
│  └────────┘  │                                          │
│              │                                          │
│  [Search]    │                                          │
└──────────────┴──────────────────────────────────────────┘
```

## Considerations

- **Header Button**: Use terminal/console icon (e.g., `ConsoleSqlOutlined` or `TerminalOutlined` if available)
- **Back Button**: Add back button in WebSSH page to return to previous page
- **Asset List**: 
  - Show asset name, IP, status
  - Indicate currently connected asset
  - Support search/filter
  - Click to connect (if not connected) or switch (if connected)
- **Terminal Area**: 
  - Show current asset info
  - Connection controls (disconnect, reconnect)
  - Terminal display (xterm.js)
- **Responsive Design**: On mobile, asset list could be collapsible or overlay