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

- `name` (String) - User-defined name for the source. Must be unique within the project.
- `type` (String) - The type of source. See Supported Source Types above.
- `config` (String, JSON) - Configuration options specific to the source type as a JSON-encoded string. Use `jsonencode()` to construct this value.

### Optional

- `metadata_policy` (Block) - Controls how Dremio refreshes metadata from the source.
  - `auth_ttl_ms` (Number) - How long source permissions are cached, in milliseconds.
  - `names_refresh_ms` (Number) - How often to refresh the source's namespace, in milliseconds.
  - `dataset_refresh_after_ms` (Number) - How often to refresh dataset metadata, in milliseconds.
  - `dataset_expire_after_ms` (Number) - How long before dataset metadata expires, in milliseconds.
  - `dataset_update_mode` (String) - Metadata update policy. Valid values: `PREFETCH`, `PREFETCH_QUERIED`, `INLINE`.
  - `delete_unavailable_datasets` (Boolean) - Remove dataset definitions if underlying data is unavailable.
  - `auto_promote_datasets` (Boolean) - Automatically format files into tables when queried.

- `acceleration_grace_period_ms` (Number) - Grace period before using Reflections, in milliseconds.
- `acceleration_refresh_period_ms` (Number) - How often to refresh Reflections, in milliseconds.
- `acceleration_active_policy_type` (String) - Active policy type. Valid values: `PERIOD`, `SCHEDULE`, `NEVER`.
- `acceleration_refresh_schedule` (String) - Cron expression for refresh schedule (UTC).
- `acceleration_refresh_on_data_changes` (Boolean) - Refresh Reflections when source data changes.

- `access_control_list` (Block) - User and role access settings.
  - `users` (Block List) - List of user access controls.
    - `id` (String) - User ID.
    - `permissions` (List of String) - List of permissions.
  - `roles` (Block List) - List of role access controls.
    - `id` (String) - Role ID.
    - `permissions` (List of String) - List of permissions.

### Read-Only

- `id` (String) - Unique identifier of the source.
- `entity_type` (String) - Type of catalog object (always `source`).
- `tag` (String) - Version tag for optimistic concurrency control. Required for updates.

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

