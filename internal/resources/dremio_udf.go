package resources

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
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &dremioUDF{}
	_ resource.ResourceWithConfigure   = &dremioUDF{}
	_ resource.ResourceWithImportState = &dremioUDF{}
)

type dremioUDF struct {
	client *dremioClient.Client
}

func NewDremioUDFResource() resource.Resource {
	return &dremioUDF{}
}

// Metadata returns the resource type name.
func (r *dremioUDF) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_udf"
}

func (r *dremioUDF) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dremioClient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *dremioClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *dremioUDF) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *dremioUDF) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioUDFModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.RequestToDremio("DELETE", fmt.Sprintf("/catalog/%s", id), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to delete UDF, got error: %s", err),
		)
		return
	}
}

func (r *dremioUDF) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio User-Defined Function (UDF) resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the UDF",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"entity_type": schema.StringAttribute{
				MarkdownDescription: "Entity type (always 'function')",
				Computed:            true,
				Default:             stringdefault.StaticString("function"),
			},
			"path": schema.ListAttribute{
				MarkdownDescription: "Full path to the UDF, including the function name as the last element",
				Required:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[^/:[\]]*$`),
							"path elements must not contain the characters: /, :, [, ]",
						),
					),
				},
			},
			"is_scalar": schema.BoolAttribute{
				MarkdownDescription: "If true, the UDF is a scalar function. If false, the UDF is a tabular function",
				Required:            true,
			},
			"function_arg_list": schema.StringAttribute{
				MarkdownDescription: "The name of each argument in the UDF and the argument's data type. Separate the name and data type with a single space. If the function includes multiple arguments, separate the arguments with a comma. Example: 'domain VARCHAR, orderdate DATE'",
				Required:            true,
			},
			"function_body": schema.StringAttribute{
				MarkdownDescription: "The SQL statement that the UDF should execute",
				Required:            true,
			},
			"return_type": schema.StringAttribute{
				MarkdownDescription: "The data type of the result that the function returns (for scalar functions) or of each column that the function returns, separated by commas (for tabular functions). Example: 'name VARCHAR, email VARCHAR, order_date DATE'",
				Required:            true,
			},
			"access_control_list": schema.SingleNestedAttribute{
				MarkdownDescription: "User and role access settings",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"users": schema.ListNestedAttribute{
						MarkdownDescription: "List of user access controls",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "User ID",
									Required:            true,
								},
								"permissions": schema.ListAttribute{
									MarkdownDescription: "List of permissions",
									Required:            true,
									ElementType:         types.StringType,
								},
							},
						},
					},
					"roles": schema.ListNestedAttribute{
						MarkdownDescription: "List of role access controls",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "Role ID",
									Required:            true,
								},
								"permissions": schema.ListAttribute{
									MarkdownDescription: "List of permissions",
									Required:            true,
									ElementType:         types.StringType,
								},
							},
						},
					},
				},
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Version tag for optimistic concurrency control. This value changes with every update.",
				Computed:            true,
			},
		},
	}
}

// Create a new resource.
func (r *dremioUDF) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioUDFModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := r.parseResourceToRequestBody(ctx, &data, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	// Make API request
	api_resp, err := r.client.RequestToDremio("POST", "/catalog", reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to create UDF, got error: %s", err),
		)
		return
	}

	// Parse response
	body, err := io.ReadAll(api_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}

	var udfResp models.UDFResponse
	if err := json.Unmarshal(body, &udfResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &udfResp, &data, &resp.Diagnostics)
	api_resp.Body.Close()

	// The create API does not support ACLs, so we need to update the UDF to set them after creation
	if !data.AccessControlList.IsNull() {
		reqBody := r.parseResourceToRequestBody(ctx, &data, &resp.Diagnostics)
		if reqBody == nil {
			return
		}

		// Set ID and Tag for update
		reqBody.ID = data.ID.ValueString()
		reqBody.Tag = data.Tag.ValueString()

		// Make API request
		api_resp, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/catalog/%s", data.ID.ValueString()), reqBody)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error", fmt.Sprintf("Unable to set ACL on UDF, got error: %s", err),
			)
			return
		}
		body, err := io.ReadAll(api_resp.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Read Error",
				fmt.Sprintf("Unable to read response body: %s", err),
			)
			return
		}

		var udfUpdateResp models.UDFResponse
		if err := json.Unmarshal(body, &udfUpdateResp); err != nil {
			resp.Diagnostics.AddError(
				"Parse Error",
				fmt.Sprintf("Unable to parse response: %s", err),
			)
			return
		}
		// We need to update the tag with the new value from the update response
		r.fromResponseToState(ctx, &udfUpdateResp, &data, &resp.Diagnostics)
		api_resp.Body.Close()
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a UDF resource")
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read resource information.
func (r *dremioUDF) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.DremioUDFModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	var udfResp models.UDFResponse
	udf_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s", id), nil)
	if err != nil {
		// If resource is not found (404), remove it from state so Terraform will recreate it
		if strings.Contains(err.Error(), "status 404") {
			tflog.Warn(ctx, fmt.Sprintf("UDF %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read UDF, got error: %s", err),
		)
		return
	}
	defer udf_resp.Body.Close()

	resp_body, err := io.ReadAll(udf_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}
	if err := json.Unmarshal(resp_body, &udfResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &udfResp, &state, &resp.Diagnostics)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dremioUDF) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioUDFModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the tag (computed field)
	var state models.DremioUDFModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := r.parseResourceToRequestBody(ctx, &plan, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	// Set ID and Tag for optimistic concurrency control
	// Tag comes from state (not plan) because it's a computed field
	id := plan.ID.ValueString()
	reqBody.ID = id
	reqBody.Tag = state.Tag.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("UDF update request with ID: %s, and Tag: %s", id, reqBody.Tag))

	api_resp, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/catalog/%s", id), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to update UDF, got error: %s", err),
		)
		return
	}
	defer api_resp.Body.Close()

	// Parse response
	body, err := io.ReadAll(api_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}

	var udfResp models.UDFResponse
	if err := json.Unmarshal(body, &udfResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &udfResp, &plan, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated a UDF resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dremioUDF) parseResourceToRequestBody(ctx context.Context, data *models.DremioUDFModel, diags *diag.Diagnostics) *models.UDFRequest {
	// Build the request body
	reqBody := &models.UDFRequest{
		EntityType: "function",
	}

	// Handle Path
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		var path []string
		diagsL := data.Path.ElementsAs(ctx, &path, false)
		if diagsL.HasError() {
			diags.Append(diagsL...)
			return nil
		}
		reqBody.Path = path
	}

	// Handle IsScalar
	reqBody.IsScalar = data.IsScalar.ValueBool()

	// Handle FunctionArgList
	reqBody.FunctionArgList = data.FunctionArgList.ValueString()

	// Handle FunctionBody
	reqBody.FunctionBody = data.FunctionBody.ValueString()

	// Handle ReturnType
	reqBody.ReturnType = data.ReturnType.ValueString()

	// Handle AccessControlList - use helper function
	var aclDiags diag.Diagnostics
	reqBody.AccessControlList, aclDiags = helpers.ConvertACLFromTerraform(ctx, data.AccessControlList)
	if aclDiags.HasError() {
		diags.Append(aclDiags...)
		return nil
	}

	return reqBody
}

func (r *dremioUDF) fromResponseToState(ctx context.Context, udfResp *models.UDFResponse, state *models.DremioUDFModel, diags *diag.Diagnostics) {
	if udfResp.ID != nil {
		state.ID = types.StringValue(*udfResp.ID)
	}
	if udfResp.Tag != nil {
		state.Tag = types.StringValue(*udfResp.Tag)
	}

	// Access control list block - use helper function
	var aclDiags diag.Diagnostics
	state.AccessControlList, aclDiags = helpers.ConvertACLToTerraform(ctx, udfResp.AccessControlList, state.AccessControlList)
	diags.Append(aclDiags...)
}
