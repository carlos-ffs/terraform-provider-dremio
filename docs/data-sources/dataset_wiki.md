# dremio_dataset_wiki (Data Source)

Retrieves wiki documentation for an existing dataset in Dremio.

## Example Usage

```hcl
data "dremio_dataset_wiki" "orders_wiki" {
  dataset_id = data.dremio_table.orders.id
}

output "wiki_content" {
  value = data.dremio_dataset_wiki.orders_wiki.text
}
```

## Schema

### Required

- `dataset_id` (String) - UUID of the dataset to retrieve wiki for.

### Read-Only

- `text` (String) - Wiki content in GitHub-flavored Markdown.
- `version` (String) - Version identifier for the wiki.

## Notes

- The `dataset_id` must reference an existing dataset (table, view, or UDF).
- The `text` attribute contains Markdown content that can be rendered.
- If no wiki exists for the dataset, `text` will be empty.

## Example with View

```hcl
data "dremio_view" "customer_360" {
  path = ["Analytics", "customer_360"]
}

data "dremio_dataset_wiki" "customer_360_wiki" {
  dataset_id = data.dremio_view.customer_360.id
}

output "has_documentation" {
  value = length(data.dremio_dataset_wiki.customer_360_wiki.text) > 0
}

output "wiki_version" {
  value = data.dremio_dataset_wiki.customer_360_wiki.version
}
```

## Example: Append to Existing Wiki

```hcl
data "dremio_view" "report" {
  path = ["Reports", "daily_report"]
}

data "dremio_dataset_wiki" "existing" {
  dataset_id = data.dremio_view.report.id
}

resource "dremio_dataset_wiki" "updated" {
  dataset_id = data.dremio_view.report.id
  text       = <<-EOT
    ${data.dremio_dataset_wiki.existing.text}

    ---

    ## Terraform Managed
    Last updated by Terraform.
  EOT
}
```

## Example: Check Multiple Datasets

```hcl
locals {
  datasets = {
    orders   = data.dremio_table.orders.id
    products = data.dremio_table.products.id
  }
}

data "dremio_dataset_wiki" "all" {
  for_each   = local.datasets
  dataset_id = each.value
}

output "documented_datasets" {
  value = [
    for k, v in data.dremio_dataset_wiki.all : k
    if length(v.text) > 0
  ]
}

output "undocumented_datasets" {
  value = [
    for k, v in data.dremio_dataset_wiki.all : k
    if length(v.text) == 0
  ]
}
```

