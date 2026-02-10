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

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | UUID of the engine. Either `id` or `name` must be specified. |
| `name` | String | Name of the engine. Either `id` or `name` must be specified. |

### Read-Only

| Attribute | Type | Description |
|-----------|------|-------------|
| `size` | String | Size of the engine. Values: `XX_SMALL_V1`, `X_SMALL_V1`, `SMALL_V1`, `MEDIUM_V1`, `LARGE_V1`, `X_LARGE_V1`, `XX_LARGE_V1`, `XXX_LARGE_V1`. |
| `min_replicas` | Number | Minimum number of engine replicas. |
| `max_replicas` | Number | Maximum number of engine replicas. |
| `auto_stop_delay_seconds` | Number | Time (in seconds) that auto-stop is delayed after the last query completes. |
| `queue_time_limit_seconds` | Number | Maximum time (in seconds) a query will wait in the engine's queue. |
| `runtime_limit_seconds` | Number | Maximum time (in seconds) a query can run. |
| `drain_time_limit_seconds` | Number | Maximum time (in seconds) an engine replica will continue to run after resize/disable/delete. |
| `max_concurrency` | Number | Maximum number of concurrent queries per replica. |
| `description` | String | Description of the engine. |
| `state` | String | Current state of the engine. Values: `DELETING`, `DISABLED`, `DISABLING`, `ENABLED`, `ENABLING`, `INVALID`. |
| `active_replicas` | Number | Number of engine replicas currently active. |
| `queried_at` | String | Date and time (ISO 8601) the engine was last used to execute a query. |
| `status_changed_at` | String | Date and time (UTC) that the state of the engine changed. |
| `additional_engine_state_info` | String | Additional engine state information (typically `NONE`). |

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

