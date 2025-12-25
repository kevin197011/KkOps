## Context
Currently, SSH key management is implemented as a tab within the WebSSH management page. Users access SSH key management by navigating to the WebSSH page and clicking the "SSH密钥" tab. This design was chosen to keep SSH-related functionality together, but it creates UI complexity and mixes different concerns.

## Goals
- Separate SSH key management into its own dedicated page
- Remove SSH key management UI from WebSSH page
- Maintain all existing SSH key functionality (CRUD operations, authentication)
- Provide independent navigation to SSH key management
- Keep WebSSH page focused solely on terminal access

## Non-Goals
- Changing SSH key backend API or data model
- Modifying SSH key authentication flow in WebSSH
- Adding new SSH key features
- Changing SSH key security or encryption

## Decisions

### Decision: Create Dedicated SSH Keys Page
**What**: Create a new `/ssh-keys` route with a dedicated page component for SSH key management.

**Why**: 
- Separation of concerns: SSH key management is independent from terminal access
- Better UX: Users can manage keys without navigating through WebSSH page
- Cleaner UI: WebSSH page becomes simpler, focused only on terminal access
- First-class feature: SSH key management deserves its own dedicated space

**Alternatives considered**:
- Keep SSH key management in WebSSH but move to a separate route - **Rejected**: Still mixes concerns in the same component
- Move SSH key management to user settings - **Rejected**: Keys are used for WebSSH, should be easily accessible for terminal access
- Keep current tab-based design - **Rejected**: User explicitly requested separation

### Decision: Menu Item Placement
**What**: Place SSH Keys menu item after "WebSSH管理" in the navigation menu.

**Why**: 
- Logical grouping: SSH keys are related to WebSSH functionality
- Easy access: Users managing terminals will find key management nearby
- Clear hierarchy: Shows the relationship between WebSSH and SSH keys

**Alternatives considered**:
- Place in settings section - **Rejected**: Keys are used actively for terminal access, not just configuration
- Place before WebSSH - **Rejected**: WebSSH is the primary feature, keys are supporting

### Decision: Reuse Existing Service
**What**: Keep using `sshKeyService` from `frontend/src/services/sshKey.ts` without modification.

**Why**: 
- No API changes needed
- Service is already well-designed and functional
- Maintains consistency with existing code

### Decision: Extract Code from WebSSH
**What**: Move SSH key management UI code (state, handlers, forms, tables) from WebSSH component to new SSHKeys component.

**Why**: 
- Code reuse: Avoid duplicating functionality
- Consistency: Maintain the same UI/UX as before
- Simplicity: Extract working code rather than rewriting

