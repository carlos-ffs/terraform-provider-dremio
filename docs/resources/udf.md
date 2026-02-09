# dremio_udf (Resource)

Creates and manages a User-Defined Function (UDF) in Dremio. UDFs allow you to create reusable SQL functions that can be called in queries.

## Example Usage

```hcl
resource "dremio_udf" "calculate_fare" {
  path      = ["Samples", "samples.dremio.com", "terraform_top_folder", "calculate_fare"]
  is_scalar = true
  
  function_arg_list = [
    {
      name     = "base_fare"
      data_type = "DOUBLE"
    },
    {
      name     = "tip"
      data_type = "DOUBLE"
    }
  ]
  
  function_body = "SELECT base_fare + tip + (base_fare * 0.08)"
  
  return_type = {
    data_type = "DOUBLE"
  }
}
```

## Schema

### Required

- `path` (List of String) - Full path to the UDF, including the source/space name and folder hierarchy. The last element is the function name. Path elements must not contain: `/`, `:`, `[`, `]`.
- `is_scalar` (Boolean) - Whether the function is scalar (`true`) or tabular (`false`). Scalar functions return a single value; tabular functions return a table.
- `function_body` (String) - SQL expression defining the function logic.
- `return_type` (Block) - Return type specification for scalar functions.
  - `data_type` (String) - SQL data type of the return value (e.g., `DOUBLE`, `VARCHAR`, `INT`, `BOOLEAN`).

### Optional

- `function_arg_list` (Block List) - List of function arguments.
  - `name` (String) - Argument name.
  - `data_type` (String) - SQL data type of the argument.

- `access_control_list` (Block) - User and role access settings.
  - `users` (Block List) - List of user access controls.
    - `id` (String) - User ID.
    - `permissions` (List of String) - List of permissions.
  - `roles` (Block List) - List of role access controls.
    - `id` (String) - Role ID.
    - `permissions` (List of String) - List of permissions.

### Read-Only

- `id` (String) - Unique identifier of the UDF.
- `entity_type` (String) - Type of catalog object (always `function`).
- `tag` (String) - Version tag for optimistic concurrency control.

## Import

UDFs can be imported using their ID:

```bash
terraform import dremio_udf.example udf-uuid-here
```

## Notes

- **Scalar vs Tabular**: Scalar functions return a single value and can be used in SELECT, WHERE, etc. Tabular functions return a result set and are used in FROM clauses.
- **Parent folders must exist**: Ensure all parent folders exist before creating the UDF.
- **Function arguments**: Arguments are referenced by name in the function body.
- **Return type**: Required for scalar functions; for tabular functions, the return schema is inferred.

## Tabular Function Example

```hcl
resource "dremio_udf" "get_high_fares" {
  path      = ["Samples", "samples.dremio.com", "analytics", "get_high_fares"]
  is_scalar = false
  
  function_arg_list = [
    {
      name      = "min_fare"
      data_type = "DOUBLE"
    }
  ]
  
  function_body = <<-EOT
    SELECT *
    FROM "NYC-taxi-trips"
    WHERE fare_amount > min_fare
  EOT
  
  return_type = {
    data_type = "TABLE"
  }
}
```

## Usage in Queries

After creating a UDF, you can use it in SQL queries:

```sql
-- Scalar function
SELECT calculate_fare(base_fare, tip) AS total_fare
FROM my_table;

-- Tabular function
SELECT * FROM TABLE(get_high_fares(50.0));
```

## Supported Data Types

Common data types for arguments and return values:
- `INT`, `BIGINT`, `DOUBLE`, `FLOAT`, `DECIMAL`
- `VARCHAR`, `CHAR`
- `BOOLEAN`
- `DATE`, `TIME`, `TIMESTAMP`
- `TABLE` (for tabular function return types)

