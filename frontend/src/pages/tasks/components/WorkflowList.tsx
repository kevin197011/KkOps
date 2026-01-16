// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { Empty } from 'antd'
import WorkflowCard from './WorkflowCard'
import type { WorkflowWithExecutions } from '../TaskManagementPage'

interface WorkflowListProps {
  workflows: WorkflowWithExecutions[]
  selectedWorkflowId: number | null
  onSelect: (id: number) => void
}

const WorkflowList = ({ workflows, selectedWorkflowId, onSelect }: WorkflowListProps) => {
  if (workflows.length === 0) {
    return (
      <div style={{ padding: 40 }}>
        <Empty description="暂无工作流" />
      </div>
    )
  }

  return (
    <div style={{ padding: '8px 0' }}>
      {workflows.map((workflow) => (
        <WorkflowCard
          key={workflow.id}
          workflow={workflow}
          isSelected={selectedWorkflowId === workflow.id}
          onClick={() => onSelect(workflow.id)}
        />
      ))}
    </div>
  )
}

export default WorkflowList
