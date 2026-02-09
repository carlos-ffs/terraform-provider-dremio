terraform {
  required_providers {
    dremio = {
      source = "registry.terraform.io/carlos-ffs/dremio"
    }
  }
}

provider "dremio" {
  //personal_access_token = ""
  project_id = "07c43507-ebad-417d-8d22-148bf2408c66"
  type = "cloud"
  host = "https://api.dremio.cloud"
}

resource "dremio_source" "samples_bucket" {
    type = "S3"
    name = "Samples_123"
    config = jsonencode({
        externalBucketList = ["samples.dremio.com"]
        secure = false
        propertyList = []
        credentialType = "NONE"
    })
    acceleration_refresh_period_ms = 3600018
    acceleration_grace_period_ms = 10800000
    acceleration_active_policy_type = "PERIOD"
    acceleration_refresh_schedule = "0 0 8 * * *"
    acceleration_refresh_on_data_changes = false

    metadata_policy = {
        auth_ttl_ms = 86400000
        auto_promote_datasets = false
        dataset_expire_after_ms = 259200000
        dataset_refresh_after_ms = 86400000
        dataset_update_mode = "PREFETCH_QUERIED"
        delete_unavailable_datasets = true
        names_refresh_ms = 86400000
    }
}

output "resource_dremio_source_id" {
  value = dremio_source.samples_bucket.id
  description = "ID of the source"
}

output "resource_dremio_source_tag" {
  value = dremio_source.samples_bucket.tag
  description = "Tag of the source"
}

output "resource_dremio_source_name" {
    value = dremio_source.samples_bucket.name
    description = "Name of the source"
}

data "dremio_source" "datasource_dremio_source_example" {
    name = dremio_source.samples_bucket.name
}

output "datasource_dremio_source_example_id" {
  value = data.dremio_source.datasource_dremio_source_example.id
  description = "ID of the source"
}

output "datasource_dremio_source_example_metadataPolicy" {
  value = data.dremio_source.datasource_dremio_source_example.metadata_policy
  description = "metadataPolicy of the source"
}


# Folder datasource examples
resource "dremio_folder" "resource_folder_carlos_santos" {
  path = ["carlos_santos"]
}

resource "dremio_folder" "resource_folder_example" {
  path = concat(dremio_folder.resource_folder_carlos_santos.path, ["test-folder3"])
}

data "dremio_folder" "datasource_folder_example" {
  path = [dremio_source.samples_bucket.name, "samples.dremio.com", "NYC-taxi-trips"]
}

output "resource_folder_example_id" {
  value = dremio_folder.resource_folder_example.id
  description = "ID of the folder"
}

output "data_folder_id" {
  value = data.dremio_folder.datasource_folder_example.id
  description = "ID of the folder"
}

output "data_folder_path" {
  value = data.dremio_folder.datasource_folder_example.path
  description = "Path of the folder"
}


# File Datasource examples

data "dremio_file" "datasource_file_example" {
  path = ["${dremio_source.samples_bucket.name}", "samples.dremio.com", "NYC-taxi-trips.csv"]
}

output "data_file_id" {
  value = data.dremio_file.datasource_file_example.id
  description = "ID of the file"
}

# Table Resource examples

resource "dremio_table" "resource_table_example" {
    path = ["${dremio_source.samples_bucket.name}", "samples.dremio.com", "NYC-taxi-trips.csv"]
    file_or_folder_id = "dremio:/${join("/", data.dremio_file.datasource_file_example.path)}"
    
    format = {
        type                        = "Text"
        field_delimiter             = ","
        skip_first_line             = false
        extract_header              = true
        quote                       = "\""
        comment                     = "#"
        escape                      = "\""
        line_delimiter              = "\\r\\n"
        auto_generate_column_names  = true
        trim_header                 = false
    }

    depends_on = [ dremio_source.samples_bucket ]
}

data "dremio_table" "datasource_table_example" {
  id = dremio_table.resource_table_example.id
}

output "dremio_table_id" {
  value = data.dremio_table.datasource_table_example.id
  description = "ID of the table"
}

output "dremio_table_tag" {
  value = data.dremio_table.datasource_table_example.tag
  description = "Tag of the table"
}

output "dremio_table_format" {
  value = data.dremio_table.datasource_table_example.format
  description = "Format of the table"
}

