# dremio_folder (Data Source)

Retrieves information about an existing folder in Dremio.

## Example Usage

### By Path

```hcl
data "dremio_folder" "analytics" {
  path = ["Samples", "samples.dremio.com", "analytics"]
}

output "folder_id" {
  value = data.dremio_folder.analytics.id
}
```

### By ID

```hcl
data "dremio_folder" "by_id" {
  id = "folder-uuid-here"
}
```

## Schema

### Optional (One Required)

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | UUID of the folder. Either `id` or `path` must be specified. |
| `path` | List of String | Full path to the folder, including the source/space name. Either `id` or `path` must be specified. |

### Read-Only

| Attribute | Type | Description |
|-----------|------|-------------|
| `entity_type` | String | Type of catalog object (always `folder`). |
| `tag` | String | Version tag for optimistic concurrency control. |

#### access_control_list (Object)

User and role access settings.

**users** (List of Object):

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | UUID of the user. |
| `permissions` | List of String | List of permissions granted. |

**roles** (List of Object):

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | UUID of the role. |
| `permissions` | List of String | List of permissions granted. |

## Notes

- Specify either `id` or `path`, but not both.
- The path is an ordered list from source/space to the folder name.
- Use this data source to reference existing folders when creating child resources.

## Example with Dependent Resource

```hcl
data "dremio_folder" "existing" {
  path = ["Samples", "samples.dremio.com", "existing_folder"]
}

resource "dremio_folder" "child" {
  path = [
    "Samples",
    "samples.dremio.com",
    "existing_folder",
    "new_child_folder"
  ]
}

resource "dremio_grants" "folder_access" {
  catalog_object_id = data.dremio_folder.existing.id
  grants = [
    {
      id           = "role-uuid"
      grantee_type = "ROLE"
      privileges   = ["SELECT"]
    }
  ]
}
```

