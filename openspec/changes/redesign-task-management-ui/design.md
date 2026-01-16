# Design: GitHub Actions-Style Task Management UI

## Context

The current task management system uses a traditional table-based layout that separates task listing, execution records, and logs into different pages. This creates friction for users who need to monitor and manage multiple workflow executions. The new design consolidates these views into a unified, GitHub Actions-inspired interface.

### Stakeholders
- Operations engineers who create and execute tasks
- Administrators who monitor task executions across projects
- Developers who need to debug failed executions

### Constraints
- Must work with existing backend APIs
- Must support real-time log streaming via WebSocket
- Must maintain responsive design for mobile/tablet

## Goals / Non-Goals

### Goals
- Create a dual-column layout with workflow list and execution details
- Enable inline real-time log viewing without page navigation
- Provide clear visual status indicators for workflow states
- Support quick actions (edit, execute, cancel) from the main view
- Allow host selection by project or environment

### Non-Goals
- Complex DAG-based workflow orchestration (out of scope)
- Step-level dependencies within a single workflow (future enhancement)
- Parallel execution visualization (all hosts execute simultaneously)

## UI Layout Design

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Header: Task Management                              [+ New Task] [Filters] │
├────────────────────────┬────────────────────────────────────────────────────┤
│ Workflow List (300px)  │ Execution Details Panel                            │
│                        │                                                    │
│ ┌────────────────────┐ │ ┌────────────────────────────────────────────────┐ │
│ │ 🟢 Deploy App      │ │ │ Deploy App                      [Edit][Execute]│ │
│ │ Last: 2min ago ✓   │ │ │ Template: deploy-script.sh                     │ │
│ │ #12 success        │ │ │ Status: Success                                │ │
│ └────────────────────┘ │ └────────────────────────────────────────────────┘ │
│                        │                                                    │
│ ┌────────────────────┐ │ ┌─ Execution History ─────────────────────────────┐│
│ │ 🔵 Backup DB       │ │ │ ▶ Run #12 - Success - 2min ago    [View Logs]  ││
│ │ Running...         │ │ │ ▷ Run #11 - Failed - 1hr ago      [View Logs]  ││
│ │ #8 running         │ │ │ ▷ Run #10 - Success - 2hr ago     [View Logs]  ││
│ └────────────────────┘ │ └─────────────────────────────────────────────────┘│
│                        │                                                    │
│ ┌────────────────────┐ │ ┌─ Inline Logs (Expanded) ────────────────────────┐│
│ │ 🔴 Health Check    │ │ │ [node01] Connecting to 192.168.1.8...          ││
│ │ Failed 10min ago   │ │ │ [node01] Running script...                      ││
│ │ #5 failed          │ │ │ [node02] Connecting to 192.168.1.9...          ││
│ └────────────────────┘ │ │ [node01] ✓ Command completed (exit: 0)          ││
│                        │ │ [node02] ✓ Command completed (exit: 0)          ││
│ [Show More...]         │ │ ▼ Auto-scroll                                    ││
│                        │ └─────────────────────────────────────────────────┘│
└────────────────────────┴────────────────────────────────────────────────────┘
```

## Component Architecture

### New Components

1. **TaskManagementPage** (container)
   - Manages layout and state coordination
   - Handles WebSocket connections for real-time updates

2. **WorkflowList** (left panel)
   - Displays workflow cards with status indicators
   - Supports filtering by status, project, environment
   - Handles workflow selection

3. **WorkflowCard**
   - Shows workflow name, last execution status, run count
   - Status indicator: 🟢 success, 🔵 running, 🔴 failed, ⚪ pending

4. **ExecutionDetailsPanel** (right panel)
   - Shows selected workflow details and actions
   - Contains execution history list
   - Manages inline log expansion

5. **ExecutionHistoryList**
   - Lists recent executions with status and timestamps
   - Expandable rows for inline log viewing

6. **InlineLogViewer**
   - Real-time log streaming via WebSocket
   - Per-host log separation with color coding
   - Auto-scroll with manual scroll lock

### State Management

```typescript
interface TaskManagementState {
  workflows: Workflow[]
  selectedWorkflowId: string | null
  executions: Map<string, Execution[]>
  expandedExecutionId: string | null
  logConnections: Map<string, WebSocket>
  filters: {
    status: string[]
    projectId: number | null
    environmentId: number | null
  }
}
```

## Decisions

### Decision 1: Single-page dual-column layout
**What**: Consolidate task list, executions, and logs into one page
**Why**: Reduces navigation overhead, provides immediate context
**Alternatives**: Keep separate pages with enhanced navigation (rejected - too many clicks)

### Decision 2: Inline expandable log viewer
**What**: Logs expand within the execution history row
**Why**: Quick access without losing context of other executions
**Alternatives**: Modal popup (rejected - blocks view), separate panel (considered - may implement later)

### Decision 3: WebSocket per expanded execution
**What**: Only connect WebSocket for the expanded execution's logs
**Why**: Conserves resources, only active log view needs streaming
**Alternatives**: Single WebSocket for all logs (complex multiplexing)

### Decision 4: Template as workflow
**What**: Each template represents a complete workflow definition
**Why**: Matches user mental model, simple execution model
**Alternatives**: Multi-step workflow builder (too complex for current needs)

## Host Selection Design

### Selection Flow
1. User creates/edits task
2. Modal shows:
   - Filter by Project (optional)
   - Filter by Environment (optional)
   - Asset list (filterable, multi-select)
3. Selected assets stored with task

### UI Component
```
┌─ Select Execution Targets ──────────────────────┐
│ Project: [All Projects    ▼]                    │
│ Environment: [All Environments ▼]               │
│ ┌─────────────────────────────────────────────┐ │
│ │ ☑ node01 (192.168.1.8) - prod/g01           │ │
│ │ ☑ node02 (192.168.1.9) - prod/g01           │ │
│ │ ☐ node03 (192.168.1.10) - dev/g02           │ │
│ └─────────────────────────────────────────────┘ │
│ Selected: 2 hosts                               │
└─────────────────────────────────────────────────┘
```

## Status Indicators

| Status    | Color   | Icon | Animation |
|-----------|---------|------|-----------|
| pending   | gray    | ○    | none      |
| running   | blue    | ●    | pulse     |
| success   | green   | ✓    | none      |
| failed    | red     | ✗    | none      |
| cancelled | orange  | ⊘    | none      |

## Risks / Trade-offs

### Risk 1: WebSocket connection management
- **Issue**: Multiple expanded logs could create many connections
- **Mitigation**: Limit to one expanded execution at a time, auto-close after timeout

### Risk 2: Large log volumes
- **Issue**: Long-running tasks generate massive logs
- **Mitigation**: Virtual scrolling, limit displayed lines, "Load more" button

### Risk 3: Mobile responsiveness
- **Issue**: Dual-column layout doesn't work on mobile
- **Mitigation**: Stack columns vertically, collapsible panels on mobile

## Migration Plan

1. Create new TaskManagementPage component alongside existing pages
2. Implement WorkflowList and WorkflowCard first (read-only)
3. Add ExecutionDetailsPanel with execution history
4. Implement InlineLogViewer with WebSocket integration
5. Add task creation/edit modal with enhanced host selection
6. Replace existing TaskList route with new page
7. Keep old components temporarily for fallback

## Open Questions

- ~~Should logs show per-host or merged view?~~ → Per-host with collapsible sections
- ~~How many executions to show in history by default?~~ → 5 recent, with "Load more"
