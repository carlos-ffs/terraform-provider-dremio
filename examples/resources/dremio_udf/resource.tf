# =============================================================================
# Dremio UDF (User-Defined Function) Resource Example
# =============================================================================
# This example is based on the working configuration from main.tf
# Usage:
#   DESCRIBE FUNCTION carlos_santos."test-folder3".count_high_passenger_trips;
#   SELECT carlos_santos."test-folder3".count_high_passenger_trips(5);
# =============================================================================

resource "dremio_udf" "example" {
  path = concat(dremio_folder.resource_folder_example.path, ["count_high_passenger_trips"])

  is_scalar         = true
  function_arg_list = "min_passengers BIGINT"
  function_body     = "SELECT count(*) FROM \"${join("\".\"", dremio_table.resource_table_example.path)}\" WHERE passenger_count > min_passengers"
  return_type       = "BIGINT"

  #   access_control_list = {
  #     users = [
  #       {
  #         id          = "user-id-123"
  #         permissions = ["VIEW", "MODIFY"]
  #       }
  #     ]
  #   }

  depends_on = [dremio_folder.resource_folder_example]
}

output "dremio_udf_id" {
  value       = dremio_udf.example.id
  description = "ID of the UDF"
}

