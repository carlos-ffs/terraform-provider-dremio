# =============================================================================
# Dremio Grants Resource Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

resource "dremio_grants" "example_grants" {
  catalog_object_id = dremio_table.resource_table_example.id
  grants = [
    {
      id           = "d8d24ff7-1b4c-4f8c-b712-91e3801182a1" # user@example.com
      grantee_type = "USER"
      privileges   = ["SELECT", "ALTER"]
    }
  ]
}

