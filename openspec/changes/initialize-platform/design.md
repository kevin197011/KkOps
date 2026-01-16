# Design: Platform Architecture and Technical Decisions

## Context

Building a new intelligent operations management platform from scratch. The platform needs to support:
- Multi-user management with role-based access control
- Comprehensive IT asset tracking and management
- Automated operation task execution
- Secure WebSSH terminal access
- Modern, professional frontend interface

## Goals / Non-Goals

### Goals
- Provide a unified platform for IT operations management
- Support both web UI and programmatic API access
- Ensure security through RBAC, encrypted storage, and audit logging
- Enable scalable architecture for future expansion
- Deliver professional, efficient user experience

### Non-Goals
- Real-time monitoring and alerting (planned for future)
- Workflow engine (planned for future)
- Mobile applications (web-first approach)
- Multi-tenant isolation (single-tenant initial design)

## Decisions

### Backend Architecture
- **Language**: Golang 1.21+ for performance and concurrency
- **Web Framework**: Gin (recommended) for performance and ecosystem
- **ORM**: GORM for database operations
- **Database**: PostgreSQL for relational data, Redis for caching/sessions
- **Authentication**: JWT for web sessions, API tokens for programmatic access

**Rationale**: Golang provides excellent performance for concurrent operations (SSH connections, task execution). PostgreSQL offers robust data integrity and JSON support. Redis enables fast session management and caching.

### Frontend Architecture
- **Framework**: React 18+ with TypeScript
- **UI Library**: Ant Design 5+ for professional components
- **State Management**: Redux Toolkit or Zustand
- **Build Tool**: Vite for fast development
- **Design System**: Swiss Modernism 2.0 + Minimalism (professional, data-dense)

**Rationale**: React + Ant Design provides a mature, professional UI framework. TypeScript ensures type safety. Vite offers excellent developer experience.

### Security Architecture
- **Password Storage**: bcrypt hashing
- **API Tokens**: SHA-256 or bcrypt hashing (similar to passwords)
- **SSH Private Keys**: AES-256 encryption
- **Token Display**: Full token shown only once, then prefix-only
- **Permission Model**: RBAC (User → Role → Permission)

**Rationale**: Industry-standard encryption ensures security. RBAC provides flexible, scalable permission management.

### Database Design
- **Schema**: Normalized relational design
- **JSON Fields**: Used for flexible config data (asset_configs, auth_config)
- **Indexing**: Strategic indexes on frequently queried fields
- **Audit Fields**: created_at, updated_at, last_used_at for tracking

**Rationale**: Normalized design ensures data integrity. JSON fields provide flexibility where needed.

## Risks / Trade-offs

### Risks
- **Complexity**: Multiple integrated systems (asset management, task execution, WebSSH)
  - **Mitigation**: Modular design, phased implementation
- **Security**: Handling sensitive data (SSH keys, credentials)
  - **Mitigation**: Encryption at rest, secure transmission, audit logging
- **Performance**: Concurrent SSH connections and task execution
  - **Mitigation**: Goroutine pools, connection pooling, Redis caching

### Trade-offs
- **Database**: PostgreSQL chosen over NoSQL for data integrity despite JSON fields being used
- **Frontend**: Ant Design provides professional UI but adds bundle size
- **SSH Implementation**: xterm.js provides terminal functionality but requires WebSocket infrastructure

## Migration Plan

N/A - This is an initial implementation, no migration required.

## Open Questions

- Specific deployment environment requirements?
- Expected user scale for initial launch?
- Integration requirements with existing systems?
