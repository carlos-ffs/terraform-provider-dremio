# dremio_dataset_tags (Resource)

Manages tags (labels) for a dataset in Dremio. Tags help organize and categorize datasets, views, and other catalog objects.

## Example Usage

```hcl
resource "dremio_dataset_tags" "nyc_trips_tags" {
  dataset_id = dremio_view.nyc_trips.id
  tags       = ["production", "analytics", "taxi-data"]
}
```

## Schema

### Required

- `dataset_id` (String) - UUID of the dataset to tag.
- `tags` (List of String) - List of tags to apply to the dataset. Tags are case-insensitive and must not contain: `/`, `:`, `[`, `]`.

### Read-Only

- `version` (String) - Version identifier for optimistic concurrency control.

## Import

Dataset tags can be imported using the dataset ID:

```bash
terraform import dremio_dataset_tags.example dataset-uuid-here
```

## Notes

- **Case insensitivity**: Tags are stored and compared case-insensitively.
- **Character restrictions**: Tag names cannot contain `/`, `:`, `[`, or `]`.
- **Replaces all tags**: This resource manages all tags for a dataset. Any existing tags not in the `tags` list will be removed.
- **Dataset must exist**: The referenced dataset must exist before applying tags.
- **Version control**: The `version` attribute changes with each update and is used for optimistic concurrency.

## Example with View

```hcl
resource "dremio_view" "sales_report" {
  path = ["MySpace", "reports", "sales_report"]
  sql  = "SELECT * FROM sales"
}

resource "dremio_dataset_tags" "sales_report_tags" {
  dataset_id = dremio_view.sales_report.id
  tags       = [
    "finance",
    "quarterly-reports",
    "approved"
  ]
}

output "tags_version" {
  value = dremio_dataset_tags.sales_report_tags.version
}
```

## Example with Table

```hcl
resource "dremio_table" "orders" {
  path              = ["Samples", "samples.dremio.com", "orders"]
  file_or_folder_id = data.dremio_file.orders_csv.id
  format = {
    type = "Text"
    field_delimiter = ","
    extract_header = true
  }
}

resource "dremio_dataset_tags" "orders_tags" {
  dataset_id = dremio_table.orders.id
  tags       = ["raw-data", "orders", "e-commerce"]
}
```

## Tag Naming Best Practices

- Use lowercase for consistency
- Use hyphens or underscores for multi-word tags
- Create a consistent tagging taxonomy across your organization
- Common tag categories:
  - Environment: `production`, `staging`, `development`
  - Data domain: `sales`, `finance`, `marketing`
  - Data quality: `verified`, `raw`, `transformed`
  - Ownership: `team-analytics`, `dept-finance`

