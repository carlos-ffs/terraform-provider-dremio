# dremio_udf (Data Source)

Retrieves information about an existing User-Defined Function (UDF) in Dremio.

## Example Usage

### By Path

```hcl
data "dremio_udf" "calculate_fare" {
  path = ["Analytics", "functions", "calculate_fare"]
}

output "udf_id" {
  value = data.dremio_udf.calculate_fare.id
}

output "udf_body" {
  value = data.dremio_udf.calculate_fare.function_body
}
```

### By ID

```hcl
data "dremio_udf" "by_id" {
  id = "udf-uuid-here"
}
```

## Schema

### Optional (One Required)

- `id` (String) - UUID of the UDF. Either `id` or `path` must be specified.
- `path` (List of String) - Full path to the UDF. Either `id` or `path` must be specified.

### Read-Only

- `entity_type` (String) - Type of catalog object (always `function`).
- `is_scalar` (Boolean) - Whether the function is scalar (`true`) or tabular (`false`).
- `function_body` (String) - SQL expression defining the function.
- `tag` (String) - Version tag for the UDF.

- `function_arg_list` (List of Object) - List of function arguments.
  - `name` (String) - Argument name.
  - `data_type` (String) - SQL data type.

- `return_type` (Object) - Return type specification.
  - `data_type` (String) - SQL data type of return value.

- `access_control_list` (Object) - ACL settings.
  - `users` (List of Object) - User access controls.
  - `roles` (List of Object) - Role access controls.

## Notes

- Specify either `id` or `path`, but not both.
- The function body contains the SQL logic of the UDF.
- Use `is_scalar` to determine if the function returns a single value or a table.

## Example with Grants

```hcl
data "dremio_udf" "format_currency" {
  path = ["Utils", "format_currency"]
}

resource "dremio_grants" "udf_access" {
  catalog_object_id = data.dremio_udf.format_currency.id
  grants = [
    {
      id           = "all-users-role-uuid"
      grantee_type = "ROLE"
      privileges   = ["SELECT"]
    }
  ]
}

output "function_args" {
  value = data.dremio_udf.format_currency.function_arg_list
}

output "return_type" {
  value = data.dremio_udf.format_currency.return_type
}
```

## Example: Reference UDF in View

```hcl
data "dremio_udf" "calculate_tax" {
  path = ["Utils", "calculate_tax"]
}

resource "dremio_view" "orders_with_tax" {
  path = ["Reports", "orders_with_tax"]
  sql  = <<-EOT
    SELECT 
      order_id,
      subtotal,
      Utils.calculate_tax(subtotal) as tax,
      subtotal + Utils.calculate_tax(subtotal) as total
    FROM orders
  EOT
}
```

