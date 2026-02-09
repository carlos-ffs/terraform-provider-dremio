package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetOwnerAttrTypes returns the attribute type definitions for Owner structures.
func GetOwnerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"owner_id":   types.StringType,
		"owner_type": types.StringType,
	}
}

// ConvertOwnerToTerraform converts API Owner response to Terraform state.
//
// Parameters:
//   - ctx: Context for the operation
//   - apiOwner: The owner from the API response (can be nil)
//
// Returns:
//   - types.Object: The converted owner as a Terraform object
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertOwnerToTerraform(
	ctx context.Context,
	apiOwner *models.Owner,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := GetOwnerAttrTypes()

	if apiOwner == nil {
		return types.ObjectNull(attrTypes), diags
	}

	ownerModel := models.OwnerModel{
		OwnerID:   types.StringValue(apiOwner.OwnerID),
		OwnerType: types.StringValue(apiOwner.OwnerType),
	}

	ownerObj, d := types.ObjectValueFrom(ctx, attrTypes, ownerModel)
	diags.Append(d...)
	return ownerObj, diags
}

// ConvertOwnerFromTerraform converts Terraform Owner state to API request format.
//
// Parameters:
//   - ctx: Context for the operation
//   - ownerObj: The owner from Terraform state/plan
//
// Returns:
//   - *models.Owner: The converted owner for API requests (nil if input is null/unknown)
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertOwnerFromTerraform(
	ctx context.Context,
	ownerObj types.Object,
) (*models.Owner, diag.Diagnostics) {
	var diags diag.Diagnostics

	if ownerObj.IsNull() || ownerObj.IsUnknown() {
		return nil, diags
	}

	var ownerModel models.OwnerModel
	diagsL := ownerObj.As(ctx, &ownerModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return nil, diags
	}

	result := &models.Owner{
		OwnerID:   ownerModel.OwnerID.ValueString(),
		OwnerType: ownerModel.OwnerType.ValueString(),
	}

	return result, diags
}

