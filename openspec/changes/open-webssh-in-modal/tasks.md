## 1. Extract WebSSH Content Component
- [x] 1.1 Create reusable WebSSHContent component
  - Extract WebSSH page content into a new component `WebSSHContent.tsx`
  - Move all state management, tree logic, and terminal tab logic to the new component
  - Ensure component is self-contained and doesn't depend on routing context
  - Component should accept optional props for customization (if needed)

- [x] 1.2 Update WebSSH page to use extracted component
  - Modify `frontend/src/pages/WebSSH.tsx` to use the new `WebSSHContent` component
  - Ensure page route still works correctly
  - Test direct page access functionality

## 2. Create WebSSH Modal Component
- [x] 2.1 Create WebSSHModal component
  - Create new component `frontend/src/components/WebSSHModal.tsx`
  - Use Ant Design `Modal` component with full-screen configuration
  - Configure modal props: width, height, centered, closable, mask, keyboard
  - Accept `open` and `onClose` props to control visibility

- [x] 2.2 Integrate WebSSHContent into modal
  - Render `WebSSHContent` component inside the modal
  - Ensure proper sizing and layout within modal
  - Handle scroll behavior if content exceeds modal size
  - Test modal rendering and layout

## 3. Update MainLayout to Use Modal
- [x] 3.1 Add modal state management
  - Add `websshModalVisible` state to MainLayout component
  - Add handlers for opening and closing modal

- [x] 3.2 Update header button
  - Change button from `Link` component to regular `Button`
  - Update `onClick` handler to open modal instead of navigating
  - Remove Link import if no longer needed
  - Test button click behavior

- [x] 3.3 Integrate WebSSHModal component
  - Import and render `WebSSHModal` component in MainLayout
  - Pass `open` and `onClose` props to control modal visibility
  - Ensure modal appears above all other content (z-index)

## 4. Modal Configuration and Styling
- [x] 4.1 Configure modal dimensions
  - Set modal width to `95vw` or `1200px` (whichever is smaller)
  - Set modal height to `90vh` or `800px` (whichever is smaller)
  - Ensure modal is centered on screen
  - Test on different screen sizes

- [x] 4.2 Configure modal behavior
  - Enable close button (X) in modal header
  - Enable ESC key to close modal
  - Configure mask (backdrop) behavior (clickable or not)
  - Add appropriate modal title (e.g., "WebSSH 终端")

- [x] 4.3 Style modal content
  - Ensure WebSSHContent fits properly within modal
  - Handle overflow and scrolling if needed
  - Test layout on various screen sizes
  - Ensure terminal component works correctly in modal context

## 5. Testing and Validation
- [ ] 5.1 Test modal functionality
  - Verify modal opens when header button is clicked
  - Verify modal closes via close button, ESC key, and backdrop click
  - Test modal behavior on different screen sizes
  - Verify WebSSH functionality works correctly in modal

- [ ] 5.2 Test backward compatibility
  - Verify direct page access via `/webssh` route still works
  - Verify WebSSH page component renders correctly
  - Test that both access methods work independently

- [ ] 5.3 Test terminal functionality in modal
  - Verify terminal connections work correctly
  - Test terminal tab management (open, close, clone, etc.)
  - Verify tree view and host selection work correctly
  - Test terminal authentication flow

- [x] 5.4 Code validation
  - Run `npm run build` to verify frontend compiles without errors
  - Fix any linting errors
  - Ensure TypeScript types are correct
  - Verify no console errors or warnings

## 6. Documentation
- [ ] 6.1 Update component documentation
  - Add JSDoc comments to new components
  - Document component props and usage
  - Update any relevant README or documentation

