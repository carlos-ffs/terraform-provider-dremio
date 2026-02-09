# =============================================================================
# Dremio Terraform Provider Configuration Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

terraform {
  required_providers {
    dremio = {
      source = "registry.terraform.io/carlos-ffs/dremio"
    }
  }
}

provider "dremio" {
  # Personal Access Token for authentication
  # Can also be set via DREMIO_PAT environment variable
  //personal_access_token = ""
  project_id = "07c43507-ebad-417d-8d22-148bf2408c66"
  type       = "cloud"
  host       = "https://api.dremio.cloud"
}

