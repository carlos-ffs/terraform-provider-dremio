# =============================================================================
# Dremio Engine Data Source Example (Dremio Cloud Only)
# =============================================================================
# This example is based on the working configuration from main.tf
# Note: This resource is only available for Dremio Cloud deployments.
# =============================================================================

data "dremio_engine" "datasource_engine_example" {
  name = "preview"
}

output "dremio_engine_id" {
  value       = data.dremio_engine.datasource_engine_example.id
  description = "ID of the engine"
}

