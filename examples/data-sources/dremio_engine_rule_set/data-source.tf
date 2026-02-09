# =============================================================================
# Dremio Engine Rule Set Data Source Example (Dremio Cloud Only)
# =============================================================================
# This example is based on the working configuration from main.tf
# Note: This resource is only available for Dremio Cloud deployments.
# =============================================================================

data "dremio_engine_rule_set" "example" {
  depends_on = [dremio_engine_rule_set.example_engine_rule_set]
}

output "engine_rules" {
  value = data.dremio_engine_rule_set.example.rule_infos
}

