# =============================================================================
# Dremio View Data Source Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

data "dremio_view" "datasource_view_example" {
  path = dremio_view.example.path
}

output "dremio_datasource_view_id" {
  value       = data.dremio_view.datasource_view_example.id
  description = "ID of the view from datasource"
}

output "dremio_datasource_view_sql" {
  value       = data.dremio_view.datasource_view_example.sql
  description = "SQL query of the view from datasource"
}

output "dremio_datasource_view_fields" {
  value       = data.dremio_view.datasource_view_example.fields
  description = "Fields of the view from datasource as JSON"
}

