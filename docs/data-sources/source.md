# dremio_source (Data Source)

Retrieves information about an existing data source in Dremio.

## Example Usage

### By Name

```hcl
data "dremio_source" "samples" {
  name = "Samples"
}

output "source_id" {
  value = data.dremio_source.samples.id
}
```

### By ID

```hcl
data "dremio_source" "by_id" {
  id = "source-uuid-here"
}
```

## Schema

### Optional (One Required)

- `id` (String) - UUID of the source. Either `id` or `name` must be specified.
- `name` (String) - Name of the source. Either `id` or `name` must be specified.

### Read-Only

- `entity_type` (String) - Type of catalog object (always `source`).
- `type` (String) - The source type (e.g., `S3`, `MYSQL`, `POSTGRES`).
- `config` (String, JSON) - Source configuration as a JSON string.
- `tag` (String) - Version tag for the source.

- `metadata_policy` (Object) - Metadata refresh policy.
  - `auth_ttl_ms` (Number) - Auth TTL in milliseconds.
  - `names_refresh_ms` (Number) - Names refresh interval.
  - `dataset_refresh_after_ms` (Number) - Dataset refresh interval.
  - `dataset_expire_after_ms` (Number) - Dataset expiration time.
  - `dataset_update_mode` (String) - Update mode.
  - `delete_unavailable_datasets` (Boolean) - Delete unavailable datasets.
  - `auto_promote_datasets` (Boolean) - Auto-promote datasets.

- `acceleration_grace_period_ms` (Number) - Acceleration grace period.
- `acceleration_refresh_period_ms` (Number) - Acceleration refresh period.
- `acceleration_active_policy_type` (String) - Active policy type.
- `acceleration_refresh_schedule` (String) - Refresh schedule (cron).
- `acceleration_refresh_on_data_changes` (Boolean) - Refresh on data changes.

- `access_control_list` (Object) - ACL settings.
  - `users` (List of Object) - User access controls.
  - `roles` (List of Object) - Role access controls.

## Notes

- Specify either `id` or `name`, but not both.
- If `name` is specified, it's resolved to an ID via the catalog API.
- The `config` attribute contains sensitive information and should be handled carefully.

