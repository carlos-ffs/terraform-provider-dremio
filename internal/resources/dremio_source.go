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
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &dremioSource{}
	_ resource.ResourceWithConfigure   = &dremioSource{}
	_ resource.ResourceWithImportState = &dremioSource{}
)

type dremioSource struct {
	client *dremioClient.Client
}

func NewDremioSourceResource() resource.Resource {
	return &dremioSource{}
}

// Metadata returns the resource type name.
func (r *dremioSource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *dremioSource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Schema defines the schema for the resource.
func (r *dremioSource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Source resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the source",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"entity_type": schema.StringAttribute{
				MarkdownDescription: "Entity type (always 'source')",
				Computed:            true,
				Default:             stringdefault.StaticString("source"),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Source type (e.g., ARCTIC, S3, SNOWFLAKE, MYSQL, POSTGRES, etc.)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ARCTIC", "S3", "SNOWFLAKE", "MYSQL", "POSTGRES", "BIGQUERY", "REDSHIFT", "ORACLE", "MSSQL", "AZURE_STORAGE", "AWS_GLUE", "DB2", "ICEBERG_REST_CATALOG", "AZURE_SYNAPSE", "SAPHANA", "SNOWFLAKE_OPEN_CATALOG", "UNITY_CATALOG", "VERTICA"}...),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "User-defined name of the source",
				Required:            true,
			},
			"config": schema.StringAttribute{
				MarkdownDescription: "Configuration options specific to the source type as a JSON string. The schema varies based on the source type. See https://docs.dremio.com/cloud/reference/api/catalog/source/source-config for available configuration options for each source type.",
				Required:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"metadata_policy": schema.SingleNestedAttribute{
				MarkdownDescription: "Metadata refresh policy",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"auth_ttl_ms": schema.Int64Attribute{
						MarkdownDescription: "Sets the length of time, in milliseconds, that source permissions are cached. Default is 24 hours (86400000 milliseconds). Minimum is one minute (60000 milliseconds).",
						Optional:            true,
					},
					"names_refresh_ms": schema.Int64Attribute{
						MarkdownDescription: "Sets when to run a refresh of a source, in milliseconds. Default is one hour (3600000 milliseconds). Minimum is one minute (60000 milliseconds).",
						Optional:            true,
					},
					"dataset_refresh_after_ms": schema.Int64Attribute{
						MarkdownDescription: "Determines how often the metadata in the dataset is refreshed, in milliseconds. Default is one hour (3600000 milliseconds). Minimum is one minute (60000 milliseconds).",
						Optional:            true,
					},
					"dataset_expire_after_ms": schema.Int64Attribute{
						MarkdownDescription: "Sets the amount of time, in milliseconds, to keep the metadata before it expires. Default is one hour (3600000 milliseconds). Minimum is one minute (60000 milliseconds).",
						Optional:            true,
					},
					"dataset_update_mode": schema.StringAttribute{
						MarkdownDescription: "Sets the metadata policy for when a dataset is updated. Use PREFETCH_QUERIED to update the details for previously queried objects in a source.",
						Optional:            true,
					},
					"delete_unavailable_datasets": schema.BoolAttribute{
						MarkdownDescription: "Option to remove dataset definitions if the underlying data is unavailable to Dremio. Default is true. Set to false to keep the dataset definitions.",
						Optional:            true,
					},
					"auto_promote_datasets": schema.BoolAttribute{
						MarkdownDescription: "Option to automatically format files into tables when a query is issued. Default is false. Set to true to enable this capability. This attribute applies only to metastore and object storage sources.",
						Optional:            true,
					},
				},
			},
			"acceleration_grace_period_ms": schema.Int64Attribute{
				MarkdownDescription: "Time to keep Reflections before expiration (milliseconds)",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"acceleration_refresh_period_ms": schema.Int64Attribute{
				MarkdownDescription: "Refresh frequency for Reflections (milliseconds)",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"acceleration_active_policy_type": schema.StringAttribute{
				MarkdownDescription: "Policy for refreshing Reflections (NEVER, PERIOD, SCHEDULE)",
				Optional:            true,
			},
			"acceleration_refresh_schedule": schema.StringAttribute{
				MarkdownDescription: "Cron expression for Reflection refresh schedule (UTC)",
				Optional:            true,
			},
			"acceleration_refresh_on_data_changes": schema.BoolAttribute{
				MarkdownDescription: "Refresh Reflections when Iceberg table snapshots change",
				Optional:            true,
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
func (r *dremioSource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioSourceModel

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
			"Client Error", fmt.Sprintf("Unable to create source, got error: %s", err),
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

	var sourceResp models.SourceResponse
	if err := json.Unmarshal(body, &sourceResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Validate response data
	if sourceResp.ID == nil || sourceResp.Tag == nil {
		resp.Diagnostics.AddError(
			"Invalid Response",
			fmt.Sprintf("API response missing required fields. Response body: %s", string(body)),
		)
		return
	}

	// Update state with response data
	// Preserve the plan's config to avoid drift from API-added defaults
	r.fromResponseToState(ctx, &sourceResp, &data, &resp.Diagnostics)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read resource information.
func (r *dremioSource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.DremioSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	var sourceResp models.SourceResponse
	source_resp, err := r.client.RequestToDremio("GET", fmt.Sprintf("/catalog/%s", id), nil)
	if err != nil {
		// If resource is not found (404), remove it from state so Terraform will recreate it
		if strings.Contains(err.Error(), "status 404") {
			tflog.Warn(ctx, fmt.Sprintf("Source %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read source, got error: %s", err),
		)
		return
	}
	defer source_resp.Body.Close()

	resp_body, err := io.ReadAll(source_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}
	if err := json.Unmarshal(resp_body, &sourceResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	// Use API response to detect actual changes during refresh
	r.fromResponseToState(ctx, &sourceResp, &state, &resp.Diagnostics)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dremioSource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the tag (computed field)
	var state models.DremioSourceModel
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

	tflog.Debug(ctx, fmt.Sprintf("Source update request with ID: %s, and Tag: %s", id, reqBody.Tag))

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

	var sourceResp models.SourceResponse
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

func (r *dremioSource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioSourceModel
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

func (r *dremioSource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// fromResponseToState updates the state with values from the API response.
// The state parameter contains the current state, which is used to preserve null values
// for optional nested objects (metadata_policy, access_control_list) to avoid drift.
// The preserveConfig parameter controls whether to keep the plan's config (true for Create/Update)
// or use the API response config (false for Read). This prevents drift from API-added default values.
func (r *dremioSource) fromResponseToState(ctx context.Context, sourceResp *models.SourceResponse, state *models.DremioSourceModel, diags *diag.Diagnostics) {

	// Basic scalar fields
	if sourceResp.ID == nil {
		tflog.Error(ctx, "Source response ID is nil")
		diags.AddError(
			"Invalid API Response",
			"API response is missing the ID field",
		)
		return
	}
	state.ID = types.StringValue(*sourceResp.ID)
	state.EntityType = types.StringValue("source")
	state.Type = types.StringValue(*sourceResp.Type)
	state.Name = types.StringValue(*sourceResp.Name)
	state.Tag = types.StringValue(*sourceResp.Tag)

	tflog.Debug(ctx, fmt.Sprintf("fromResponseToState: ID=%s, Name=%s", *sourceResp.ID, *sourceResp.Name))

	// We dont change the config since we dont want to drift from the user's config
	// from API-added defaults.

	// Acceleration settings
	if sourceResp.AccelerationGracePeriodMs != nil {
		state.AccelerationGracePeriodMs = types.Int64Value(*sourceResp.AccelerationGracePeriodMs)
	} else {
		state.AccelerationGracePeriodMs = types.Int64Null()
	}

	if sourceResp.AccelerationRefreshPeriodMs != nil {
		state.AccelerationRefreshPeriodMs = types.Int64Value(*sourceResp.AccelerationRefreshPeriodMs)
	} else {
		state.AccelerationRefreshPeriodMs = types.Int64Null()
	}

	if sourceResp.AccelerationActivePolicyType != nil {
		state.AccelerationActivePolicyType = types.StringValue(*sourceResp.AccelerationActivePolicyType)
	} else {
		state.AccelerationActivePolicyType = types.StringNull()
	}

	if sourceResp.AccelerationRefreshSchedule != nil {
		state.AccelerationRefreshSchedule = types.StringValue(*sourceResp.AccelerationRefreshSchedule)
	} else {
		state.AccelerationRefreshSchedule = types.StringNull()
	}

	metadataPolicyAttrTypes := map[string]attr.Type{
		"auth_ttl_ms":                 types.Int64Type,
		"names_refresh_ms":            types.Int64Type,
		"dataset_refresh_after_ms":    types.Int64Type,
		"dataset_expire_after_ms":     types.Int64Type,
		"dataset_update_mode":         types.StringType,
		"delete_unavailable_datasets": types.BoolType,
		"auto_promote_datasets":       types.BoolType,
	}

	// If the user didn't specify metadata_policy in the plan (it's null in state),
	// keep it null instead of populating with API defaults to avoid drift.
	if state.MetadataPolicy.IsNull() && sourceResp.MetadataPolicy != nil {
		// Keep the null value - don't populate from API response
	} else if sourceResp.MetadataPolicy == nil {
		state.MetadataPolicy = types.ObjectNull(metadataPolicyAttrTypes)
	} else if !state.MetadataPolicy.IsNull() {
		mp := sourceResp.MetadataPolicy

		mpModel := models.MetadataPolicyModel{
			AuthTTLMs:                 types.Int64Null(),
			NamesRefreshMs:            types.Int64Null(),
			DatasetRefreshAfterMs:     types.Int64Null(),
			DatasetExpireAfterMs:      types.Int64Null(),
			DatasetUpdateMode:         types.StringNull(),
			DeleteUnavailableDatasets: types.BoolNull(),
			AutoPromoteDatasets:       types.BoolNull(),
		}

		if mp.AuthTTLMs != nil {
			mpModel.AuthTTLMs = types.Int64Value(*mp.AuthTTLMs)
		}
		if mp.NamesRefreshMs != nil {
			mpModel.NamesRefreshMs = types.Int64Value(*mp.NamesRefreshMs)
		}
		if mp.DatasetRefreshAfterMs != nil {
			mpModel.DatasetRefreshAfterMs = types.Int64Value(*mp.DatasetRefreshAfterMs)
		}
		if mp.DatasetExpireAfterMs != nil {
			mpModel.DatasetExpireAfterMs = types.Int64Value(*mp.DatasetExpireAfterMs)
		}
		if mp.DatasetUpdateMode != nil {
			mpModel.DatasetUpdateMode = types.StringValue(*mp.DatasetUpdateMode)
		}
		if mp.DeleteUnavailableDatasets != nil {
			mpModel.DeleteUnavailableDatasets = types.BoolValue(*mp.DeleteUnavailableDatasets)
		}
		if mp.AutoPromoteDatasets != nil {
			mpModel.AutoPromoteDatasets = types.BoolValue(*mp.AutoPromoteDatasets)
		}

		metadataPolicyObj, d := types.ObjectValueFrom(ctx, metadataPolicyAttrTypes, mpModel)
		diags.Append(d...)
		state.MetadataPolicy = metadataPolicyObj
	}

	// Access control list block - use helper function
	var aclDiags diag.Diagnostics
	state.AccessControlList, aclDiags = helpers.ConvertACLToTerraform(ctx, sourceResp.AccessControlList, state.AccessControlList)
	diags.Append(aclDiags...)
}

func (r *dremioSource) parseResourceToRequestBody(ctx context.Context, data *models.DremioSourceModel, diags *diag.Diagnostics) *models.SourceRequest {
	// Build the request body
	reqBody := &models.SourceRequest{
		EntityType: "source",
		Type:       data.Type.ValueString(),
		Name:       data.Name.ValueString(),
	}

	// Parse config JSON string based on source type
	if !data.Config.IsNull() && !data.Config.IsUnknown() {
		configInterface, err := parseConfigByType(data.Type.ValueString(), data.Config.ValueString())
		if err != nil {
			diags.AddError(
				"Invalid Config",
				fmt.Sprintf("Unable to parse config JSON for source type %s: %s", data.Type.ValueString(), err),
			)
			return nil
		}

		// Convert the typed config back to map[string]interface{} for the API
		configBytes, err := json.Marshal(configInterface)
		if err != nil {
			diags.AddError(
				"Config Marshal Error",
				fmt.Sprintf("Unable to marshal config: %s", err),
			)
			return nil
		}

		var configMap map[string]interface{}
		if err := json.Unmarshal(configBytes, &configMap); err != nil {
			diags.AddError(
				"Config Unmarshal Error",
				fmt.Sprintf("Unable to unmarshal config to map: %s", err),
			)
			return nil
		}
		reqBody.Config = configMap
	}

	// Handle optional fields
	if !data.AccelerationGracePeriodMs.IsNull() && !data.AccelerationGracePeriodMs.IsUnknown() {
		reqBody.AccelerationGracePeriodMs = data.AccelerationGracePeriodMs.ValueInt64()
	}

	if !data.AccelerationRefreshPeriodMs.IsNull() && !data.AccelerationRefreshPeriodMs.IsUnknown() {
		reqBody.AccelerationRefreshPeriodMs = data.AccelerationRefreshPeriodMs.ValueInt64()
	}

	if !data.AccelerationActivePolicyType.IsNull() && !data.AccelerationActivePolicyType.IsUnknown() {
		reqBody.AccelerationActivePolicyType = data.AccelerationActivePolicyType.ValueString()
	}

	if !data.AccelerationRefreshSchedule.IsNull() && !data.AccelerationRefreshSchedule.IsUnknown() {
		reqBody.AccelerationRefreshSchedule = data.AccelerationRefreshSchedule.ValueString()
	}

	if !data.AccelerationRefreshOnDataChanges.IsNull() && !data.AccelerationRefreshOnDataChanges.IsUnknown() {
		reqBody.AccelerationRefreshOnDataChanges = data.AccelerationRefreshOnDataChanges.ValueBool()
	}

	// Handle MetadataPolicy
	if !data.MetadataPolicy.IsNull() && !data.MetadataPolicy.IsUnknown() {
		var metadataPolicy models.MetadataPolicyModel
		diags := data.MetadataPolicy.As(ctx, &metadataPolicy, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			diags.Append(diags...)
			return nil
		}

		reqBody.MetadataPolicy = &models.MetadataPolicy{}
		if !metadataPolicy.AuthTTLMs.IsNull() && !metadataPolicy.AuthTTLMs.IsUnknown() {
			v := metadataPolicy.AuthTTLMs.ValueInt64()
			reqBody.MetadataPolicy.AuthTTLMs = &v
		}
		if !metadataPolicy.NamesRefreshMs.IsNull() && !metadataPolicy.NamesRefreshMs.IsUnknown() {
			v := metadataPolicy.NamesRefreshMs.ValueInt64()
			reqBody.MetadataPolicy.NamesRefreshMs = &v
		}
		if !metadataPolicy.DatasetRefreshAfterMs.IsNull() && !metadataPolicy.DatasetRefreshAfterMs.IsUnknown() {
			v := metadataPolicy.DatasetRefreshAfterMs.ValueInt64()
			reqBody.MetadataPolicy.DatasetRefreshAfterMs = &v
		}
		if !metadataPolicy.DatasetExpireAfterMs.IsNull() && !metadataPolicy.DatasetExpireAfterMs.IsUnknown() {
			v := metadataPolicy.DatasetExpireAfterMs.ValueInt64()
			reqBody.MetadataPolicy.DatasetExpireAfterMs = &v
		}
		if !metadataPolicy.DatasetUpdateMode.IsNull() && !metadataPolicy.DatasetUpdateMode.IsUnknown() {
			v := metadataPolicy.DatasetUpdateMode.ValueString()
			reqBody.MetadataPolicy.DatasetUpdateMode = &v
		}
		if !metadataPolicy.DeleteUnavailableDatasets.IsNull() && !metadataPolicy.DeleteUnavailableDatasets.IsUnknown() {
			v := metadataPolicy.DeleteUnavailableDatasets.ValueBool()
			reqBody.MetadataPolicy.DeleteUnavailableDatasets = &v
		}
		if !metadataPolicy.AutoPromoteDatasets.IsNull() && !metadataPolicy.AutoPromoteDatasets.IsUnknown() {
			v := metadataPolicy.AutoPromoteDatasets.ValueBool()
			reqBody.MetadataPolicy.AutoPromoteDatasets = &v
		}
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

// parseConfigByType parses the config JSON string based on the source type
// and returns the appropriate typed struct
func parseConfigByType(sourceType, configJSON string) (interface{}, error) {
	var config interface{}

	switch sourceType {
	case "ARCTIC":
		config = &models.ArcticConfig{}
	case "S3":
		config = &models.S3Config{}
	case "SNOWFLAKE":
		config = &models.SnowflakeConfig{}
	case "MYSQL":
		config = &models.MySQLConfig{}
	case "POSTGRES":
		config = &models.PostgreSQLConfig{}
	case "BIGQUERY":
		config = &models.BigQueryConfig{}
	case "REDSHIFT":
		config = &models.RedshiftConfig{}
	case "ORACLE":
		config = &models.OracleConfig{}
	case "MSSQL":
		config = &models.MSSQLConfig{}
	case "AZURE_STORAGE":
		config = &models.AzureStorageConfig{}
	case "AWS_GLUE":
		config = &models.AWSGlueConfig{}
	case "DB2":
		config = &models.Db2Config{}
	case "ICEBERG_REST_CATALOG":
		config = &models.IcebergRESTCatalogConfig{}
	case "AZURE_SYNAPSE":
		config = &models.AzureSynapseConfig{}
	case "SAPHANA":
		config = &models.SAPHANAConfig{}
	case "SNOWFLAKE_OPEN_CATALOG":
		config = &models.SnowflakeOpenCatalogConfig{}
	case "UNITY_CATALOG":
		config = &models.UnityCatalogConfig{}
	case "VERTICA":
		config = &models.VerticaConfig{}
	default:
		return nil, fmt.Errorf("unsupported source type: %s", sourceType)
	}

	// Use a decoder with DisallowUnknownFields to strictly validate the config
	decoder := json.NewDecoder(strings.NewReader(configJSON))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return config, nil
}
