# =============================================================================
# Dremio File Data Source Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

data "dremio_file" "datasource_file_example" {
  path = ["${dremio_source.samples_bucket.name}", "samples.dremio.com", "NYC-taxi-trips.csv"]
}

output "data_file_id" {
  value       = data.dremio_file.datasource_file_example.id
  description = "ID of the file"
}

