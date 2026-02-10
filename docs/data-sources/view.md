# dremio_view (Data Source)

Retrieves information about an existing view (virtual dataset) in Dremio.

## Example Usage

### By Path

```hcl
data "dremio_view" "sales_summary" {
  path = ["Analytics", "reports", "sales_summary"]
}

output "view_id" {
  value = data.dremio_view.sales_summary.id
}

output "view_sql" {
  value = data.dremio_view.sales_summary.sql
}
```

### By ID

```hcl
data "dremio_view" "by_id" {
  id = "view-uuid-here"
}
```

## Schema

### Optional (One Required)

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | UUID of the view. Either `id` or `path` must be specified. |
| `path` | List of String | Full path to the view, including the source/space name. Either `id` or `path` must be specified. |

### Read-Only

| Attribute | Type | Description |
|-----------|------|-------------|
| `entity_type` | String | Type of catalog object (always `dataset`). |
| `type` | String | Dataset type (always `VIRTUAL_DATASET`). |
| `sql` | String | SQL query defining the view. |
| `sql_context` | List of String | Default schema context for the SQL query. |
| `tag` | String | Version tag for optimistic concurrency control. |
| `fields` | String (JSON) | JSON representation of the view's field schema, including column names and data types. Use `jsondecode()` to parse this value. |

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
- The `sql` attribute contains the complete view definition.
- The `fields` attribute provides schema information as JSON.

## Example with Dependencies

```hcl
# Reference an existing view
data "dremio_view" "base_view" {
  path = ["Analytics", "base_metrics"]
}

# Create a new view based on the existing one
resource "dremio_view" "extended_metrics" {
  path = ["Analytics", "extended_metrics"]
  sql  = "SELECT *, CURRENT_TIMESTAMP as refresh_time FROM ${join(".", formatlist("\"%s\"", data.dremio_view.base_view.path))}"
  sql_context = data.dremio_view.base_view.sql_context
}

# Add tags to existing view
resource "dremio_dataset_tags" "base_view_tags" {
  dataset_id = data.dremio_view.base_view.id
  tags       = ["core", "metrics"]
}
```

## Example with Grants

```hcl
data "dremio_view" "customer_report" {
  path = ["Reports", "customer_360"]
}

resource "dremio_grants" "report_access" {
  catalog_object_id = data.dremio_view.customer_report.id
  grants = [
    {
      id           = "sales-team-role-uuid"
      grantee_type = "ROLE"
      privileges   = ["SELECT"]
    },
    {
      id           = "admin-role-uuid"
      grantee_type = "ROLE"
      privileges   = ["SELECT", "ALTER", "DROP"]
    }
  ]
}
```

