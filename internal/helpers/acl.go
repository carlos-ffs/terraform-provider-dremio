package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetACLAttrTypes returns the attribute type definitions for ACL structures.
// These are used for creating typed null values and converting between Go structs and Terraform values.
func GetACLAttrTypes() (userPermission, rolePermission, accessControl map[string]attr.Type) {
	userPermission = map[string]attr.Type{
		"id":          types.StringType,
		"permissions": types.ListType{ElemType: types.StringType},
	}
	rolePermission = map[string]attr.Type{
		"id":          types.StringType,
		"permissions": types.ListType{ElemType: types.StringType},
	}
	accessControl = map[string]attr.Type{
		"users": types.ListType{ElemType: types.ObjectType{AttrTypes: userPermission}},
		"roles": types.ListType{ElemType: types.ObjectType{AttrTypes: rolePermission}},
	}
	return
}

// ConvertACLToTerraform converts API ACL response to Terraform state.
// It handles the conversion from models.AccessControlList (API response) to types.Object (Terraform state).
//
// Parameters:
//   - ctx: Context for the operation
//   - apiACL: The access control list from the API response (can be nil)
//   - planACL: The access control list from the Terraform plan/state
//
// Returns:
//   - types.Object: The converted ACL as a Terraform object
//   - diag.Diagnostics: Any diagnostics encountered during conversion
//
// Behavior:
//   - If planACL is null and apiACL is not nil, returns null (avoids drift from API defaults)
//   - If apiACL is nil, returns a typed null object
//   - If planACL is null, returns a typed null object
//   - Otherwise, converts the API ACL to Terraform format
func ConvertACLToTerraform(
	ctx context.Context,
	apiACL *models.AccessControlList,
	planACL types.Object,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	userPermissionAttrTypes, rolePermissionAttrTypes, accessControlAttrTypes := GetACLAttrTypes()

	// If the user didn't specify access_control_list in the plan (it's null in state),
	// keep it null instead of populating with API defaults to avoid drift.
	if planACL.IsNull() && apiACL != nil {
		return types.ObjectNull(accessControlAttrTypes), diags
	}

	if apiACL == nil {
		return types.ObjectNull(accessControlAttrTypes), diags
	}

	if planACL.IsNull() {
		return types.ObjectNull(accessControlAttrTypes), diags
	}

	aclModel := models.AccessControlListModel{
		Users: types.ListNull(types.ObjectType{AttrTypes: userPermissionAttrTypes}),
		Roles: types.ListNull(types.ObjectType{AttrTypes: rolePermissionAttrTypes}),
	}

	// Convert users
	if len(apiACL.Users) > 0 {
		userModels := make([]models.UserPermissionModel, 0, len(apiACL.Users))
		for _, u := range apiACL.Users {
			permsList, d := types.ListValueFrom(ctx, types.StringType, u.Permissions)
			diags.Append(d...)
			userModels = append(userModels, models.UserPermissionModel{
				ID:          types.StringValue(u.ID),
				Permissions: permsList,
			})
		}
		usersList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: userPermissionAttrTypes}, userModels)
		diags.Append(d...)
		aclModel.Users = usersList
	}

	// Convert roles
	if len(apiACL.Roles) > 0 {
		roleModels := make([]models.RolePermissionModel, 0, len(apiACL.Roles))
		for _, r := range apiACL.Roles {
			permsList, d := types.ListValueFrom(ctx, types.StringType, r.Permissions)
			diags.Append(d...)
			roleModels = append(roleModels, models.RolePermissionModel{
				ID:          types.StringValue(r.ID),
				Permissions: permsList,
			})
		}
		rolesList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: rolePermissionAttrTypes}, roleModels)
		diags.Append(d...)
		aclModel.Roles = rolesList
	}

	aclObj, d := types.ObjectValueFrom(ctx, accessControlAttrTypes, aclModel)
	diags.Append(d...)
	return aclObj, diags
}

// ConvertACLFromTerraform converts Terraform ACL state to API request format.
// It handles the conversion from types.Object (Terraform state) to *models.AccessControlList (API request).
//
// Parameters:
//   - ctx: Context for the operation
//   - aclObj: The access control list from Terraform state/plan
//
// Returns:
//   - *models.AccessControlList: The converted ACL for API requests (nil if input is null/unknown)
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertACLFromTerraform(
	ctx context.Context,
	aclObj types.Object,
) (*models.AccessControlList, diag.Diagnostics) {
	var diags diag.Diagnostics

	if aclObj.IsNull() || aclObj.IsUnknown() {
		return nil, diags
	}

	var acl models.AccessControlListModel
	diagsL := aclObj.As(ctx, &acl, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return nil, diags
	}

	result := &models.AccessControlList{}

	// Handle Users
	if !acl.Users.IsNull() && !acl.Users.IsUnknown() {
		var users []models.UserPermissionModel
		diagsL := acl.Users.ElementsAs(ctx, &users, false)
		if diagsL.HasError() {
			diags.Append(diagsL...)
			return nil, diags
		}

		for _, user := range users {
			var permissions []string
			diagsL := user.Permissions.ElementsAs(ctx, &permissions, false)
			if diagsL.HasError() {
				diags.Append(diagsL...)
				return nil, diags
			}

			result.Users = append(result.Users, models.UserPermission{
				ID:          user.ID.ValueString(),
				Permissions: permissions,
			})
		}
	}

	// Handle Roles
	if !acl.Roles.IsNull() && !acl.Roles.IsUnknown() {
		var roles []models.RolePermissionModel
		diagsL := acl.Roles.ElementsAs(ctx, &roles, false)
		if diagsL.HasError() {
			diags.Append(diagsL...)
			return nil, diags
		}

		for _, role := range roles {
			var permissions []string
			diagsL := role.Permissions.ElementsAs(ctx, &permissions, false)
			if diagsL.HasError() {
				diags.Append(diagsL...)
				return nil, diags
			}

			result.Roles = append(result.Roles, models.RolePermission{
				ID:          role.ID.ValueString(),
				Permissions: permissions,
			})
		}
	}

	return result, diags
}
