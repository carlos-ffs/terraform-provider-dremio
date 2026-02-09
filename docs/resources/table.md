# dremio_table (Resource)

Promotes a file or folder to a queryable table (physical dataset) in Dremio. This resource allows you to configure format options and acceleration settings for promoted tables.

## Example Usage

### Promoting a CSV File

```hcl
data "dremio_file" "csv_file" {
  path = ["Samples", "samples.dremio.com", "NYC-taxi-trips", "trips_pickupdate"]
}

resource "dremio_table" "nyc_trips" {
  path            = ["Samples", "samples.dremio.com", "NYC-taxi-trips", "trips_pickupdate"]
  file_or_folder_id = data.dremio_file.csv_file.id
  
  format = {
    type              = "Text"
    field_delimiter   = ","
    skip_first_line   = false
    extract_header    = true
    trim_header       = true
  }
}
```

## Schema

### Required

- `path` (List of String) - Full path to the table, including the source name. Each element represents a level in the hierarchy.
- `file_or_folder_id` (String) - The ID of the file or folder to promote. Use the `dremio_file` data source to look up file IDs by path.

### Optional

- `format` (Block) - Format configuration for the promoted table.
  - `type` (String) - Type of data. Valid values: `Delta`, `Excel`, `Iceberg`, `JSON`, `Parquet`, `Text`, `Unknown`, `XLS`.
  - `ignore_other_file_formats` (Boolean) - For Parquet folders, ignore non-Parquet files.
  - `skip_first_line` (Boolean) - Skip first line when creating table (Excel/Text).
  - `extract_header` (Boolean) - Extract column names from first line (Excel/Text).
  - `has_merged_cells` (Boolean) - Expand merged cells (Excel).
  - `sheet_name` (String) - Sheet name for Excel files with multiple sheets.
  - `field_delimiter` (String) - Field delimiter character (Text). Default: `,`.
  - `quote` (String) - Quote character (Text). Default: `"`.
  - `comment` (String) - Comment character (Text). Default: `#`.
  - `escape` (String) - Escape character (Text). Default: `"`.
  - `line_delimiter` (String) - Line delimiter character (Text). Default: `\n`.
  - `auto_generate_column_names` (Boolean) - Auto-generate column names (Text).
  - `trim_header` (Boolean) - Trim header whitespace (Text).

- `acceleration_refresh_policy` (Block) - Acceleration refresh policy for the table.
  - `active_policy_type` (String) - Policy for refreshing Reflections. Valid values: `NEVER`, `PERIOD`, `SCHEDULE`, `REFRESH_ON_DATA_CHANGES`.
  - `refresh_period_ms` (Number) - Refresh period in milliseconds. Minimum: 3600000 (1 hour).
  - `refresh_schedule` (String) - Cron expression for refresh schedule (UTC). Example: `0 0 8 * * ?`.
  - `grace_period_ms` (Number) - Maximum age for Reflection data in milliseconds.
  - `method` (String) - Method for refreshing Reflections. Valid values: `AUTO`, `FULL`, `INCREMENTAL`.
  - `refresh_field` (String) - Field to use for incremental refresh.
  - `never_expire` (Boolean) - Whether Reflections never expire.

- `access_control_list` (Block) - User and role access settings.
  - `users` (Block List) - List of user access controls.
    - `id` (String) - User ID.
    - `permissions` (List of String) - List of permissions.
  - `roles` (Block List) - List of role access controls.
    - `id` (String) - Role ID.
    - `permissions` (List of String) - List of permissions.

### Read-Only

- `id` (String) - Unique identifier of the table.
- `entity_type` (String) - Type of catalog object.
- `type` (String) - Dataset type (always `PHYSICAL_DATASET`).
- `tag` (String) - Version tag for optimistic concurrency control.

## Import

Tables can be imported using their ID:

```bash
terraform import dremio_table.example table-uuid-here
```

## Notes

- **File lookup**: Use the `dremio_file` data source to look up the `file_or_folder_id` by path.
- **Format detection**: If format is not specified, Dremio will attempt to auto-detect the format.
- **Deletion**: Deleting this resource will unpromote the table, reverting it to a file.
- **Text format**: Common delimiters include `,` (CSV), `\t` (TSV), and `|` (pipe-delimited).

## Parquet Table Example

```hcl
data "dremio_file" "parquet_folder" {
  path = ["Samples", "samples.dremio.com", "NYC-taxi-trips"]
}

resource "dremio_table" "parquet_table" {
  path              = ["Samples", "samples.dremio.com", "NYC-taxi-trips"]
  file_or_folder_id = data.dremio_file.parquet_folder.id
  
  format = {
    type                      = "Parquet"
    ignore_other_file_formats = true
  }
}
```

