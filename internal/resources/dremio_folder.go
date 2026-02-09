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
	_ resource.Resource                = &dremioFolder{}
	_ resource.ResourceWithConfigure   = &dremioFolder{}
	_ resource.ResourceWithImportState = &dremioFolder{}
)

type dremioFolder struct {
	client *dremioClient.Client
}

func NewDremioFolderResource() resource.Resource {
	return &dremioFolder{}
}

// Metadata returns the resource type name.
func (r *dremioFolder) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (r *dremioFolder) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dremioFolder) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *dremioFolder) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioFolderModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.RequestToDremio("DELETE", fmt.Sprintf("/catalog/%s", id), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to delete source, got error: %s", err),
		)
		return
	}
}

func (r *dremioFolder) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Folder resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the folder",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"entity_type": schema.StringAttribute{
				MarkdownDescription: "Entity type (always 'folder')",
				Computed:            true,
				Default:             stringdefault.StaticString("folder"),
			},
			"path": schema.ListAttribute{
				MarkdownDescription: "Full path to the folder",
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
func (r *dremioFolder) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioFolderModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := r.parseResourceToRequestBodyCreate(ctx, &data, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	// Make API request
	api_resp, err := r.client.RequestToDremio("POST", "/catalog", reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to create source, got error: %s", err),
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

	var folderResp models.FolderResponse
	if err := json.Unmarshal(body, &folderResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	// Preserve the plan's config to avoid drift from API-added defaults
	r.fromResponseToState(ctx, &folderResp, &data, &resp.Diagnostics)
	api_resp.Body.Close()

	// The create API does not support ACLs, so we need to update the folder to set them after creation
	if !data.AccessControlList.IsNull() {
		reqBody := r.parseResourceToRequestBodyUpdate(ctx, &data, &resp.Diagnostics)
		if reqBody == nil {
			return
		}

		// Make API request
		api_resp, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/catalog/%s", data.ID.ValueString()), reqBody)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error", fmt.Sprintf("Unable to create source, got error: %s", err),
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

		var folderUpdateResp models.FolderResponse
		if err := json.Unmarshal(body, &folderUpdateResp); err != nil {
			resp.Diagnostics.AddError(
				"Parse Error",
				fmt.Sprintf("Unable to parse response: %s", err),
			)
			return
		}
		// We need to update the tag with the new value from the update response
		r.fromResponseToState(ctx, &folderUpdateResp, &data, &resp.Diagnostics)
		api_resp.Body.Close()
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read resource information.
func (r *dremioFolder) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.DremioFolderModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	var folderResp models.FolderResponse
	folder_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s", id), nil)
	if err != nil {
		// If resource is not found (404), remove it from state so Terraform will recreate it
		if strings.Contains(err.Error(), "status 404") {
			tflog.Warn(ctx, fmt.Sprintf("Folder %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read source, got error: %s", err),
		)
		return
	}
	defer folder_resp.Body.Close()

	resp_body, err := io.ReadAll(folder_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}
	if err := json.Unmarshal(resp_body, &folderResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	// Use API response to detect actual changes during refresh
	r.fromResponseToState(ctx, &folderResp, &state, &resp.Diagnostics)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dremioFolder) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioFolderModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the tag (computed field)
	var state models.DremioFolderModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := r.parseResourceToRequestBodyUpdate(ctx, &plan, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	// Set ID and Tag for optimistic concurrency control
	// Tag comes from state (not plan) because it's a computed field
	id := plan.ID.ValueString()
	reqBody.Tag = state.Tag.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Folder update request with ID: %s, and Tag: %s", id, reqBody.Tag))

	api_resp, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/catalog/%s", id), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to update source, got error: %s", err),
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

	var sourceResp models.FolderResponse
	if err := json.Unmarshal(body, &sourceResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	// Preserve the plan's config to avoid drift from API-added defaults
	r.fromResponseToState(ctx, &sourceResp, &plan, &resp.Diagnostics)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "updated a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dremioFolder) parseResourceToRequestBodyCreate(ctx context.Context, data *models.DremioFolderModel, diags *diag.Diagnostics) *models.FolderCreateRequest {
	// Build the request body
	reqBody := &models.FolderCreateRequest{
		EntityType: "folder",
	}

	// Handle Path
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		var path []string
		diags := data.Path.ElementsAs(ctx, &path, false)
		if diags.HasError() {
			diags.Append(diags...)
			return nil
		}
		reqBody.Path = path
	}

	return reqBody
}

func (r *dremioFolder) parseResourceToRequestBodyUpdate(ctx context.Context, data *models.DremioFolderModel, diags *diag.Diagnostics) *models.FolderUpdateRequest {
	// Build the request body
	reqBody := &models.FolderUpdateRequest{
		EntityType: "folder",
		ID:         data.ID.ValueString(),
		Tag:        data.Tag.ValueString(),
	}

	// Handle Path
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		var path []string
		diags := data.Path.ElementsAs(ctx, &path, false)
		if diags.HasError() {
			diags.Append(diags...)
			return nil
		}
		reqBody.Path = path
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

func (r *dremioFolder) fromResponseToState(ctx context.Context, folderResp *models.FolderResponse, state *models.DremioFolderModel, diags *diag.Diagnostics) {
	state.ID = types.StringValue(folderResp.ID)
	state.Tag = types.StringValue(folderResp.Tag)

	// Access control list block - use helper function
	var aclDiags diag.Diagnostics
	state.AccessControlList, aclDiags = helpers.ConvertACLToTerraform(ctx, folderResp.AccessControlList, state.AccessControlList)
	diags.Append(aclDiags...)
}
