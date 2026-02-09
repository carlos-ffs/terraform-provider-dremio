package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetRuleInfoAttrTypes returns the attribute type definitions for RuleInfo structures.
func GetRuleInfoAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":           types.StringType,
		"condition":      types.StringType,
		"engine_name":    types.StringType,
		"action":         types.StringType,
		"reject_message": types.StringType,
	}
}

// GetRuleSetAttrTypes returns the attribute type definitions for RuleSet structures.
func GetRuleSetAttrTypes() map[string]attr.Type {
	ruleInfoAttrTypes := GetRuleInfoAttrTypes()
	return map[string]attr.Type{
		"rule_infos":        types.ListType{ElemType: types.ObjectType{AttrTypes: ruleInfoAttrTypes}},
		"rule_info_default": types.ObjectType{AttrTypes: ruleInfoAttrTypes},
		"tag":               types.StringType,
	}
}

// ConvertRuleInfoFromTerraform converts Terraform RuleInfo state to API request format.
// This is a request-only structure, so there is no ToTerraform conversion.
//
// Parameters:
//   - ctx: Context for the operation
//   - ruleObj: The rule info from Terraform state/plan
//
// Returns:
//   - *models.RuleInfo: The converted rule info for API requests
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertRuleInfoFromTerraform(
	ctx context.Context,
	ruleObj types.Object,
) (*models.RuleInfo, diag.Diagnostics) {
	var diags diag.Diagnostics

	if ruleObj.IsNull() || ruleObj.IsUnknown() {
		return nil, diags
	}

	var ruleModel models.RuleInfoModel
	diagsL := ruleObj.As(ctx, &ruleModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return nil, diags
	}

	result := &models.RuleInfo{
		Name:          ruleModel.Name.ValueString(),
		Condition:     ruleModel.Condition.ValueString(),
		EngineName:    ruleModel.EngineName.ValueString(),
		Action:        ruleModel.Action.ValueString(),
		RejectMessage: ruleModel.RejectMessage.ValueString(),
	}

	return result, diags
}

// ConvertRuleSetFromTerraform converts Terraform RuleSet state to API request format.
// This is a request-only structure, so there is no ToTerraform conversion.
//
// Parameters:
//   - ctx: Context for the operation
//   - ruleSetObj: The rule set from Terraform state/plan
//
// Returns:
//   - *models.RuleSet: The converted rule set for API requests (nil if input is null/unknown)
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertRuleSetFromTerraform(
	ctx context.Context,
	ruleSetObj types.Object,
) (*models.RuleSet, diag.Diagnostics) {
	var diags diag.Diagnostics

	if ruleSetObj.IsNull() || ruleSetObj.IsUnknown() {
		return nil, diags
	}

	var ruleSetModel models.RuleSetModel
	diagsL := ruleSetObj.As(ctx, &ruleSetModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return nil, diags
	}

	result := &models.RuleSet{}

	// Convert RuleInfos list
	if !ruleSetModel.RuleInfos.IsNull() && !ruleSetModel.RuleInfos.IsUnknown() {
		var ruleInfoObjs []types.Object
		diagsL := ruleSetModel.RuleInfos.ElementsAs(ctx, &ruleInfoObjs, false)
		if diagsL.HasError() {
			diags.Append(diagsL...)
			return nil, diags
		}

		for _, ruleInfoObj := range ruleInfoObjs {
			ruleInfo, d := ConvertRuleInfoFromTerraform(ctx, ruleInfoObj)
			diags.Append(d...)
			if ruleInfo != nil {
				result.RuleInfos = append(result.RuleInfos, ruleInfo)
			}
		}
	}

	// Convert RuleInfoDefault
	if !ruleSetModel.RuleInfoDefault.IsNull() && !ruleSetModel.RuleInfoDefault.IsUnknown() {
		ruleInfoDefault, d := ConvertRuleInfoFromTerraform(ctx, ruleSetModel.RuleInfoDefault)
		diags.Append(d...)
		result.RuleInfoDefault = ruleInfoDefault
	}

	if !ruleSetModel.Tag.IsNull() && !ruleSetModel.Tag.IsUnknown() {
		result.Tag = ruleSetModel.Tag.ValueString()
	}

	return result, diags
}

// ConvertRuleInfoToTerraform converts API RuleInfo response to Terraform state format.
//
// Parameters:
//   - ctx: Context for the operation
//   - ruleInfo: The rule info from API response
//
// Returns:
//   - types.Object: The converted rule info for Terraform state
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertRuleInfoToTerraform(
	ctx context.Context,
	ruleInfo *models.RuleInfo,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if ruleInfo == nil {
		return types.ObjectNull(GetRuleInfoAttrTypes()), diags
	}

	ruleInfoAttrTypes := GetRuleInfoAttrTypes()

	ruleObj, d := types.ObjectValue(ruleInfoAttrTypes, map[string]attr.Value{
		"name":           types.StringValue(ruleInfo.Name),
		"condition":      types.StringValue(ruleInfo.Condition),
		"engine_name":    types.StringValue(ruleInfo.EngineName),
		"action":         types.StringValue(ruleInfo.Action),
		"reject_message": types.StringValue(ruleInfo.RejectMessage),
	})
	diags.Append(d...)

	return ruleObj, diags
}

// ConvertRuleInfosToTerraform converts a list of API RuleInfo responses to Terraform list format.
//
// Parameters:
//   - ctx: Context for the operation
//   - ruleInfos: The list of rule infos from API response
//
// Returns:
//   - types.List: The converted rule infos for Terraform state
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertRuleInfosToTerraform(
	ctx context.Context,
	ruleInfos []*models.RuleInfo,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	ruleInfoAttrTypes := GetRuleInfoAttrTypes()

	if ruleInfos == nil {
		return types.ListNull(types.ObjectType{AttrTypes: ruleInfoAttrTypes}), diags
	}

	var ruleObjects []attr.Value
	for _, ruleInfo := range ruleInfos {
		ruleObj, d := ConvertRuleInfoToTerraform(ctx, ruleInfo)
		diags.Append(d...)
		ruleObjects = append(ruleObjects, ruleObj)
	}

	rulesList, d := types.ListValue(types.ObjectType{AttrTypes: ruleInfoAttrTypes}, ruleObjects)
	diags.Append(d...)

	return rulesList, diags
}
