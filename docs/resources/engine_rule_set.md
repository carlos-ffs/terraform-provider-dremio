# dremio_engine_rule_set (Resource)

Manages query routing rules for engines in Dremio Cloud. Rules determine which engine processes specific queries based on conditions.

> [!WARNING]
> Only one `dremio_engine_rule_set` resource should be defined per Terraform configuration (i.e., per project). This resource manages ALL routing rules for the project.

## Example Usage

```hcl
resource "dremio_engine_rule_set" "routing" {
  rule_infos = [
    {
      name        = "route-analytics"
      condition   = "user_name = 'analyst@company.com'"
      engine_name = "analytics-engine"
      action      = "ROUTE"
    },
    {
      name           = "reject-expensive"
      condition      = "query_cost > 1000000"
      action         = "REJECT"
      reject_message = "Query too expensive. Please optimize."
    }
  ]
}
```

## Schema

### Required

- `rule_infos` (List of Object) - Ordered list of routing rules. Rules are evaluated in order.
  - `name` (String) - User-defined name for the rule.
  - `condition` (String) - SQL-like condition for matching queries.
  - `action` (String) - Action when condition matches. Valid values: `ROUTE`, `REJECT`.
  - `engine_name` (String) - Name of engine to route to (required for `ROUTE` action).
  - `reject_message` (String) - Message shown when rejecting (optional for `REJECT` action).

### Read-Only

- `rule_info_default` (Object) - The default rule applied when no other rules match.
  - `name` (String) - Rule name.
  - `condition` (String) - Condition (always matches for default).
  - `engine_name` (String) - Default engine name.
  - `action` (String) - Action (always `ROUTE` for default).
  - `reject_message` (String) - Not applicable for default rule.
- `tag` (String) - Version tag for optimistic concurrency control.

## Import

Engine rule sets can be imported (no ID required as there's only one per project):

```bash
terraform import dremio_engine_rule_set.example rules
```

## Rule Conditions

Conditions use SQL-like syntax with these available fields:

| Field | Type | Description |
|-------|------|-------------|
| `user_name` | String | Username of query submitter |
| `role_name` | String | Role of query submitter |
| `query_cost` | Number | Estimated query cost |
| `query_type` | String | Type of query (SELECT, DDL, etc.) |

### Condition Examples

```hcl
# Match by user
condition = "user_name = 'analyst@company.com'"

# Match by role
condition = "role_name = 'data-engineers'"

# Match by query cost
condition = "query_cost > 500000"

# Compound conditions
condition = "role_name = 'analysts' AND query_cost < 100000"
```

## Notes

- **Single resource**: Only one `dremio_engine_rule_set` should exist per configuration.
- **Order matters**: Rules are evaluated in order; first match wins.
- **Default rule**: A default rule always exists and routes to the default engine.
- **Replaces all rules**: This resource replaces ALL existing routing rules.
- **Cloud only**: This resource is only available for Dremio Cloud.

## Complete Example

```hcl
# Create engines
resource "dremio_engine" "small" {
  name = "small-engine"
  size = "SMALL_V1"
}

resource "dremio_engine" "large" {
  name = "large-engine"
  size = "LARGE_V1"
}

# Configure routing rules
resource "dremio_engine_rule_set" "routing" {
  rule_infos = [
    {
      name        = "expensive-to-large"
      condition   = "query_cost > 100000"
      engine_name = dremio_engine.large.name
      action      = "ROUTE"
    },
    {
      name        = "analysts-small"
      condition   = "role_name = 'analysts'"
      engine_name = dremio_engine.small.name
      action      = "ROUTE"
    },
    {
      name           = "reject-huge"
      condition      = "query_cost > 10000000"
      action         = "REJECT"
      reject_message = "Query exceeds maximum allowed cost"
    }
  ]
  
  depends_on = [dremio_engine.small, dremio_engine.large]
}

output "default_engine" {
  value = dremio_engine_rule_set.routing.rule_info_default.engine_name
}
```

