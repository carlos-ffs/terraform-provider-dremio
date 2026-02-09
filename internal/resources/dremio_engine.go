package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &dremioEngine{}
	_ resource.ResourceWithConfigure   = &dremioEngine{}
	_ resource.ResourceWithImportState = &dremioEngine{}
)

type dremioEngine struct {
	client *dremioClient.Client
}

func NewDremioEngineResource() resource.Resource {
	return &dremioEngine{}
}

func (r *dremioEngine) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engine"
}

func (r *dremioEngine) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dremioEngine) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Engine resource - manages compute engines in Dremio Cloud",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the engine (UUID)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "User-defined name for the engine",
				Required:            true,
			},
			"size": schema.StringAttribute{
				MarkdownDescription: "Size of the engine (XX_SMALL_V1, X_SMALL_V1, SMALL_V1, MEDIUM_V1, LARGE_V1, X_LARGE_V1, XX_LARGE_V1, XXX_LARGE_V1)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("XX_SMALL_V1", "X_SMALL_V1", "SMALL_V1", "MEDIUM_V1", "LARGE_V1", "X_LARGE_V1", "XX_LARGE_V1", "XXX_LARGE_V1"),
				},
			},
			"min_replicas": schema.Int64Attribute{
				MarkdownDescription: "Minimum number of engine replicas that will be enabled at any given time",
				Required:            true,
			},
			"max_replicas": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of engine replicas that will be enabled at any given time",
				Required:            true,
			},
			"auto_stop_delay_seconds": schema.Int64Attribute{
				MarkdownDescription: "Time (in seconds) that auto-stop is delayed",
				Required:            true,
			},
			"queue_time_limit_seconds": schema.Int64Attribute{
				MarkdownDescription: "Maximum time (in seconds) a query will wait in the engine's queue before being canceled. Should be >= 120 seconds",
				Required:            true,
			},
			"runtime_limit_seconds": schema.Int64Attribute{
				MarkdownDescription: "Maximum time (in seconds) a query can run before being terminated. Set to 0 for no limit",
				Required:            true,
			},
			"drain_time_limit_seconds": schema.Int64Attribute{
				MarkdownDescription: "Maximum time (in seconds) an engine replica will continue to run after resize/disable/delete before termination",
				Required:            true,
			},
			"max_concurrency": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of concurrent queries that an engine replica can run",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for the engine",
				Optional:            true,
			},
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Whether the engine is enabled. Defaults to true",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			// Computed fields (read-only)
			"state": schema.StringAttribute{
				MarkdownDescription: "Current state of the engine (DELETING, DISABLED, DISABLING, ENABLED, ENABLING, INVALID)",
				Computed:            true,
			},
			"active_replicas": schema.Int64Attribute{
				MarkdownDescription: "Number of engine replicas currently active",
				Computed:            true,
			},
			"queried_at": schema.StringAttribute{
				MarkdownDescription: "Date and time the engine was last used to execute a query",
				Computed:            true,
			},
			"status_changed_at": schema.StringAttribute{
				MarkdownDescription: "Date and time (UTC) that the state of the engine changed",
				Computed:            true,
			},
			"instance_family": schema.StringAttribute{
				MarkdownDescription: "Instance family (M5D, M6ID, M6GD, DDV4, DDV5)",
				Computed:            true,
			},
			"additional_engine_state_info": schema.StringAttribute{
				MarkdownDescription: "Additional engine state information (typically NONE)",
				Computed:            true,
			},
		},
	}
}

// Create a new engine resource.
func (r *dremioEngine) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioEngineModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := r.parseResourceToRequestBody(ctx, &data)
	// Generate a unique request ID for idempotency
	reqBody.RequestID = uuid.New().String()

	api_resp, err := r.client.RequestToDremio("POST", "/engines", reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to create engine, got error: %s", err),
		)
		return
	}
	defer api_resp.Body.Close()

	body, err := io.ReadAll(api_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}

	// Create response only returns the ID
	var createResp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &createResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	data.ID = types.StringValue(createResp.ID)

	// If enable is false, disable the engine after creation (engines are created enabled by default)
	if !data.Enable.ValueBool() {
		_, err = r.client.RequestToDremio("PUT", fmt.Sprintf("/engines/%s/disable", createResp.ID), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error", fmt.Sprintf("Unable to disable engine after creation, got error: %s", err),
			)
			return
		}
	}

	// Read the full engine state to populate computed fields
	r.readEngineState(ctx, &data, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created an engine resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read resource information.
func (r *dremioEngine) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.DremioEngineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readEngineState(ctx, &state, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dremioEngine) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioEngineModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.DremioEngineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()

	// Update engine configuration
	reqBody := r.parseResourceToRequestBody(ctx, &plan)

	// Name is not allowed to be updated, set to empty string (omitted from json)
	reqBody.Name = ""

	_, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/engines/%s", id), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to update engine, got error: %s", err),
		)
		return
	}

	// Handle enable/disable state changes
	planEnabled := plan.Enable.ValueBool()
	stateEnabled := state.Enable.ValueBool()

	if planEnabled != stateEnabled {
		var enablePath string
		if planEnabled {
			enablePath = fmt.Sprintf("/engines/%s/enable", id)
		} else {
			enablePath = fmt.Sprintf("/engines/%s/disable", id)
		}

		_, err = r.client.RequestToDremio("PUT", enablePath, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error", fmt.Sprintf("Unable to change engine enable state, got error: %s", err),
			)
			return
		}
	}

	// Read the full engine state to populate computed fields
	r.readEngineState(ctx, &plan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updated an engine resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dremioEngine) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioEngineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	api_resp, err := r.client.RequestToDremio("DELETE", fmt.Sprintf("/engines/%s", id), nil)
	if err != nil {
		// If engine doesn't exist (404/400), treat as successful delete
		if api_resp.StatusCode == 404 || api_resp.StatusCode == 400 {
			tflog.Warn(ctx, fmt.Sprintf("Engine %s already deleted or doesn't exist", id))
			return
		}
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to delete engine, got error: %s", err),
		)
		return
	}
}

