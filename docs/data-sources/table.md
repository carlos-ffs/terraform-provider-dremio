# dremio_table (Data Source)

Retrieves information about an existing promoted table (physical dataset) in Dremio.

## Example Usage

### By Path

```hcl
data "dremio_table" "trips" {
  path = ["Samples", "samples.dremio.com", "NYC-taxi-trips"]
}

output "table_id" {
  value = data.dremio_table.trips.id
}

output "table_format" {
  value = data.dremio_table.trips.format
}
```

### By ID

```hcl
data "dremio_table" "by_id" {
  id = "table-uuid-here"
}
```

## Schema

### Optional (One Required)

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | UUID of the table. Either `id` or `path` must be specified. |
| `path` | List of String | Full path to the table, including the source name. Either `id` or `path` must be specified. |

### Read-Only

| Attribute | Type | Description |
|-----------|------|-------------|
| `entity_type` | String | Type of catalog object (always `dataset`). |
| `type` | String | Dataset type (always `PHYSICAL_DATASET`). |
| `tag` | String | Version tag for optimistic concurrency control. |

#### format (Object)

Format configuration of the promoted table.

| Attribute | Type | Description |
|-----------|------|-------------|
| `type` | String | Format type. Values: `Text`, `Parquet`, `JSON`, `Delta`, `Iceberg`, `Excel`, `XLS`, `Unknown`. |
| `field_delimiter` | String | Field delimiter character (Text). |
| `quote` | String | Quote character (Text). |
| `comment` | String | Comment character (Text). |
| `escape` | String | Escape character (Text). |
| `line_delimiter` | String | Line delimiter (Text). |
| `skip_first_line` | Boolean | Skip first line when reading (Excel/Text). |
| `extract_header` | Boolean | Extract column names from first line (Excel/Text). |
| `trim_header` | Boolean | Trim whitespace from column names (Text). |
| `ignore_other_file_formats` | Boolean | For Parquet folders, ignore non-Parquet files. |
| `sheet_name` | String | Sheet name for Excel files with multiple sheets. |
| `has_merged_cells` | Boolean | Expand merged cells (Excel). |
| `auto_generate_column_names` | Boolean | Auto-generate column names if no header (Text). |

#### acceleration_refresh_policy (Object)

Acceleration (Reflection) refresh policy for the table.

| Attribute | Type | Description |
|-----------|------|-------------|
| `active_policy_type` | String | Policy type. Values: `NEVER`, `PERIOD`, `SCHEDULE`, `REFRESH_ON_DATA_CHANGES`. |
| `refresh_period_ms` | Number | Refresh period in milliseconds. |
| `refresh_schedule` | String | Cron expression for refresh schedule (UTC). |
| `grace_period_ms` | Number | Maximum age for Reflection data in milliseconds. |
| `method` | String | Refresh method. Values: `AUTO`, `FULL`, `INCREMENTAL`. |
| `refresh_field` | String | Field to use for incremental refresh. |
| `never_expire` | Boolean | Whether Reflections never expire. |

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
- Only promoted tables can be retrieved; unpromoted files use `dremio_file`.
- Format settings reflect how the table was configured during promotion.

## Example with Grants

```hcl
data "dremio_table" "orders" {
  path = ["Samples", "samples.dremio.com", "orders"]
}

resource "dremio_grants" "orders_access" {
  catalog_object_id = data.dremio_table.orders.id
  grants = [
    {
      id           = "analysts-role-uuid"
      grantee_type = "ROLE"
      privileges   = ["SELECT"]
    }
  ]
}

resource "dremio_dataset_tags" "orders_tags" {
  dataset_id = data.dremio_table.orders.id
  tags       = ["production", "sales"]
}
```

