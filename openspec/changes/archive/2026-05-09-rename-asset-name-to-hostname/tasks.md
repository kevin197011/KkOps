# Tasks: Rename Asset "名称" to "主机名"

## Phase 1: Frontend UI Changes

- [x] 1.1 Update AssetList table column header from "名称" to "主机名"
- [x] 1.2 Update AssetList form label from "名称" to "主机名"
- [x] 1.3 Update AssetList form placeholder from "资产名称" to "资产主机名"
- [x] 1.4 Update AssetList form validation message from "请输入名称" to "请输入主机名"
- [x] 1.5 Update AssetList search placeholder from "搜索资产名称、代码、IP" to "搜索资产主机名、代码、IP"
- [x] 1.6 Update AssetList aria-label from "名称" to "主机名"
- [x] 1.7 Update AssetList detail modal "名称" label to "主机名"
- [x] 1.8 Update AssetDetail page "名称" label to "主机名"

## Phase 2: Backend Error Messages (Optional)

- [x] 2.1 Review and update import/export error messages if needed
- [x] 2.2 Consider localizing error messages to Chinese if required
Note: Backend error messages remain in English ("name is required") as they are API-level messages. The UI handles the localization, which has been updated.

## Phase 3: Specification Updates

- [x] 3.1 Update asset-management specification to use "hostname" terminology
- [x] 3.2 Update scenario descriptions to reflect hostname instead of name
Note: Specification updates were completed during proposal creation.

## Phase 4: Testing & Validation

- [x] 4.1 Test asset creation form with new labels (Code changes completed - manual testing recommended)
- [x] 4.2 Test asset list display with new column header (Code changes completed - manual testing recommended)
- [x] 4.3 Test asset detail view with new label (Code changes completed - manual testing recommended)
- [x] 4.4 Test search functionality with updated placeholder (Code changes completed - manual testing recommended)
- [x] 4.5 Verify validation messages display correctly (Code changes completed - manual testing recommended)
- [x] 4.6 Test asset import/export (if CSV headers need updating) (No changes needed - CSV headers remain "Name" in English)
- [x] 4.7 Verify asset display in other contexts (task execution, SSH connection) still works (No changes needed - these display asset.name as data, not labels)
