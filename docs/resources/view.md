# dremio_view (Resource)

Creates and manages a virtual dataset (view) in Dremio. Views are SQL-based abstractions that allow you to define reusable queries without materializing data.

## Example Usage

```hcl
resource "dremio_view" "nyc_trips" {
  path = [
    "Samples",
    "samples.dremio.com",
    "terraform_top_folder",
    "terraform_nested_folder",
    "NYC-taxi-trips_view"
  ]
  sql = "SELECT * FROM \"Samples\".\"samples.dremio.com\".\"NYC-taxi-trips\""
  sql_context = ["Samples", "samples.dremio.com"]
}
```

## Schema

### Required

- `path` (List of String) - Full path to the view, including the source/space name and folder hierarchy. The last element is the view name. Path elements must not contain: `/`, `:`, `[`, `]`.
- `sql` (String) - SQL query defining the view.

### Optional

- `sql_context` (List of String) - Default schema context for the SQL query. Objects referenced without full paths are resolved relative to this context.

- `access_control_list` (Block) - User and role access settings.
  - `users` (Block List) - List of user access controls.
    - `id` (String) - User ID.
    - `permissions` (List of String) - List of permissions.
  - `roles` (Block List) - List of role access controls.
    - `id` (String) - Role ID.
    - `permissions` (List of String) - List of permissions.

### Read-Only

- `id` (String) - Unique identifier of the view.
- `entity_type` (String) - Type of catalog object (always `dataset`).
- `type` (String) - Dataset type (always `VIRTUAL_DATASET`).
- `tag` (String) - Version tag for optimistic concurrency control.
- `fields` (String, JSON) - JSON representation of the view's field schema, including column names and data types.

## Import

Views can be imported using their ID:

```bash
terraform import dremio_view.example view-uuid-here
```

## Notes

- **Path structure**: The path includes the full hierarchy from source/space to the view name.
- **SQL context**: Using `sql_context` simplifies SQL queries by establishing a default schema.
- **Parent folders must exist**: Ensure all parent folders exist before creating the view.
- **Fields are computed**: The `fields` attribute is populated after creation based on the SQL query's output schema.

## Example with Dependencies

```hcl
resource "dremio_folder" "analytics" {
  path = ["Samples", "samples.dremio.com", "analytics"]
}

resource "dremio_view" "daily_summary" {
  path = [
    "Samples",
    "samples.dremio.com",
    "analytics",
    "daily_summary"
  ]
  sql = <<-EOT
    SELECT 
      pickup_date,
      COUNT(*) as trip_count,
      AVG(fare_amount) as avg_fare
    FROM "NYC-taxi-trips"
    GROUP BY pickup_date
  EOT
  sql_context = ["Samples", "samples.dremio.com"]
  
  depends_on = [dremio_folder.analytics]
}

output "view_id" {
  value = dremio_view.daily_summary.id
}

output "view_fields" {
  value = dremio_view.daily_summary.fields
}
```

## SQL Context Example

The `sql_context` attribute allows you to write simpler SQL:

```hcl
# Without sql_context - must use fully qualified names
resource "dremio_view" "without_context" {
  path = ["MySpace", "my_view"]
  sql  = "SELECT * FROM \"Samples\".\"samples.dremio.com\".\"NYC-taxi-trips\""
}

# With sql_context - can use relative names
resource "dremio_view" "with_context" {
  path        = ["MySpace", "my_view"]
  sql         = "SELECT * FROM \"NYC-taxi-trips\""
  sql_context = ["Samples", "samples.dremio.com"]
}
```

