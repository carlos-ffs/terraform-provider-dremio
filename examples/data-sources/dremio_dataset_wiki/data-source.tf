# =============================================================================
# Dremio Dataset Wiki Data Source Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

data "dremio_dataset_wiki" "datasource_dataset_wiki_example" {
  dataset_id = dremio_table.resource_table_example.id
}

output "data_dataset_wiki_text" {
  value       = data.dremio_dataset_wiki.datasource_dataset_wiki_example.text
  description = "Text of the dataset wiki from datasource"
}

output "data_dataset_wiki_version" {
  value       = data.dremio_dataset_wiki.datasource_dataset_wiki_example.version
  description = "Version of the dataset wiki from datasource"
}

