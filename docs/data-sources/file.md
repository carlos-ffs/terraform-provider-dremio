# dremio_file (Data Source)

Retrieves information about a file in Dremio. This data source is commonly used to look up file IDs for table promotion.

## Example Usage

### By Path

```hcl
data "dremio_file" "trips_csv" {
  path = ["Samples", "samples.dremio.com", "NYC-taxi-trips", "trips_pickupdate"]
}

output "file_id" {
  value = data.dremio_file.trips_csv.id
}
```

### By ID

```hcl
data "dremio_file" "by_id" {
  id = "file-uuid-here"
}
```

## Schema

### Optional (One Required)

- `id` (String) - UUID of the file. Either `id` or `path` must be specified.
- `path` (List of String) - Full path to the file. Either `id` or `path` must be specified.

### Read-Only

- `entity_type` (String) - Type of catalog object (always `file`).
- `tag` (String) - Version tag for the file.

## Notes

- Specify either `id` or `path`, but not both.
- Files represent unpromoted data files within sources.
- Use this data source to get file IDs for the `dremio_table` resource.

## Example with Table Promotion

```hcl
# Look up the file by path
data "dremio_file" "csv_data" {
  path = ["Samples", "samples.dremio.com", "NYC-taxi-trips", "trips_pickupdate"]
}

# Promote the file to a table
resource "dremio_table" "trips" {
  path              = ["Samples", "samples.dremio.com", "NYC-taxi-trips", "trips_pickupdate"]
  file_or_folder_id = data.dremio_file.csv_data.id
  
  format = {
    type              = "Text"
    field_delimiter   = ","
    extract_header    = true
    trim_header       = true
  }
}

output "table_id" {
  value = dremio_table.trips.id
}
```

## Looking Up Folder IDs

For promoting folders (e.g., Parquet folders), use the same data source:

```hcl
data "dremio_file" "parquet_folder" {
  path = ["Samples", "samples.dremio.com", "NYC-taxi-trips"]
}

resource "dremio_table" "parquet_data" {
  path              = ["Samples", "samples.dremio.com", "NYC-taxi-trips"]
  file_or_folder_id = data.dremio_file.parquet_folder.id
  
  format = {
    type                      = "Parquet"
    ignore_other_file_formats = true
  }
}
```

