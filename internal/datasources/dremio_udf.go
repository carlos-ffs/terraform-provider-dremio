package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/helpers"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                     = &dremioUDFDataSource{}
	_ datasource.DataSourceWithConfigure        = &dremioUDFDataSource{}
	_ datasource.DataSourceWithConfigValidators = &dremioUDFDataSource{}
)

func NewDremioUDFDataSource() datasource.DataSource {
	return &dremioUDFDataSource{}
}

type dremioUDFDataSource struct {
	client *dremioClient.Client
}

// ConfigValidators returns a list of functions which will all be performed during validation.
func (d *dremioUDFDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("path"),
		),
	}
}

// Metadata returns the data source type name.
func (d *dremioUDFDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_udf"
}

func (d *dremioUDFDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dremioUDFDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio UDF data source - retrieves information about an existing user-defined function",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the UDF. Exactly one of `id` or `path` must be specified.",
				Computed:            true,
				Optional:            true,
			},
			"path": schema.ListAttribute{
				MarkdownDescription: "Full path to the UDF. Exactly one of `id` or `path` must be specified.",
				Computed:            true,
				Optional:            true,
				ElementType:         types.StringType,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Version tag for optimistic concurrency control",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Date and time the UDF was created (UTC)",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "Date and time the UDF was last modified (UTC)",
				Computed:            true,
			},
			"is_scalar": schema.BoolAttribute{
				MarkdownDescription: "Whether the function is scalar",
				Computed:            true,
			},
			"function_arg_list": schema.StringAttribute{
				MarkdownDescription: "Function arguments as a string",
				Computed:            true,
			},
			"function_body": schema.StringAttribute{
				MarkdownDescription: "SQL body of the function",
				Computed:            true,
			},
			"return_type": schema.StringAttribute{
				MarkdownDescription: "Return type as a string",
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
				MarkdownDescription: "User's permissions on the UDF",
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
func (d *dremioUDFDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioUDFDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var udfID string
	if !data.ID.IsNull() {
		udfID = data.ID.ValueString()
	}

	var udf_path []string
	if !data.Path.IsNull() {
		diags := data.Path.ElementsAs(ctx, &udf_path, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if udfID == "" && len(udf_path) == 0 {
		resp.Diagnostics.AddError(
			"Missing UDF ID or Path",
			"Either `id` or `path` must be specified for Dremio UDF data source.",
		)
		return
	}
	if udfID != "" && len(udf_path) > 0 {
		resp.Diagnostics.AddError(
			"Both UDF ID and Path specified",
			"Only one of `id` or `path` must be specified for Dremio UDF data source.",
		)
		return
	}

	var path string
	if udfID != "" { // Read UDF by ID
		path = fmt.Sprintf("/catalog/%s", udfID)
	} else { // Lookup UDF ID by path
		udf_path_str := "/" + strings.Join(udf_path, "/")
		path = fmt.Sprintf("/catalog/by-path/%s", udf_path_str)
	}

	api_resp, err := d.client.RequestToDremio("GET", path, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to request UDF: %s", err),
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

	var udfResp models.UDFResponse
	if err := json.Unmarshal(api_resp_body, &udfResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, &udfResp, &data, &resp.Diagnostics)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioUDFDataSource) mapResponseToState(ctx context.Context, udfResp *models.UDFResponse, data *models.DremioUDFDataSourceModel, diags *diag.Diagnostics) {

	// Map basic fields
	if udfResp.ID != nil {
		data.ID = types.StringValue(*udfResp.ID)
	} else {
		data.ID = types.StringNull()
	}

	if udfResp.Tag != nil {
		data.Tag = types.StringValue(*udfResp.Tag)
	} else {
		data.Tag = types.StringNull()
	}

	if len(udfResp.Path) == 0 {
		data.Path = types.ListNull(types.StringType)
	} else {
		pathFromAPI, diagsTemp := types.ListValueFrom(ctx, types.StringType, udfResp.Path)
		diags.Append(diagsTemp...)
		data.Path = pathFromAPI
	}

	if udfResp.CreatedAt != nil {
		data.CreatedAt = types.StringValue(*udfResp.CreatedAt)
	} else {
		data.CreatedAt = types.StringNull()
	}

	if udfResp.LastModified != nil {
		data.LastModified = types.StringValue(*udfResp.LastModified)
	} else {
		data.LastModified = types.StringNull()
	}

	if udfResp.IsScalar != nil {
		data.IsScalar = types.BoolValue(*udfResp.IsScalar)
	} else {
		data.IsScalar = types.BoolNull()
	}

	if udfResp.FunctionArgList != nil {
		data.FunctionArgList = types.StringValue(*udfResp.FunctionArgList)
	} else {
		data.FunctionArgList = types.StringNull()
	}

	if udfResp.FunctionBody != nil {
		data.FunctionBody = types.StringValue(*udfResp.FunctionBody)
	} else {
		data.FunctionBody = types.StringNull()
	}

	if udfResp.ReturnType != nil {
		data.ReturnType = types.StringValue(*udfResp.ReturnType)
	} else {
		data.ReturnType = types.StringNull()
	}

	// Map access control list - use helper function
	// For datasources, we always populate from API (no plan to compare against)
	// So we pass an unknown object as the plan parameter to force conversion
	_, _, accessControlAttrTypes := helpers.GetACLAttrTypes()
	if udfResp.AccessControlList == nil {
		data.AccessControlList = types.ObjectNull(accessControlAttrTypes)
	} else {
		// Pass a non-null plan to force conversion
		data.AccessControlList, *diags = helpers.ConvertACLToTerraform(ctx, udfResp.AccessControlList, types.ObjectUnknown(accessControlAttrTypes))
	}
	if diags.HasError() {
		return
	}

	// Map permissions
	if len(udfResp.Permissions) == 0 {
		data.Permissions = types.ListNull(types.StringType)
	} else {
		permsList, d := types.ListValueFrom(ctx, types.StringType, udfResp.Permissions)
		diags.Append(d...)
		data.Permissions = permsList
	}
	if diags.HasError() {
		return
	}

	// Map owner - use helper function
	var ownerDiags diag.Diagnostics
	data.Owner, ownerDiags = helpers.ConvertOwnerToTerraform(ctx, udfResp.Owner)
	diags.Append(ownerDiags...)
}
