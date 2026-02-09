package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                     = &dremioEngineDataSource{}
	_ datasource.DataSourceWithConfigure        = &dremioEngineDataSource{}
	_ datasource.DataSourceWithConfigValidators = &dremioEngineDataSource{}
)

func NewDremioEngineDataSource() datasource.DataSource {
	return &dremioEngineDataSource{}
}

type dremioEngineDataSource struct {
	client *dremioClient.Client
}

// Metadata returns the data source type name.
func (d *dremioEngineDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engine"
}

func (d *dremioEngineDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *dremioEngineDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

// Schema defines the schema for the data source.
func (d *dremioEngineDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Engine data source - retrieves information about an existing engine in Dremio Cloud",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the engine (UUID). Exactly one of `id` or `name` must be specified.",
				Computed:            true,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the engine. Exactly one of `id` or `name` must be specified.",
				Computed:            true,
				Optional:            true,
			},
			"size": schema.StringAttribute{
				MarkdownDescription: "Size of the engine (XX_SMALL_V1, X_SMALL_V1, SMALL_V1, MEDIUM_V1, LARGE_V1, X_LARGE_V1, XX_LARGE_V1, XXX_LARGE_V1)",
				Computed:            true,
			},
			"min_replicas": schema.Int64Attribute{
				MarkdownDescription: "Minimum number of engine replicas",
				Computed:            true,
			},
			"max_replicas": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of engine replicas",
				Computed:            true,
			},
			"auto_stop_delay_seconds": schema.Int64Attribute{
				MarkdownDescription: "Time (in seconds) that auto-stop is delayed",
				Computed:            true,
			},
			"queue_time_limit_seconds": schema.Int64Attribute{
				MarkdownDescription: "Maximum time (in seconds) a query will wait in the engine's queue",
				Computed:            true,
			},
			"runtime_limit_seconds": schema.Int64Attribute{
				MarkdownDescription: "Maximum time (in seconds) a query can run",
				Computed:            true,
			},
			"drain_time_limit_seconds": schema.Int64Attribute{
				MarkdownDescription: "Maximum time (in seconds) an engine replica will continue to run after resize/disable/delete",
				Computed:            true,
			},
			"max_concurrency": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of concurrent queries per replica",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the engine",
				Computed:            true,
			},
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
			"additional_engine_state_info": schema.StringAttribute{
				MarkdownDescription: "Additional engine state information",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioEngineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioEngineDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	engineID := data.ID.ValueString()
	engineName := data.Name.ValueString()

	var apiPath string
	if engineName != "" {
		// Need to list engines and find by name
		apiPath = "/engines"
	} else {
		apiPath = fmt.Sprintf("/engines/%s", engineID)
	}

	api_resp, err := d.client.RequestToDremio("GET", apiPath, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read engine: %s", err),
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

	if engineName != "" {
		// Parse list response and find engine by name
		var engines []models.EngineResponse
		if err := json.Unmarshal(body, &engines); err != nil {
			resp.Diagnostics.AddError(
				"Parse Error",
				fmt.Sprintf("Unable to parse engines list response: %s", err),
			)
			return
		}

		var foundEngine *models.EngineResponse
		for _, engine := range engines {
			if engine.Name == engineName {
				foundEngine = &engine
				break
			}
		}

		if foundEngine == nil {
			resp.Diagnostics.AddError(
				"Not Found",
				fmt.Sprintf("Engine with name '%s' not found", engineName),
			)
			return
		}

		d.mapResponseToState(foundEngine, &data)
	} else {
		// Parse single engine response
		var engineResp models.EngineResponse
		if err := json.Unmarshal(body, &engineResp); err != nil {
			resp.Diagnostics.AddError(
				"Parse Error",
				fmt.Sprintf("Unable to parse engine response: %s", err),
			)
			return
		}

		d.mapResponseToState(&engineResp, &data)
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioEngineDataSource) mapResponseToState(engineResp *models.EngineResponse, data *models.DremioEngineDataSourceModel) {
	data.ID = types.StringValue(engineResp.ID)
	data.Name = types.StringValue(engineResp.Name)
	data.Size = types.StringValue(engineResp.Size)
	data.MinReplicas = types.Int64Value(int64(engineResp.MinReplicas))
	data.MaxReplicas = types.Int64Value(int64(engineResp.MaxReplicas))
	data.AutoStopDelaySeconds = types.Int64Value(int64(engineResp.AutoStopDelaySeconds))
	data.QueueTimeLimitSeconds = types.Int64Value(int64(engineResp.QueueTimeLimitSeconds))
	data.RuntimeLimitSeconds = types.Int64Value(int64(engineResp.RuntimeLimitSeconds))
	data.DrainTimeLimitSeconds = types.Int64Value(int64(engineResp.DrainTimeLimitSeconds))
	data.MaxConcurrency = types.Int64Value(int64(engineResp.MaxConcurrency))
	data.Description = types.StringValue(engineResp.Description)
	data.QueriedAt = types.StringValue(engineResp.QueriedAt)
	data.StatusChangedAt = types.StringValue(engineResp.StatusChangedAt)
	data.State = types.StringValue(engineResp.State)
	data.ActiveReplicas = types.Int64Value(int64(engineResp.ActiveReplicas))
	data.AdditionalEngineStateInfo = types.StringValue(engineResp.AdditionalEngineStateInfo)
}
