# dremio_udf (Resource)

Creates and manages a User-Defined Function (UDF) in Dremio. UDFs allow you to create reusable SQL functions that can be called in queries.

## Example Usage

```hcl
resource "dremio_udf" "calculate_fare" {
  path              = ["Samples", "samples.dremio.com", "terraform_top_folder", "calculate_fare"]
  is_scalar         = true
  function_arg_list = "base_fare DOUBLE, tip DOUBLE"
  function_body     = "SELECT base_fare + tip + (base_fare * 0.08)"
  return_type       = "DOUBLE"
}
```

## Schema

### Required

| Attribute | Type | Description |
|-----------|------|-------------|
| `path` | List of String | Full path to the UDF, including the source/space name and folder hierarchy. The last element is the function name. Path elements must not contain: `/`, `:`, `[`, `]`. |
| `is_scalar` | Boolean | If `true`, the UDF is a scalar function (returns a single value). If `false`, the UDF is a tabular function (returns a result set). |
| `function_arg_list` | String | The name and data type of each argument, separated by space. Multiple arguments separated by commas. Example: `"domain VARCHAR, orderdate DATE"` |
| `function_body` | String | The SQL statement that the UDF should execute. |
| `return_type` | String | The data type of the result (scalar) or column definitions (tabular). Examples: `"DOUBLE"` for scalar, `"name VARCHAR, email VARCHAR"` for tabular. |

### Optional

#### access_control_list (Block)

User and role access settings.

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
| `id` | String | Unique identifier of the UDF (UUID). |
| `entity_type` | String | Type of catalog object (always `function`). |
| `tag` | String | Version tag for optimistic concurrency control. This value changes with every update. |

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
  path              = ["Samples", "samples.dremio.com", "analytics", "get_high_fares"]
  is_scalar         = false
  function_arg_list = "min_fare DOUBLE"
  function_body     = <<-EOT
    SELECT *
    FROM "NYC-taxi-trips"
    WHERE fare_amount > min_fare
  EOT
  return_type       = "pickup_datetime TIMESTAMP, fare_amount DOUBLE, trip_distance DOUBLE"
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

