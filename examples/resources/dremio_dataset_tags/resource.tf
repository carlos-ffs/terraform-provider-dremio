# =============================================================================
# Dremio Dataset Tags Resource Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

resource "dremio_dataset_tags" "dataset_tags_example" {
  dataset_id = dremio_table.resource_table_example.id
  tags       = ["carlos-santos", "terraform", "SRE"]
}

output "dremio_dataset_tags_resources_version" {
  value       = dremio_dataset_tags.dataset_tags_example.version
  description = "Version of the dataset tags"
}

