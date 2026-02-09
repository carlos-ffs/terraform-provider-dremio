# =============================================================================
# Dremio View Resource Example
# =============================================================================
# This example is based on the working configuration from main.tf
# Usage: SELECT * FROM carlos_santos."test-folder3".high_passenger_trips;
# =============================================================================

locals {
  # NYC-taxi-trips.csv
  last_element                    = dremio_table.resource_table_example.path[length(dremio_table.resource_table_example.path) - 1]
  table_path_without_last_element = slice(dremio_table.resource_table_example.path, 0, length(dremio_table.resource_table_example.path) - 1)
}

resource "dremio_view" "example" {
  path = concat(dremio_folder.resource_folder_example.path, ["high_passenger_trips"])

  sql = "SELECT passenger_count, trip_distance_mi FROM \"${local.last_element}\" WHERE passenger_count > 5"

  sql_context = local.table_path_without_last_element

  #   access_control_list = {
  #     users = [
  #       {
  #         id          = "user-id-123"
  #         permissions = ["VIEW", "MODIFY"]
  #       }
  #     ]
  #   }

  depends_on = [dremio_folder.resource_folder_example, dremio_table.resource_table_example]
}

output "dremio_view_id" {
  value       = dremio_view.example.id
  description = "ID of the view"
}

output "dremio_view_fields" {
  value       = jsondecode(dremio_view.example.fields)
  description = "Fields of the view as JSON"
}

