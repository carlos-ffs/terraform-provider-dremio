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
	_ datasource.DataSource              = &dremioDatasetTagsDataSource{}
	_ datasource.DataSourceWithConfigure = &dremioDatasetTagsDataSource{}
)

type dremioDatasetTagsDataSource struct {
	client *dremioClient.Client
}

func NewDremioDatasetTagsDataSource() datasource.DataSource {
	return &dremioDatasetTagsDataSource{}
}

// Metadata returns the data source type name.
func (d *dremioDatasetTagsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataset_tags"
}

func (d *dremioDatasetTagsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dremioDatasetTagsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Dataset Tags data source - retrieves tags for an existing dataset",
		Attributes: map[string]schema.Attribute{
			"dataset_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the dataset",
				Required:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "List of tags applied to the dataset. Tags are case-insensitive labels used for organization and discovery.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version identifier for the current set of tags. Used for optimistic concurrency control.",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioDatasetTagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioDatasetTagsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	datasetID := data.DatasetID.ValueString()

	// Make API request
	api_resp, err := d.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/collaboration/tag", datasetID), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read dataset tags: %s", err),
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

	var tagResp models.TagResponse
	if err := json.Unmarshal(api_resp_body, &tagResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, &tagResp, &data, &resp.Diagnostics)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioDatasetTagsDataSource) mapResponseToState(ctx context.Context, tagResp *models.TagResponse, data *models.DremioDatasetTagsDataSourceModel, diags *diag.Diagnostics) {
	// Map tags
	if len(tagResp.Tags) == 0 {
		data.Tags = types.ListNull(types.StringType)
	} else {
		tagsList, d := types.ListValueFrom(ctx, types.StringType, tagResp.Tags)
		diags.Append(d...)
		data.Tags = tagsList
	}

	// Map version
	data.Version = types.StringValue(tagResp.Version)
}

