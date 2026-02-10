# dremio_source (Resource)

Manages a data source connection in Dremio. Sources define connections to external data stores such as S3, Snowflake, MySQL, PostgreSQL, and many others.

## Example Usage

```hcl
# S3 Source
resource "dremio_source" "s3_samples" {
  name = "Samples"
  type = "S3"
  config = jsonencode({
    accessKey          = ""
    accessSecret       = ""
    rootPath           = "/"
    secure             = true
    externalBucketList = ["samples.dremio.com"]
  })

  metadata_policy = {
    auth_ttl_ms              = 86400000
    names_refresh_ms         = 3600000
    dataset_refresh_after_ms = 3600000
    dataset_expire_after_ms  = 10800000
    dataset_update_mode      = "PREFETCH_QUERIED"
  }
}
```

## Supported Source Types

- `ARCTIC` - Dremio Arctic (Nessie catalog)
- `S3` - Amazon S3
- `SNOWFLAKE` - Snowflake
- `MYSQL` - MySQL
- `POSTGRES` - PostgreSQL
- `BIGQUERY` - Google BigQuery
- `REDSHIFT` - Amazon Redshift
- `ORACLE` - Oracle Database
- `MSSQL` - Microsoft SQL Server
- `AZURE_STORAGE` - Azure Blob Storage / ADLS Gen2
- `AWS_GLUE` - AWS Glue Data Catalog
- `DB2` - IBM DB2
- `ICEBERG_REST_CATALOG` - Iceberg REST Catalog
- `AZURE_SYNAPSE` - Azure Synapse Analytics
- `SAPHANA` - SAP HANA
- `SNOWFLAKE_OPEN_CATALOG` - Snowflake Open Catalog
- `UNITY_CATALOG` - Databricks Unity Catalog
- `VERTICA` - Vertica

## Schema

### Required

| Attribute | Type | Description |
|-----------|------|-------------|
| `name` | String | User-defined name for the source. Must be unique within the project. |
| `type` | String | The type of source. See Supported Source Types above. |
| `config` | String (JSON) | Configuration options specific to the source type as a JSON-encoded string. Use `jsonencode()` to construct this value. See [Dremio API documentation](https://docs.dremio.com/cloud/reference/api/catalog/source/source-config) for available options. |

### Optional

#### metadata_policy (Block)

Controls how Dremio refreshes metadata from the source.

| Attribute | Type | Default | Description |
|-----------|------|---------|-------------|
| `auth_ttl_ms` | Number | `86400000` (24 hours) | How long source permissions are cached, in milliseconds. Minimum: 60000 (1 minute). |
| `names_refresh_ms` | Number | `3600000` (1 hour) | How often to refresh the source's namespace, in milliseconds. Minimum: 60000 (1 minute). |
| `dataset_refresh_after_ms` | Number | `3600000` (1 hour) | How often to refresh dataset metadata, in milliseconds. Minimum: 60000 (1 minute). |
| `dataset_expire_after_ms` | Number | `3600000` (1 hour) | How long before dataset metadata expires, in milliseconds. Minimum: 60000 (1 minute). |
| `dataset_update_mode` | String | `PREFETCH_QUERIED` | Metadata update policy. Valid values: `PREFETCH` (update all datasets), `PREFETCH_QUERIED` (update only previously queried datasets), `INLINE` (update on query). |
| `delete_unavailable_datasets` | Boolean | `true` | Remove dataset definitions if underlying data is unavailable to Dremio. |
| `auto_promote_datasets` | Boolean | `false` | Automatically format files into tables when queried. Applies only to metastore and object storage sources. |

#### Acceleration Settings

| Attribute | Type | Description |
|-----------|------|-------------|
| `acceleration_grace_period_ms` | Number | Time to keep Reflections before expiration (milliseconds). |
| `acceleration_refresh_period_ms` | Number | Refresh frequency for Reflections (milliseconds). |
| `acceleration_active_policy_type` | String | Policy for refreshing Reflections. Valid values: `NEVER`, `PERIOD`, `SCHEDULE`. |
| `acceleration_refresh_schedule` | String | Cron expression for Reflection refresh schedule (UTC). Example: `0 0 8 * * ?`. |
| `acceleration_refresh_on_data_changes` | Boolean | Refresh Reflections when Iceberg table snapshots change. |

#### access_control_list (Block)

User and role access settings. Can only be set via update (not on initial creation).

**users** (List of Object):

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | Yes | UUID of the user. |
| `permissions` | List of String | Yes | List of permissions to grant. |

**roles** (List of Object):

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | Yes | UUID of the role. |
| `permissions` | List of String | Yes | List of permissions to grant. |

### Read-Only

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | Unique identifier of the source (UUID). |
| `entity_type` | String | Type of catalog object (always `source`). |
| `tag` | String | Version tag for optimistic concurrency control. This value changes with every update. |

## Import

Sources can be imported using their ID:

```bash
terraform import dremio_source.example source-uuid-here
```

Or by name:

```bash
terraform import dremio_source.example source-name-here
```

## Notes

- The `config` attribute must be a valid JSON string for the specified source type. Refer to the [Dremio API documentation](https://docs.dremio.com/cloud/reference/api/) for the specific configuration options required for each source type.
- Changes to the `name` attribute will force recreation of the resource.
- Access control lists can only be set after initial creation (via update operation).
- The `tag` attribute is used for optimistic concurrency control. It changes with each update.

## Source-Specific Configuration Examples

### MySQL Source

```hcl
resource "dremio_source" "mysql" {
  name = "my-mysql"
  type = "MYSQL"
  config = jsonencode({
    hostname = "mysql-server.example.com"
    port     = "3306"
    username = "dremio_user"
    password = "secret"
  })
}
```

### PostgreSQL Source

```hcl
resource "dremio_source" "postgres" {
  name = "my-postgres"
  type = "POSTGRES"
  config = jsonencode({
    hostname = "postgres-server.example.com"
    port     = "5432"
    databaseName = "mydb"
    username = "dremio_user"
    password = "secret"
  })
}
```

