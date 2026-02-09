package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/helpers"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                     = &dremioViewDataSource{}
	_ datasource.DataSourceWithConfigure        = &dremioViewDataSource{}
	_ datasource.DataSourceWithConfigValidators = &dremioViewDataSource{}
)

func NewDremioViewDataSource() datasource.DataSource {
	return &dremioViewDataSource{}
}

type dremioViewDataSource struct {
	client *dremioClient.Client
}

// ConfigValidators returns a list of functions which will all be performed during validation.
func (d *dremioViewDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("path"),
		),
	}
}

// Metadata returns the data source type name.
func (d *dremioViewDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_view"
}

func (d *dremioViewDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dremioViewDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio View data source - retrieves information about an existing view/virtual dataset",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the view. Exactly one of `id` or `path` must be specified.",
				Computed:            true,
				Optional:            true,
			},
			"path": schema.ListAttribute{
				MarkdownDescription: "Full path to the view",
				Optional:            true,
				Computed:            true,
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
			"type": schema.StringAttribute{
				MarkdownDescription: "Dataset type (VIRTUAL_DATASET)",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Date and time the view was created (UTC)",
				Computed:            true,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Version tag for optimistic concurrency control",
				Computed:            true,
			},
			"sql": schema.StringAttribute{
				MarkdownDescription: "SQL query defining the view",
				Computed:            true,
			},
			"sql_context": schema.ListAttribute{
				MarkdownDescription: "Context for SQL query execution",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"fields": schema.StringAttribute{
				MarkdownDescription: "View fields/columns as JSON string. Due to the recursive nature of view schemas (STRUCT and LIST types can be arbitrarily nested), fields are represented as a JSON string. Use jsondecode() to parse this value in Terraform configurations.",
				Computed:            true,
			},
			"access_control_list": schema.SingleNestedAttribute{
				MarkdownDescription: "User and role access settings",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"users": schema.ListNestedAttribute{
						MarkdownDescription: "List of user access controls",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "User ID",
									Computed:            true,
								},
								"permissions": schema.ListAttribute{
									MarkdownDescription: "List of permissions",
									Computed:            true,
									ElementType:         types.StringType,
								},
							},
						},
					},
					"roles": schema.ListNestedAttribute{
						MarkdownDescription: "List of role access controls",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "Role ID",
									Computed:            true,
								},
								"permissions": schema.ListAttribute{
									MarkdownDescription: "List of permissions",
									Computed:            true,
									ElementType:         types.StringType,
								},
							},
						},
					},
				},
			},
			"permissions": schema.ListAttribute{
				MarkdownDescription: "User's permissions on the view",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "Owner information",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"owner_id": schema.StringAttribute{
						MarkdownDescription: "Owner ID",
						Computed:            true,
					},
					"owner_type": schema.StringAttribute{
						MarkdownDescription: "Owner type (USER or ROLE)",
						Computed:            true,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioViewDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioViewDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var viewID string
	if !data.ID.IsNull() {
		viewID = data.ID.ValueString()
	}

	var view_path []string
	if !data.Path.IsNull() {
		diags := data.Path.ElementsAs(ctx, &view_path, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if viewID == "" && len(view_path) == 0 {
		resp.Diagnostics.AddError(
			"Missing View ID or Path",
			"Either `id` or `path` must be specified for Dremio View data source.",
		)
		return
	}
	if viewID != "" && len(view_path) > 0 {
		resp.Diagnostics.AddError(
			"Both View ID and Path specified",
			"Only one of `id` or `path` must be specified for Dremio View data source.",
		)
		return
	}

	var path string
	if viewID != "" { // Read view by ID
		path = fmt.Sprintf("/catalog/%s", viewID)
	} else { // Lookup view ID by path
		view_path_str := "/" + strings.Join(view_path, "/")
		path = fmt.Sprintf("/catalog/by-path/%s", view_path_str)
	}

	api_resp, err := d.client.RequestToDremio("GET", path, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to request view: %s", err),
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

	var viewResp models.ViewResponse
	if err := json.Unmarshal(api_resp_body, &viewResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, &viewResp, &data, &resp.Diagnostics)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioViewDataSource) mapResponseToState(ctx context.Context, viewResp *models.ViewResponse, data *models.DremioViewDataSourceModel, diags *diag.Diagnostics) {
	// Map basic fields
	data.ID = types.StringValue(viewResp.ID)
	data.Type = types.StringValue(viewResp.Type)
	data.Tag = types.StringValue(viewResp.Tag)

	if viewResp.CreatedAt != "" {
		data.CreatedAt = types.StringValue(viewResp.CreatedAt)
	} else {
		data.CreatedAt = types.StringNull()
	}

	// Map path
	if len(viewResp.Path) == 0 {
		data.Path = types.ListNull(types.StringType)
	} else {
		pathFromAPI, diagsTemp := types.ListValueFrom(ctx, types.StringType, viewResp.Path)
		diags.Append(diagsTemp...)
		data.Path = pathFromAPI
	}

	// Map SQL
	data.SQL = types.StringValue(viewResp.SQL)

	// Map SQL Context
	if len(viewResp.SQLContext) == 0 {
		data.SQLContext = types.ListNull(types.StringType)
	} else {
		sqlContextFromAPI, diagsTemp := types.ListValueFrom(ctx, types.StringType, viewResp.SQLContext)
		diags.Append(diagsTemp...)
		data.SQLContext = sqlContextFromAPI
	}

	// Map fields - convert to JSON string
	// View schemas can be arbitrarily deep with nested STRUCT and LIST types,
	// so we use JSON representation instead of trying to model the recursive structure
	fieldsJSON, fieldsDiags := helpers.ConvertTableFieldsToJSON(ctx, viewResp.Fields)
	diags.Append(fieldsDiags...)
	data.Fields = fieldsJSON

	// Map access control list - use helper function
	// For datasources, we always populate from API (no plan to compare against)
	// So we pass an unknown object as the plan parameter to force conversion
	_, _, accessControlAttrTypes := helpers.GetACLAttrTypes()
	if viewResp.AccessControlList == nil {
		data.AccessControlList = types.ObjectNull(accessControlAttrTypes)
	} else {
		// Pass a non-null plan to force conversion
		data.AccessControlList, *diags = helpers.ConvertACLToTerraform(ctx, viewResp.AccessControlList, types.ObjectUnknown(accessControlAttrTypes))
	}
	if diags.HasError() {
		return
	}

	// Map permissions
	if len(viewResp.Permissions) == 0 {
		data.Permissions = types.ListNull(types.StringType)
	} else {
		permsList, d := types.ListValueFrom(ctx, types.StringType, viewResp.Permissions)
		diags.Append(d...)
		data.Permissions = permsList
	}
	if diags.HasError() {
		return
	}

	// Map owner - use helper function
	var ownerDiags diag.Diagnostics
	data.Owner, ownerDiags = helpers.ConvertOwnerToTerraform(ctx, viewResp.Owner)
	diags.Append(ownerDiags...)
}
