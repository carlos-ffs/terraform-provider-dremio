package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetGranteeRequestAttrTypes returns the attribute type definitions for GranteeRequest structures.
// Uses SetType for privileges to ensure order-independent comparison.
func GetGranteeRequestAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"privileges":   types.SetType{ElemType: types.StringType},
		"grantee_type": types.StringType,
		"id":           types.StringType,
	}
}

// GetScriptGranteeAttrTypes returns the attribute type definitions for ScriptGrantee structures.
func GetScriptGranteeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"grantee_id": types.StringType,
		"privileges": types.ListType{ElemType: types.StringType},
	}
}

// ConvertGranteeRequestFromTerraform converts Terraform GranteeRequest state to API request format.
// This is a request-only structure, so there is no ToTerraform conversion.
//
// Parameters:
//   - ctx: Context for the operation
//   - granteeObj: The grantee request from Terraform state/plan
//
// Returns:
//   - models.GranteeRequest: The converted grantee request for API requests
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertGranteeRequestFromTerraform(
	ctx context.Context,
	granteeObj types.Object,
) (models.GranteeRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var granteeModel models.GranteeRequestModel
	diagsL := granteeObj.As(ctx, &granteeModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return models.GranteeRequest{}, diags
	}

	var privileges []string
	diagsL = granteeModel.Privileges.ElementsAs(ctx, &privileges, false)
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return models.GranteeRequest{}, diags
	}

	result := models.GranteeRequest{
		Privileges:  privileges,
		GranteeType: granteeModel.GranteeType.ValueString(),
		ID:          granteeModel.ID.ValueString(),
	}

	return result, diags
}

// ConvertScriptGranteeFromTerraform converts Terraform ScriptGrantee state to API request format.
// This is a request-only structure, so there is no ToTerraform conversion.
//
// Parameters:
//   - ctx: Context for the operation
//   - granteeObj: The script grantee from Terraform state/plan
//
// Returns:
//   - models.ScriptGrantee: The converted script grantee for API requests
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertScriptGranteeFromTerraform(
	ctx context.Context,
	granteeObj types.Object,
) (models.ScriptGrantee, diag.Diagnostics) {
	var diags diag.Diagnostics

	var granteeModel models.ScriptGranteeModel
	diagsL := granteeObj.As(ctx, &granteeModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return models.ScriptGrantee{}, diags
	}

	var privileges []string
	diagsL = granteeModel.Privileges.ElementsAs(ctx, &privileges, false)
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return models.ScriptGrantee{}, diags
	}

	result := models.ScriptGrantee{
		GranteeID:  granteeModel.GranteeID.ValueString(),
		Privileges: privileges,
	}

	return result, diags
}

// ConvertGranteesToTerraform converts API GranteesResponse list to Terraform set.
// This is used for the grants resource to convert API response to Terraform state.
// Uses sets for both grants and privileges to ensure order-independent comparison.
//
// Parameters:
//   - ctx: Context for the operation
//   - grants: The list of grantees from the API response
//
// Returns:
//   - types.Set: The converted grants as a Terraform set
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertGranteesToTerraform(
	ctx context.Context,
	grants []models.GranteesResponse,
) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	granteeAttrTypes := GetGranteeRequestAttrTypes()

	var grantObjects []attr.Value
	for _, grant := range grants {
		privsSet, d := types.SetValueFrom(ctx, types.StringType, grant.Privileges)
		diags.Append(d...)

		grantObj, d := types.ObjectValue(granteeAttrTypes, map[string]attr.Value{
			"id":           types.StringValue(grant.ID),
			"grantee_type": types.StringValue(grant.GranteeType),
			"privileges":   privsSet,
		})
		diags.Append(d...)
		grantObjects = append(grantObjects, grantObj)
	}

	grantsSet, d := types.SetValue(types.ObjectType{AttrTypes: granteeAttrTypes}, grantObjects)
	diags.Append(d...)

	return grantsSet, diags
}

// ConvertGranteeSetFromTerraform converts Terraform grants set to API request format.
//
// Parameters:
//   - ctx: Context for the operation
//   - grantsSet: The grants set from Terraform state/plan
//
// Returns:
//   - []models.GranteeRequest: The converted grantees for API requests
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertGranteeSetFromTerraform(
	ctx context.Context,
	grantsSet types.Set,
) ([]models.GranteeRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var granteeObjects []types.Object
	d := grantsSet.ElementsAs(ctx, &granteeObjects, false)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	var grants []models.GranteeRequest
	for _, granteeObj := range granteeObjects {
		grant, d := ConvertGranteeRequestFromTerraform(ctx, granteeObj)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		grants = append(grants, grant)
	}

	return grants, diags
}
