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
	_ datasource.DataSource                     = &dremioFolderDataSource{}
	_ datasource.DataSourceWithConfigure        = &dremioFolderDataSource{}
	_ datasource.DataSourceWithConfigValidators = &dremioFolderDataSource{}
)

func NewDremioFolderDataSource() datasource.DataSource {
	return &dremioFolderDataSource{}
}

type dremioFolderDataSource struct {
	client *dremioClient.Client
}

// ConfigValidators returns a list of functions which will all be performed during validation.
func (d *dremioFolderDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("path"),
		),
	}
}

// Metadata returns the data source type name.
func (d *dremioFolderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (d *dremioFolderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dremioFolderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Folder data source - retrieves information about an existing folder",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the folder. Exactly one of `id` or `path` must be specified.",
				Computed:            true,
				Optional:            true,
			},
			"path": schema.ListAttribute{
				MarkdownDescription: "Full path to the folder",
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
			"entity_type": schema.StringAttribute{
				MarkdownDescription: "Type of catalog object (always 'folder')",
				Computed:            true,
			},
			"max_children": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of child objects to include in results. Default is 25.",
				Optional:            true,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Version tag for optimistic concurrency control",
				Computed:            true,
			},
			"children": schema.ListNestedAttribute{
				MarkdownDescription: "Child entities in the folder",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier of the child object",
							Computed:            true,
						},
						"path": schema.ListAttribute{
							MarkdownDescription: "Full path to the child object",
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
							MarkdownDescription: "Container type (FOLDER if type is CONTAINER)",
							Computed:            true,
						},
						"dataset_type": schema.StringAttribute{
							MarkdownDescription: "Dataset type (VIRTUAL or PROMOTED if type is DATASET)",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "Date and time the child object was created",
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
				MarkdownDescription: "User's permissions on the folder",
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
			"storage_uri": schema.StringAttribute{
				MarkdownDescription: "Indicates the location of the Open Catalog folder in object storage.",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioFolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioFolderDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var folderID string
	if !data.ID.IsNull() {
		folderID = data.ID.ValueString()
	}

	var folder_path []string
	if !data.Path.IsNull() {
		diags := data.Path.ElementsAs(ctx, &folder_path, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if folderID == "" && len(folder_path) == 0 {
		resp.Diagnostics.AddError(
			"Missing Folder ID or Path",
			"Either `id` or `path` must be specified for Dremio Folder data source.",
		)
		return
	}
	if folderID != "" && len(folder_path) > 0 {
		resp.Diagnostics.AddError(
			"Both Folder ID and Path specified",
			"Only one of `id` or `path` must be specified for Dremio Folder data source.",
		)
		return
	}

	var path string
	if folderID != "" { // Read folder by ID
		path = fmt.Sprintf("/catalog/%s", folderID)
	} else { // Lookup folder ID by name
		folder_path_str := "/" + strings.Join(folder_path, "/")
		path = fmt.Sprintf("/catalog/by-path/%s", folder_path_str)
	}

	// We will not support page token.
	if !data.MaxChildren.IsNull() {
		path += fmt.Sprintf("?maxChildren=%d", data.MaxChildren.ValueInt64())
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

	var folderResp models.FolderResponse
	if err := json.Unmarshal(api_resp_body, &folderResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, &folderResp, &data, &resp.Diagnostics)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioFolderDataSource) mapResponseToState(ctx context.Context, folderResp *models.FolderResponse, data *models.DremioFolderDataSourceModel, diags *diag.Diagnostics) {

	// Map basic fields
	data.ID = types.StringValue(folderResp.ID)
	data.Tag = types.StringValue(folderResp.Tag)

	if folderResp.EntityType != "" {
		data.EntityType = types.StringValue(folderResp.EntityType)
	} else {
		data.EntityType = types.StringNull()
	}

	if len(folderResp.Path) == 0 {
		data.Path = types.ListNull(types.StringType)
	} else {
		pathFromAPI, diagsTemp := types.ListValueFrom(ctx, types.StringType, folderResp.Path)
		diags.Append(diagsTemp...)
		data.Path = pathFromAPI
	}

	// Map children - use helper function
	folderChildAttrTypes := helpers.GetFolderChildAttrTypes()
	if len(folderResp.Children) == 0 {
		data.Children = types.ListNull(types.ObjectType{AttrTypes: folderChildAttrTypes})
	} else {
		childObjects := make([]types.Object, 0, len(folderResp.Children))
		for _, child := range folderResp.Children {
			childObj, diagsTemp := helpers.ConvertFolderChildToTerraform(ctx, child)
			diags.Append(diagsTemp...)
			childObjects = append(childObjects, childObj)
		}

		childrenList, diagsTemp := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: folderChildAttrTypes}, childObjects)
		diags.Append(diagsTemp...)
		data.Children = childrenList
	}

	// Map access control list - use helper function
	// For datasources, we always populate from API (no plan to compare against)
	// So we pass an unknown object as the plan parameter to force conversion
	_, _, accessControlAttrTypes := helpers.GetACLAttrTypes()
	if folderResp.AccessControlList == nil {
		data.AccessControlList = types.ObjectNull(accessControlAttrTypes)
	} else {
		// Pass a non-null plan to force conversion
		data.AccessControlList, *diags = helpers.ConvertACLToTerraform(ctx, folderResp.AccessControlList, types.ObjectUnknown(accessControlAttrTypes))
	}
	if diags.HasError() {
		return
	}

	// Map permissions
	if len(folderResp.Permissions) == 0 {
		data.Permissions = types.ListNull(types.StringType)
	} else {
		permsList, d := types.ListValueFrom(ctx, types.StringType, folderResp.Permissions)
		diags.Append(d...)
		data.Permissions = permsList
	}
	if diags.HasError() {
		return
	}

	if folderResp.StorageURI != "" {
		data.StorageURI = types.StringValue(folderResp.StorageURI)
	} else {
		data.StorageURI = types.StringNull()
	}

	// Map owner - use helper function
	var ownerDiags diag.Diagnostics
	data.Owner, ownerDiags = helpers.ConvertOwnerToTerraform(ctx, folderResp.Owner)
	diags.Append(ownerDiags...)
}
