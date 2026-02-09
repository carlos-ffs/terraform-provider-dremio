# dremio_engine (Data Source)

Retrieves information about an existing compute engine in Dremio Cloud.

> [!NOTE]
> This data source is only available for Dremio Cloud deployments.

## Example Usage

### By Name

```hcl
data "dremio_engine" "analytics" {
  name = "analytics-engine"
}

output "engine_id" {
  value = data.dremio_engine.analytics.id
}

output "engine_state" {
  value = data.dremio_engine.analytics.state
}
```

### By ID

```hcl
data "dremio_engine" "by_id" {
  id = "engine-uuid-here"
}
```

## Schema

### Optional (One Required)

- `id` (String) - UUID of the engine. Either `id` or `name` must be specified.
- `name` (String) - Name of the engine. Either `id` or `name` must be specified.

### Read-Only

- `size` (String) - Engine size (e.g., `SMALL_V1`, `MEDIUM_V1`).
- `min_replicas` (Number) - Minimum number of replicas.
- `max_replicas` (Number) - Maximum number of replicas.
- `auto_stop_delay_seconds` (Number) - Auto-stop delay in seconds.
- `queue_time_limit_seconds` (Number) - Queue time limit.
- `runtime_limit_seconds` (Number) - Runtime limit.
- `drain_time_limit_seconds` (Number) - Drain time limit.
- `max_concurrency` (Number) - Maximum concurrent queries.
- `description` (String) - Engine description.
- `enable` (Boolean) - Whether engine is enabled.
- `state` (String) - Current state (`ENABLED`, `DISABLED`, `ENABLING`, `DISABLING`, `DELETING`, `INVALID`).
- `active_replicas` (Number) - Current running replicas.
- `queried_at` (String) - Timestamp of last query.
- `status_changed_at` (String) - Timestamp of last status change.
- `instance_family` (String) - Cloud instance family.
- `additional_engine_state_info` (String) - Additional state info.

## Notes

- Specify either `id` or `name`, but not both.
- The `state` attribute shows the current operational status.
- The `active_replicas` shows how many replicas are currently running.

## Example: Monitor Engine Status

```hcl
data "dremio_engine" "production" {
  name = "production-engine"
}

output "is_healthy" {
  value = data.dremio_engine.production.state == "ENABLED"
}

output "active_capacity" {
  value = {
    current = data.dremio_engine.production.active_replicas
    max     = data.dremio_engine.production.max_replicas
  }
}

output "last_query" {
  value = data.dremio_engine.production.queried_at
}
```

## Example: Reference Engine in Rules

```hcl
data "dremio_engine" "default" {
  name = "default-engine"
}

data "dremio_engine" "high_priority" {
  name = "high-priority-engine"
}

resource "dremio_engine_rule_set" "routing" {
  rule_infos = [
    {
      name        = "vip-users"
      condition   = "role_name = 'vip'"
      engine_name = data.dremio_engine.high_priority.name
      action      = "ROUTE"
    }
  ]
}
```

