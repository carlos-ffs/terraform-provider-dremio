package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/helpers"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dremioGrants{}
	_ resource.ResourceWithConfigure   = &dremioGrants{}
	_ resource.ResourceWithImportState = &dremioGrants{}
)

type dremioGrants struct {
	client *dremioClient.Client
}

func NewDremioGrantsResource() resource.Resource {
	return &dremioGrants{}
}

// Metadata returns the resource type name.
func (r *dremioGrants) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_grants"
}

// Configure adds the provider configured client to the resource.
func (r *dremioGrants) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Schema defines the schema for the resource.
func (r *dremioGrants) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages grants (privileges) on a Dremio catalog object. This resource allows you to grant privileges to users and roles on catalog objects such as sources, spaces, folders, datasets, views, and UDFs.\n\n**Important:** This resource manages ALL grants on the catalog object. When this resource is created, it will **overwrite** any existing grants on the object. When destroyed, it will remove all grants from the object.",

		Attributes: map[string]schema.Attribute{
			"catalog_object_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the Dremio catalog object to manage grants for.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"grants": schema.SetNestedAttribute{
				MarkdownDescription: "Set of grants to apply to the catalog object. Each grant specifies a user or role and the privileges to grant. If empty, all explicit grants will be removed from the object.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "UUID of the user or role to grant privileges to.",
							Required:            true,
						},
						"grantee_type": schema.StringAttribute{
							MarkdownDescription: "Type of grantee. Must be 'USER' or 'ROLE'.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("USER", "ROLE"),
							},
						},
						"privileges": schema.SetAttribute{
							MarkdownDescription: "Set of privileges to grant. Available privileges depend on the catalog object type. Common privileges include: ALTER, SELECT, MANAGE_GRANTS, DELETE, INSERT, TRUNCATE, UPDATE, DROP, CREATE_TABLE, MODIFY, READ_METADATA, ALTER_REFLECTION, VIEW_REFLECTION.",
							Required:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
			"available_privileges": schema.ListAttribute{
				MarkdownDescription: "List of available privileges for this catalog object type. This is computed from the API response.",
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *dremioGrants) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to catalog_object_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("catalog_object_id"), req, resp)
}

func (r *dremioGrants) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioGrantsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogObjectID := state.CatalogObjectID.ValueString()

	// Delete by setting grants to empty array
	// There is no delete operation in the API, the correspondent delete operation is setting the grants to an empty array.
	reqBody := models.GrantsRequest{
		Grants: []models.GranteeRequest{},
	}

	tflog.Debug(ctx, fmt.Sprintf("Deleting grants for catalog object: %s (setting to empty array)", catalogObjectID))

	_, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/catalog/%s/grants", catalogObjectID), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to delete grants, got error: %s", err),
		)
		return
	}
}

