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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dremioEngineRuleSetDataSource{}
	_ datasource.DataSourceWithConfigure = &dremioEngineRuleSetDataSource{}
)

type dremioEngineRuleSetDataSource struct {
	client *dremioClient.Client
}

func NewDremioEngineRuleSetDataSource() datasource.DataSource {
	return &dremioEngineRuleSetDataSource{}
}

// Metadata returns the data source type name.
func (d *dremioEngineRuleSetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engine_rule_set"
}

func (d *dremioEngineRuleSetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Schema defines the schema for the data source.
func (d *dremioEngineRuleSetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ruleInfoSchema := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "User-defined name for the rule",
				Computed:            true,
			},
			"condition": schema.StringAttribute{
				MarkdownDescription: "Routing condition using SQL syntax",
				Computed:            true,
			},
			"engine_name": schema.StringAttribute{
				MarkdownDescription: "Name of the engine to route jobs to",
				Computed:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Rule type: ROUTE or REJECT",
				Computed:            true,
			},
			"reject_message": schema.StringAttribute{
				MarkdownDescription: "Message displayed to the user if the rule rejects jobs",
				Computed:            true,
			},
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Dremio Engine Rule Set data source - retrieves engine routing rules for a Dremio project",

		Attributes: map[string]schema.Attribute{
			"rule_infos": schema.ListNestedAttribute{
				MarkdownDescription: "List of routing rules. Rules are evaluated in order.",
				Computed:            true,
				NestedObject:        ruleInfoSchema,
			},
			"rule_info_default": schema.SingleNestedAttribute{
				MarkdownDescription: "The default rule that applies to jobs without a matching rule",
				Computed:            true,
				Attributes:          ruleInfoSchema.Attributes,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "UUID of a tag that routes JDBC queries to a particular session",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dremioEngineRuleSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DremioEngineRuleSetModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Make API request
	api_resp, err := d.client.RequestToDremio("GET", "/rules", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read engine rules: %s", err),
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

	var rulesResp models.EngineRulesResponse
	if err := json.Unmarshal(api_resp_body, &rulesResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, &rulesResp, &data, &resp.Diagnostics)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state model
func (d *dremioEngineRuleSetDataSource) mapResponseToState(ctx context.Context, rulesResp *models.EngineRulesResponse, data *models.DremioEngineRuleSetModel, diags *diag.Diagnostics) {
	if rulesResp.RuleSet == nil {
		return
	}

	// Convert RuleInfos
	ruleInfosList, diagsConv := helpers.ConvertRuleInfosToTerraform(ctx, rulesResp.RuleSet.RuleInfos)
	diags.Append(diagsConv...)
	data.RuleInfos = ruleInfosList

	// Convert RuleInfoDefault
	ruleInfoDefault, diagsConv := helpers.ConvertRuleInfoToTerraform(ctx, rulesResp.RuleSet.RuleInfoDefault)
	diags.Append(diagsConv...)
	data.RuleInfoDefault = ruleInfoDefault

	// Set Tag
	data.Tag = types.StringValue(rulesResp.RuleSet.Tag)
}

