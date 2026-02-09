package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dremioDatasetWiki{}
	_ resource.ResourceWithConfigure   = &dremioDatasetWiki{}
	_ resource.ResourceWithImportState = &dremioDatasetWiki{}
)

type dremioDatasetWiki struct {
	client *dremioClient.Client
}

func NewDremioDatasetWikiResource() resource.Resource {
	return &dremioDatasetWiki{}
}

// Metadata returns the resource type name.
func (r *dremioDatasetWiki) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataset_wiki"
}

// Configure adds the provider configured client to the resource.
func (r *dremioDatasetWiki) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dremioDatasetWiki) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages wiki content for a Dremio dataset. Wiki content uses GitHub-flavored Markdown for formatting.",

		Attributes: map[string]schema.Attribute{
			"dataset_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the source, folder, or dataset for which to manage the wiki",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"text": schema.StringAttribute{
				MarkdownDescription: "Text to display in the wiki. Use GitHub-flavored Markdown for wiki formatting and \\n for new lines and blank lines. Each wiki may have a maximum of 100,000 characters.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Number for the most recent version of the wiki, starting with 0. This value changes with every update.",
				Computed:            true,
			},
		},
	}
}

func (r *dremioDatasetWiki) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to dataset_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("dataset_id"), req, resp)
}

func (r *dremioDatasetWiki) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioDatasetWikiModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	datasetID := state.DatasetID.ValueString()
	version := int(state.Version.ValueInt64())

	// Delete requires sending empty text with version
	reqBody := models.WikiRequest{
		Text:    "",
		Version: &version,
	}

	_, err := r.client.RequestToDremio("POST", fmt.Sprintf("/catalog/%s/collaboration/wiki", datasetID), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to delete dataset wiki, got error: %s", err),
		)
		return
	}
}

// Read resource information.
func (r *dremioDatasetWiki) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.DremioDatasetWikiModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	datasetID := state.DatasetID.ValueString()

	var wikiResp models.WikiResponse
	wiki_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/collaboration/wiki", datasetID), nil)
	if err != nil {
		// If resource is not found (404), remove it from state so Terraform will recreate it
		if strings.Contains(err.Error(), "status 404") {
			tflog.Warn(ctx, fmt.Sprintf("Dataset wiki for %s not found, removing from state", datasetID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read dataset wiki, got error: %s", err),
		)
		return
	}
	defer wiki_resp.Body.Close()

	resp_body, err := io.ReadAll(wiki_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}
	if err := json.Unmarshal(resp_body, &wikiResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &wikiResp, &state, &resp.Diagnostics)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource.
func (r *dremioDatasetWiki) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioDatasetWikiModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	datasetID := data.DatasetID.ValueString()

	// First, we try to GET existing wiki to check if it already exists
	// This is necessary because a dataset might already have a wiki with a version
	tflog.Debug(ctx, fmt.Sprintf("Checking for existing wiki on dataset: %s", datasetID))

	existing_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s/collaboration/wiki", datasetID), nil)
	var existingVersion *int

	if err == nil {
		// Wiki exists, read the version
		defer existing_resp.Body.Close()
		existing_body, read_err := io.ReadAll(existing_resp.Body)
		if read_err == nil {
			var existingWikiResp models.WikiResponse
			if json.Unmarshal(existing_body, &existingWikiResp) == nil {
				existingVersion = &existingWikiResp.Version
				tflog.Debug(ctx, fmt.Sprintf("Found existing wiki with version: %d", *existingVersion))
			}
		}
	} else {
		// No existing wiki found, we'll create new one
		tflog.Debug(ctx, fmt.Sprintf("No existing wiki found, will create new: %s", err))
	}

	reqBody := r.parseResourceToRequestBody(ctx, &data, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	var api_resp *http.Response

	// If we found an existing version, include it in the request (update)
	if existingVersion != nil {
		reqBody.Version = existingVersion
		tflog.Debug(ctx, fmt.Sprintf("Updating existing wiki with version: %d", *existingVersion))
	} else {
		tflog.Debug(ctx, "Creating new wiki")
	}

	api_resp, err = r.client.RequestToDremio("POST", fmt.Sprintf("/catalog/%s/collaboration/wiki", datasetID), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to create dataset wiki, got error: %s", err),
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

	var wikiResp models.WikiResponse
	if err := json.Unmarshal(body, &wikiResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &wikiResp, &data, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created dataset wiki resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dremioDatasetWiki) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioDatasetWikiModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the version (computed field)
	var state models.DremioDatasetWikiModel
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
	version := int(state.Version.ValueInt64())
	reqBody.Version = &version

	tflog.Debug(ctx, fmt.Sprintf("Dataset wiki update request for dataset: %s, with Version: %d", datasetID, version))

	api_resp, err := r.client.RequestToDremio("POST", fmt.Sprintf("/catalog/%s/collaboration/wiki", datasetID), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to update dataset wiki, got error: %s", err),
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

	var wikiResp models.WikiResponse
	if err := json.Unmarshal(body, &wikiResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &wikiResp, &plan, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated dataset wiki resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// fromResponseToState updates the state with values from the API response.
func (r *dremioDatasetWiki) fromResponseToState(ctx context.Context, wikiResp *models.WikiResponse, state *models.DremioDatasetWikiModel, diags *diag.Diagnostics) {
	state.Text = types.StringValue(wikiResp.Text)
	state.Version = types.Int64Value(int64(wikiResp.Version))

	tflog.Debug(ctx, fmt.Sprintf("fromResponseToState: Version=%d, Text length=%d", wikiResp.Version, len(wikiResp.Text)))
}

// parseResourceToRequestBody converts Terraform state/plan to API request body.
func (r *dremioDatasetWiki) parseResourceToRequestBody(ctx context.Context, data *models.DremioDatasetWikiModel, diags *diag.Diagnostics) *models.WikiRequest {
	// Build the request body
	reqBody := &models.WikiRequest{
		Text: data.Text.ValueString(),
		// Version is omitted for create, set from state for update
	}

	return reqBody
}
