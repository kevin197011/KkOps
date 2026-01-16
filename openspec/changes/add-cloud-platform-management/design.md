# Design: Cloud Platform Management

## Architecture Overview

The cloud platform management feature follows the same pattern as Environment and Project management:

1. **Backend**: Model → Service → Handler → Router
2. **Frontend**: API Client → Page Component → Route → Menu Item

## Data Model

### CloudPlatform Model

```go
type CloudPlatform struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"not null;size:100;uniqueIndex" json:"name"`
    Description string         `gorm:"type:text" json:"description"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

    // Relationships
    Assets []Asset `gorm:"foreignKey:CloudPlatformID" json:"assets,omitempty"`
}
```

### Asset Model Changes

**Before**:
```go
CloudPlatform string `gorm:"size:50" json:"cloud_platform"`
```

**After**:
```go
CloudPlatformID *uint          `json:"cloud_platform_id"`
CloudPlatform   *CloudPlatform `gorm:"foreignKey:CloudPlatformID" json:"cloud_platform,omitempty"`
```

## Database Migration

### Step 1: Create cloud_platforms table

```sql
CREATE TABLE cloud_platforms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_cloud_platforms_name ON cloud_platforms(name) WHERE deleted_at IS NULL;
```

### Step 2: Migrate existing data

```sql
-- Extract unique cloud platform values and create entries
INSERT INTO cloud_platforms (name, created_at, updated_at)
SELECT DISTINCT cloud_platform, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
FROM assets
WHERE cloud_platform IS NOT NULL AND cloud_platform != '';
```

### Step 3: Add foreign key column

```sql
ALTER TABLE assets ADD COLUMN cloud_platform_id INTEGER REFERENCES cloud_platforms(id);
```

### Step 4: Map existing data

```sql
UPDATE assets
SET cloud_platform_id = (
    SELECT id FROM cloud_platforms 
    WHERE cloud_platforms.name = assets.cloud_platform
    LIMIT 1
)
WHERE cloud_platform IS NOT NULL AND cloud_platform != '';
```

### Step 5: Drop old column

```sql
ALTER TABLE assets DROP COLUMN cloud_platform;
```

## API Design

### Endpoints

- `GET /cloud-platforms` - List all cloud platforms
- `GET /cloud-platforms/:id` - Get cloud platform details
- `POST /cloud-platforms` - Create cloud platform
- `PUT /cloud-platforms/:id` - Update cloud platform
- `DELETE /cloud-platforms/:id` - Delete cloud platform

### Request/Response Formats

```typescript
// Create
POST /cloud-platforms
{
  "name": "阿里云",
  "description": "Alibaba Cloud"
}

// Response
{
  "id": 1,
  "name": "阿里云",
  "description": "Alibaba Cloud",
  "created_at": "2025-01-27T10:00:00Z",
  "updated_at": "2025-01-27T10:00:00Z"
}
```

## Frontend Design

### Page Structure

Similar to `EnvironmentList.tsx`:
- Table view with columns: ID, Name, Description, Actions
- Create/Edit modal form
- Delete confirmation dialog

### Asset Form Changes

**Before**:
```tsx
<Form.Item name="cloud_platform" label="云平台">
  <Input placeholder="如：阿里云、腾讯云、AWS" />
</Form.Item>
```

**After**:
```tsx
<Form.Item name="cloud_platform_id" label="云平台">
  <Select placeholder="选择云平台" allowClear>
    {cloudPlatforms.map((platform) => (
      <Select.Option key={platform.id} value={platform.id}>
        {platform.name}
      </Select.Option>
    ))}
  </Select>
</Form.Item>
```

### Menu Item

Add to `MainLayout.tsx` menu items:
```tsx
{
  key: '/cloud-platforms',
  icon: <CloudOutlined />,
  label: '云平台管理',
}
```

Place it after "环境管理" for logical grouping.

## Asset API Changes

### Asset Response

**Before**:
```typescript
{
  cloud_platform: string
}
```

**After**:
```typescript
{
  cloud_platform_id?: number,
  cloud_platform?: {
    id: number,
    name: string,
    description: string
  }
}
```

### Asset Request

**Before**:
```typescript
{
  cloud_platform?: string
}
```

**After**:
```typescript
{
  cloud_platform_id?: number
}
```

## Data Migration Strategy

1. **Extract unique values**: Query all distinct `cloud_platform` values from assets
2. **Create entries**: For each unique value, create a cloud platform entry
3. **Map relationships**: Update assets to reference cloud platform IDs
4. **Handle nulls**: Assets with null/empty cloud_platform remain with null `cloud_platform_id`

## Default Cloud Platforms

Optionally pre-populate common platforms:
- 阿里云 (Alibaba Cloud)
- 腾讯云 (Tencent Cloud)
- 华为云 (Huawei Cloud)
- AWS
- Azure
- 本地/私有云 (On-Premise)
