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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dremioDataMaintenance{}
	_ resource.ResourceWithConfigure   = &dremioDataMaintenance{}
	_ resource.ResourceWithImportState = &dremioDataMaintenance{}
)

type dremioDataMaintenance struct {
	client *dremioClient.Client
}

func NewDremioDataMaintenanceResource() resource.Resource {
	return &dremioDataMaintenance{}
}

// Metadata returns the resource type name.
func (r *dremioDataMaintenance) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_maintenance"
}

// Configure adds the provider configured client to the resource.
func (r *dremioDataMaintenance) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dremioDataMaintenance) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Dremio Cloud data maintenance task. Data maintenance tasks automate OPTIMIZE and EXPIRE_SNAPSHOTS operations on tables in Open Catalog.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier (UUID) of the maintenance task",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of maintenance task. Valid values are `OPTIMIZE` (run OPTIMIZE on the table) or `EXPIRE_SNAPSHOTS` (run VACUUM on the table).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("OPTIMIZE", "EXPIRE_SNAPSHOTS"),
				},
			},
			"level": schema.StringAttribute{
				MarkdownDescription: "The scope of the maintenance task. Currently only `TABLE` is supported.",
				Computed:            true,
			},
			"source_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Open Catalog source where the table resides.",
				Computed:            true,
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the maintenance task is enabled. When enabled, the task runs automatically based on Dremio logic.",
				Required:            true,
			},
			"table_id": schema.StringAttribute{
				MarkdownDescription: "Fully qualified table name in the format `folder1.folder2.table_name` (without source name).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *dremioDataMaintenance) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource.
func (r *dremioDataMaintenance) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioDataMaintenanceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, diags := helpers.ConvertMaintenanceTaskFromTerraform(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Make API request
	api_resp, err := r.client.RequestToDremio("POST", "/maintenance/tasks", reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to create data maintenance task, got error: %s", err),
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

	var taskResp models.MaintenanceTaskResponse
	if err := json.Unmarshal(body, &taskResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	diags = helpers.ConvertMaintenanceTaskToTerraform(ctx, &taskResp, &data)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "created data maintenance task resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read resource information.
func (r *dremioDataMaintenance) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.DremioDataMaintenanceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	var taskResp models.MaintenanceTaskResponse
	api_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/maintenance/tasks/%s", id), nil)
	if err != nil {
		// If resource is not found (404), remove it from state so Terraform will recreate it
		if strings.Contains(err.Error(), "status 404") {
			tflog.Warn(ctx, fmt.Sprintf("Data maintenance task %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read data maintenance task, got error: %s", err),
		)
		return
	}
	defer api_resp.Body.Close()

	resp_body, err := io.ReadAll(api_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}
	if err := json.Unmarshal(resp_body, &taskResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	diags = helpers.ConvertMaintenanceTaskToTerraform(ctx, &taskResp, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *dremioDataMaintenance) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioDataMaintenanceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state models.DremioDataMaintenanceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, diags := helpers.ConvertMaintenanceTaskFromTerraform(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Data maintenance task update request with ID: %s", id))

	api_resp, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/maintenance/tasks/%s", id), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to update data maintenance task, got error: %s", err),
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

	var taskResp models.MaintenanceTaskResponse
	if err := json.Unmarshal(body, &taskResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	diags = helpers.ConvertMaintenanceTaskToTerraform(ctx, &taskResp, &plan)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "updated data maintenance task resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dremioDataMaintenance) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioDataMaintenanceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, err := r.client.RequestToDremio("DELETE", fmt.Sprintf("/maintenance/tasks/%s", id), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to delete data maintenance task, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "deleted data maintenance task resource")
}
