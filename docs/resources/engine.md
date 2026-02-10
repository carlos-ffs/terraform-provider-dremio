# dremio_engine (Resource)

Manages a compute engine in Dremio Cloud. Engines provide the compute resources for running queries.

## Example Usage

```hcl
resource "dremio_engine" "analytics" {
  name                      = "analytics-engine"
  size                      = "SMALL_V1"
  min_replicas              = 1
  max_replicas              = 3
  auto_stop_delay_seconds   = 3600
  queue_time_limit_seconds  = 300
  runtime_limit_seconds     = 3600
  drain_time_limit_seconds  = 120
  max_concurrency           = 10
  description               = "Engine for analytics workloads"
  enable                    = true
}
```

## Schema

### Required

| Attribute | Type | Description |
|-----------|------|-------------|
| `name` | String | User-defined name for the engine. Must be unique within the project. |
| `size` | String | Size of the engine. Valid values: `XX_SMALL_V1`, `X_SMALL_V1`, `SMALL_V1`, `MEDIUM_V1`, `LARGE_V1`, `X_LARGE_V1`, `XX_LARGE_V1`, `XXX_LARGE_V1`. |
| `min_replicas` | Number | Minimum number of engine replicas that will be enabled at any given time. |
| `max_replicas` | Number | Maximum number of engine replicas that will be enabled at any given time. |
| `auto_stop_delay_seconds` | Number | Time (in seconds) that auto-stop is delayed after the last query completes. |
| `queue_time_limit_seconds` | Number | Maximum time (in seconds) a query will wait in the engine's queue before being canceled. Should be >= 120 seconds. |
| `runtime_limit_seconds` | Number | Maximum time (in seconds) a query can run before being terminated. Set to 0 for no limit. |
| `drain_time_limit_seconds` | Number | Maximum time (in seconds) an engine replica will continue to run after resize/disable/delete before termination. |
| `max_concurrency` | Number | Maximum number of concurrent queries that an engine replica can run. |

### Optional

| Attribute | Type | Default | Description |
|-----------|------|---------|-------------|
| `description` | String | `null` | Human-readable description for the engine. |
| `enable` | Boolean | `true` | Whether the engine is enabled. Set to `false` to disable the engine. |

### Read-Only

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | Unique identifier of the engine (UUID). |
| `state` | String | Current state of the engine. Values: `DELETING`, `DISABLED`, `DISABLING`, `ENABLED`, `ENABLING`, `INVALID`. |
| `active_replicas` | Number | Number of engine replicas currently active. |
| `queried_at` | String | Date and time (ISO 8601) the engine was last used to execute a query. |
| `status_changed_at` | String | Date and time (UTC) that the state of the engine changed. |
| `instance_family` | String | Cloud instance family used (e.g., `M5D`, `M6ID`, `M6GD`, `DDV4`, `DDV5`). |
| `additional_engine_state_info` | String | Additional engine state information (typically `NONE`). |

## Import

Engines can be imported using their ID or name:

```bash
terraform import dremio_engine.example engine-uuid-here
terraform import dremio_engine.example analytics-engine
```

## Engine Sizes

| Size | Description |
|------|-------------|
| `XX_SMALL_V1` | Extra extra small |
| `X_SMALL_V1` | Extra small |
| `SMALL_V1` | Small |
| `MEDIUM_V1` | Medium |
| `LARGE_V1` | Large |
| `X_LARGE_V1` | Extra large |
| `XX_LARGE_V1` | Extra extra large |
| `XXX_LARGE_V1` | Triple extra large |

## Notes

- **Cloud only**: This resource is only available for Dremio Cloud.
- **Auto-stop**: Engines automatically stop after `auto_stop_delay_seconds` of inactivity.
- **Scaling**: Engines scale between `min_replicas` and `max_replicas` based on load.
- **State transitions**: State changes (enabling/disabling) are asynchronous.

## Auto-Scaling Example

```hcl
resource "dremio_engine" "autoscale" {
  name                      = "autoscale-engine"
  size                      = "MEDIUM_V1"
  min_replicas              = 1
  max_replicas              = 5
  auto_stop_delay_seconds   = 1800  # 30 minutes
  queue_time_limit_seconds  = 600   # 10 minutes
  runtime_limit_seconds     = 7200  # 2 hours
  max_concurrency           = 20
  enable                    = true
}

output "engine_state" {
  value = dremio_engine.autoscale.state
}

output "active_replicas" {
  value = dremio_engine.autoscale.active_replicas
}
```

