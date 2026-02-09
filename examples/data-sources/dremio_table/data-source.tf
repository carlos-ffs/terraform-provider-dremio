# =============================================================================
# Dremio Table Data Source Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

data "dremio_table" "datasource_table_example" {
  id = dremio_table.resource_table_example.id
}

output "dremio_table_id" {
  value       = data.dremio_table.datasource_table_example.id
  description = "ID of the table"
}

output "dremio_table_tag" {
  value       = data.dremio_table.datasource_table_example.tag
  description = "Tag of the table"
}

output "dremio_table_format" {
  value       = data.dremio_table.datasource_table_example.format
  description = "Format of the table"
}