// UDF Resource examples
// DESCRIBE FUNCTION carlos_santos."test-folder3".count_high_passenger_trips;
// SELECT carlos_santos."test-folder3".count_high_passenger_trips(5);
resource "dremio_udf" "example" {
  path = concat(dremio_folder.resource_folder_example.path, ["count_high_passenger_trips"])

  is_scalar          = true
  function_arg_list  = "min_passengers BIGINT"
  function_body      = "SELECT count(*) FROM \"${join("\".\"", dremio_table.resource_table_example.path)}\" WHERE passenger_count > min_passengers"
  return_type        = "BIGINT"
  
#   access_control_list = {
#     users = [
#       {
#         id          = "user-id-123"
#         permissions = ["VIEW", "MODIFY"]
#       }
#     ]
#   }

  depends_on = [ dremio_folder.resource_folder_example ]
}

data "dremio_udf" "datasource_udf_example" {
  path = dremio_udf.example.path
}

output "dremio_udf_id" {
  value = data.dremio_udf.datasource_udf_example.id
  description = "ID of the UDF"
}

// Dremio Dataset Tags examples
resource "dremio_dataset_tags" "dataset_tags_example" {
    dataset_id = dremio_table.resource_table_example.id
    tags = ["carlos-santos", "terraform", "SRE"]
}

output "dremio_dataset_tags_resources_version" {
  value = dremio_dataset_tags.dataset_tags_example.version
  description = "Version of the dataset tags"
}

data "dremio_dataset_tags" "datasource_dataset_tags_example" {
  dataset_id = dremio_table.resource_table_example.id
}

output "dremio_datasource_dataset_tags" {
  value = data.dremio_dataset_tags.datasource_dataset_tags_example.tags
  description = "Tags of the dataset"
}

// View Resource examples
// SELECT * FROM carlos_santos."test-folder3".high_passenger_trips;

locals {
  # NYC-taxi-trips.csv
  last_element                      = dremio_table.resource_table_example.path[length(dremio_table.resource_table_example.path) - 1]
  table_path_without_last_element   = slice(dremio_table.resource_table_example.path, 0, length(dremio_table.resource_table_example.path) - 1)
}

resource "dremio_view" "example" {
  path = concat(dremio_folder.resource_folder_example.path, ["high_passenger_trips"])

  sql = "SELECT passenger_count, trip_distance_mi FROM \"${local.last_element}\" WHERE passenger_count > 5"

  sql_context = local.table_path_without_last_element

#   access_control_list = {
#     users = [
#       {
#         id          = "user-id-123"
#         permissions = ["VIEW", "MODIFY"]
#       }
#     ]
#   }

  depends_on = [ dremio_folder.resource_folder_example, dremio_table.resource_table_example ]
}

output "dremio_view_id" {
  value = dremio_view.example.id
  description = "ID of the view"
}

output "dremio_view_fields" {
  value = jsondecode(dremio_view.example.fields)
  description = "Fields of the view as JSON"
}

data "dremio_view" "datasource_view_example" {
  path = dremio_view.example.path
}

output "dremio_datasource_view_id" {
  value = data.dremio_view.datasource_view_example.id
  description = "ID of the view from datasource"
}

output "dremio_datasource_view_sql" {
  value = data.dremio_view.datasource_view_example.sql
  description = "SQL query of the view from datasource"
}

output "dremio_datasource_view_fields" {
  value = data.dremio_view.datasource_view_example.fields
  description = "Fields of the view from datasource as JSON"
}

// Dremio Dataset Wiki examples

