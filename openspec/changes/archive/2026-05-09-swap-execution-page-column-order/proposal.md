# Change: 运维执行页左右两栏显示顺序交换

## Why

当前 `/executions` 页面第一行为「模式选择 | 执行选项」、第二行为「脚本编辑 | 主机选择」。用户希望**左右两边的显示选项交换**，使布局更符合使用习惯（例如先选执行选项/主机，左侧；模式与脚本右侧）。

## What Changes

- **第一行**：左右交换 —— 左侧改为「执行选项」，右侧改为「模式选择」。
- **第二行**：左右交换 —— 左侧改为「主机选择」，右侧改为「脚本编辑」。
- 第三行（执行按钮）、第四行（执行结果）不变；单栏宽度与高度一屏约束不变。

## Impact

- **受影响规格**：execution-operator（Execution Page Layout 中块顺序描述）
- **受影响代码**：`frontend/src/pages/executions/ExecutionOperatorPage.tsx` 中两行 Row 内 Col 的先后顺序（组件左右对调）
