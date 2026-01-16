# Implementation Tasks

## Backend Implementation

- [x] 1.1 Add Logout handler method in `backend/internal/handler/auth/handler.go`
  - Add `Logout` handler function with Swagger annotations
  - Return 200 OK response (JWT is stateless, so logout is primarily client-side)

- [x] 1.2 Register logout route in `backend/cmd/server/main.go`
  - Add `POST /api/v1/auth/logout` route (protected, requires authentication)
  - Wire up the logout handler

- [x] 1.3 Add logout API call in `frontend/src/api/auth.ts`
  - Add `logout` method to `authApi` object
  - Call `POST /api/v1/auth/logout` endpoint

## Frontend Implementation

- [x] 2.1 Update MainLayout to display current user information
  - Add user information display in header (username/real name)
  - Create user menu/dropdown using Ant Design `Dropdown` component
  - Position in header next to theme toggle

- [x] 2.2 Implement logout functionality in MainLayout
  - Add logout handler that:
    - Calls `authApi.logout()` (optional, for API consistency)
    - Calls `clearAuth()` from auth store
    - Navigates to `/login` page
    - Shows success message

- [x] 2.3 Add logout menu item to user dropdown
  - Add logout option with icon (LogoutOutlined from Ant Design)
  - Ensure proper styling and accessibility (ARIA labels, keyboard navigation)

- [x] 2.4 Update ProtectedRoute to handle logout properly
  - Ensure ProtectedRoute redirects to login when user logs out
  - Verify authentication state is properly cleared
  - Note: ProtectedRoute already correctly handles this - it redirects to login when user is not authenticated

## Testing & Validation

- [ ] 3.1 Test logout flow manually
  - Verify logout button is visible when user is authenticated
  - Verify clicking logout clears token and user info from localStorage
  - Verify user is redirected to login page
  - Verify user cannot access protected routes after logout

- [ ] 3.2 Test logout API endpoint (if implemented)
  - Verify logout endpoint returns 200 OK
  - Verify endpoint requires authentication (401 if not authenticated)

- [ ] 3.3 Test user menu accessibility
  - Verify keyboard navigation works (Tab, Enter, Escape)
  - Verify screen reader announces user menu and logout option
  - Verify ARIA labels are properly set

## Documentation

- [x] 4.1 Update API documentation (Swagger)
  - Ensure logout endpoint is documented in Swagger UI
  - Add request/response examples
  - Note: Swagger annotations added to handler, documentation will be generated when `swag init` is run

- [x] 4.2 Update user documentation (if applicable)
  - Document logout functionality in user guide
  - Note: User documentation can be updated separately if needed