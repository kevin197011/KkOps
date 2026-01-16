# Change: Redesign Task Management UI with GitHub Actions Style

## Why

The current task management interface uses a basic table layout that lacks visual workflow representation and real-time execution feedback. Users need a more intuitive, GitHub Actions-inspired interface that:

- Provides clear visual representation of workflow execution status
- Enables quick navigation between workflows and their execution history
- Shows real-time execution logs inline without page navigation
- Supports template-based workflow definition with multi-host execution

## What Changes

### Frontend Changes

- **New dual-column layout**: Left panel shows workflow list, right panel shows execution details
- **GitHub Actions-style workflow cards**: Visual cards showing workflow name, status, and recent runs
- **Inline execution logs**: Expandable log viewer within the execution details panel
- **Real-time status updates**: WebSocket-based live updates for running workflows
- **Enhanced execution controls**: Edit, Execute, Cancel, and View Logs actions in a unified interface
- **Project/Environment-based host selection**: Filter and select execution targets by project or environment

### Template Enhancement

- **Templates as complete workflow definitions**: Each template defines a self-contained workflow
- **Multi-host execution support**: Execute template content across selected assets simultaneously

### UI/UX Improvements

- Status indicators with colors and animations for running states
- Collapsible execution history per workflow
- Quick-action buttons for common operations
- Responsive design for various screen sizes

## Impact

- **Affected specs**: task-management (new capability)
- **Affected code**:
  - `frontend/src/pages/tasks/TaskList.tsx` - Major redesign
  - `frontend/src/pages/tasks/TaskExecutionList.tsx` - Integrate into new layout
  - `frontend/src/pages/tasks/TaskExecutionLogs.tsx` - Convert to inline component
  - `frontend/src/api/task.ts` - May need additional endpoints
  - Backend WebSocket handler for real-time log streaming (existing)
- **No breaking API changes**: Backend APIs remain compatible
