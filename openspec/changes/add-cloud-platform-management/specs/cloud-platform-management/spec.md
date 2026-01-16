# Cloud Platform Management

## ADDED Requirements

### Requirement: Manage Cloud Platforms
The system SHALL provide cloud platform management functionality to allow administrators to create, read, update, and delete cloud platforms. Cloud platforms SHALL be standardized entries that can be associated with assets.

#### Scenario: List Cloud Platforms
- Given I am logged in as an administrator
- When I navigate to "云平台管理" (Cloud Platform Management) in the sidebar
- Then I should see a table listing all cloud platforms with columns: ID, Name, Description, and Actions
- And the table should support pagination

#### Scenario: Create Cloud Platform
- Given I am on the cloud platform management page
- When I click "新增云平台" (Create Cloud Platform) button
- Then a modal form should appear
- When I fill in the name (required) and description (optional) fields and submit
- Then a new cloud platform should be created
- And I should see a success message
- And the new cloud platform should appear in the table

#### Scenario: Edit Cloud Platform
- Given I am on the cloud platform management page
- When I click "编辑" (Edit) button for a cloud platform
- Then a modal form should appear with the current values pre-filled
- When I modify the name or description and submit
- Then the cloud platform should be updated
- And I should see a success message
- And the updated values should be reflected in the table

#### Scenario: Delete Cloud Platform
- Given I am on the cloud platform management page
- When I click "删除" (Delete) button for a cloud platform
- Then a confirmation dialog should appear
- When I confirm the deletion
- Then the cloud platform should be deleted (soft delete)
- And I should see a success message
- And the cloud platform should be removed from the table

#### Scenario: Cloud Platform Name Uniqueness
- Given I am creating or editing a cloud platform
- When I enter a name that already exists
- Then I should see an error message indicating the name already exists
- And the form should not be submitted

### Requirement: Associate Assets with Cloud Platforms
Assets MUST be able to be associated with cloud platforms through a foreign key relationship. When creating or editing assets, users SHALL select cloud platforms from a dropdown list instead of entering free-text values.

#### Scenario: Select Cloud Platform in Asset Form
- Given I am creating or editing an asset
- When I view the "云平台" (Cloud Platform) field
- Then I should see a dropdown select field instead of a text input
- And the dropdown should list all available cloud platforms
- And I should be able to select a cloud platform or leave it empty
- When I submit the asset form
- Then the asset should be associated with the selected cloud platform

#### Scenario: Display Cloud Platform in Asset List
- Given I am viewing the asset list
- When the asset has an associated cloud platform
- Then the cloud platform name should be displayed in the "云平台" column
- And it should be displayed as a tag or badge

#### Scenario: Filter Assets by Cloud Platform
- Given I am viewing the asset list
- When I use the filter/search functionality
- Then I should be able to filter assets by cloud platform
- And the filter should use cloud platform ID matching

## MODIFIED Requirements

### Requirement: Asset Management - Cloud Platform Field
The cloud platform field in asset management MUST use a foreign key relationship (`cloud_platform_id`) instead of a free-text string field (`cloud_platform`). The system SHALL migrate existing string values to cloud platform entries and update asset records to reference them.

#### Scenario: Asset Creation with Cloud Platform
- Given I am creating a new asset
- When I fill in the asset form
- Then the "云平台" field should be a dropdown showing all available cloud platforms
- And I should be able to select from existing cloud platforms
- And the selected cloud platform ID should be saved with the asset

#### Scenario: Asset Display with Cloud Platform
- Given an asset is associated with a cloud platform
- When I view the asset in the list or detail page
- Then the cloud platform name should be displayed (not an ID or code)
