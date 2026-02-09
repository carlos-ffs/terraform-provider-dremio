# =============================================================================
# Dremio Source Resource Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

resource "dremio_source" "samples_bucket" {
  type = "S3"
  name = "Samples"

  config = jsonencode({
    externalBucketList = ["samples.dremio.com"]
    secure             = false
    propertyList       = []
    credentialType     = "NONE"
  })

  # Acceleration settings
  acceleration_refresh_period_ms        = 3600000
  acceleration_grace_period_ms          = 10800000
  acceleration_active_policy_type       = "PERIOD"
  acceleration_refresh_schedule         = "0 0 8 * * *"
  acceleration_refresh_on_data_changes  = false

  # Metadata policy settings
  metadata_policy = {
    auth_ttl_ms                 = 86400000
    auto_promote_datasets       = false
    dataset_expire_after_ms     = 259200000
    dataset_refresh_after_ms    = 86400000
    dataset_update_mode         = "PREFETCH_QUERIED"
    delete_unavailable_datasets = true
    names_refresh_ms            = 86400000
  }
}

output "source_id" {
  value       = dremio_source.samples_bucket.id
  description = "ID of the source"
}

output "source_tag" {
  value       = dremio_source.samples_bucket.tag
  description = "Tag of the source"
}

output "source_name" {
  value       = dremio_source.samples_bucket.name
  description = "Name of the source"
}

