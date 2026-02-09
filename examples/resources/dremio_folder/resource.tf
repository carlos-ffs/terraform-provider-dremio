# =============================================================================
# Dremio Folder Resource Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

# Create a top-level folder
resource "dremio_folder" "resource_folder_carlos_santos" {
  path = ["carlos_santos"]
}

# Create a nested folder using concat
resource "dremio_folder" "resource_folder_example" {
  path = concat(dremio_folder.resource_folder_carlos_santos.path, ["test-folder3"])
}

output "resource_folder_example_id" {
  value       = dremio_folder.resource_folder_example.id
  description = "ID of the folder"
}

