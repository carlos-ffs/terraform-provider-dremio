package datasources

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
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                     = &dremioTableDataSource{}
	_ datasource.DataSourceWithConfigure        = &dremioTableDataSource{}
	_ datasource.DataSourceWithConfigValidators = &dremioTableDataSource{}
)

func NewDremioTableDataSource() datasource.DataSource {
	return &dremioTableDataSource{}
}

type dremioTableDataSource struct {
	client *dremioClient.Client
}

// ConfigValidators returns a list of functions which will all be performed during validation.
func (d *dremioTableDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("path"),
		),
	}
}

// Metadata returns the data source type name.
func (d *dremioTableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_table"
}

func (d *dremioTableDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dremioTableDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Table data source - retrieves information about an existing table/physical dataset",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the table. Exactly one of `id` or `path` must be specified.",
				Computed:            true,
				Optional:            true,
			},
			"path": schema.ListAttribute{
				MarkdownDescription: "Full path to the table",
				Optional:            true,
				Computed:            true,
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
			"type": schema.StringAttribute{
				MarkdownDescription: "Dataset type (PHYSICAL_DATASET)",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Date and time the table was created (UTC)",
				Computed:            true,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Version tag for optimistic concurrency control",
				Computed:            true,
			},
			"acceleration_refresh_policy": schema.SingleNestedAttribute{
				MarkdownDescription: "Acceleration refresh policy for the table",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"active_policy_type": schema.StringAttribute{
						MarkdownDescription: "Policy for refreshing Reflections (NEVER, PERIOD, SCHEDULE, REFRESH_ON_DATA_CHANGES)",
						Computed:            true,
					},
					"refresh_period_ms": schema.Int64Attribute{
						MarkdownDescription: "Refresh period in milliseconds (minimum 3600000, default 3600000)",
						Computed:            true,
					},
					"refresh_schedule": schema.StringAttribute{
						MarkdownDescription: "Cron expression for refresh schedule (UTC), e.g., '0 0 8 * * ?'",
						Computed:            true,
					},
					"grace_period_ms": schema.Int64Attribute{
						MarkdownDescription: "Maximum age for Reflection data in milliseconds",
						Computed:            true,
					},
					"method": schema.StringAttribute{
						MarkdownDescription: "Method for refreshing Reflections (AUTO, FULL, INCREMENTAL)",
						Computed:            true,
					},
					"refresh_field": schema.StringAttribute{
						MarkdownDescription: "Field to use for incremental refresh",
						Computed:            true,
					},
					"never_expire": schema.BoolAttribute{
						MarkdownDescription: "Whether Reflections never expire",
						Computed:            true,
					},
				},
			},
			"format": schema.SingleNestedAttribute{
				MarkdownDescription: "Table format information",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Type of data in the table (Delta, Excel, Iceberg, JSON, Parquet, Text, Unknown, XLS)",
						Computed:            true,
					},
					"ignore_other_file_formats": schema.BoolAttribute{
						MarkdownDescription: "For Parquet folders, ignore non-Parquet files",
						Computed:            true,
					},
					"skip_first_line": schema.BoolAttribute{
						MarkdownDescription: "Skip first line when creating table (Excel/Text)",
						Computed:            true,
					},
					"extract_header": schema.BoolAttribute{
						MarkdownDescription: "Extract column names from first line (Excel/Text)",
						Computed:            true,
					},
					"has_merged_cells": schema.BoolAttribute{
						MarkdownDescription: "Expand merged cells (Excel)",
						Computed:            true,
					},
					"sheet_name": schema.StringAttribute{
						MarkdownDescription: "Sheet name for Excel files with multiple sheets",
						Computed:            true,
					},
					"field_delimiter": schema.StringAttribute{
						MarkdownDescription: "Field delimiter character (Text), default: ','",
						Computed:            true,
					},
					"quote": schema.StringAttribute{
						MarkdownDescription: "Quote character (Text), default: '\"'",
						Computed:            true,
					},
					"comment": schema.StringAttribute{
						MarkdownDescription: "Comment character (Text), default: '#'",
						Computed:            true,
					},
					"escape": schema.StringAttribute{
						MarkdownDescription: "Escape character (Text), default: '\"'",
						Computed:            true,
					},
					"line_delimiter": schema.StringAttribute{
						MarkdownDescription: "Line delimiter character (Text), default: '\\n'",
						Computed:            true,
					},
					"auto_generate_column_names": schema.BoolAttribute{
						MarkdownDescription: "Auto-generate column names (Text)",
						Computed:            true,
					},
					"trim_header": schema.BoolAttribute{
						MarkdownDescription: "Trim header whitespace (Text)",
						Computed:            true,
					},
					"auto_correct_corrupt_dates": schema.BoolAttribute{
						MarkdownDescription: "Auto-correct corrupted date fields (read-only)",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Table name (read-only)",
						Computed:            true,
					},
					"full_path": schema.ListAttribute{
						MarkdownDescription: "Full path to the table (read-only)",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"ctime": schema.Int64Attribute{
						MarkdownDescription: "Creation time (read-only)",
						Computed:            true,
					},
					"is_folder": schema.BoolAttribute{
						MarkdownDescription: "Whether the table was created from a folder (read-only)",
						Computed:            true,
					},
					"location": schema.StringAttribute{
						MarkdownDescription: "Location where table metadata is stored (read-only)",
						Computed:            true,
					},
				},
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
			"fields": schema.StringAttribute{
				MarkdownDescription: "Table fields/columns as JSON string. Due to the recursive nature of table schemas (STRUCT and LIST types can be arbitrarily nested), fields are represented as a JSON string. Use jsondecode() to parse this value in Terraform configurations.",
				Computed:            true,
			},
			"approximate_statistics_allowed": schema.BoolAttribute{
				MarkdownDescription: "Whether approximate statistics are allowed",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioTableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioTableDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tableID string
	if !data.ID.IsNull() {
		tableID = data.ID.ValueString()
	}

	var table_path []string
	if !data.Path.IsNull() {
		diags := data.Path.ElementsAs(ctx, &table_path, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if tableID == "" && len(table_path) == 0 {
		resp.Diagnostics.AddError(
			"Missing Table ID or Path",
			"Either `id` or `path` must be specified for Dremio Table data source.",
		)
		return
	}
	if tableID != "" && len(table_path) > 0 {
		resp.Diagnostics.AddError(
			"Both Table ID and Path specified",
			"Only one of `id` or `path` must be specified for Dremio Table data source.",
		)
		return
	}

	var path string
	if tableID != "" { // Read table by ID
		path = fmt.Sprintf("/catalog/%s", tableID)
	} else { // Lookup table ID by path
		table_path_str := "/" + strings.Join(table_path, "/")
		path = fmt.Sprintf("/catalog/by-path/%s", table_path_str)
	}

	api_resp, err := d.client.RequestToDremio("GET", path, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to request table: %s", err),
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

	var tableResp models.TableResponse
	if err := json.Unmarshal(api_resp_body, &tableResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, &tableResp, &data, &resp.Diagnostics)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioTableDataSource) mapResponseToState(ctx context.Context, tableResp *models.TableResponse, data *models.DremioTableDataSourceModel, diags *diag.Diagnostics) {
	// Map basic fields
	data.ID = types.StringValue(tableResp.ID)
	data.Type = types.StringValue(tableResp.Type)
	data.Tag = types.StringValue(tableResp.Tag)

	if tableResp.CreatedAt != "" {
		data.CreatedAt = types.StringValue(tableResp.CreatedAt)
	} else {
		data.CreatedAt = types.StringNull()
	}

	// Map path
	if len(tableResp.Path) == 0 {
		data.Path = types.ListNull(types.StringType)
	} else {
		pathFromAPI, diagsTemp := types.ListValueFrom(ctx, types.StringType, tableResp.Path)
		diags.Append(diagsTemp...)
		data.Path = pathFromAPI
	}

	// Map acceleration refresh policy - use helper function
	// For datasources, we always populate from API (no plan to compare against)
	// So we pass an unknown object as the plan parameter to force conversion
	accelerationRefreshPolicyAttrTypes := helpers.GetAccelerationRefreshPolicyAttrTypes()
	data.AccelerationRefreshPolicy, *diags = helpers.ConvertAccelerationRefreshPolicyToTerraform(
		ctx,
		tableResp.AccelerationRefreshPolicy,
		types.ObjectUnknown(accelerationRefreshPolicyAttrTypes),
	)

	// Map format - use datasource-specific helper function that includes read-only fields
	data.Format, *diags = helpers.ConvertTableFormatToTerraformDatasource(
		ctx,
		tableResp.Format,
	)

	// Map access control list - use helper function
	_, _, accessControlAttrTypes := helpers.GetACLAttrTypes()
	if tableResp.AccessControlList == nil {
		data.AccessControlList = types.ObjectNull(accessControlAttrTypes)
	} else {
		// Pass a non-null plan to force conversion
		data.AccessControlList, *diags = helpers.ConvertACLToTerraform(ctx, tableResp.AccessControlList, types.ObjectUnknown(accessControlAttrTypes))
	}

	// Map owner - use helper function
	var ownerDiags diag.Diagnostics
	data.Owner, ownerDiags = helpers.ConvertOwnerToTerraform(ctx, tableResp.Owner)
	diags.Append(ownerDiags...)

	// Map fields - convert to JSON string
	// Table schemas can be arbitrarily deep with nested STRUCT and LIST types,
	// so we use JSON representation instead of trying to model the recursive structure
	fieldsJSON, fieldsDiags := helpers.ConvertTableFieldsToJSON(ctx, tableResp.Fields)
	diags.Append(fieldsDiags...)
	data.Fields = fieldsJSON

	// Map approximate statistics allowed
	data.ApproximateStatisticsAllowed = types.BoolValue(tableResp.ApproximateStatisticsAllowed)
}
