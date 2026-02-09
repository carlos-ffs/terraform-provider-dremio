package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/helpers"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dremioDataMaintenanceTaskDataSource{}
	_ datasource.DataSourceWithConfigure = &dremioDataMaintenanceTaskDataSource{}
)

type dremioDataMaintenanceTaskDataSource struct {
	client *dremioClient.Client
}

func NewDremioDataMaintenanceTaskDataSource() datasource.DataSource {
	return &dremioDataMaintenanceTaskDataSource{}
}

// Metadata returns the data source type name.
func (d *dremioDataMaintenanceTaskDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_maintenance_task"
}

func (d *dremioDataMaintenanceTaskDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dremioDataMaintenanceTaskDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Data Maintenance Task data source - retrieves information about an existing data maintenance task",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier (UUID) of the maintenance task",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of maintenance task. Valid values are `OPTIMIZE` or `EXPIRE_SNAPSHOTS`.",
				Computed:            true,
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
				MarkdownDescription: "Whether the maintenance task is enabled.",
				Computed:            true,
			},
			"table_id": schema.StringAttribute{
				MarkdownDescription: "Fully qualified table name in the format `folder1.folder2.table_name` (without source name).",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioDataMaintenanceTaskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioDataMaintenanceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	taskID := data.ID.ValueString()

	// Make API request
	api_resp, err := d.client.RequestToDremio("GET", fmt.Sprintf("/maintenance/tasks/%s", taskID), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read data maintenance task: %s", err),
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

	var taskResp models.MaintenanceTaskResponse
	if err := json.Unmarshal(api_resp_body, &taskResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state using helper
	diags := helpers.ConvertMaintenanceTaskToDataSource(ctx, &taskResp, &data)
	resp.Diagnostics.Append(diags...)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

