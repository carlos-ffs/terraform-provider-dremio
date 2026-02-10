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

| Attribute | Type | Description |
|-----------|------|-------------|
| `path` | List of String | Full path to the table, including the source name. Each element represents a level in the hierarchy. Path elements must not contain: `/`, `:`, `[`, `]`. |
| `file_or_folder_id` | String | Unique identifier of the source file or folder to format as a table. Use the `dremio_file` data source to look up file IDs by path. |
| `format` | Block | Format configuration for the promoted table. See format block below. |

### Optional

#### format (Block) - Required

Defines the file format configuration when promoting files to tables.

| Attribute | Type | Default | Description |
|-----------|------|---------|-------------|
| `type` | String | **Required** | Type of data. Valid values: `Delta`, `Excel`, `Iceberg`, `JSON`, `Parquet`, `Text`, `Unknown`, `XLS`. |
| `ignore_other_file_formats` | Boolean | `false` | For Parquet folders, ignore non-Parquet files. |
| `skip_first_line` | Boolean | `false` | Skip first line when creating table (Excel/Text). |
| `extract_header` | Boolean | `false` | Extract column names from first line (Excel/Text). |
| `has_merged_cells` | Boolean | `false` | Expand merged cells (Excel). |
| `sheet_name` | String | `null` | Sheet name for Excel files with multiple sheets. |
| `field_delimiter` | String | `,` | Field delimiter character (Text). Common values: `,` (CSV), `\t` (TSV), `\|` (pipe). |
| `quote` | String | `"` | Quote character (Text). |
| `comment` | String | `#` | Comment character (Text). |
| `escape` | String | `"` | Escape character (Text). |
| `line_delimiter` | String | `\r\n` | Line delimiter (Text). |
| `auto_generate_column_names` | Boolean | `false` | Auto-generate column names if no header (Text). |
| `trim_header` | Boolean | `false` | Trim whitespace from column names (Text). |

#### acceleration_refresh_policy (Block)

Defines the acceleration (Reflection) refresh policy for the table.

| Attribute | Type | Default | Description |
|-----------|------|---------|-------------|
| `active_policy_type` | String | `PERIOD` | Policy for refreshing Reflections. Valid values: `NEVER`, `PERIOD`, `SCHEDULE`, `REFRESH_ON_DATA_CHANGES`. |
| `refresh_period_ms` | Number | `3600000` | Refresh period in milliseconds. Minimum: 3600000 (1 hour). |
| `refresh_schedule` | String | `null` | Cron expression for refresh schedule (UTC). Example: `0 0 8 * * ?` (daily at 8 AM UTC). |
| `grace_period_ms` | Number | `null` | Maximum age for Reflection data in milliseconds before it's considered stale. |
| `method` | String | `AUTO` | Method for refreshing Reflections. Valid values: `AUTO` (Dremio decides), `FULL` (complete refresh), `INCREMENTAL` (only new data). |
| `refresh_field` | String | `null` | Field to use for incremental refresh. Required when `method` is `INCREMENTAL`. |
| `never_expire` | Boolean | `false` | Whether Reflections never expire. |

#### access_control_list (Block)

User and role access settings.

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
| `id` | String | Unique identifier of the table (UUID). |
| `entity_type` | String | Type of catalog object (always `dataset`). |
| `type` | String | Dataset type (always `PHYSICAL_DATASET`). |
| `tag` | String | Version tag for optimistic concurrency control. This value changes with every update. |

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

