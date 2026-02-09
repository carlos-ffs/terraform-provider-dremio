package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/helpers"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dremioGrantsDataSource{}
	_ datasource.DataSourceWithConfigure = &dremioGrantsDataSource{}
)

type dremioGrantsDataSource struct {
	client *dremioClient.Client
}

func NewDremioGrantsDataSource() datasource.DataSource {
	return &dremioGrantsDataSource{}
}

// Metadata returns the data source type name.
func (d *dremioGrantsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_grants"
}

func (d *dremioGrantsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dremioGrantsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Grants data source - retrieves grants (privileges) for an existing catalog object such as sources, spaces, folders, datasets, views, and UDFs.",
		Attributes: map[string]schema.Attribute{
			"catalog_object_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the Dremio catalog object to retrieve grants for.",
				Required:            true,
			},
			"grants": schema.SetNestedAttribute{
				MarkdownDescription: "Set of grants on the catalog object. Each grant specifies a user or role and their privileges.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "UUID of the user or role.",
							Computed:            true,
						},
						"grantee_type": schema.StringAttribute{
							MarkdownDescription: "Type of grantee. Either 'USER' or 'ROLE'.",
							Computed:            true,
						},
						"privileges": schema.SetAttribute{
							MarkdownDescription: "Set of privileges granted.",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
			"available_privileges": schema.ListAttribute{
				MarkdownDescription: "List of available privileges for this catalog object type.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioGrantsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioGrantsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogObjectID := data.CatalogObjectID.ValueString()

	// Make API request
	apiResp, err := d.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/grants", catalogObjectID), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read grants: %s", err),
		)
		return
	}
	defer apiResp.Body.Close()

	apiRespBody, err := io.ReadAll(apiResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}

	var grantsResp models.GrantsResponse
	if err := json.Unmarshal(apiRespBody, &grantsResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, &grantsResp, &data, &resp.Diagnostics)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioGrantsDataSource) mapResponseToState(ctx context.Context, grantsResp *models.GrantsResponse, data *models.DremioGrantsDataSourceModel, diags *diag.Diagnostics) {
	// Map catalog object ID from response
	data.CatalogObjectID = types.StringValue(grantsResp.ID)

	// Map available privileges
	availablePrivsList, diagsL := types.ListValueFrom(ctx, types.StringType, grantsResp.AvailablePrivileges)
	diags.Append(diagsL...)
	data.AvailablePrivileges = availablePrivsList

	// Map grants using helper
	grantsList, diagsL := helpers.ConvertGranteesToTerraform(ctx, grantsResp.Grants)
	diags.Append(diagsL...)
	data.Grants = grantsList
}
