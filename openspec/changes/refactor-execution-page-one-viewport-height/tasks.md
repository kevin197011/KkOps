# Tasks: 运维执行页高度一屏、宽度不变、功能块重规划

## 1. 恢复单栏宽度

- [x] 1.1 撤销 ExecutionOperatorPage 的左右分栏（Row/Col 配置区与结果区左右分列）
- [x] 1.2 恢复内容区为单栏：maxWidth: 1200（或原值）、margin: 0 auto、width: 100%
- [x] 1.3 根容器保持或设置 minHeight: calc(100vh - 120px)，保证一屏高度

## 2. 功能块布局重规划

- [x] 2.1 第一行：模式选择（ExecutionModeSelector）与执行选项（ExecutionOptions）并排（Row + Col）
- [x] 2.2 第二行：脚本编辑（ScriptEditor）与主机选择（HostSelector）并排（Row + Col），各设 maxHeight + overflow: auto
- [x] 2.3 第三行：执行按钮（及「查看任务历史」链接）
- [x] 2.4 第四行：执行结果（ExecutionResults），占剩余高度，容器 maxHeight 或 flex + overflow: auto
- [x] 2.5 宽度 &lt; 992px 时：脚本区与主机区改为上下堆叠（Col xs={24}）

## 3. 高度约束

- [x] 3.1 脚本编辑区：maxHeight（如 22vh 或 200px），内部滚动
- [x] 3.2 主机列表容器：maxHeight（如 220px 或 24vh），内部滚动
- [x] 3.3 执行结果区：占剩余高度，内部滚动，不撑开整页

## 4. 样式与验收

- [x] 4.1 统一间距与内边距，保持与现有风格一致
- [x] 4.2 在典型桌面分辨率下验证：宽度为单栏居中，高度一屏内可见，功能块按新布局展示
- [x] 4.3 验证脚本/主机/结果区长内容时仅内部滚动
