# dremio_data_maintenance_task (Data Source)

Retrieves information about an existing data maintenance task in Dremio Cloud.

> [!NOTE]
> This data source is only available for Dremio Cloud deployments.

## Example Usage

```hcl
data "dremio_data_maintenance_task" "optimize_orders" {
  id = "task-uuid-here"
}

output "task_type" {
  value = data.dremio_data_maintenance_task.optimize_orders.type
}

output "is_enabled" {
  value = data.dremio_data_maintenance_task.optimize_orders.is_enabled
}
```

## Schema

### Required

- `id` (String) - UUID of the maintenance task.

### Read-Only

- `type` (String) - Type of maintenance task (`OPTIMIZE` or `EXPIRE_SNAPSHOTS`).
- `level` (String) - Scope of the task (currently only `TABLE`).
- `source_name` (String) - Name of the Open Catalog source.
- `is_enabled` (Boolean) - Whether the task is enabled.
- `table_id` (String) - Fully qualified table name (`folder1.folder2.table_name`).

## Notes

- The `id` must be a valid maintenance task UUID.
- Tasks are associated with Open Catalog tables.
- The `table_id` uses dot notation without the source name prefix.

## Maintenance Task Types

| Type | Description |
|------|-------------|
| `OPTIMIZE` | Compacts small files for better query performance |
| `EXPIRE_SNAPSHOTS` | Removes old snapshots to reclaim storage |

## Example: Check Task Status

```hcl
data "dremio_data_maintenance_task" "orders_optimize" {
  id = var.optimize_task_id
}

output "task_status" {
  value = {
    type       = data.dremio_data_maintenance_task.orders_optimize.type
    enabled    = data.dremio_data_maintenance_task.orders_optimize.is_enabled
    table      = data.dremio_data_maintenance_task.orders_optimize.table_id
    source     = data.dremio_data_maintenance_task.orders_optimize.source_name
  }
}
```

## Example: Reference Task in Outputs

```hcl
# Get existing maintenance tasks
data "dremio_data_maintenance_task" "optimize" {
  id = var.optimize_task_id
}

data "dremio_data_maintenance_task" "expire" {
  id = var.expire_task_id
}

# Output maintenance status
output "maintenance_overview" {
  value = {
    optimize = {
      table   = data.dremio_data_maintenance_task.optimize.table_id
      enabled = data.dremio_data_maintenance_task.optimize.is_enabled
    }
    expire = {
      table   = data.dremio_data_maintenance_task.expire.table_id
      enabled = data.dremio_data_maintenance_task.expire.is_enabled
    }
  }
}
```

## Example: Conditional Management

```hcl
data "dremio_data_maintenance_task" "existing" {
  id = var.task_id
}

# Create or update task based on current state
resource "dremio_data_maintenance" "managed" {
  type       = data.dremio_data_maintenance_task.existing.type
  table_id   = data.dremio_data_maintenance_task.existing.table_id
  is_enabled = var.enable_maintenance
}
```

