# dremio_engine_rule_set (Data Source)

Retrieves the engine routing rules for a Dremio Cloud project.

> [!NOTE]
> This data source is only available for Dremio Cloud deployments.

## Example Usage

```hcl
data "dremio_engine_rule_set" "current" {}

output "rules" {
  value = data.dremio_engine_rule_set.current.rule_infos
}

output "default_rule" {
  value = data.dremio_engine_rule_set.current.rule_info_default
}
```

## Schema

### Read-Only

- `rule_infos` (List of Object) - Ordered list of routing rules.
  - `name` (String) - Rule name.
  - `condition` (String) - SQL-like condition for matching.
  - `engine_name` (String) - Target engine name.
  - `action` (String) - Action (`ROUTE` or `REJECT`).
  - `reject_message` (String) - Message for rejected queries.

- `rule_info_default` (Object) - The default rule.
  - `name` (String) - Default rule name.
  - `condition` (String) - Default condition.
  - `engine_name` (String) - Default engine name.
  - `action` (String) - Always `ROUTE`.
  - `reject_message` (String) - Not applicable.

- `tag` (String) - Version tag for optimistic concurrency control.

## Notes

- No input parameters required; returns all rules for the project.
- Rules are evaluated in order; first match wins.
- The default rule is always present and handles unmatched queries.

## Example: Display Current Rules

```hcl
data "dremio_engine_rule_set" "rules" {}

output "rule_summary" {
  value = [
    for rule in data.dremio_engine_rule_set.rules.rule_infos : {
      name      = rule.name
      condition = rule.condition
      action    = rule.action
      target    = rule.action == "ROUTE" ? rule.engine_name : "REJECTED"
    }
  ]
}

output "default_engine" {
  value = data.dremio_engine_rule_set.rules.rule_info_default.engine_name
}
```

## Example: Check Existing Rules Before Modification

```hcl
data "dremio_engine_rule_set" "current" {}

locals {
  has_reject_rule = anytrue([
    for rule in data.dremio_engine_rule_set.current.rule_infos :
    rule.action == "REJECT"
  ])
  
  routing_engines = distinct([
    for rule in data.dremio_engine_rule_set.current.rule_infos :
    rule.engine_name
    if rule.action == "ROUTE"
  ])
}

output "has_reject_rules" {
  value = local.has_reject_rule
}

output "engines_in_use" {
  value = local.routing_engines
}
```

## Example: Preserve Existing Rules and Add New

```hcl
data "dremio_engine_rule_set" "existing" {}

resource "dremio_engine_rule_set" "updated" {
  rule_infos = concat(
    # Keep existing rules
    [for rule in data.dremio_engine_rule_set.existing.rule_infos : {
      name           = rule.name
      condition      = rule.condition
      engine_name    = rule.engine_name
      action         = rule.action
      reject_message = rule.reject_message
    }],
    # Add new rule
    [
      {
        name        = "new-rule"
        condition   = "role_name = 'new-team'"
        engine_name = "team-engine"
        action      = "ROUTE"
      }
    ]
  )
}
```

