# Dremio Terraform Provider

The Dremio Terraform Provider enables you to manage [Dremio Cloud](https://www.dremio.com/) resources using Infrastructure as Code.

> [!IMPORTANT]
> This provider currently only supports Dremio Cloud.

## Documentation

Full documentation is available in the [docs](docs/index.md) folder, including:

- [Provider Configuration](docs/index.md)
- [Resources](docs/resources/)
- [Data Sources](docs/data-sources/)

## Features

- **Sources** - Manage data source connections (S3, Snowflake, MySQL, PostgreSQL, and more)
- **Folders** - Create and organize folders within sources
- **Tables** - Promote files to queryable tables
- **Views** - Create virtual datasets with SQL
- **UDFs** - Create user-defined functions
- **Dataset Tags & Wiki** - Add metadata and documentation
- **Grants** - Manage access control and permissions
- **Engines** - Configure compute engines
- **Engine Rules** - Set up query routing rules
- **Data Maintenance** - Automate table optimization tasks

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (for building from source)

## Installation

### From Terraform Registry

```hcl
terraform {
  required_providers {
    dremio = {
      source = "registry.terraform.io/carlos-ffs/dremio"
    }
  }
}
```

### Building from Source

1. Clone the repository:

```shell
git clone https://github.com/carlos-ffs/terraform-provider-dremio.git
cd dremio-terraform-provider
```

2. Build and install the provider:

```shell
make install
```

This will build the provider binary and install it to `$GOPATH/bin`.

3. Configure Terraform to use the local provider by creating a `~/.terraformrc` file:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/carlos-ffs/dremio" = "/path/to/your/gopath/bin"
  }
  direct {}
}
```

## Quick Start

1. Configure the provider with your Dremio Cloud credentials:

```hcl
provider "dremio" {
  host                  = "https://api.dremio.cloud"
  personal_access_token = var.dremio_pat
  type                  = "cloud"
  project_id            = var.dremio_project_id
}
```

2. Create resources:

```hcl
resource "dremio_source" "s3_data" {
  name = "my-s3-source"
  type = "S3"
  config = jsonencode({
    accessKey          = var.aws_access_key
    accessSecret       = var.aws_secret_key
    rootPath           = "/"
    secure             = true
    externalBucketList = ["my-bucket"]
  })
}
```

See the [examples](examples/) folder for more usage examples.

## Development

### Prerequisites

- Go 1.24+
- Make

### Available Make Targets

| Target | Description |
| ------ | ----------- |
| `make build` | Build the provider binary |
| `make install` | Build and install to `$GOPATH/bin` |
| `make fmt` | Format Go source code |
| `make lint` | Run linter |
| `make test` | Run unit tests |
| `make testacc` | Run acceptance tests (creates real resources) |
| `make generate` | Generate documentation |

### Running Tests

```shell
# Unit tests
make test

# Acceptance tests (requires Dremio Cloud credentials)
export DREMIO_PAT="your-personal-access-token"
export DREMIO_PROJECT_ID="your-project-id"
make testacc
```

> [!WARNING]
> Acceptance tests create real resources in Dremio Cloud and may incur costs.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.
