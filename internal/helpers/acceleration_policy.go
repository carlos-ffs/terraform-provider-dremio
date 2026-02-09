package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetAccelerationRefreshPolicyAttrTypes returns the attribute type definitions for AccelerationRefreshPolicy structures.
func GetAccelerationRefreshPolicyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"active_policy_type": types.StringType,
		"refresh_period_ms":  types.Int64Type,
		"refresh_schedule":   types.StringType,
		"grace_period_ms":    types.Int64Type,
		"method":             types.StringType,
		"refresh_field":      types.StringType,
		"never_expire":       types.BoolType,
	}
}

// ConvertAccelerationRefreshPolicyToTerraform converts API AccelerationRefreshPolicy response to Terraform state.
//
// Parameters:
//   - ctx: Context for the operation
//   - apiPolicy: The acceleration refresh policy from the API response (can be nil)
//   - planPolicy: The acceleration refresh policy from the Terraform plan/state
//
// Returns:
//   - types.Object: The converted acceleration refresh policy as a Terraform object
//   - diag.Diagnostics: Any diagnostics encountered during conversion
//
// Behavior:
//   - If planPolicy is null and apiPolicy is not nil, returns null (avoids drift from API defaults)
//   - If apiPolicy is nil, returns a typed null object
//   - Otherwise, converts the API acceleration refresh policy to Terraform format
func ConvertAccelerationRefreshPolicyToTerraform(
	ctx context.Context,
	apiPolicy *models.AccelerationRefreshPolicy,
	planPolicy types.Object,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := GetAccelerationRefreshPolicyAttrTypes()

	// If the user didn't specify acceleration_refresh_policy in the plan (it's null in state),
	// keep it null instead of populating with API defaults to avoid drift.
	if planPolicy.IsNull() && apiPolicy != nil {
		return types.ObjectNull(attrTypes), diags
	}

	if apiPolicy == nil {
		return types.ObjectNull(attrTypes), diags
	}

	policyModel := models.AccelerationRefreshPolicyModel{}

	// Handle ActivePolicyType
	if apiPolicy.ActivePolicyType != nil {
		policyModel.ActivePolicyType = types.StringValue(*apiPolicy.ActivePolicyType)
	} else {
		policyModel.ActivePolicyType = types.StringNull()
	}

	// Handle RefreshPeriodMs
	if apiPolicy.RefreshPeriodMs != nil {
		policyModel.RefreshPeriodMs = types.Int64Value(*apiPolicy.RefreshPeriodMs)
	} else {
		policyModel.RefreshPeriodMs = types.Int64Null()
	}

	// Handle RefreshSchedule
	if apiPolicy.RefreshSchedule != nil {
		policyModel.RefreshSchedule = types.StringValue(*apiPolicy.RefreshSchedule)
	} else {
		policyModel.RefreshSchedule = types.StringNull()
	}

	// Handle GracePeriodMs
	if apiPolicy.GracePeriodMs != nil {
		policyModel.GracePeriodMs = types.Int64Value(*apiPolicy.GracePeriodMs)
	} else {
		policyModel.GracePeriodMs = types.Int64Null()
	}

	// Handle Method
	if apiPolicy.Method != nil {
		policyModel.Method = types.StringValue(*apiPolicy.Method)
	} else {
		policyModel.Method = types.StringNull()
	}

	// Handle RefreshField
	if apiPolicy.RefreshField != nil {
		policyModel.RefreshField = types.StringValue(*apiPolicy.RefreshField)
	} else {
		policyModel.RefreshField = types.StringNull()
	}

	// Handle NeverExpire (pointer to bool)
	if apiPolicy.NeverExpire != nil {
		policyModel.NeverExpire = types.BoolValue(*apiPolicy.NeverExpire)
	} else {
		policyModel.NeverExpire = types.BoolNull()
	}

	policyObj, d := types.ObjectValueFrom(ctx, attrTypes, policyModel)
	diags.Append(d...)
	return policyObj, diags
}

// ConvertAccelerationRefreshPolicyFromTerraform converts Terraform AccelerationRefreshPolicy state to API request format.
//
// Parameters:
//   - ctx: Context for the operation
//   - policyObj: The acceleration refresh policy from Terraform state/plan
//
// Returns:
//   - *models.AccelerationRefreshPolicy: The converted acceleration refresh policy for API requests (nil if input is null/unknown)
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertAccelerationRefreshPolicyFromTerraform(
	ctx context.Context,
	policyObj types.Object,
) (*models.AccelerationRefreshPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	if policyObj.IsNull() || policyObj.IsUnknown() {
		return nil, diags
	}

	var policyModel models.AccelerationRefreshPolicyModel
	diagsL := policyObj.As(ctx, &policyModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return nil, diags
	}

	result := &models.AccelerationRefreshPolicy{}

	// Handle ActivePolicyType
	if !policyModel.ActivePolicyType.IsNull() && !policyModel.ActivePolicyType.IsUnknown() {
		activePolicyType := policyModel.ActivePolicyType.ValueString()
		result.ActivePolicyType = &activePolicyType
	}

	// Handle RefreshPeriodMs
	if !policyModel.RefreshPeriodMs.IsNull() && !policyModel.RefreshPeriodMs.IsUnknown() {
		refreshPeriodMs := policyModel.RefreshPeriodMs.ValueInt64()
		result.RefreshPeriodMs = &refreshPeriodMs
	}

	// Handle RefreshSchedule
	if !policyModel.RefreshSchedule.IsNull() && !policyModel.RefreshSchedule.IsUnknown() {
		refreshSchedule := policyModel.RefreshSchedule.ValueString()
		result.RefreshSchedule = &refreshSchedule
	}

	// Handle GracePeriodMs
	if !policyModel.GracePeriodMs.IsNull() && !policyModel.GracePeriodMs.IsUnknown() {
		gracePeriodMs := policyModel.GracePeriodMs.ValueInt64()
		result.GracePeriodMs = &gracePeriodMs
	}

	// Handle Method
	if !policyModel.Method.IsNull() && !policyModel.Method.IsUnknown() {
		method := policyModel.Method.ValueString()
		result.Method = &method
	}

	// Handle RefreshField
	if !policyModel.RefreshField.IsNull() && !policyModel.RefreshField.IsUnknown() {
		refreshField := policyModel.RefreshField.ValueString()
		result.RefreshField = &refreshField
	}

	// Handle NeverExpire (pointer to bool)
	if !policyModel.NeverExpire.IsNull() && !policyModel.NeverExpire.IsUnknown() {
		neverExpire := policyModel.NeverExpire.ValueBool()
		result.NeverExpire = &neverExpire
	}

	return result, diags
}
