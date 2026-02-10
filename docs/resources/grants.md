# dremio_grants (Resource)

Manages access control grants (privileges) for a catalog object in Dremio. This resource controls which users and roles can access sources, spaces, folders, datasets, views, and UDFs.

## Example Usage

```hcl
resource "dremio_grants" "source_access" {
  catalog_object_id = dremio_source.s3_samples.id
  
  grants = [
    {
      id           = "user-uuid-here"
      grantee_type = "USER"
      privileges   = ["SELECT", "ALTER"]
    },
    {
      id           = "role-uuid-here"
      grantee_type = "ROLE"
      privileges   = ["SELECT"]
    }
  ]
}
```

## Schema

### Required

| Attribute | Type | Description |
|-----------|------|-------------|
| `catalog_object_id` | String | UUID of the catalog object (source, folder, dataset, etc.) to manage grants for. |

#### grants (Set of Object) - Required

Set of grants to apply to the catalog object. If empty, all explicit grants will be removed from the object.

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | Yes | UUID of the user or role to grant privileges to. |
| `grantee_type` | String | Yes | Type of grantee. Valid values: `USER`, `ROLE`. |
| `privileges` | Set of String | Yes | Set of privileges to grant. Available privileges depend on the catalog object type. See Available Privileges below. |

### Read-Only

| Attribute | Type | Description |
|-----------|------|-------------|
| `available_privileges` | List of String | List of available privileges for this catalog object type. This is computed from the API response. |

## Import

Grants can be imported using the catalog object ID:

```bash
terraform import dremio_grants.example catalog-object-uuid-here
```

## Available Privileges

Privileges vary by object type:

| Privilege | Description |
|-----------|-------------|
| `SELECT` | Query data from the object |
| `ALTER` | Modify object settings |
| `MANAGE_GRANTS` | Grant/revoke privileges |
| `DELETE` | Delete rows (tables) |
| `INSERT` | Insert rows (tables) |
| `TRUNCATE` | Remove all rows (tables) |
| `UPDATE` | Update rows (tables) |
| `DROP` | Delete the object |
| `CREATE_TABLE` | Create tables (folders) |
| `MODIFY` | Modify object content |
| `READ_METADATA` | Read object metadata |
| `ALTER_REFLECTION` | Modify Reflections |
| `VIEW_REFLECTION` | View Reflection status |

## Notes

- **Replaces all grants**: This resource manages ALL grants for the object. Existing grants not in the configuration are removed.
- **One resource per object**: Only one `dremio_grants` resource should exist per catalog object.
- **Available privileges computed**: The `available_privileges` attribute shows which privileges are valid for this object type.
- **User/Role IDs**: You must use UUIDs, not names, for user and role identifiers.

## Example with View

```hcl
resource "dremio_view" "sales_report" {
  path = ["Analytics", "reports", "sales_report"]
  sql  = "SELECT * FROM sales"
}

resource "dremio_grants" "sales_report_access" {
  catalog_object_id = dremio_view.sales_report.id
  
  grants = [
    {
      id           = "analysts-role-uuid"
      grantee_type = "ROLE"
      privileges   = ["SELECT"]
    },
    {
      id           = "data-engineers-role-uuid"
      grantee_type = "ROLE"
      privileges   = ["SELECT", "ALTER", "DROP"]
    }
  ]
}

output "available_privileges" {
  value = dremio_grants.sales_report_access.available_privileges
}
```

## Best Practices

1. **Use roles over users**: Assign privileges to roles, then add users to roles.
2. **Principle of least privilege**: Grant only the minimum required privileges.
3. **Document access patterns**: Use tags and wikis to document why access is granted.
4. **Manage grants with Terraform**: Avoid mixing manual grants with Terraform-managed grants.

