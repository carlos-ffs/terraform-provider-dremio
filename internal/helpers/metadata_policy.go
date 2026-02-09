package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetMetadataPolicyAttrTypes returns the attribute type definitions for MetadataPolicy structures.
func GetMetadataPolicyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"auth_ttl_ms":                 types.Int64Type,
		"names_refresh_ms":            types.Int64Type,
		"dataset_refresh_after_ms":    types.Int64Type,
		"dataset_expire_after_ms":     types.Int64Type,
		"dataset_update_mode":         types.StringType,
		"delete_unavailable_datasets": types.BoolType,
		"auto_promote_datasets":       types.BoolType,
	}
}

// ConvertMetadataPolicyToTerraform converts API MetadataPolicy response to Terraform state.
//
// Parameters:
//   - ctx: Context for the operation
//   - apiPolicy: The metadata policy from the API response (can be nil)
//   - planPolicy: The metadata policy from the Terraform plan/state
//
// Returns:
//   - types.Object: The converted metadata policy as a Terraform object
//   - diag.Diagnostics: Any diagnostics encountered during conversion
//
// Behavior:
//   - If planPolicy is null and apiPolicy is not nil, returns null (avoids drift from API defaults)
//   - If apiPolicy is nil, returns a typed null object
//   - Otherwise, converts the API metadata policy to Terraform format
func ConvertMetadataPolicyToTerraform(
	ctx context.Context,
	apiPolicy *models.MetadataPolicy,
	planPolicy types.Object,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := GetMetadataPolicyAttrTypes()

	// If the user didn't specify metadata_policy in the plan (it's null in state),
	// keep it null instead of populating with API defaults to avoid drift.
	if planPolicy.IsNull() && apiPolicy != nil {
		return types.ObjectNull(attrTypes), diags
	}

	if apiPolicy == nil {
		return types.ObjectNull(attrTypes), diags
	}

	policyModel := models.MetadataPolicyModel{}

	if apiPolicy.AuthTTLMs != nil {
		policyModel.AuthTTLMs = types.Int64Value(*apiPolicy.AuthTTLMs)
	} else {
		policyModel.AuthTTLMs = types.Int64Null()
	}

	if apiPolicy.NamesRefreshMs != nil {
		policyModel.NamesRefreshMs = types.Int64Value(*apiPolicy.NamesRefreshMs)
	} else {
		policyModel.NamesRefreshMs = types.Int64Null()
	}

	if apiPolicy.DatasetRefreshAfterMs != nil {
		policyModel.DatasetRefreshAfterMs = types.Int64Value(*apiPolicy.DatasetRefreshAfterMs)
	} else {
		policyModel.DatasetRefreshAfterMs = types.Int64Null()
	}

	if apiPolicy.DatasetExpireAfterMs != nil {
		policyModel.DatasetExpireAfterMs = types.Int64Value(*apiPolicy.DatasetExpireAfterMs)
	} else {
		policyModel.DatasetExpireAfterMs = types.Int64Null()
	}

	if apiPolicy.DatasetUpdateMode != nil {
		policyModel.DatasetUpdateMode = types.StringValue(*apiPolicy.DatasetUpdateMode)
	} else {
		policyModel.DatasetUpdateMode = types.StringNull()
	}

	if apiPolicy.DeleteUnavailableDatasets != nil {
		policyModel.DeleteUnavailableDatasets = types.BoolValue(*apiPolicy.DeleteUnavailableDatasets)
	} else {
		policyModel.DeleteUnavailableDatasets = types.BoolNull()
	}

	if apiPolicy.AutoPromoteDatasets != nil {
		policyModel.AutoPromoteDatasets = types.BoolValue(*apiPolicy.AutoPromoteDatasets)
	} else {
		policyModel.AutoPromoteDatasets = types.BoolNull()
	}

	policyObj, d := types.ObjectValueFrom(ctx, attrTypes, policyModel)
	diags.Append(d...)
	return policyObj, diags
}

// ConvertMetadataPolicyFromTerraform converts Terraform MetadataPolicy state to API request format.
//
// Parameters:
//   - ctx: Context for the operation
//   - policyObj: The metadata policy from Terraform state/plan
//
// Returns:
//   - *models.MetadataPolicy: The converted metadata policy for API requests (nil if input is null/unknown)
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertMetadataPolicyFromTerraform(
	ctx context.Context,
	policyObj types.Object,
) (*models.MetadataPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	if policyObj.IsNull() || policyObj.IsUnknown() {
		return nil, diags
	}

	var policyModel models.MetadataPolicyModel
	diagsL := policyObj.As(ctx, &policyModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return nil, diags
	}

	result := &models.MetadataPolicy{}

	if !policyModel.AuthTTLMs.IsNull() && !policyModel.AuthTTLMs.IsUnknown() {
		v := policyModel.AuthTTLMs.ValueInt64()
		result.AuthTTLMs = &v
	}

	if !policyModel.NamesRefreshMs.IsNull() && !policyModel.NamesRefreshMs.IsUnknown() {
		v := policyModel.NamesRefreshMs.ValueInt64()
		result.NamesRefreshMs = &v
	}

	if !policyModel.DatasetRefreshAfterMs.IsNull() && !policyModel.DatasetRefreshAfterMs.IsUnknown() {
		v := policyModel.DatasetRefreshAfterMs.ValueInt64()
		result.DatasetRefreshAfterMs = &v
	}

	if !policyModel.DatasetExpireAfterMs.IsNull() && !policyModel.DatasetExpireAfterMs.IsUnknown() {
		v := policyModel.DatasetExpireAfterMs.ValueInt64()
		result.DatasetExpireAfterMs = &v
	}

	if !policyModel.DatasetUpdateMode.IsNull() && !policyModel.DatasetUpdateMode.IsUnknown() {
		v := policyModel.DatasetUpdateMode.ValueString()
		result.DatasetUpdateMode = &v
	}

	if !policyModel.DeleteUnavailableDatasets.IsNull() && !policyModel.DeleteUnavailableDatasets.IsUnknown() {
		v := policyModel.DeleteUnavailableDatasets.ValueBool()
		result.DeleteUnavailableDatasets = &v
	}

	if !policyModel.AutoPromoteDatasets.IsNull() && !policyModel.AutoPromoteDatasets.IsUnknown() {
		v := policyModel.AutoPromoteDatasets.ValueBool()
		result.AutoPromoteDatasets = &v
	}

	return result, diags
}
