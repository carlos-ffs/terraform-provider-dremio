# dremio_grants (Data Source)

Retrieves grants (privileges) for an existing catalog object in Dremio.

## Example Usage

```hcl
data "dremio_grants" "source_grants" {
  catalog_object_id = data.dremio_source.samples.id
}

output "grants" {
  value = data.dremio_grants.source_grants.grants
}

output "available_privileges" {
  value = data.dremio_grants.source_grants.available_privileges
}
```

## Schema

### Required

- `catalog_object_id` (String) - UUID of the catalog object (source, folder, dataset, etc.).

### Read-Only

- `grants` (Set of Object) - Set of grants on the catalog object.
  - `id` (String) - UUID of the user or role.
  - `grantee_type` (String) - Type of grantee (`USER` or `ROLE`).
  - `privileges` (Set of String) - Set of privileges granted.

- `available_privileges` (List of String) - List of privileges that can be granted on this object type.

## Notes

- Use this data source to view existing grants before modifying them.
- The `available_privileges` attribute shows valid privileges for the object type.
- Grants include both user and role assignments.

## Example with Source

```hcl
data "dremio_source" "s3_data" {
  name = "S3-Data"
}

data "dremio_grants" "s3_grants" {
  catalog_object_id = data.dremio_source.s3_data.id
}

output "s3_privileges" {
  value = data.dremio_grants.s3_grants.available_privileges
}

output "current_grants" {
  value = [
    for grant in data.dremio_grants.s3_grants.grants : {
      type       = grant.grantee_type
      id         = grant.id
      privileges = grant.privileges
    }
  ]
}
```

## Example: Check Access for Specific Roles

```hcl
data "dremio_view" "report" {
  path = ["Analytics", "sales_report"]
}

data "dremio_grants" "report_grants" {
  catalog_object_id = data.dremio_view.report.id
}

locals {
  roles_with_select = [
    for grant in data.dremio_grants.report_grants.grants :
    grant.id
    if grant.grantee_type == "ROLE" && contains(grant.privileges, "SELECT")
  ]
}

output "roles_with_select_access" {
  value = local.roles_with_select
}
```

## Example: Audit Access Across Objects

```hcl
locals {
  objects_to_audit = {
    source1 = data.dremio_source.data_lake.id
    folder1 = data.dremio_folder.analytics.id
    view1   = data.dremio_view.dashboard.id
  }
}

data "dremio_grants" "audit" {
  for_each          = local.objects_to_audit
  catalog_object_id = each.value
}

output "access_audit" {
  value = {
    for k, v in data.dremio_grants.audit : k => {
      grant_count = length(v.grants)
      privileges  = v.available_privileges
    }
  }
}
```

