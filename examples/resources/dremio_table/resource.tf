# =============================================================================
# Dremio Table Resource Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

resource "dremio_table" "resource_table_example" {
  path              = ["${dremio_source.samples_bucket.name}", "samples.dremio.com", "NYC-taxi-trips.csv"]
  file_or_folder_id = "dremio:/${join("/", data.dremio_file.datasource_file_example.path)}"

  format = {
    type                       = "Text"
    field_delimiter            = ","
    skip_first_line            = false
    extract_header             = true
    quote                      = "\""
    comment                    = "#"
    escape                     = "\""
    line_delimiter             = "\\r\\n"
    auto_generate_column_names = true
    trim_header                = false
  }

  depends_on = [dremio_source.samples_bucket]
}

