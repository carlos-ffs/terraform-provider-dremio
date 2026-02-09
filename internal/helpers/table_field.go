package helpers

import (
	"context"
	"encoding/json"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ConvertTableFieldsToJSON converts API TableField slice to a JSON string for Terraform state.
// This approach is used because Terraform's schema system does not support truly recursive schemas.
// Table schemas can be arbitrarily deep with nested STRUCT and LIST types.
//
// Parameters:
//   - ctx: Context for the operation
//   - apiFields: The table fields from the API response
//
// Returns:
//   - types.String: The fields as a JSON string
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertTableFieldsToJSON(
	ctx context.Context,
	apiFields []models.TableField,
) (basetypes.StringValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(apiFields) == 0 {
		return types.StringNull(), diags
	}

	// Marshal the fields to JSON
	jsonBytes, err := json.Marshal(apiFields)
	if err != nil {
		diags.AddError(
			"JSON Conversion Error",
			"Unable to convert table fields to JSON: "+err.Error(),
		)
		return types.StringNull(), diags
	}

	return types.StringValue(string(jsonBytes)), diags
}
