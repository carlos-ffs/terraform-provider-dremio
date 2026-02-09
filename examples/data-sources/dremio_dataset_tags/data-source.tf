# =============================================================================
# Dremio Dataset Tags Data Source Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

data "dremio_dataset_tags" "datasource_dataset_tags_example" {
  dataset_id = dremio_table.resource_table_example.id
}

output "dremio_datasource_dataset_tags" {
  value       = data.dremio_dataset_tags.datasource_dataset_tags_example.tags
  description = "Tags of the dataset"
}