func (r *dremioEngine) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// parseResourceToRequestBody converts the Terraform model to an API request body
func (r *dremioEngine) parseResourceToRequestBody(_ context.Context, data *models.DremioEngineModel) *models.EngineRequest {
	reqBody := &models.EngineRequest{
		Name:                  data.Name.ValueString(),
		Size:                  data.Size.ValueString(),
		MinReplicas:           int(data.MinReplicas.ValueInt64()),
		MaxReplicas:           int(data.MaxReplicas.ValueInt64()),
		AutoStopDelaySeconds:  int(data.AutoStopDelaySeconds.ValueInt64()),
		QueueTimeLimitSeconds: int(data.QueueTimeLimitSeconds.ValueInt64()),
		RuntimeLimitSeconds:   int(data.RuntimeLimitSeconds.ValueInt64()),
		DrainTimeLimitSeconds: int(data.DrainTimeLimitSeconds.ValueInt64()),
		MaxConcurrency:        int(data.MaxConcurrency.ValueInt64()),
	}

	// API requires description to not be null, so set to empty string if not provided
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		reqBody.Description = data.Description.ValueString()
	} else {
		reqBody.Description = ""
	}

	return reqBody
}

// readEngineState reads the engine from the API and updates the Terraform state
func (r *dremioEngine) readEngineState(ctx context.Context, state *models.DremioEngineModel, resp interface{}) {
	id := state.ID.ValueString()

	api_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/engines/%s", id), nil)
	if err != nil {
		if api_resp.StatusCode == 404 || api_resp.StatusCode == 400 {
			tflog.Warn(ctx, fmt.Sprintf("Engine %s not found, removing from state", id))
			if readResp, ok := resp.(*resource.ReadResponse); ok {
				readResp.State.RemoveResource(ctx)
			}
			return
		}
		// Add error for all response types
		switch v := resp.(type) {
		case *resource.CreateResponse:
			v.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read engine, got error: %s", err))
		case *resource.ReadResponse:
			v.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read engine, got error: %s", err))
		case *resource.UpdateResponse:
			v.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read engine, got error: %s", err))
		}
		return
	}
	defer api_resp.Body.Close()

	body, err := io.ReadAll(api_resp.Body)
	if err != nil {
		switch v := resp.(type) {
		case *resource.CreateResponse:
			v.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body: %s", err))
		case *resource.ReadResponse:
			v.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body: %s", err))
		case *resource.UpdateResponse:
			v.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body: %s", err))
		}
		return
	}

	var engineResp models.EngineResponse
	if err := json.Unmarshal(body, &engineResp); err != nil {
		switch v := resp.(type) {
		case *resource.CreateResponse:
			v.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err))
		case *resource.ReadResponse:
			v.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err))
		case *resource.UpdateResponse:
			v.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err))
		}
		return
	}

	r.fromResponseToState(&engineResp, state)
}

// fromResponseToState maps the API response to the Terraform state
func (r *dremioEngine) fromResponseToState(engineResp *models.EngineResponse, state *models.DremioEngineModel) {
	state.ID = types.StringValue(engineResp.ID)
	state.Name = types.StringValue(engineResp.Name)
	state.Size = types.StringValue(engineResp.Size)
	state.MinReplicas = types.Int64Value(int64(engineResp.MinReplicas))
	state.MaxReplicas = types.Int64Value(int64(engineResp.MaxReplicas))
	state.AutoStopDelaySeconds = types.Int64Value(int64(engineResp.AutoStopDelaySeconds))
	state.QueueTimeLimitSeconds = types.Int64Value(int64(engineResp.QueueTimeLimitSeconds))
	state.RuntimeLimitSeconds = types.Int64Value(int64(engineResp.RuntimeLimitSeconds))
	state.DrainTimeLimitSeconds = types.Int64Value(int64(engineResp.DrainTimeLimitSeconds))
	state.MaxConcurrency = types.Int64Value(int64(engineResp.MaxConcurrency))

	// Handle optional description
	if engineResp.Description != "" {
		state.Description = types.StringValue(engineResp.Description)
	} else {
		state.Description = types.StringNull()
	}

	// Set enable based on state
	isEnabled := engineResp.State != "DISABLED"
	state.Enable = types.BoolValue(isEnabled)

	// Computed fields
	state.State = types.StringValue(engineResp.State)
	state.ActiveReplicas = types.Int64Value(int64(engineResp.ActiveReplicas))

	if engineResp.QueriedAt != "" {
		state.QueriedAt = types.StringValue(engineResp.QueriedAt)
	} else {
		state.QueriedAt = types.StringNull()
	}

	if engineResp.StatusChangedAt != "" {
		state.StatusChangedAt = types.StringValue(engineResp.StatusChangedAt)
	} else {
		state.StatusChangedAt = types.StringNull()
	}

	if engineResp.InstanceFamily != "" {
		state.InstanceFamily = types.StringValue(engineResp.InstanceFamily)
	} else {
		state.InstanceFamily = types.StringNull()
	}

	if engineResp.AdditionalEngineStateInfo != "" {
		state.AdditionalEngineStateInfo = types.StringValue(engineResp.AdditionalEngineStateInfo)
	} else {
		state.AdditionalEngineStateInfo = types.StringNull()
	}
}
