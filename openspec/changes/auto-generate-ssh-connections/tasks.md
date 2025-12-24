## 1. Backend Implementation

- [x] 1.1 Add `AutoGenerateSSHConnections` method to `SSHConnectionService`
  - Query all hosts with Linux/Unix OS types (e.g., "Linux", "Ubuntu", "CentOS", "RedHat", "Debian", "FreeBSD", "OpenBSD")
  - Filter out hosts that already have SSH connections (check by host_id)
  - For each eligible host, create SSH connection with:
    - hostname: host.IPAddress or host.Hostname
    - port: host.SSHPort (default 22)
    - project_id: host.ProjectID
    - host_id: host.ID
    - username: configurable default (e.g., "root", "ubuntu", "admin")
    - auth_type: configurable default (e.g., "password" or "key")
    - name: auto-generated name (e.g., "Auto: {hostname}:{port}")
    - status: "active"

- [x] 1.2 Add `AutoGenerateSSHConnectionsForProject` method
  - Similar to above but scoped to a specific project
  - Accept project_id as parameter

- [x] 1.3 Add `AutoGenerateSSHConnectionsForHost` method
  - Generate SSH connection for a single host
  - Useful for on-demand generation when a new host is added

- [x] 1.4 Add handler endpoint `POST /api/v1/ssh/connections/auto-generate`
  - Accept optional parameters: project_id, default_username, default_auth_type
  - Return summary: total hosts processed, connections created, skipped (already exists)

- [x] 1.5 Add handler endpoint `POST /api/v1/ssh/connections/auto-generate/:host_id`
  - Generate SSH connection for a specific host
  - Return created connection or error if already exists

- [x] 1.6 Add validation logic
  - Ensure host has valid IP address or hostname
  - Ensure host has valid SSH port
  - Ensure host OS type is Linux/Unix (skip Windows, macOS, etc.)

## 2. Frontend Implementation

- [x] 2.1 Add `autoGenerateConnections` method to `sshService` in `frontend/src/services/ssh.ts`
  - Call `POST /api/v1/ssh/connections/auto-generate` endpoint
  - Accept optional parameters: project_id, default_username, default_auth_type

- [x] 2.2 Add "自动生成 SSH 连接" button to SSH management page
  - Place in SSH Connections tab, next to "新建连接" button
  - Show modal with options:
    - Project filter (optional, default: all projects)
    - Default username (optional, default: "root")
    - Default authentication type (optional, default: "password")
    - Preview: estimated number of hosts eligible

- [x] 2.3 Add auto-generation result display
  - Show success message with summary (e.g., "成功为 15 台主机创建 SSH 连接")
  - Show warning for skipped hosts (e.g., "3 台主机已存在连接，已跳过")
  - Refresh connection list after generation

- [x] 2.4 Add "自动生成" action to host management page
  - Add button/action in host row for Linux/Unix hosts
  - Quick action to generate SSH connection for single host

- [x] 2.5 Add visual indicator for auto-generated connections
  - Show badge or icon in connection list for auto-generated connections
  - Differentiate from manually created connections

## 3. Testing

- [ ] 3.1 Backend unit tests
  - Test `AutoGenerateSSHConnections` with various OS types
  - Test skipping hosts with existing connections
  - Test validation logic (invalid IP, port, etc.)
  - Test project-scoped generation

- [ ] 3.2 Backend integration tests
  - Test auto-generation API endpoints
  - Test concurrent generation (idempotency)

- [ ] 3.3 Frontend unit tests
  - Test auto-generation UI components
  - Test API service calls

- [ ] 3.4 Manual testing
  - Test with real Linux hosts
  - Test with mixed OS types (Linux, Windows, etc.)
  - Test with hosts that already have connections
  - Test with hosts missing required information

## 4. Documentation

- [ ] 4.1 Update API documentation
  - Document new auto-generation endpoints
  - Document parameters and response format

- [ ] 4.2 Update user manual
  - Document auto-generation feature
  - Document configuration options (default username, auth type)

