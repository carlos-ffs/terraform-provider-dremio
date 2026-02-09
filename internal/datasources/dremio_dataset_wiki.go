package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dremioDatasetWikiDataSource{}
	_ datasource.DataSourceWithConfigure = &dremioDatasetWikiDataSource{}
)

type dremioDatasetWikiDataSource struct {
	client *dremioClient.Client
}

func NewDremioDatasetWikiDataSource() datasource.DataSource {
	return &dremioDatasetWikiDataSource{}
}

// Metadata returns the data source type name.
func (d *dremioDatasetWikiDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataset_wiki"
}

func (d *dremioDatasetWikiDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dremioClient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *dremioClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

func (d *dremioDatasetWikiDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Dataset Wiki data source - retrieves wiki content for an existing dataset",
		Attributes: map[string]schema.Attribute{
			"dataset_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the source, folder, or dataset",
				Required:            true,
			},
			"text": schema.StringAttribute{
				MarkdownDescription: "Text displayed in the wiki, formatted with GitHub-flavored Markdown",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Number for the most recent version of the wiki, starting with 0",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioDatasetWikiDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioDatasetWikiDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	datasetID := data.DatasetID.ValueString()

	// Make API request
	api_resp, err := d.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/collaboration/wiki", datasetID), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read dataset wiki: %s", err),
		)
		return
	}
	defer api_resp.Body.Close()

	api_resp_body, err := io.ReadAll(api_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}

	var wikiResp models.WikiResponse
	if err := json.Unmarshal(api_resp_body, &wikiResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, &wikiResp, &data, &resp.Diagnostics)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioDatasetWikiDataSource) mapResponseToState(ctx context.Context, wikiResp *models.WikiResponse, data *models.DremioDatasetWikiDataSourceModel, diags *diag.Diagnostics) {
	// Map text
	data.Text = types.StringValue(wikiResp.Text)

	// Map version
	data.Version = types.Int64Value(int64(wikiResp.Version))
}

