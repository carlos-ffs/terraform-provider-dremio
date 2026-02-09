package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &dremioFileDataSource{}
	_ datasource.DataSourceWithConfigure = &dremioFileDataSource{}
)

func NewDremioFileDataSource() datasource.DataSource {
	return &dremioFileDataSource{}
}

type dremioFileDataSource struct {
	client *dremioClient.Client
}

// Metadata returns the data source type name.
func (d *dremioFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *dremioFileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dremioFileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio File data source - retrieves information about an existing file",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the file",
				Computed:            true,
			},
			"path": schema.ListAttribute{
				MarkdownDescription: "Full path to the file",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[^/:[\]]*$`),
							"path elements must not contain the characters: /, :, [, ]",
						),
					),
				},
			},
			"entity_type": schema.StringAttribute{
				MarkdownDescription: "Type of catalog object (always 'file')",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioFileDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var file_path []string
	if !data.Path.IsNull() {
		diags := data.Path.ElementsAs(ctx, &file_path, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if len(file_path) == 0 {
		resp.Diagnostics.AddError(
			"Missing File Path",
			"The `path` must be specified for Dremio File data source.",
		)
		return
	}

	// Construct the API path using by-path endpoint
	file_path_str := "/" + strings.Join(file_path, "/")
	path := fmt.Sprintf("/catalog/by-path/%s", file_path_str)

	api_resp, err := d.client.RequestToDremio("GET", path, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to request file: %s", err),
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

	var fileResp models.FileResponse
	if err := json.Unmarshal(api_resp_body, &fileResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, &fileResp, &data, &resp.Diagnostics)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioFileDataSource) mapResponseToState(ctx context.Context, fileResp *models.FileResponse, data *models.DremioFileDataSourceModel, diags *diag.Diagnostics) {

	// Map basic fields
	data.ID = types.StringValue(fileResp.ID)

	if fileResp.EntityType != "" {
		data.EntityType = types.StringValue(fileResp.EntityType)
	} else {
		data.EntityType = types.StringNull()
	}

	if len(fileResp.Path) == 0 {
		data.Path = types.ListNull(types.StringType)
	} else {
		pathFromAPI, diagsTemp := types.ListValueFrom(ctx, types.StringType, fileResp.Path)
		diags.Append(diagsTemp...)
		data.Path = pathFromAPI
	}
}
