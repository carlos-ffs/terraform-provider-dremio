# dremio_folder (Resource)

Manages a folder in Dremio. Folders are used to organize datasets, views, and other resources within sources and spaces.

## Example Usage

### Top-Level Folder

```hcl
resource "dremio_folder" "top_level" {
  path = ["Samples", "samples.dremio.com", "terraform_top_folder"]
}
```

### Nested Folder

```hcl
resource "dremio_folder" "nested" {
  path = [
    "Samples",
    "samples.dremio.com",
    "terraform_top_folder",
    "terraform_nested_folder"
  ]
  depends_on = [dremio_folder.top_level]
}
```

## Schema

### Required

| Attribute | Type | Description |
|-----------|------|-------------|
| `path` | List of String | Full path to the folder, including the source or space name. Each element represents a level in the hierarchy. Path elements must not contain: `/`, `:`, `[`, `]`. |

### Optional

#### access_control_list (Block)

User and role access settings. Can only be set via update (not on initial creation).

**users** (List of Object):

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | Yes | UUID of the user. |
| `permissions` | List of String | Yes | List of permissions to grant. |

**roles** (List of Object):

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | Yes | UUID of the role. |
| `permissions` | List of String | Yes | List of permissions to grant. |

### Read-Only

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | Unique identifier of the folder (UUID). |
| `entity_type` | String | Type of catalog object (always `folder`). |
| `tag` | String | Version tag for optimistic concurrency control. This value changes with every update. |

## Import

Folders can be imported using their ID:

```bash
terraform import dremio_folder.example folder-uuid-here
```

## Notes

- **Path requires replacement**: Changing the `path` attribute will force recreation of the folder.
- **Parent folders must exist**: Ensure parent folders exist before creating nested folders. Use `depends_on` to enforce ordering.
- **Path validation**: Path elements cannot contain `/`, `:`, `[`, or `]` characters.
- **Access control**: ACLs can only be set after the folder is created (in an update operation).
- **Deletion**: Deleting a folder will fail if it contains child objects. Delete children first.

## Path Structure

The path is an ordered list of strings representing the hierarchy:

```hcl
path = [
  "SourceName",       # First element: source or space name
  "TopLevelFolder",   # Second element: first folder level
  "SubFolder",        # Third element: nested folder
  "DeepFolder"        # Additional nesting as needed
]
```

## Example with Outputs

```hcl
resource "dremio_folder" "example" {
  path = ["Samples", "samples.dremio.com", "my_folder"]
}

output "folder_id" {
  value = dremio_folder.example.id
}

output "folder_tag" {
  value = dremio_folder.example.tag
}
```

## Using with Other Resources

Folders are commonly used as locations for views and UDFs:

```hcl
resource "dremio_folder" "analytics" {
  path = ["Samples", "samples.dremio.com", "analytics"]
}

resource "dremio_view" "sales_summary" {
  path = [
    "Samples",
    "samples.dremio.com",
    "analytics",
    "sales_summary"
  ]
  sql = "SELECT * FROM my_table"
  depends_on = [dremio_folder.analytics]
}
```

