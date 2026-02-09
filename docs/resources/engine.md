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

- `name` (String) - Name of the engine. Must be unique within the project.
- `size` (String) - Engine size. Valid values: `XX_SMALL_V1`, `X_SMALL_V1`, `SMALL_V1`, `MEDIUM_V1`, `LARGE_V1`, `X_LARGE_V1`, `XX_LARGE_V1`, `XXX_LARGE_V1`.

### Optional

- `min_replicas` (Number) - Minimum number of engine replicas. Default: 0.
- `max_replicas` (Number) - Maximum number of engine replicas.
- `auto_stop_delay_seconds` (Number) - Seconds of inactivity before auto-stopping. Default: 3600 (1 hour).
- `queue_time_limit_seconds` (Number) - Maximum seconds a query can wait in queue.
- `runtime_limit_seconds` (Number) - Maximum seconds a query can run.
- `drain_time_limit_seconds` (Number) - Seconds to wait for queries to complete before scaling down.
- `max_concurrency` (Number) - Maximum concurrent queries per replica.
- `description` (String) - Human-readable description of the engine.
- `enable` (Boolean) - Whether the engine is enabled. Default: true.

### Read-Only

- `id` (String) - Unique identifier of the engine.
- `state` (String) - Current engine state. Values: `DELETING`, `DISABLED`, `DISABLING`, `ENABLED`, `ENABLING`, `INVALID`.
- `active_replicas` (Number) - Current number of running replicas.
- `queried_at` (String) - Timestamp of the last query.
- `status_changed_at` (String) - Timestamp of the last status change.
- `instance_family` (String) - Cloud instance family used.
- `additional_engine_state_info` (String) - Additional state information.

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

