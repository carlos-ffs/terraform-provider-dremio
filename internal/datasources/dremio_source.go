package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/helpers"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                     = &dremioSourceDataSource{}
	_ datasource.DataSourceWithConfigure        = &dremioSourceDataSource{}
	_ datasource.DataSourceWithConfigValidators = &dremioSourceDataSource{}
)

func NewDremioSourceDataSource() datasource.DataSource {
	return &dremioSourceDataSource{}
}

type dremioSourceDataSource struct {
	client *dremioClient.Client
}

// Metadata returns the data source type name.
func (d *dremioSourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (d *dremioSourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// ConfigValidators returns a list of functions which will all be performed during validation.
func (d *dremioSourceDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

// Schema defines the schema for the data source.
func (d *dremioSourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Source data source - retrieves information about an existing source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the source. Exactly one of `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "User-defined name of the source. Exactly one of `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Version tag for optimistic concurrency control",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Source type (e.g., ARCTIC, S3, SNOWFLAKE, MYSQL, POSTGRES, etc.)",
				Computed:            true,
			},
			"config": schema.StringAttribute{
				MarkdownDescription: "Configuration options specific to the source type as a JSON string",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"metadata_policy": schema.SingleNestedAttribute{
				MarkdownDescription: "Metadata refresh policy",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"auth_ttl_ms": schema.Int64Attribute{
						MarkdownDescription: "Length of time, in milliseconds, that source permissions are cached",
						Computed:            true,
					},
					"names_refresh_ms": schema.Int64Attribute{
						MarkdownDescription: "When to run a refresh of a source, in milliseconds",
						Computed:            true,
					},
					"dataset_refresh_after_ms": schema.Int64Attribute{
						MarkdownDescription: "How often the metadata in the dataset is refreshed, in milliseconds",
						Computed:            true,
					},
					"dataset_expire_after_ms": schema.Int64Attribute{
						MarkdownDescription: "Amount of time, in milliseconds, to keep the metadata before it expires",
						Computed:            true,
					},
					"dataset_update_mode": schema.StringAttribute{
						MarkdownDescription: "Metadata policy for when a dataset is updated",
						Computed:            true,
					},
					"delete_unavailable_datasets": schema.BoolAttribute{
						MarkdownDescription: "Option to remove dataset definitions if the underlying data is unavailable",
						Computed:            true,
					},
					"auto_promote_datasets": schema.BoolAttribute{
						MarkdownDescription: "Option to automatically format files into tables when a query is issued",
						Computed:            true,
					},
				},
			},
			"acceleration_grace_period_ms": schema.Int64Attribute{
				MarkdownDescription: "Grace period before using Reflections (milliseconds)",
				Computed:            true,
			},
			"acceleration_refresh_period_ms": schema.Int64Attribute{
				MarkdownDescription: "Refresh period for Reflections (milliseconds)",
				Computed:            true,
			},
			"acceleration_never_expire": schema.BoolAttribute{
				MarkdownDescription: "Whether Reflections never expire",
				Computed:            true,
			},
			"acceleration_never_refresh": schema.BoolAttribute{
				MarkdownDescription: "Whether Reflections never refresh",
				Computed:            true,
			},
			"acceleration_active_policy_type": schema.StringAttribute{
				MarkdownDescription: "Active policy type (PERIOD or NEVER)",
				Computed:            true,
			},
			"acceleration_refresh_schedule": schema.StringAttribute{
				MarkdownDescription: "Cron expression for refresh schedule",
				Computed:            true,
			},
			"children": schema.ListNestedAttribute{
				MarkdownDescription: "Child entities in the source",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier of the entity",
							Computed:            true,
						},
						"path": schema.ListAttribute{
							MarkdownDescription: "Full path to the entity",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"tag": schema.StringAttribute{
							MarkdownDescription: "Version tag",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Entity type (CONTAINER or DATASET)",
							Computed:            true,
						},
						"container_type": schema.StringAttribute{
							MarkdownDescription: "Container type (SPACE, SOURCE, FOLDER, HOME)",
							Computed:            true,
						},
						"dataset_type": schema.StringAttribute{
							MarkdownDescription: "Dataset type (VIRTUAL_DATASET or PHYSICAL_DATASET)",
							Computed:            true,
						},
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
			"permissions": schema.ListAttribute{
				MarkdownDescription: "User's permissions on the source",
				Computed:            true,
				ElementType:         types.StringType,
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
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioSourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioSourceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sourceId string
	if !data.ID.IsNull() {
		sourceId = data.ID.ValueString()
	}
	var sourceName string
	if !data.Name.IsNull() {
		sourceName = data.Name.ValueString()
	}

	if sourceId == "" && sourceName == "" {
		resp.Diagnostics.AddError(
			"Missing Source ID or Name",
			"Either `id` or `name` must be specified for Dremio source data source.",
		)
		return
	}
	if sourceId != "" && sourceName != "" {
		resp.Diagnostics.AddError(
			"Both Source ID and Name specified",
			"Only one of `id` or `name` must be specified for Dremio source data source.",
		)
		return
	}

	var path string
	if sourceName != "" { // Lookup source ID by name
		path = fmt.Sprintf("/catalog/by-path/%s", sourceName)
	} else { // Read source by ID
		path = fmt.Sprintf("/catalog/%s", sourceId)
	}

	api_resp, err := d.client.RequestToDremio("GET", path, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to request source: %s", err),
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

	var sourceResp models.SourceResponse
	if err := json.Unmarshal(api_resp_body, &sourceResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, &sourceResp, &data, &resp.Diagnostics)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioSourceDataSource) mapResponseToState(ctx context.Context, sourceResp *models.SourceResponse, data *models.DremioSourceDataSourceModel, diags *diag.Diagnostics) {
	// Map basic fields
	if sourceResp.ID != nil {
		data.ID = types.StringValue(*sourceResp.ID)
	} else {
		data.ID = types.StringNull()
	}

	if sourceResp.Name != nil {
		data.Name = types.StringValue(*sourceResp.Name)
	} else {
		data.Name = types.StringNull()
	}

	if sourceResp.Tag != nil {
		data.Tag = types.StringValue(*sourceResp.Tag)
	} else {
		data.Tag = types.StringNull()
	}

	if sourceResp.Type != nil {
		data.Type = types.StringValue(*sourceResp.Type)
	} else {
		data.Type = types.StringNull()
	}

	// Map config
	if sourceResp.Config == nil {
		data.Config = jsontypes.NewNormalizedNull()
	} else {
		configBytes, err := json.Marshal(sourceResp.Config)
		if err != nil {
			diags.AddError(
				"Config Marshal Error",
				fmt.Sprintf("Unable to marshal source config from API response: %s", err),
			)
			return
		}
		data.Config = jsontypes.NewNormalizedValue(string(configBytes))
	}

	// Map acceleration fields
	if sourceResp.AccelerationGracePeriodMs != nil {
		data.AccelerationGracePeriodMs = types.Int64Value(*sourceResp.AccelerationGracePeriodMs)
	} else {
		data.AccelerationGracePeriodMs = types.Int64Null()
	}

	if sourceResp.AccelerationRefreshPeriodMs != nil {
		data.AccelerationRefreshPeriodMs = types.Int64Value(*sourceResp.AccelerationRefreshPeriodMs)
	} else {
		data.AccelerationRefreshPeriodMs = types.Int64Null()
	}

	if sourceResp.AccelerationNeverExpire != nil {
		data.AccelerationNeverExpire = types.BoolValue(*sourceResp.AccelerationNeverExpire)
	} else {
		data.AccelerationNeverExpire = types.BoolNull()
	}

	if sourceResp.AccelerationNeverRefresh != nil {
		data.AccelerationNeverRefresh = types.BoolValue(*sourceResp.AccelerationNeverRefresh)
	} else {
		data.AccelerationNeverRefresh = types.BoolNull()
	}

	if sourceResp.AccelerationActivePolicyType != nil {
		data.AccelerationActivePolicyType = types.StringValue(*sourceResp.AccelerationActivePolicyType)
	} else {
		data.AccelerationActivePolicyType = types.StringNull()
	}

	if sourceResp.AccelerationRefreshSchedule != nil {
		data.AccelerationRefreshSchedule = types.StringValue(*sourceResp.AccelerationRefreshSchedule)
	} else {
		data.AccelerationRefreshSchedule = types.StringNull()
	}

	// Map metadata policy - use helper function
	// For datasources, we always populate from API (no plan to compare against)
	// So we pass an unknown object as the plan parameter to force conversion
	metadataPolicyAttrTypes := helpers.GetMetadataPolicyAttrTypes()
	data.MetadataPolicy, *diags = helpers.ConvertMetadataPolicyToTerraform(
		ctx,
		sourceResp.MetadataPolicy,
		types.ObjectUnknown(metadataPolicyAttrTypes),
	)

	// Map children - use helper function
	catalogEntityAttrTypes := helpers.GetCatalogEntityAttrTypes()
	if sourceResp.Children == nil || len(*sourceResp.Children) == 0 {
		data.Children = types.ListNull(types.ObjectType{AttrTypes: catalogEntityAttrTypes})
	} else {
		childObjects := make([]types.Object, 0, len(*sourceResp.Children))
		for _, child := range *sourceResp.Children {
			childObj, d := helpers.ConvertCatalogEntityToTerraform(ctx, child)
			diags.Append(d...)
			childObjects = append(childObjects, childObj)
		}

		childrenList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: catalogEntityAttrTypes}, childObjects)
		diags.Append(d...)
		data.Children = childrenList
	}

	// Map access control list - use helper function
	// For datasources, we always populate from API (no plan to compare against)
	// So we pass an unknown object as the plan parameter to force conversion
	_, _, accessControlAttrTypes := helpers.GetACLAttrTypes()

	var aclDiags diag.Diagnostics
	if sourceResp.AccessControlList == nil {
		data.AccessControlList = types.ObjectNull(accessControlAttrTypes)
	} else {
		// Pass a non-null plan to force conversion
		data.AccessControlList, aclDiags = helpers.ConvertACLToTerraform(ctx, sourceResp.AccessControlList, types.ObjectUnknown(accessControlAttrTypes))
		diags.Append(aclDiags...)
	}

	// Map permissions
	if len(sourceResp.Permissions) == 0 {
		data.Permissions = types.ListNull(types.StringType)
	} else {
		permsList, d := types.ListValueFrom(ctx, types.StringType, sourceResp.Permissions)
		diags.Append(d...)
		data.Permissions = permsList
	}

	// Map owner - use helper function
	var ownerDiags diag.Diagnostics
	data.Owner, ownerDiags = helpers.ConvertOwnerToTerraform(ctx, sourceResp.Owner)
	diags.Append(ownerDiags...)
}
