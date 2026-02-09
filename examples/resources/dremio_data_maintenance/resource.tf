# =============================================================================
# Dremio Data Maintenance Resource Example (Dremio Cloud Only)
# =============================================================================
# This example is based on the working configuration from main.tf
# Note: This resource is only available for Dremio Cloud deployments.
# =============================================================================

resource "dremio_data_maintenance" "example_data_maintenance" {
  type       = "OPTIMIZE"
  is_enabled = true
  table_id   = "${join("\".\"", dremio_view.example.path)}"
}