// Read resource information.
func (r *dremioGrants) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.DremioGrantsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogObjectID := state.CatalogObjectID.ValueString()

	var grantsResp models.GrantsResponse
	apiResp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/grants", catalogObjectID), nil)
	if err != nil {
		// If resource is not found (404), remove it from state so Terraform will recreate it
		if strings.Contains(err.Error(), "status 404") {
			tflog.Warn(ctx, fmt.Sprintf("Grants for catalog object %s not found, removing from state", catalogObjectID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read grants, got error: %s", err),
		)
		return
	}
	defer apiResp.Body.Close()

	respBody, err := io.ReadAll(apiResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}
	if err := json.Unmarshal(respBody, &grantsResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &grantsResp, &state, &resp.Diagnostics)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource.
func (r *dremioGrants) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioGrantsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	catalogObjectID := data.CatalogObjectID.ValueString()

	// First, check for existing grants and warn the user if they exist
	tflog.Debug(ctx, fmt.Sprintf("Checking for existing grants on catalog object: %s", catalogObjectID))

	existingResp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/grants", catalogObjectID), nil)
	if err == nil {
		defer existingResp.Body.Close()
		existingBody, readErr := io.ReadAll(existingResp.Body)
		if readErr == nil {
			var existingGrantsResp models.GrantsResponse
			if json.Unmarshal(existingBody, &existingGrantsResp) == nil {
				if len(existingGrantsResp.Grants) > 0 {
					// Warn the user that existing grants will be overwritten
					resp.Diagnostics.AddWarning(
						"Existing Grants Will Be Overwritten",
						fmt.Sprintf("The catalog object %s already has %d existing grant(s). These grants will be overwritten with the grants specified in this resource. Existing grants: %v",
							catalogObjectID, len(existingGrantsResp.Grants), formatExistingGrants(existingGrantsResp.Grants)),
					)
				}
			}
		}
	}

	reqBody := r.parseResourceToRequestBody(ctx, &data, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating grants for catalog object: %s", catalogObjectID))

	apiResp, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/catalog/%s/grants", catalogObjectID), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to create grants, got error: %s", err),
		)
		return
	}
	defer apiResp.Body.Close()

	// PUT returns 204 No Content, so we need to GET to retrieve the current state
	getResp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/grants", catalogObjectID), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read grants after creation, got error: %s", err),
		)
		return
	}
	defer getResp.Body.Close()

	body, err := io.ReadAll(getResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}

	var grantsResp models.GrantsResponse
	if err := json.Unmarshal(body, &grantsResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &grantsResp, &data, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created grants resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dremioGrants) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioGrantsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogObjectID := plan.CatalogObjectID.ValueString()

	reqBody := r.parseResourceToRequestBody(ctx, &plan, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating grants for catalog object: %s", catalogObjectID))

	apiResp, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/catalog/%s/grants", catalogObjectID), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to update grants, got error: %s", err),
		)
		return
	}
	defer apiResp.Body.Close()

	// PUT returns 204 No Content, so we need to GET to retrieve the current state
	getResp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/grants", catalogObjectID), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read grants after update, got error: %s", err),
		)
		return
	}
	defer getResp.Body.Close()

	body, err := io.ReadAll(getResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}

	var grantsResp models.GrantsResponse
	if err := json.Unmarshal(body, &grantsResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &grantsResp, &plan, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated grants resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// formatExistingGrants formats existing grants for display in warning messages.
func formatExistingGrants(grants []models.GranteesResponse) string {
	var parts []string
	for _, g := range grants {
		parts = append(parts, fmt.Sprintf("%s (%s): %v", g.Name, g.GranteeType, g.Privileges))
	}
	return strings.Join(parts, "; ")
}

// fromResponseToState updates the state with values from the API response.
func (r *dremioGrants) fromResponseToState(ctx context.Context, grantsResp *models.GrantsResponse, state *models.DremioGrantsModel, diags *diag.Diagnostics) {
	// Set the catalog object ID from response
	state.CatalogObjectID = types.StringValue(grantsResp.ID)

	// Convert available privileges to Terraform list
	availablePrivsList, d := types.ListValueFrom(ctx, types.StringType, grantsResp.AvailablePrivileges)
	diags.Append(d...)
	state.AvailablePrivileges = availablePrivsList

	// Convert grants to Terraform set using helper
	grantsList, d := helpers.ConvertGranteesToTerraform(ctx, grantsResp.Grants)
	diags.Append(d...)
	state.Grants = grantsList

	tflog.Debug(ctx, fmt.Sprintf("fromResponseToState: CatalogObjectID=%s, Grants count=%d, AvailablePrivileges=%v",
		grantsResp.ID, len(grantsResp.Grants), grantsResp.AvailablePrivileges))
}

// parseResourceToRequestBody converts Terraform state/plan to API request body.
func (r *dremioGrants) parseResourceToRequestBody(ctx context.Context, data *models.DremioGrantsModel, diags *diag.Diagnostics) *models.GrantsRequest {
	// Convert grants set to GranteeRequest slice using helper
	grants, d := helpers.ConvertGranteeSetFromTerraform(ctx, data.Grants)
	diags.Append(d...)
	if diags.HasError() {
		return nil
	}

	// Build the request body
	reqBody := &models.GrantsRequest{
		Grants: grants,
	}

	return reqBody
}
