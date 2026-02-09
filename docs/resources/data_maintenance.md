# dremio_data_maintenance (Resource)

Manages data maintenance tasks for tables in Dremio Cloud. Maintenance tasks automate table optimization and snapshot cleanup for Open Catalog tables.

## Example Usage

### Table Optimization Task

```hcl
resource "dremio_data_maintenance" "optimize_orders" {
  type       = "OPTIMIZE"
  table_id   = "folder1.folder2.orders"
  is_enabled = true
}
```

### Snapshot Expiration Task

```hcl
resource "dremio_data_maintenance" "expire_snapshots" {
  type       = "EXPIRE_SNAPSHOTS"
  table_id   = "folder1.folder2.orders"
  is_enabled = true
}
```

## Schema

### Required

- `type` (String) - Type of maintenance task. Valid values:
  - `OPTIMIZE` - Compacts small files into larger files for better query performance.
  - `EXPIRE_SNAPSHOTS` - Removes old table snapshots to reclaim storage.
- `table_id` (String) - Fully qualified table name in format `folder1.folder2.table_name` (without source name).
- `is_enabled` (Boolean) - Whether the maintenance task is enabled.

### Read-Only

- `id` (String) - Unique identifier of the maintenance task.
- `level` (String) - Scope of the maintenance task (currently only `TABLE`).
- `source_name` (String) - Name of the Open Catalog source where the table resides.

## Import

Data maintenance tasks can be imported using their ID:

```bash
terraform import dremio_data_maintenance.example task-uuid-here
```

## Maintenance Task Types

### OPTIMIZE

The OPTIMIZE task compacts small data files into larger ones, which:
- Improves query performance by reducing file I/O
- Reduces metadata overhead
- Should be run periodically on tables with many small files

### EXPIRE_SNAPSHOTS

The EXPIRE_SNAPSHOTS task removes old table snapshots, which:
- Reclaims storage space
- Cleans up historical data no longer needed for time travel
- Should be run after data retention requirements are met

## Notes

- **Cloud only**: This resource is only available for Dremio Cloud.
- **Open Catalog tables**: These tasks only apply to tables in Open Catalog sources.
- **Table path format**: The `table_id` uses dot notation without the source name.
- **Automatic scheduling**: Dremio automatically schedules enabled maintenance tasks.
- **One task per type per table**: You can have one OPTIMIZE and one EXPIRE_SNAPSHOTS task per table.

## Example with Both Task Types

```hcl
# Create both maintenance tasks for a table
resource "dremio_data_maintenance" "orders_optimize" {
  type       = "OPTIMIZE"
  table_id   = "sales.orders"
  is_enabled = true
}

resource "dremio_data_maintenance" "orders_expire" {
  type       = "EXPIRE_SNAPSHOTS"
  table_id   = "sales.orders"
  is_enabled = true
}

output "optimize_task_id" {
  value = dremio_data_maintenance.orders_optimize.id
}

output "expire_task_source" {
  value = dremio_data_maintenance.orders_expire.source_name
}
```

## Managing Task State

You can enable or disable tasks without destroying them:

```hcl
# Disable optimization during migration
resource "dremio_data_maintenance" "orders_optimize" {
  type       = "OPTIMIZE"
  table_id   = "sales.orders"
  is_enabled = false  # Temporarily disabled
}
```

## Best Practices

1. **Enable both tasks**: For most tables, enable both OPTIMIZE and EXPIRE_SNAPSHOTS.
2. **Monitor performance**: Check query performance to determine if OPTIMIZE is needed.
3. **Consider retention**: Set EXPIRE_SNAPSHOTS based on your time-travel requirements.
4. **Disable during migrations**: Temporarily disable tasks during large data migrations.

