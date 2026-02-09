# =============================================================================
# Dremio Folder Data Source Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

data "dremio_folder" "datasource_folder_example" {
  path = [dremio_source.samples_bucket.name, "samples.dremio.com", "NYC-taxi-trips"]
}

output "data_folder_id" {
  value       = data.dremio_folder.datasource_folder_example.id
  description = "ID of the folder"
}

output "data_folder_path" {
  value       = data.dremio_folder.datasource_folder_example.path
  description = "Path of the folder"
}

