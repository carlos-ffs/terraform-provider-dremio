# =============================================================================
# Dremio Engine Rule Set Resource Example (Dremio Cloud Only)
# =============================================================================
# This example is based on the working configuration from main.tf
# Note: This resource is only available for Dremio Cloud deployments.
#
# IMPORTANT CONSIDERATIONS:
# - Only ONE engine rule set resource should be defined per Terraform configuration
# - Multiple resources will override each other (API replaces all rules on update)
# - When applied, any rules not defined in the resource will be DELETED
# - If you remove this resource, all routing rules will be deleted
# =============================================================================

# The removed block prevents the resource from being destroyed when the
# Terraform configuration is removed. This way you can keep the engine rules
# set in Dremio without managing it with Terraform.
# Comment dremio_engine_rule_set.example_engine_rule_set block and uncomment
# the removed block.
# removed {
#   from = dremio_engine_rule_set.example_engine_rule_set
#   lifecycle {
#     destroy = false
#   }
# }

resource "dremio_engine_rule_set" "example_engine_rule_set" {
  rule_infos = [
    {
      name        = "UI to Preview"
      condition   = "query_type() = 'UI Preview' OR query_type() = 'Internal Preview'"
      engine_name = "preview"
      action      = "ROUTE"
    },
    {
      name        = "Reflections"
      condition   = "query_type() = 'Reflections'"
      engine_name = "preview"
      action      = "ROUTE"
    },
    {
      name        = "Metadata Refresh"
      condition   = "query_type() = 'Metadata Refresh'"
      engine_name = "preview"
      action      = "ROUTE"
    }
  ]
  tag = ""
}

