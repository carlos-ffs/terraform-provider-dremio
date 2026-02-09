# dremio_dataset_wiki (Resource)

Manages wiki documentation for a dataset in Dremio. Wikis provide rich documentation using GitHub-flavored Markdown.

## Example Usage

```hcl
resource "dremio_dataset_wiki" "nyc_trips_wiki" {
  dataset_id = dremio_view.nyc_trips.id
  text       = <<-EOT
    # NYC Taxi Trips Dataset

    This view contains New York City taxi trip data.

    ## Columns
    - **pickup_datetime**: Timestamp of pickup
    - **dropoff_datetime**: Timestamp of dropoff
    - **passenger_count**: Number of passengers
    - **trip_distance**: Distance in miles
    - **fare_amount**: Fare in USD

    ## Usage Notes
    - Data is updated daily
    - Use for analytics and reporting
  EOT
}
```

## Schema

### Required

- `dataset_id` (String) - UUID of the dataset to document.
- `text` (String) - Wiki content in GitHub-flavored Markdown. Maximum 100,000 characters.

### Read-Only

- `version` (String) - Version identifier for optimistic concurrency control.

## Import

Dataset wikis can be imported using the dataset ID:

```bash
terraform import dremio_dataset_wiki.example dataset-uuid-here
```

## Notes

- **Markdown support**: Content supports GitHub-flavored Markdown including headings, lists, tables, code blocks, and links.
- **Character limit**: Maximum content length is 100,000 characters.
- **Replaces content**: This resource manages the entire wiki. Any existing content is replaced.
- **Dataset must exist**: The referenced dataset must exist before adding documentation.

## Example with View

```hcl
resource "dremio_view" "customer_360" {
  path = ["Analytics", "customer_360"]
  sql  = "SELECT * FROM customers JOIN orders USING (customer_id)"
}

resource "dremio_dataset_wiki" "customer_360_docs" {
  dataset_id = dremio_view.customer_360.id
  text       = <<-EOT
    # Customer 360 View

    A unified view of customer data combining customer profiles with order history.

    ## Data Sources
    | Source | Description |
    |--------|-------------|
    | customers | Customer master data |
    | orders | Order transaction history |

    ## Key Metrics
    - **total_orders**: Lifetime order count
    - **total_spend**: Lifetime order value
    - **avg_order_value**: Average order amount

    ## Refresh Schedule
    Updated every 6 hours from source systems.

    ## Owner
    Data Analytics Team - analytics@company.com
  EOT
}
```

## Markdown Features

GitHub-flavored Markdown supports:

- **Headings**: `# H1`, `## H2`, `### H3`
- **Emphasis**: `*italic*`, `**bold**`, `~~strikethrough~~`
- **Lists**: Ordered (`1.`) and unordered (`-` or `*`)
- **Links**: `[text](url)`
- **Code**: Inline `` `code` `` and fenced code blocks
- **Tables**: Pipe-delimited tables with headers
- **Blockquotes**: `> quote`
- **Task lists**: `- [ ] task` and `- [x] completed`

## Best Practices

1. **Start with a clear title** describing the dataset
2. **Document columns** with descriptions and data types
3. **Include data lineage** - where does the data come from?
4. **Note refresh schedules** and update frequency
5. **Specify owners and contacts** for questions
6. **Add usage examples** with sample queries
7. **Document known issues** or caveats

