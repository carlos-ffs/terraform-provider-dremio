# =============================================================================
# Dremio Dataset Wiki Resource Example
# =============================================================================
# This example is based on the working configuration from main.tf
# =============================================================================

resource "dremio_dataset_wiki" "dataset_wiki_example" {
  dataset_id = dremio_table.resource_table_example.id
  text       = <<-EOT
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

