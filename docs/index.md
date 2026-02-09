# Dremio Provider

The Dremio provider enables Terraform to manage resources in [Dremio](https://www.dremio.com/), a data lakehouse platform. 

> [!IMPORTANT]
> This provider currently only supports Dremio Cloud.

## Features

- **Sources**: Manage data source connections (S3, Snowflake, MySQL, PostgreSQL, BigQuery, and more)
- **Folders**: Create and organize folders within spaces and sources
- **Tables**: Promote files to queryable tables with format configuration
- **Views**: Create virtual datasets with SQL queries
- **User-Defined Functions (UDFs)**: Create reusable SQL functions
- **Dataset Tags & Wiki**: Add metadata and documentation to datasets
- **Grants**: Manage access control and permissions
- **Engines** (Cloud only): Configure compute engines
- **Engine Rules** (Cloud only): Set up query routing rules
- **Data Maintenance** (Cloud only): Automate table optimization tasks

## Example Usage

```hcl
terraform {
  required_providers {
    dremio = {
      source = "registry.terraform.io/carlos-ffs/dremio"
    }
  }
}

provider "dremio" {
  host                   = "https://api.dremio.cloud"
  personal_access_token  = var.dremio_pat
  type                   = "cloud"
  project_id             = var.dremio_project_id
}

# Create a source
resource "dremio_source" "example" {
  name = "my-s3-source"
  type = "S3"
  config = jsonencode({
    accessKey         = "your-access-key"
    accessSecret      = "your-secret-key"
    rootPath          = "/"
    secure            = true
    externalBucketList = ["my-bucket"]
  })
}
```

## Authentication

The Dremio provider requires a Personal Access Token (PAT) for authentication. You can configure authentication in several ways:

### Configuration File

```hcl
provider "dremio" {
  personal_access_token = "your-personal-access-token"
  host                  = "https://api.dremio.cloud"
  type                  = "cloud"
  project_id            = "your-project-id"
}
```

### Environment Variables

```bash
export DREMIO_PAT="your-personal-access-token"
export DREMIO_HOST="https://api.dremio.cloud"
export DREMIO_TYPE="cloud"
export DREMIO_PROJECT_ID="your-project-id"
```

Configuration values specified in the provider block take precedence over environment variables.

## Schema

### Optional

- `host` (String) - Dremio API Host. Defaults to `https://api.dremio.cloud`. For Dremio Software, use your instance URL (e.g., `http://localhost:9047`).
- `personal_access_token` (String, Sensitive) - Dremio Personal Access Token. Can also be set via the `DREMIO_PAT` environment variable.
- `type` (String) - Dremio account type. Valid values are `cloud` or `software`. Defaults to `cloud`.
- `project_id` (String) - Dremio Project ID. Required for Dremio Cloud. Can also be set via the `DREMIO_PROJECT_ID` environment variable.

## Generating a Personal Access Token

### Dremio Cloud

1. Log in to your Dremio Cloud account
2. Click on your profile icon and select **Account Settings**
3. Navigate to **Personal Access Tokens**
4. Click **Create Token** and provide a name
5. Copy the token value (it will only be shown once)

## Resources

- [dremio_source](resources/source) - Manage data source connections
- [dremio_folder](resources/folder) - Manage folders
- [dremio_table](resources/table) - Promote files to tables
- [dremio_view](resources/view) - Create virtual datasets
- [dremio_udf](resources/udf) - Create user-defined functions
- [dremio_dataset_tags](resources/dataset_tags) - Manage dataset tags
- [dremio_dataset_wiki](resources/dataset_wiki) - Manage dataset documentation
- [dremio_grants](resources/grants) - Manage access control
- [dremio_engine](resources/engine) - Manage compute engines (Cloud only)
- [dremio_engine_rule_set](resources/engine_rule_set) - Manage routing rules (Cloud only)
- [dremio_data_maintenance](resources/data_maintenance) - Manage maintenance tasks (Cloud only)

## Data Sources

- [dremio_source](data-sources/source) - Read source information
- [dremio_folder](data-sources/folder) - Read folder information
- [dremio_file](data-sources/file) - Read file information
- [dremio_table](data-sources/table) - Read table information
- [dremio_view](data-sources/view) - Read view information
- [dremio_udf](data-sources/udf) - Read UDF information
- [dremio_dataset_tags](data-sources/dataset_tags) - Read dataset tags
- [dremio_dataset_wiki](data-sources/dataset_wiki) - Read dataset wiki
- [dremio_grants](data-sources/grants) - Read grants information
- [dremio_engine](data-sources/engine) - Read engine information (Cloud only)
- [dremio_engine_rule_set](data-sources/engine_rule_set) - Read routing rules (Cloud only)
- [dremio_data_maintenance_task](data-sources/data_maintenance_task) - Read maintenance tasks (Cloud only)

