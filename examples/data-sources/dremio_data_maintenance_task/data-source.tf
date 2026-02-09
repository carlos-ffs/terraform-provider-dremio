# =============================================================================
# Dremio Data Maintenance Task Data Source Example (Dremio Cloud Only)
# =============================================================================
# This example is based on the working configuration from main.tf
# Note: This resource is only available for Dremio Cloud deployments.
# =============================================================================

data "dremio_data_maintenance_task" "example" {
  id = dremio_data_maintenance.example_data_maintenance.id
}

output "data_maintenance_task_is_enabled" {
  value       = data.dremio_data_maintenance_task.example.is_enabled
  description = "Is enabled"
}

