## ADDED Requirements

### Requirement: User Logout
The system SHALL provide user logout functionality that allows authenticated users to end their session and clear their authentication state.

#### Scenario: User logs out via UI
- **WHEN** an authenticated user clicks the logout button in the user menu
- **THEN** the system clears the authentication token and user information from client storage (localStorage)
- **AND** the user is redirected to the login page
- **AND** all authenticated API requests are invalidated

#### Scenario: Logout API endpoint
- **WHEN** an authenticated user calls the logout API endpoint (`POST /api/v1/auth/logout`)
- **THEN** the system returns a successful response (200 OK)
- **AND** the endpoint requires authentication (returns 401 if not authenticated)
- **NOTE**: For JWT-based authentication, logout is primarily client-side (token deletion), but the endpoint provides API consistency

#### Scenario: User session cleanup after logout
- **WHEN** a user logs out
- **THEN** the authentication token is removed from localStorage
- **AND** user information is removed from localStorage
- **AND** the authentication state in the application store is cleared
- **AND** the user is no longer considered authenticated

#### Scenario: Protected route access after logout
- **WHEN** a user logs out
- **THEN** the user cannot access protected routes
- **AND** any attempt to access protected routes redirects to the login page
- **AND** API requests without authentication return 401 Unauthorized

### Requirement: User Menu in Header
The system SHALL display user information and provide a logout option in the application header.

#### Scenario: Display user information
- **WHEN** a user is authenticated and viewing the application
- **THEN** the header displays user information (username or real name)
- **AND** the user information is displayed in a user menu/dropdown component
- **AND** the user menu is positioned appropriately in the header (typically on the right side)

#### Scenario: User menu contains logout option
- **WHEN** an authenticated user opens the user menu
- **THEN** the menu displays a logout option
- **AND** the logout option is clearly labeled (e.g., "登出" or "Logout")
- **AND** the logout option has an appropriate icon (e.g., LogoutOutlined)
- **AND** clicking the logout option triggers the logout process

#### Scenario: User menu accessibility
- **WHEN** a user interacts with the user menu
- **THEN** the menu is keyboard accessible (Tab, Enter, Escape keys)
- **AND** the menu has proper ARIA labels for screen readers
- **AND** the logout option is announced correctly by screen readers