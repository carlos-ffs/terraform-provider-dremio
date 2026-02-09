package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dremioDatasetTags{}
	_ resource.ResourceWithConfigure   = &dremioDatasetTags{}
	_ resource.ResourceWithImportState = &dremioDatasetTags{}
)

type dremioDatasetTags struct {
	client *dremioClient.Client
}

func NewDremioDatasetTagsResource() resource.Resource {
	return &dremioDatasetTags{}
}

// Metadata returns the resource type name.
func (r *dremioDatasetTags) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataset_tags"
}

// Configure adds the provider configured client to the resource.
func (r *dremioDatasetTags) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dremioDatasetTags) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages tags for a Dremio dataset. Tags are case-insensitive labels that can be applied to datasets for organization and discovery.",

		Attributes: map[string]schema.Attribute{
			"dataset_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the dataset",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "List of tags to apply to the dataset. Tags are case-insensitive. Each tag can be listed only once. Tags cannot include special characters (/, :, [, ]). Send empty array to delete all tags.",
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
			"version": schema.StringAttribute{
				MarkdownDescription: "Version identifier for optimistic concurrency control. This value changes with every update.",
				Computed:            true,
			},
		},
	}
}

func (r *dremioDatasetTags) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to dataset_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("dataset_id"), req, resp)
}

func (r *dremioDatasetTags) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioDatasetTagsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	datasetID := state.DatasetID.ValueString()
	version := state.Version.ValueString()

	// Delete requires sending empty tags array with version
	reqBody := models.TagRequest{
		// Current Limitation, if other tags are applied to the dataset outside of Terraform, they will also be removed.
		// Future Improvement: Add support for managing multiple tag resources per dataset.
		// We can get the current tags from the datasource and ONLY remove the ones that are managed by Terraform before sending the request.
		// It may not be possible to manage multiple tag resources per dataset, since the version can be out of sync and lead to conflicts.
		Tags:    []string{},
		Version: version,
	}

	_, err := r.client.RequestToDremio("POST", fmt.Sprintf("/catalog/%s/collaboration/tag", datasetID), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to delete dataset tags, got error: %s", err),
		)
		return
	}
}

// Read resource information.
func (r *dremioDatasetTags) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.DremioDatasetTagsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	datasetID := state.DatasetID.ValueString()

	var tagResp models.TagResponse
	tag_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/collaboration/tag", datasetID), nil)
	if err != nil {
		// If resource is not found (404), remove it from state so Terraform will recreate it
		if strings.Contains(err.Error(), "status 404") {
			tflog.Warn(ctx, fmt.Sprintf("Dataset tags for %s not found, removing from state", datasetID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read dataset tags, got error: %s", err),
		)
		return
	}
	defer tag_resp.Body.Close()

	resp_body, err := io.ReadAll(tag_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}
	if err := json.Unmarshal(resp_body, &tagResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &tagResp, &state, &resp.Diagnostics)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource.
func (r *dremioDatasetTags) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioDatasetTagsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	datasetID := data.DatasetID.ValueString()

	// First, we try to GET existing tags to check if they already exist (even if empty)
	// This is necessary because a dataset might already have an empty tags array with a version
	tflog.Debug(ctx, fmt.Sprintf("Checking for existing tags on dataset: %s", datasetID))

	existing_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/collaboration/tag", datasetID), nil)
	var existingVersion string

	if err == nil {
		// Tags exist, read the version
		defer existing_resp.Body.Close()
		existing_body, read_err := io.ReadAll(existing_resp.Body)
		if read_err == nil {
			var existingTagResp models.TagResponse
			if json.Unmarshal(existing_body, &existingTagResp) == nil {
				existingVersion = existingTagResp.Version
				tflog.Debug(ctx, fmt.Sprintf("Found existing tags with version: %s", existingVersion))
			}
		}
	} else {
		// No existing tags found, we'll create new ones
		tflog.Debug(ctx, fmt.Sprintf("No existing tags found, will create new: %s", err))
	}

	reqBody := r.parseResourceToRequestBody(ctx, &data, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	var api_resp *http.Response
	var method string

	// If we found an existing version, use PUT (update), otherwise use POST (create)
	if existingVersion != "" {
		reqBody.Version = existingVersion
		tflog.Debug(ctx, fmt.Sprintf("Updating existing tags with version: %s", existingVersion))
	} else {
		tflog.Debug(ctx, "Creating new tags")
	}

	api_resp, err = r.client.RequestToDremio("POST", fmt.Sprintf("/catalog/%s/collaboration/tag", datasetID), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to %s dataset tags, got error: %s", strings.ToLower(method), err),
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

	var tagResp models.TagResponse
	if err := json.Unmarshal(body, &tagResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &tagResp, &data, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created dataset tags resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dremioDatasetTags) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioDatasetTagsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the version (computed field)
	var state models.DremioDatasetTagsModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := r.parseResourceToRequestBody(ctx, &plan, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	// Set Version for optimistic concurrency control
	// Version comes from state (not plan) because it's a computed field
	datasetID := plan.DatasetID.ValueString()
	reqBody.Version = state.Version.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Dataset tags update request for dataset: %s, with Version: %s", datasetID, reqBody.Version))

	api_resp, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/catalog/%s/collaboration/tag", datasetID), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to update dataset tags, got error: %s", err),
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

	var tagResp models.TagResponse
	if err := json.Unmarshal(body, &tagResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &tagResp, &plan, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated dataset tags resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// fromResponseToState updates the state with values from the API response.
func (r *dremioDatasetTags) fromResponseToState(ctx context.Context, tagResp *models.TagResponse, state *models.DremioDatasetTagsModel, diags *diag.Diagnostics) {
	// Convert tags to Terraform list
	tagsList, d := types.ListValueFrom(ctx, types.StringType, tagResp.Tags)
	diags.Append(d...)
	state.Tags = tagsList
	state.Version = types.StringValue(tagResp.Version)

	tflog.Debug(ctx, fmt.Sprintf("fromResponseToState: Version=%s, Tags=%v", tagResp.Version, tagResp.Tags))
}

// parseResourceToRequestBody converts Terraform state/plan to API request body.
func (r *dremioDatasetTags) parseResourceToRequestBody(ctx context.Context, data *models.DremioDatasetTagsModel, diags *diag.Diagnostics) *models.TagRequest {
	// Convert tags list to string slice
	var tags []string
	d := data.Tags.ElementsAs(ctx, &tags, false)
	diags.Append(d...)
	if diags.HasError() {
		return nil
	}

	// Build the request body
	reqBody := &models.TagRequest{
		Tags: tags,
		// Version is omitted for create, set from state for update
	}

	return reqBody
}
