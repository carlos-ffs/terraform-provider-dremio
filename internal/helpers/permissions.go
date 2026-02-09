package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetPermissionsAttrTypes returns the attribute type definitions for Permissions structures.
func GetPermissionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"can_view":                     types.BoolType,
		"can_alter":                    types.BoolType,
		"can_delete":                   types.BoolType,
		"can_manage_grants":            types.BoolType,
		"can_edit_access_control_list": types.BoolType,
		"can_create_children":          types.BoolType,
		"can_read":                     types.BoolType,
		"can_edit_format_settings":     types.BoolType,
		"can_select":                   types.BoolType,
		"can_view_reflections":         types.BoolType,
		"can_alter_reflections":        types.BoolType,
		"can_create_reflections":       types.BoolType,
		"can_drop_reflections":         types.BoolType,
	}
}

// ConvertPermissionsToTerraform converts API Permissions response to Terraform state.
// This is a response-only structure, so there is no FromTerraform conversion.
//
// Parameters:
//   - ctx: Context for the operation
//   - apiPermissions: The permissions from the API response (can be nil)
//
// Returns:
//   - types.Object: The converted permissions as a Terraform object
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertPermissionsToTerraform(
	ctx context.Context,
	apiPermissions *models.Permissions,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := GetPermissionsAttrTypes()

	if apiPermissions == nil {
		return types.ObjectNull(attrTypes), diags
	}

	permissionsModel := models.PermissionsModel{
		CanView:                  types.BoolValue(apiPermissions.CanView),
		CanAlter:                 types.BoolValue(apiPermissions.CanAlter),
		CanDelete:                types.BoolValue(apiPermissions.CanDelete),
		CanManageGrants:          types.BoolValue(apiPermissions.CanManageGrants),
		CanEditAccessControlList: types.BoolValue(apiPermissions.CanEditAccessControlList),
		CanCreateChildren:        types.BoolValue(apiPermissions.CanCreateChildren),
		CanRead:                  types.BoolValue(apiPermissions.CanRead),
		CanEditFormatSettings:    types.BoolValue(apiPermissions.CanEditFormatSettings),
		CanSelect:                types.BoolValue(apiPermissions.CanSelect),
		CanViewReflections:       types.BoolValue(apiPermissions.CanViewReflections),
		CanAlterReflections:      types.BoolValue(apiPermissions.CanAlterReflections),
		CanCreateReflections:     types.BoolValue(apiPermissions.CanCreateReflections),
		CanDropReflections:       types.BoolValue(apiPermissions.CanDropReflections),
	}

	permissionsObj, d := types.ObjectValueFrom(ctx, attrTypes, permissionsModel)
	diags.Append(d...)
	return permissionsObj, diags
}
