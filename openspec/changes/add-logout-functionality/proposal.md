# Change: Add Logout Functionality

## Why

Currently, the KkOps platform has user login functionality but lacks a logout feature. Users cannot properly end their session, which is a critical security and user experience gap. This change will:
- Allow users to securely log out and clear their authentication state
- Provide a clear UI entry point for logout (user menu/dropdown)
- Ensure proper cleanup of authentication tokens and user data
- Improve security by allowing users to terminate their sessions
- Align with standard authentication flow expectations

## What Changes

- **ADDED**: Logout API endpoint (`POST /api/v1/auth/logout`) for API consistency (optional, primarily client-side for JWT)
- **ADDED**: Frontend logout functionality that clears authentication state (token, user info)
- **ADDED**: User menu/dropdown in MainLayout header with user information and logout option
- **ADDED**: Navigation to login page after logout
- **MODIFIED**: MainLayout header to include user information display

## Impact

- Affected code: Backend auth handler/service, Frontend auth store, MainLayout component
- Affected files: 
  - `backend/internal/handler/auth/handler.go` (add Logout handler)
  - `backend/internal/service/auth/service.go` (add Logout method, optional)
  - `frontend/src/api/auth.ts` (add logout API call)
  - `frontend/src/layouts/MainLayout.tsx` (add user menu with logout)
  - `frontend/src/stores/auth.ts` (ensure clearAuth works properly)
- User experience: Users can now properly log out and end their sessions
- Security: Improved session management and token cleanup