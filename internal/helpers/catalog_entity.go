package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetCatalogEntityAttrTypes returns the attribute type definitions for CatalogEntity structures.
func GetCatalogEntityAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.StringType,
		"path":           types.ListType{ElemType: types.StringType},
		"tag":            types.StringType,
		"type":           types.StringType,
		"container_type": types.StringType,
		"dataset_type":   types.StringType,
	}
}

// GetFolderChildAttrTypes returns the attribute type definitions for FolderChild structures.
func GetFolderChildAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.StringType,
		"path":           types.ListType{ElemType: types.StringType},
		"tag":            types.StringType,
		"type":           types.StringType,
		"container_type": types.StringType,
		"dataset_type":   types.StringType,
		"created_at":     types.StringType,
	}
}

// ConvertCatalogEntityToTerraform converts API CatalogEntity response to Terraform state.
//
// Parameters:
//   - ctx: Context for the operation
//   - apiEntity: The catalog entity from the API response
//
// Returns:
//   - types.Object: The converted catalog entity as a Terraform object
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertCatalogEntityToTerraform(
	ctx context.Context,
	apiEntity models.CatalogEntity,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := GetCatalogEntityAttrTypes()

	pathList, d := types.ListValueFrom(ctx, types.StringType, apiEntity.Path)
	diags.Append(d...)

	entityModel := models.ChildEntityModel{
		ID:            types.StringValue(apiEntity.ID),
		Path:          pathList,
		Tag:           types.StringValue(apiEntity.Tag),
		Type:          types.StringValue(apiEntity.Type),
		ContainerType: types.StringValue(apiEntity.ContainerType),
		DatasetType:   types.StringValue(apiEntity.DatasetType),
	}

	entityObj, d := types.ObjectValueFrom(ctx, attrTypes, entityModel)
	diags.Append(d...)
	return entityObj, diags
}

// ConvertFolderChildToTerraform converts API FolderChild response to Terraform state.
//
// Parameters:
//   - ctx: Context for the operation
//   - apiChild: The folder child from the API response
//
// Returns:
//   - types.Object: The converted folder child as a Terraform object
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertFolderChildToTerraform(
	ctx context.Context,
	apiChild models.FolderChild,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := GetFolderChildAttrTypes()

	pathList, d := types.ListValueFrom(ctx, types.StringType, apiChild.Path)
	diags.Append(d...)

	childModel := models.FolderChildModel{
		ID:            types.StringValue(apiChild.ID),
		Path:          pathList,
		Tag:           types.StringValue(apiChild.Tag),
		Type:          types.StringValue(apiChild.Type),
		ContainerType: types.StringValue(apiChild.ContainerType),
		DatasetType:   types.StringValue(apiChild.DatasetType),
		CreatedAt:     types.StringValue(apiChild.CreatedAt),
	}

	childObj, d := types.ObjectValueFrom(ctx, attrTypes, childModel)
	diags.Append(d...)
	return childObj, diags
}

