# =============================================================================
# Dremio UDF Data Source Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

data "dremio_udf" "datasource_udf_example" {
  path = dremio_udf.example.path
}

output "dremio_udf_id" {
  value       = data.dremio_udf.datasource_udf_example.id
  description = "ID of the UDF"
}

