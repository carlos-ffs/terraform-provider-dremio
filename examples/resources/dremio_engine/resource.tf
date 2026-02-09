# =============================================================================
# Dremio Engine Resource Example (Dremio Cloud Only)
# =============================================================================
# This example is based on the working configuration from main.tf
# Note: This resource is only available for Dremio Cloud deployments.
# =============================================================================

resource "dremio_engine" "example_engine" {
  name                     = "carlos-santos-test-engine"
  description              = "Test engine created by Terraform"
  size                     = "SMALL_V1"
  min_replicas             = 0
  max_replicas             = 2
  auto_stop_delay_seconds  = 300
  queue_time_limit_seconds = 300
  runtime_limit_seconds    = 0
  drain_time_limit_seconds = 300
  max_concurrency          = 1
  enable                   = false
}

