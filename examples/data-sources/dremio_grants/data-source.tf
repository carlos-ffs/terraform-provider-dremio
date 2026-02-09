# =============================================================================
# Dremio Grants Data Source Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

data "dremio_grants" "datasource_grants_example" {
  catalog_object_id = dremio_table.resource_table_example.id
  depends_on        = [dremio_grants.example_grants]
}

output "resource_table_example_grants" {
  value       = data.dremio_grants.datasource_grants_example.grants
  description = "Grants of the table"
}

