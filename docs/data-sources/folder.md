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

- `id` (String) - UUID of the folder. Either `id` or `path` must be specified.
- `path` (List of String) - Full path to the folder. Either `id` or `path` must be specified.

### Read-Only

- `entity_type` (String) - Type of catalog object (always `folder`).
- `tag` (String) - Version tag for the folder.
- `access_control_list` (Object) - ACL settings.
  - `users` (List of Object) - User access controls.
    - `id` (String) - User ID.
    - `permissions` (List of String) - Permissions.
  - `roles` (List of Object) - Role access controls.
    - `id` (String) - Role ID.
    - `permissions` (List of String) - Permissions.

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

