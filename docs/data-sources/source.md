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

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | UUID of the source. Either `id` or `name` must be specified. |
| `name` | String | Name of the source. Either `id` or `name` must be specified. |

### Read-Only

| Attribute | Type | Description |
|-----------|------|-------------|
| `type` | String | Source type (e.g., `ARCTIC`, `S3`, `SNOWFLAKE`, `MYSQL`, `POSTGRES`, etc.). |
| `config` | String (JSON) | Configuration options specific to the source type as a JSON string. |
| `tag` | String | Version tag for optimistic concurrency control. |
| `permissions` | List of String | User's permissions on the source. |

#### metadata_policy (Object)

Metadata refresh policy settings.

| Attribute | Type | Description |
|-----------|------|-------------|
| `auth_ttl_ms` | Number | How long source permissions are cached (milliseconds). |
| `names_refresh_ms` | Number | When to run a refresh of the source namespace (milliseconds). |
| `dataset_refresh_after_ms` | Number | How often the dataset metadata is refreshed (milliseconds). |
| `dataset_expire_after_ms` | Number | Time before metadata expires (milliseconds). |
| `dataset_update_mode` | String | Metadata policy for dataset updates (`PREFETCH`, `PREFETCH_QUERIED`, `INLINE`). |
| `delete_unavailable_datasets` | Boolean | Remove dataset definitions if underlying data is unavailable. |
| `auto_promote_datasets` | Boolean | Automatically format files into tables when queried. |

#### Acceleration Settings

| Attribute | Type | Description |
|-----------|------|-------------|
| `acceleration_grace_period_ms` | Number | Grace period before using Reflections (milliseconds). |
| `acceleration_refresh_period_ms` | Number | Refresh period for Reflections (milliseconds). |
| `acceleration_never_expire` | Boolean | Whether Reflections never expire. |
| `acceleration_never_refresh` | Boolean | Whether Reflections never refresh. |
| `acceleration_active_policy_type` | String | Active policy type (`PERIOD` or `NEVER`). |
| `acceleration_refresh_schedule` | String | Cron expression for refresh schedule. |

#### children (List of Object)

Child entities in the source (folders, datasets, etc.).

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | Unique identifier of the entity. |
| `path` | List of String | Full path to the entity. |
| `tag` | String | Version tag. |
| `type` | String | Entity type (`CONTAINER` or `DATASET`). |
| `container_type` | String | Container type (`SPACE`, `SOURCE`, `FOLDER`, `HOME`). |
| `dataset_type` | String | Dataset type (`VIRTUAL_DATASET` or `PHYSICAL_DATASET`). |

#### owner (Object)

Owner information for the source.

| Attribute | Type | Description |
|-----------|------|-------------|
| `owner_id` | String | UUID of the owner. |
| `owner_type` | String | Owner type (`USER` or `ROLE`). |

#### access_control_list (Object)

User and role access settings.

**users** (List of Object):

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | UUID of the user. |
| `permissions` | List of String | List of permissions granted. |

**roles** (List of Object):

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | UUID of the role. |
| `permissions` | List of String | List of permissions granted. |

## Notes

- Specify either `id` or `name`, but not both.
- If `name` is specified, it's resolved to an ID via the catalog API.
- The `config` attribute contains sensitive information and should be handled carefully.

