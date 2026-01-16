# Change: Change Asset Detail View to Modal

## Why

Currently, when users click "查看" (View) button in the asset list, they are navigated to a separate detail page. This breaks the user's context and requires navigation back to the list. Changing to a modal dialog provides a better user experience by:
- Keeping users in context of the list page
- Faster access to asset details
- No page navigation required
- Better for quick information lookups

## What Changes

- **MODIFIED**: Asset list page "查看" button text changed to "详情"
- **MODIFIED**: Asset detail view changed from page navigation to modal dialog
- **MODIFIED**: Detail information displayed in a Modal using Descriptions component
- **REMOVED**: Navigation to asset detail page from list (but detail page route remains for direct access)

## Impact

- **UI/UX**: Improved user experience with in-page detail viewing
- **Frontend**: AssetList component updated to include detail modal
- **No Breaking Changes**: Asset detail page route remains available for direct URL access
- **Affected Specs**: `asset-management` specification may need update to reflect modal-based detail view