resource "dremio_dataset_wiki" "dataset_wiki_example" {
    dataset_id = dremio_table.resource_table_example.id
    text = <<-EOT
      # Test Wiki
      This is an example wiki for a catalog object in Dremio. Here is some text in **bold**. Here is some text in *italics*.

      Here is an example excerpt with quotation formatting:

      > Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.


      ## Heading Level 2

      Here is a bulleted list:
      * An item in a bulleted list
      * A second item in a bulleted list
      * A third item in a bulleted list


      ### Heading Level 3

      Here is a numbered list:
      1. An item in a numbered list
      1. A second item in a numbered list
      1. A third item in a numbered list


      Here is a sentence that includes an [external link to https://dremio.com](https://dremio.com).

      Here is an image:

      ![](https://www.dremio.com/wp-content/uploads/2022/03/Dremio-logo.png)

      Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
    EOT
}

data "dremio_dataset_wiki" "datasource_dataset_wiki_example" {
    dataset_id = dremio_table.resource_table_example.id
}

output "data_dataset_wiki_text" {
  value = data.dremio_dataset_wiki.datasource_dataset_wiki_example.text
  description = "Text of the dataset wiki from datasource"
}

output "data_dataset_wiki_version" {
  value = data.dremio_dataset_wiki.datasource_dataset_wiki_example.version
  description = "Version of the dataset wiki from datasource"
}

// Dremio Grants examples
resource "dremio_grants" "example_grants" {
    catalog_object_id = dremio_table.resource_table_example.id
    grants = [
        {
            id           = "d8d24ff7-1b4c-4f8c-b712-91e3801182a1" // carlos.santos@dremio.com
            grantee_type = "USER"
            privileges   = ["SELECT", "ALTER"]
        }
    ]
}

data "dremio_grants" "datasource_grants_example" {
  catalog_object_id = dremio_table.resource_table_example.id
  depends_on = [ dremio_grants.example_grants ]
}

output "resource_table_example_grants" {
  value = data.dremio_grants.datasource_grants_example.grants
  description = "Grants of the table"
}

// Dremio Engine examples
resource "dremio_engine" "example_engine" {
  name                     = "carlos-santos-test-engine"
  description              = "Test engine created by Terraform"
  size                     = "SMALL_V1"
  min_replicas             = 0
  max_replicas             = 2
  auto_stop_delay_seconds  = 300
  queue_time_limit_seconds = 300
  runtime_limit_seconds    = 0
  drain_time_limit_seconds = 300
  max_concurrency          = 1
  enable                   = false
}

data "dremio_engine" "datasource_engine_example" {
  name = "preview"
}

output "dremio_engine_id" {
  value = data.dremio_engine.datasource_engine_example.id
  description = "ID of the engine"
}

// Dremio Engine Rule Set examples
// Considerations:
// - Only one engine rule set resource should be defined per Terraform configuration. 
//   Multiple resources will override each other since the API replaces all rules on each update.
// - When this resource is applied, any existing rules not defined in the resource will be deleted. 
//   Consequently, if you create a rule with the UI, it will be deleted when the resource is applied. 
//   As well as, if you remove the resource from the terraform configuration, all rules will be deleted.

// The removed block prevents the resource from being destroyed when the Terraform configuration is removed.
// This way you can keep the engine rules set in Dremio without managing it with Terraform.
// Comment dremio_engine_rule_set.example_engine_rule_set block and uncomment the removed block.
# removed {
#   from = dremio_engine_rule_set.example_engine_rule_set
#   lifecycle {
#     destroy = false
#   }
# }

resource "dremio_engine_rule_set" "example_engine_rule_set" {
  rule_infos = [
    {
      name        = "UI to Preview"
      condition   = "query_type() = 'UI Preview' OR query_type() = 'Internal Preview'"
      engine_name = "preview"
      action      = "ROUTE"
    },
    {
      name        = "Reflections"
      condition   = "query_type() = 'Reflections'"
      engine_name = "preview"
      action      = "ROUTE"
    },
    {
      name        = "Metadata Refresh"
      condition   = "query_type() = 'Metadata Refresh'"
      engine_name = "preview"
      action      = "ROUTE"
    }
  ]
  tag = ""
}

data "dremio_engine_rule_set" "example" {
    depends_on = [ dremio_engine_rule_set.example_engine_rule_set ]
}

output "engine_rules" {
  value = data.dremio_engine_rule_set.example.rule_infos
}


# Data maintenance

resource "dremio_data_maintenance" "example_data_maintenance" {
  type       = "OPTIMIZE"
  is_enabled = true
  table_id   = "${join("\".\"", dremio_view.example.path)}"
}

data "dremio_data_maintenance_task" "example" {
  id = dremio_data_maintenance.example_data_maintenance.id
}

output "data_maintenance_task_is_enabled" {
  value = data.dremio_data_maintenance_task.example.is_enabled
  description = "Is enabled"
}
