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

- `id` (String) - UUID of the table. Either `id` or `path` must be specified.
- `path` (List of String) - Full path to the table. Either `id` or `path` must be specified.

### Read-Only

- `entity_type` (String) - Type of catalog object.
- `type` (String) - Dataset type (always `PHYSICAL_DATASET`).
- `tag` (String) - Version tag for the table.

- `format` (Object) - Format configuration of the table.
  - `type` (String) - Format type (Text, Parquet, JSON, etc.).
  - `field_delimiter` (String) - Field delimiter (Text).
  - `skip_first_line` (Boolean) - Skip first line.
  - `extract_header` (Boolean) - Extract header.
  - `trim_header` (Boolean) - Trim header whitespace.
  - `quote` (String) - Quote character.
  - `comment` (String) - Comment character.
  - `escape` (String) - Escape character.
  - `line_delimiter` (String) - Line delimiter.
  - `ignore_other_file_formats` (Boolean) - Ignore non-matching files.
  - `sheet_name` (String) - Excel sheet name.
  - `has_merged_cells` (Boolean) - Has merged cells (Excel).
  - `auto_generate_column_names` (Boolean) - Auto-generate column names.

- `acceleration_refresh_policy` (Object) - Acceleration settings.
  - `active_policy_type` (String) - Policy type.
  - `refresh_period_ms` (Number) - Refresh period.
  - `refresh_schedule` (String) - Refresh schedule (cron).
  - `grace_period_ms` (Number) - Grace period.
  - `method` (String) - Refresh method.
  - `refresh_field` (String) - Incremental refresh field.
  - `never_expire` (Boolean) - Never expire.

- `access_control_list` (Object) - ACL settings.

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

