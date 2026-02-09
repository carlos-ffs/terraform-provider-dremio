package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
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
	_ resource.Resource                = &dremioTable{}
	_ resource.ResourceWithConfigure   = &dremioTable{}
	_ resource.ResourceWithImportState = &dremioTable{}
)

type dremioTable struct {
	client *dremioClient.Client
}

func NewDremioTableResource() resource.Resource {
	return &dremioTable{}
}

// Metadata returns the resource type name.
func (r *dremioTable) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_table"
}

func (r *dremioTable) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dremioTable) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *dremioTable) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioTableModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	_, err := r.client.RequestToDremio("DELETE", fmt.Sprintf("/catalog/%s", id), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to delete table, got error: %s", err),
		)
		return
	}
}

func (r *dremioTable) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Table resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the table",
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
				MarkdownDescription: "Dataset type (always 'PHYSICAL_DATASET')",
				Computed:            true,
				Default:             stringdefault.StaticString("PHYSICAL_DATASET"),
			},
			"path": schema.ListAttribute{
				MarkdownDescription: "Full path to the table",
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
			"file_or_folder_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the source file or folder to format as a table",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"acceleration_refresh_policy": schema.SingleNestedAttribute{
				MarkdownDescription: "Acceleration refresh policy for the table",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"active_policy_type": schema.StringAttribute{
						MarkdownDescription: "Policy for refreshing Reflections (NEVER, PERIOD, SCHEDULE, REFRESH_ON_DATA_CHANGES)",
						Optional:            true,
					},
					"refresh_period_ms": schema.Int64Attribute{
						MarkdownDescription: "Refresh period in milliseconds (minimum 3600000, default 3600000)",
						Optional:            true,
					},
					"refresh_schedule": schema.StringAttribute{
						MarkdownDescription: "Cron expression for refresh schedule (UTC), e.g., '0 0 8 * * ?'",
						Optional:            true,
					},
					"grace_period_ms": schema.Int64Attribute{
						MarkdownDescription: "Maximum age for Reflection data in milliseconds",
						Optional:            true,
					},
					"method": schema.StringAttribute{
						MarkdownDescription: "Method for refreshing Reflections (AUTO, FULL, INCREMENTAL)",
						Optional:            true,
					},
					"refresh_field": schema.StringAttribute{
						MarkdownDescription: "Field to use for incremental refresh",
						Optional:            true,
					},
					"never_expire": schema.BoolAttribute{
						MarkdownDescription: "Whether Reflections never expire",
						Optional:            true,
					},
				},
			},
			"format": schema.SingleNestedAttribute{
				MarkdownDescription: "Format parameters for the table",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Type of data in the table (Delta, Excel, Iceberg, JSON, Parquet, Text, Unknown, XLS)",
						Required:            true,
					},
					"ignore_other_file_formats": schema.BoolAttribute{
						MarkdownDescription: "For Parquet folders, ignore non-Parquet files",
						Optional:            true,
					},
					"skip_first_line": schema.BoolAttribute{
						MarkdownDescription: "Skip first line when creating table (Excel/Text)",
						Optional:            true,
					},
					"extract_header": schema.BoolAttribute{
						MarkdownDescription: "Extract column names from first line (Excel/Text)",
						Optional:            true,
					},
					"has_merged_cells": schema.BoolAttribute{
						MarkdownDescription: "Expand merged cells (Excel)",
						Optional:            true,
					},
					"sheet_name": schema.StringAttribute{
						MarkdownDescription: "Sheet name for Excel files with multiple sheets",
						Optional:            true,
					},
					"field_delimiter": schema.StringAttribute{
						MarkdownDescription: "Field delimiter character (Text), default: ','",
						Optional:            true,
					},
					"quote": schema.StringAttribute{
						MarkdownDescription: "Quote character (Text), default: '\"'",
						Optional:            true,
					},
					"comment": schema.StringAttribute{
						MarkdownDescription: "Comment character (Text), default: '#'",
						Optional:            true,
					},
					"escape": schema.StringAttribute{
						MarkdownDescription: "Escape character (Text), default: '\"'",
						Optional:            true,
					},
					"line_delimiter": schema.StringAttribute{
						MarkdownDescription: "Line delimiter (Text), default: '\\r\\n'",
						Optional:            true,
					},
					"auto_generate_column_names": schema.BoolAttribute{
						MarkdownDescription: "Use existing column names (Text)",
						Optional:            true,
					},
					"trim_header": schema.BoolAttribute{
						MarkdownDescription: "Trim column names (Text)",
						Optional:            true,
					},
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
func (r *dremioTable) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioTableModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := r.parseResourceToRequestBody(ctx, &data, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	fileOrFolderID := url.QueryEscape(data.FileOrFolderID.ValueString())

	// Make API request
	api_resp, err := r.client.RequestToDremio("POST", fmt.Sprintf("/catalog/%s", fileOrFolderID), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to create table, got error: %s", err),
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

	var tableResp models.TableResponse
	if err := json.Unmarshal(body, &tableResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &tableResp, &data, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a table resource")
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read resource information.
func (r *dremioTable) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.DremioTableModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	var tableResp models.TableResponse
	table_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s", id), nil)
	if err != nil {
		// If resource is not found (404), remove it from state so Terraform will recreate it
		if strings.Contains(err.Error(), "status 404") {
			tflog.Warn(ctx, fmt.Sprintf("Table %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read table, got error: %s", err),
		)
		return
	}
	defer table_resp.Body.Close()

	resp_body, err := io.ReadAll(table_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}
	if err := json.Unmarshal(resp_body, &tableResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &tableResp, &state, &resp.Diagnostics)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dremioTable) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioTableModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the tag (computed field)
	var state models.DremioTableModel
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
	reqBody.Tag = state.Tag.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Table update request with ID: %s, and Tag: %s", id, reqBody.Tag))

	api_resp, err := r.client.RequestToDremio("PUT", fmt.Sprintf("/catalog/%s", id), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to update table, got error: %s", err),
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

	var tableResp models.TableResponse
	if err := json.Unmarshal(body, &tableResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &tableResp, &plan, &resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated a table resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dremioTable) parseResourceToRequestBody(ctx context.Context, data *models.DremioTableModel, diags *diag.Diagnostics) *models.TableRequest {
	// Build the request body
	reqBody := &models.TableRequest{
		EntityType:       "dataset",
		Type:             "PHYSICAL_DATASET",
		SourceOrFolderID: data.FileOrFolderID.ValueString(),
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

	// Handle AccelerationRefreshPolicy - use helper function
	var arpDiags diag.Diagnostics
	reqBody.AccelerationRefreshPolicy, arpDiags = helpers.ConvertAccelerationRefreshPolicyFromTerraform(ctx, data.AccelerationRefreshPolicy)
	if arpDiags.HasError() {
		diags.Append(arpDiags...)
		return nil
	}

	// Handle Format - use helper function
	var formatDiags diag.Diagnostics
	reqBody.Format, formatDiags = helpers.ConvertTableFormatFromTerraform(ctx, data.Format)
	if formatDiags.HasError() {
		diags.Append(formatDiags...)
		return nil
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

func (r *dremioTable) fromResponseToState(ctx context.Context, tableResp *models.TableResponse, state *models.DremioTableModel, diags *diag.Diagnostics) {
	state.ID = types.StringValue(tableResp.ID)
	state.EntityType = types.StringValue("dataset")
	state.Type = types.StringValue(tableResp.Type)
	state.Tag = types.StringValue(tableResp.Tag)

	// Handle Path
	if len(tableResp.Path) == 0 {
		state.Path = types.ListNull(types.StringType)
	} else {
		pathFromAPI, diagsTemp := types.ListValueFrom(ctx, types.StringType, tableResp.Path)
		diags.Append(diagsTemp...)
		state.Path = pathFromAPI
	}

	// Handle AccelerationRefreshPolicy - use helper function
	// Important: Only update if the field was null in the state to avoid drift
	var arpDiags diag.Diagnostics
	state.AccelerationRefreshPolicy, arpDiags = helpers.ConvertAccelerationRefreshPolicyToTerraform(ctx, tableResp.AccelerationRefreshPolicy, state.AccelerationRefreshPolicy)
	diags.Append(arpDiags...)

	// Handle Format - use helper function
	// Important: Only update if the field was null in the state to avoid drift
	var formatDiags diag.Diagnostics
	state.Format, formatDiags = helpers.ConvertTableFormatToTerraform(ctx, tableResp.Format, state.Format)
	diags.Append(formatDiags...)

	// Handle AccessControlList - use helper function
	var aclDiags diag.Diagnostics
	state.AccessControlList, aclDiags = helpers.ConvertACLToTerraform(ctx, tableResp.AccessControlList, state.AccessControlList)
	diags.Append(aclDiags...)
}
