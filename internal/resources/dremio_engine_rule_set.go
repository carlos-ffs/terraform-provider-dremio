package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	dremioClient "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/helpers"
	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &dremioEngineRuleSet{}
	_ resource.ResourceWithConfigure = &dremioEngineRuleSet{}
)

type dremioEngineRuleSet struct {
	client *dremioClient.Client
}

func NewDremioEngineRuleSetResource() resource.Resource {
	return &dremioEngineRuleSet{}
}

// Metadata returns the resource type name.
func (r *dremioEngineRuleSet) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engine_rule_set"
}

func (r *dremioEngineRuleSet) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dremioEngineRuleSet) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	ruleInfoSchema := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "User-defined name for the rule",
				Required:            true,
			},
			"condition": schema.StringAttribute{
				MarkdownDescription: "Routing condition using SQL syntax. See Dremio Workload Management documentation for more information.",
				Optional:            true,
				Computed:            true,
			},
			"engine_name": schema.StringAttribute{
				MarkdownDescription: "Name of the engine to route jobs to. Must be empty when action is REJECT.",
				Optional:            true,
				Computed:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Rule type: ROUTE (route to engine) or REJECT (reject the query)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ROUTE", "REJECT"),
				},
			},
			"reject_message": schema.StringAttribute{
				MarkdownDescription: "Message displayed to the user if the rule rejects jobs. Only applicable when action is REJECT.",
				Optional:            true,
				Computed:            true,
			},
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages engine routing rules for a Dremio project. Engine rules are used to route jobs to specific engines based on conditions.
**Important Notes:**
- Only one engine rule set resource should be defined per Terraform configuration. Multiple resources will override each other since the API replaces all rules on each update.
- When this resource is applied, any existing rules not defined in the resource will be deleted.`,

		Attributes: map[string]schema.Attribute{
			"rule_infos": schema.ListNestedAttribute{
				MarkdownDescription: "List of routing rules. Rules are evaluated in order. When adding rules, include all existing rules you want to retain; otherwise, they will be deleted.",
				Optional:            true,
				Computed:            true,
				NestedObject:        ruleInfoSchema,
			},
			"rule_info_default": schema.SingleNestedAttribute{
				MarkdownDescription: "The default rule that applies to jobs without a matching rule. This rule cannot be deleted and is computed from the API.",
				Computed:            true,
				Attributes:          ruleInfoSchema.Attributes,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "UUID of a tag that routes JDBC queries to a particular session. When the JDBC connection property ROUTING_TAG is set, the specified tag value is associated with all queries executed within that connection's session.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *dremioEngineRuleSet) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.DremioEngineRuleSetModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current default rule from state to preserve it
	defaultRule, d := helpers.ConvertRuleInfoFromTerraform(ctx, state.RuleInfoDefault)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete all rules except the default by sending an empty ruleInfos array
	reqBody := models.EngineRulesRequest{
		RuleSet: &models.RuleSet{
			RuleInfos:       []*models.RuleInfo{}, // Empty array to delete all non-default rules
			RuleInfoDefault: defaultRule,
			Tag:             state.Tag.ValueString(),
		},
	}

	_, err := r.client.RequestToDremio("PUT", "/rules", reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to delete engine rules, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "deleted engine rule set resource (reset to default rule only)")
}

// Read resource information.
func (r *dremioEngineRuleSet) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.DremioEngineRuleSetModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var rulesResp models.EngineRulesResponse
	rules_resp, err := r.client.RequestToDremio("GET", "/rules", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to read engine rules, got error: %s", err),
		)
		return
	}
	defer rules_resp.Body.Close()

	resp_body, err := io.ReadAll(rules_resp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read response body: %s", err),
		)
		return
	}
	if err := json.Unmarshal(resp_body, &rulesResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	// Update state with response data
	r.fromResponseToState(ctx, &rulesResp, &state, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Create a new resource.
func (r *dremioEngineRuleSet) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.DremioEngineRuleSetModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for existing rules and warn about overriding
	existingResp, err := r.client.RequestToDremio("GET", "/rules", nil)
	if err == nil {
		defer existingResp.Body.Close()
		existingBody, readErr := io.ReadAll(existingResp.Body)
		if readErr == nil {
			var existingRulesResp models.EngineRulesResponse
			if json.Unmarshal(existingBody, &existingRulesResp) == nil {
				r.warnAboutOverriddenRules(ctx, &data, existingRulesResp.RuleSet, &resp.Diagnostics)
			}
		}
	}

	reqBody := r.parseResourceToRequestBody(ctx, &data, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	api_resp, err := r.client.RequestToDremio("PUT", "/rules", reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to create engine rules, got error: %s", err),
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

	var rulesResp models.EngineRulesResponse
	if err := json.Unmarshal(body, &rulesResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	r.fromResponseToState(ctx, &rulesResp, &data, &resp.Diagnostics)

	tflog.Trace(ctx, "created engine rule set resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dremioEngineRuleSet) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DremioEngineRuleSetModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for existing rules and warn about overriding
	existingResp, err := r.client.RequestToDremio("GET", "/rules", nil)
	if err == nil {
		defer existingResp.Body.Close()
		existingBody, readErr := io.ReadAll(existingResp.Body)
		if readErr == nil {
			var existingRulesResp models.EngineRulesResponse
			if json.Unmarshal(existingBody, &existingRulesResp) == nil {
				r.warnAboutOverriddenRules(ctx, &plan, existingRulesResp.RuleSet, &resp.Diagnostics)
			}
		}
	}

	reqBody := r.parseResourceToRequestBody(ctx, &plan, &resp.Diagnostics)
	if reqBody == nil {
		return
	}

	api_resp, err := r.client.RequestToDremio("PUT", "/rules", reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error", fmt.Sprintf("Unable to update engine rules, got error: %s", err),
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

	var rulesResp models.EngineRulesResponse
	if err := json.Unmarshal(body, &rulesResp); err != nil {
		resp.Diagnostics.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse response: %s", err),
		)
		return
	}

	r.fromResponseToState(ctx, &rulesResp, &plan, &resp.Diagnostics)

	tflog.Trace(ctx, "updated engine rule set resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// fromResponseToState updates the state with values from the API response.
func (r *dremioEngineRuleSet) fromResponseToState(ctx context.Context, rulesResp *models.EngineRulesResponse, state *models.DremioEngineRuleSetModel, diags *diag.Diagnostics) {
	if rulesResp.RuleSet == nil {
		return
	}

	// Convert RuleInfos
	ruleInfosList, d := helpers.ConvertRuleInfosToTerraform(ctx, rulesResp.RuleSet.RuleInfos)
	diags.Append(d...)
	state.RuleInfos = ruleInfosList

	// Convert RuleInfoDefault
	ruleInfoDefault, d := helpers.ConvertRuleInfoToTerraform(ctx, rulesResp.RuleSet.RuleInfoDefault)
	diags.Append(d...)
	state.RuleInfoDefault = ruleInfoDefault

	// Set Tag
	state.Tag = types.StringValue(rulesResp.RuleSet.Tag)

	tflog.Debug(ctx, fmt.Sprintf("fromResponseToState: RuleInfos count=%d, Tag=%s",
		len(rulesResp.RuleSet.RuleInfos), rulesResp.RuleSet.Tag))
}

// parseResourceToRequestBody converts Terraform state/plan to API request body.
func (r *dremioEngineRuleSet) parseResourceToRequestBody(ctx context.Context, data *models.DremioEngineRuleSetModel, diags *diag.Diagnostics) *models.EngineRulesRequest {
	ruleSet := &models.RuleSet{}

	// Convert RuleInfos from plan
	if !data.RuleInfos.IsNull() && !data.RuleInfos.IsUnknown() {
		var ruleInfoObjs []types.Object
		d := data.RuleInfos.ElementsAs(ctx, &ruleInfoObjs, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil
		}

		for _, ruleInfoObj := range ruleInfoObjs {
			ruleInfo, d := helpers.ConvertRuleInfoFromTerraform(ctx, ruleInfoObj)
			diags.Append(d...)
			if ruleInfo != nil {
				// Apply defaults for optional fields
				r.applyRuleInfoDefaults(ruleInfo)
				ruleSet.RuleInfos = append(ruleSet.RuleInfos, ruleInfo)
			}
		}
	} else {
		// If not specified, use empty array
		ruleSet.RuleInfos = []*models.RuleInfo{}
	}

	// RuleInfoDefault is computed - fetch from API to include in request
	currentRulesResp, err := r.client.RequestToDremio("GET", "/rules", nil)
	if err != nil {
		diags.AddError(
			"Client Error",
			fmt.Sprintf("Unable to fetch current rules to get default rule: %s", err),
		)
		return nil
	}
	defer currentRulesResp.Body.Close()

	currentBody, err := io.ReadAll(currentRulesResp.Body)
	if err != nil {
		diags.AddError(
			"Read Error",
			fmt.Sprintf("Unable to read current rules response: %s", err),
		)
		return nil
	}

	var currentRules models.EngineRulesResponse
	if err := json.Unmarshal(currentBody, &currentRules); err != nil {
		diags.AddError(
			"Parse Error",
			fmt.Sprintf("Unable to parse current rules response: %s", err),
		)
		return nil
	}

	if currentRules.RuleSet != nil && currentRules.RuleSet.RuleInfoDefault != nil {
		ruleSet.RuleInfoDefault = currentRules.RuleSet.RuleInfoDefault
	}

	// Set Tag (optional, defaults to empty string)
	if !data.Tag.IsNull() && !data.Tag.IsUnknown() {
		ruleSet.Tag = data.Tag.ValueString()
	} else {
		ruleSet.Tag = ""
	}

	return &models.EngineRulesRequest{
		RuleSet: ruleSet,
	}
}

// applyRuleInfoDefaults sets default values for optional rule info fields.
// - Tag defaults to empty string if not specified
// - RejectMessage defaults to empty string if not specified
// - EngineName is set to empty string when action is REJECT
func (r *dremioEngineRuleSet) applyRuleInfoDefaults(ruleInfo *models.RuleInfo) {
	// EngineName should be empty when action is REJECT
	if ruleInfo.Action == "REJECT" {
		ruleInfo.EngineName = ""
	}
	// Tag and RejectMessage default to empty string (already handled by zero values)
}

// warnAboutOverriddenRules checks if existing rules will be overridden and generates warnings.
func (r *dremioEngineRuleSet) warnAboutOverriddenRules(ctx context.Context, plan *models.DremioEngineRuleSetModel, existingRuleSet *models.RuleSet, diags *diag.Diagnostics) {
	if existingRuleSet == nil || len(existingRuleSet.RuleInfos) == 0 {
		return
	}

	// Get planned rule names
	plannedRuleNames := make(map[string]bool)
	if !plan.RuleInfos.IsNull() && !plan.RuleInfos.IsUnknown() {
		var ruleInfoObjs []types.Object
		d := plan.RuleInfos.ElementsAs(ctx, &ruleInfoObjs, false)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		for _, ruleInfoObj := range ruleInfoObjs {
			ruleInfo, d := helpers.ConvertRuleInfoFromTerraform(ctx, ruleInfoObj)
			diags.Append(d...)
			if ruleInfo != nil {
				plannedRuleNames[ruleInfo.Name] = true
			}
		}
	}

	// Check for rules that will be deleted (exist in API but not in plan)
	var deletedRules []string
	for _, existingRule := range existingRuleSet.RuleInfos {
		if !plannedRuleNames[existingRule.Name] {
			deletedRules = append(deletedRules, existingRule.Name)
		}
	}

	if len(deletedRules) > 0 {
		diags.AddWarning(
			"Existing rules will be removed",
			fmt.Sprintf("The following rules exist in Dremio but are not defined in this resource and will be deleted: %v. If you want to retain these rules, add them to the rule_infos list.", deletedRules),
		)
	}

	// Check for rules that will be overridden (exist in both)
	var overriddenRules []string
	for _, existingRule := range existingRuleSet.RuleInfos {
		if plannedRuleNames[existingRule.Name] {
			overriddenRules = append(overriddenRules, existingRule.Name)
		}
	}

	if len(overriddenRules) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Rules that will be updated: %v", overriddenRules))
	}
}
