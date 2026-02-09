package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetCloudTagAttrTypes returns the attribute type definitions for CloudTag structures.
func GetCloudTagAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}

// ConvertCloudTagFromTerraform converts Terraform CloudTag state to API request format.
// This is a request-only structure, so there is no ToTerraform conversion.
//
// Parameters:
//   - ctx: Context for the operation
//   - tagObj: The cloud tag from Terraform state/plan
//
// Returns:
//   - models.CloudTag: The converted cloud tag for API requests
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertCloudTagFromTerraform(
	ctx context.Context,
	tagObj types.Object,
) (models.CloudTag, diag.Diagnostics) {
	var diags diag.Diagnostics

	var tagModel models.CloudTagModel
	diagsL := tagObj.As(ctx, &tagModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return models.CloudTag{}, diags
	}

	result := models.CloudTag{
		Key:   tagModel.Key.ValueString(),
		Value: tagModel.Value.ValueString(),
	}

	return result, diags
}

