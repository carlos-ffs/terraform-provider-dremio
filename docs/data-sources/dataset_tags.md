# dremio_dataset_tags (Data Source)

Retrieves tags (labels) for an existing dataset in Dremio.

## Example Usage

```hcl
data "dremio_dataset_tags" "orders_tags" {
  dataset_id = data.dremio_table.orders.id
}

output "tags" {
  value = data.dremio_dataset_tags.orders_tags.tags
}
```

## Schema

### Required

- `dataset_id` (String) - UUID of the dataset to retrieve tags for.

### Read-Only

- `tags` (List of String) - List of tags applied to the dataset.
- `version` (String) - Version identifier for the tags.

## Notes

- Tags are case-insensitive.
- The `dataset_id` must reference an existing dataset (table, view, or UDF).
- Use this data source to check existing tags before making modifications.

## Example with View

```hcl
data "dremio_view" "sales_report" {
  path = ["Analytics", "sales_report"]
}

data "dremio_dataset_tags" "sales_report_tags" {
  dataset_id = data.dremio_view.sales_report.id
}

output "sales_report_tags" {
  value = data.dremio_dataset_tags.sales_report_tags.tags
}

# Conditionally add more tags
resource "dremio_dataset_tags" "updated_tags" {
  dataset_id = data.dremio_view.sales_report.id
  tags       = concat(
    data.dremio_dataset_tags.sales_report_tags.tags,
    ["terraform-managed"]
  )
}
```

## Example with Multiple Datasets

```hcl
locals {
  datasets = {
    orders   = data.dremio_table.orders.id
    products = data.dremio_table.products.id
    sales    = data.dremio_view.sales.id
  }
}

data "dremio_dataset_tags" "all" {
  for_each   = local.datasets
  dataset_id = each.value
}

output "all_tags" {
  value = {
    for k, v in data.dremio_dataset_tags.all : k => v.tags
  }
}
```

