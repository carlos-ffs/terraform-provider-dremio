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
	_ resource.Resource                = &dremioView{}
	_ resource.ResourceWithConfigure   = &dremioView{}
	_ resource.ResourceWithImportState = &dremioView{}
)

type dremioView struct {
	client *dremioClient.Client
}

func NewDremioViewResource() resource.Resource {
	return &dremioView{}
}

// Metadata returns the resource type name.
func (r *dremioView) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_view"
}

func (r *dremioView) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dremioView) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *dremioView) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioViewModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.RequestToDremio("DELETE", fmt.Sprintf("/catalog/%s", id), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to delete view, got error: %s", err),
		)
		return
	}
}

func (r *dremioView) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio View resource - manages a virtual dataset (view) in Dremio",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the view",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"entity_type": schema.StringAttribute{
				MarkdownDescription: "Entity type (always 'dataset')",
				Computed:            true,
				Default:             stringdefault.StaticString("dataset"),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Dataset type (always 'VIRTUAL_DATASET')",
				Computed:            true,
				Default:             stringdefault.StaticString("VIRTUAL_DATASET"),
			},
			"path": schema.ListAttribute{
				MarkdownDescription: "Full path to the view, including the view name as the last element",
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
			"sql": schema.StringAttribute{
				MarkdownDescription: "SQL query defining the view",
				Required:            true,
			},
			"sql_context": schema.ListAttribute{
				MarkdownDescription: "Context for SQL query execution (optional)",
				Optional:            true,
				ElementType:         types.StringType,
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
			"fields": schema.StringAttribute{
				MarkdownDescription: "View fields/columns as JSON string. Due to the recursive nature of view schemas (STRUCT and LIST types can be arbitrarily nested), fields are represented as a JSON string. Use jsondecode() to parse this value in Terraform configurations.",
				Computed:            true,
			},
		},
	}
}

// Create a new resource.
func (r *dremioView) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioViewModel

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
			"Client Error", fmt.Sprintf("Unable to create view, got error: %s", err),
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

	var viewResp models.ViewResponse
	if err := json.Unmarshal(body, &viewResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &viewResp, &data, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a view resource")
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read resource information.
func (r *dremioView) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.DremioViewModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	var viewResp models.ViewResponse
	view_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s", id), nil)
	if err != nil {
		// If resource is not found (404), remove it from state so Terraform will recreate it
		if strings.Contains(err.Error(), "status 404") {
			tflog.Warn(ctx, fmt.Sprintf("View %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read view, got error: %s", err),
		)
		return
	}
	defer view_resp.Body.Close()

	resp_body, err := io.ReadAll(view_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}
	if err := json.Unmarshal(resp_body, &viewResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &viewResp, &state, &resp.Diagnostics)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dremioView) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioViewModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the tag (computed field)
	var state models.DremioViewModel
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

	tflog.Debug(ctx, fmt.Sprintf("View update request with ID: %s, and Tag: %s", id, reqBody.Tag))

	api_resp, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/catalog/%s", id), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to update view, got error: %s", err),
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

	var viewResp models.ViewResponse
	if err := json.Unmarshal(body, &viewResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &viewResp, &plan, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated a view resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dremioView) parseResourceToRequestBody(ctx context.Context, data *models.DremioViewModel, diags *diag.Diagnostics) *models.ViewRequest {
	// Build the request body
	reqBody := &models.ViewRequest{
		EntityType: "dataset",
		Type:       "VIRTUAL_DATASET",
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

	// Handle SQL
	reqBody.SQL = data.SQL.ValueString()

	// Handle SQLContext
	if !data.SQLContext.IsNull() && !data.SQLContext.IsUnknown() {
		var sqlContext []string
		diagsL := data.SQLContext.ElementsAs(ctx, &sqlContext, false)
		if diagsL.HasError() {
			diags.Append(diagsL...)
			return nil
		}
		reqBody.SQLContext = sqlContext
	}

	// Handle AccessControlList - use helper function
	var aclDiags diag.Diagnostics
	reqBody.AccessControlList, aclDiags = helpers.ConvertACLFromTerraform(ctx, data.AccessControlList)
	if aclDiags.HasError() {
		diags.Append(aclDiags...)
		return nil
	}

	return reqBody
}

func (r *dremioView) fromResponseToState(ctx context.Context, viewResp *models.ViewResponse, state *models.DremioViewModel, diags *diag.Diagnostics) {
	state.ID = types.StringValue(viewResp.ID)
	state.EntityType = types.StringValue("dataset")
	state.Type = types.StringValue("VIRTUAL_DATASET")
	state.Tag = types.StringValue(viewResp.Tag)

	// Access control list block - use helper function
	var aclDiags diag.Diagnostics
	state.AccessControlList, aclDiags = helpers.ConvertACLToTerraform(ctx, viewResp.AccessControlList, state.AccessControlList)
	diags.Append(aclDiags...)

	// Map fields - convert to JSON string
	// View schemas can be arbitrarily deep with nested STRUCT and LIST types,
	// so we use JSON representation instead of trying to model the recursive structure
	fieldsJSON, fieldsDiags := helpers.ConvertTableFieldsToJSON(ctx, viewResp.Fields)
	diags.Append(fieldsDiags...)
	state.Fields = fieldsJSON
}
