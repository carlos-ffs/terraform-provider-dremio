# =============================================================================
# Dremio Source Data Source Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

data "dremio_source" "datasource_dremio_source_example" {
  name = dremio_source.samples_bucket.name
}

output "datasource_dremio_source_example_id" {
  value       = data.dremio_source.datasource_dremio_source_example.id
  description = "ID of the source"
}

output "datasource_dremio_source_example_metadataPolicy" {
  value       = data.dremio_source.datasource_dremio_source_example.metadata_policy
  description = "metadataPolicy of the source"
}

