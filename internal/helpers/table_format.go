package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetTableFormatAttrTypes returns the attribute type definitions for TableFormat structures.
// This includes only the writable fields that can be sent in requests (TableFormatRequest).
// Use this for resources.
func GetTableFormatAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                       types.StringType,
		"ignore_other_file_formats":  types.BoolType,
		"skip_first_line":            types.BoolType,
		"extract_header":             types.BoolType,
		"has_merged_cells":           types.BoolType,
		"sheet_name":                 types.StringType,
		"field_delimiter":            types.StringType,
		"quote":                      types.StringType,
		"comment":                    types.StringType,
		"escape":                     types.StringType,
		"line_delimiter":             types.StringType,
		"auto_generate_column_names": types.BoolType,
		"trim_header":                types.BoolType,
	}
}

// GetTableFormatDatasourceAttrTypes returns the attribute type definitions for TableFormat structures in datasources.
// This includes all fields from the API response, including read-only fields.
// Use this for datasources.
func GetTableFormatDatasourceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                       types.StringType,
		"ignore_other_file_formats":  types.BoolType,
		"skip_first_line":            types.BoolType,
		"extract_header":             types.BoolType,
		"has_merged_cells":           types.BoolType,
		"sheet_name":                 types.StringType,
		"field_delimiter":            types.StringType,
		"quote":                      types.StringType,
		"comment":                    types.StringType,
		"escape":                     types.StringType,
		"line_delimiter":             types.StringType,
		"auto_generate_column_names": types.BoolType,
		"trim_header":                types.BoolType,
		"auto_correct_corrupt_dates": types.BoolType,
		"name":                       types.StringType,
		"full_path":                  types.ListType{ElemType: types.StringType},
		"ctime":                      types.Int64Type,
		"is_folder":                  types.BoolType,
		"location":                   types.StringType,
	}
}

// ConvertTableFormatToTerraform converts API TableFormatResponse to Terraform state.
//
// Parameters:
//   - ctx: Context for the operation
//   - apiFormat: The table format from the API response (can be nil)
//   - planFormat: The table format from the Terraform plan/state
//
// Returns:
//   - types.Object: The converted table format as a Terraform object
//   - diag.Diagnostics: Any diagnostics encountered during conversion
//
// Behavior:
//   - If planFormat is null and apiFormat is not nil, returns null (avoids drift from API defaults)
//   - If apiFormat is nil, returns a typed null object
//   - If planFormat is unknown (datasource case), populates all fields from API response
//   - Otherwise, converts the API table format to Terraform format
//   - For individual fields: if a field is null in planFormat, keep it null (don't update with API values)
func ConvertTableFormatToTerraform(
	ctx context.Context,
	apiFormat *models.TableFormatResponse,
	planFormat types.Object,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := GetTableFormatAttrTypes()

	// If the user didn't specify format in the plan (it's null in state),
	// keep it null instead of populating with API defaults to avoid drift.
	if planFormat.IsNull() && apiFormat != nil {
		return types.ObjectNull(attrTypes), diags
	}

	if apiFormat == nil {
		return types.ObjectNull(attrTypes), diags
	}

	// If planFormat is unknown (datasource case), populate all fields from API
	if planFormat.IsUnknown() {
		formatModel := models.TableFormatModel{
			Type:                    types.StringValue(apiFormat.Type),
			IgnoreOtherFileFormats:  types.BoolPointerValue(apiFormat.IgnoreOtherFileFormats),
			SkipFirstLine:           types.BoolPointerValue(apiFormat.SkipFirstLine),
			ExtractHeader:           types.BoolPointerValue(apiFormat.ExtractHeader),
			HasMergedCells:          types.BoolPointerValue(apiFormat.HasMergedCells),
			SheetName:               types.StringPointerValue(apiFormat.SheetName),
			FieldDelimiter:          types.StringPointerValue(apiFormat.FieldDelimiter),
			Quote:                   types.StringPointerValue(apiFormat.Quote),
			Comment:                 types.StringPointerValue(apiFormat.Comment),
			Escape:                  types.StringPointerValue(apiFormat.Escape),
			LineDelimiter:           types.StringPointerValue(apiFormat.LineDelimiter),
			AutoGenerateColumnNames: types.BoolPointerValue(apiFormat.AutoGenerateColumnNames),
			TrimHeader:              types.BoolPointerValue(apiFormat.TrimHeader),
		}

		formatObj, d := types.ObjectValueFrom(ctx, attrTypes, formatModel)
		diags.Append(d...)
		return formatObj, diags
	}

	// Get the existing plan format to preserve null values
	var planFormatModel models.TableFormatModel
	diagsL := planFormat.As(ctx, &planFormatModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return types.ObjectNull(attrTypes), diags
	}

	// Build the format model, preserving null values from the plan
	formatModel := models.TableFormatModel{
		Type: types.StringValue(apiFormat.Type),
	}

	// For each optional field, only update if it was not null in the plan
	if !planFormatModel.IgnoreOtherFileFormats.IsNull() {
		formatModel.IgnoreOtherFileFormats = types.BoolPointerValue(apiFormat.IgnoreOtherFileFormats)
	} else {
		formatModel.IgnoreOtherFileFormats = types.BoolNull()
	}

	if !planFormatModel.SkipFirstLine.IsNull() {
		formatModel.SkipFirstLine = types.BoolPointerValue(apiFormat.SkipFirstLine)
	} else {
		formatModel.SkipFirstLine = types.BoolNull()
	}

	if !planFormatModel.ExtractHeader.IsNull() {
		formatModel.ExtractHeader = types.BoolPointerValue(apiFormat.ExtractHeader)
	} else {
		formatModel.ExtractHeader = types.BoolNull()
	}

	if !planFormatModel.HasMergedCells.IsNull() {
		formatModel.HasMergedCells = types.BoolPointerValue(apiFormat.HasMergedCells)
	} else {
		formatModel.HasMergedCells = types.BoolNull()
	}

	if !planFormatModel.SheetName.IsNull() {
		formatModel.SheetName = types.StringPointerValue(apiFormat.SheetName)
	} else {
		formatModel.SheetName = types.StringNull()
	}

	if !planFormatModel.FieldDelimiter.IsNull() {
		formatModel.FieldDelimiter = types.StringPointerValue(apiFormat.FieldDelimiter)
	} else {
		formatModel.FieldDelimiter = types.StringNull()
	}

	if !planFormatModel.Quote.IsNull() {
		formatModel.Quote = types.StringPointerValue(apiFormat.Quote)
	} else {
		formatModel.Quote = types.StringNull()
	}

	if !planFormatModel.Comment.IsNull() {
		formatModel.Comment = types.StringPointerValue(apiFormat.Comment)
	} else {
		formatModel.Comment = types.StringNull()
	}

	if !planFormatModel.Escape.IsNull() {
		formatModel.Escape = types.StringPointerValue(apiFormat.Escape)
	} else {
		formatModel.Escape = types.StringNull()
	}

	if !planFormatModel.LineDelimiter.IsNull() {
		formatModel.LineDelimiter = types.StringPointerValue(apiFormat.LineDelimiter)
	} else {
		formatModel.LineDelimiter = types.StringNull()
	}

	if !planFormatModel.AutoGenerateColumnNames.IsNull() {
		formatModel.AutoGenerateColumnNames = types.BoolPointerValue(apiFormat.AutoGenerateColumnNames)
	} else {
		formatModel.AutoGenerateColumnNames = types.BoolNull()
	}

	if !planFormatModel.TrimHeader.IsNull() {
		formatModel.TrimHeader = types.BoolPointerValue(apiFormat.TrimHeader)
	} else {
		formatModel.TrimHeader = types.BoolNull()
	}

	formatObj, d := types.ObjectValueFrom(ctx, attrTypes, formatModel)
	diags.Append(d...)
	return formatObj, diags
}

// ConvertTableFormatFromTerraform converts Terraform TableFormat state to API request format.
//
// Parameters:
//   - ctx: Context for the operation
//   - formatObj: The table format from Terraform state/plan
//
// Returns:
//   - *models.TableFormatRequest: The converted table format for API requests (nil if input is null/unknown)
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertTableFormatFromTerraform(
	ctx context.Context,
	formatObj types.Object,
) (*models.TableFormatRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	if formatObj.IsNull() || formatObj.IsUnknown() {
		return nil, diags
	}

	var formatModel models.TableFormatModel
	diagsL := formatObj.As(ctx, &formatModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return nil, diags
	}

	result := &models.TableFormatRequest{
		Type: formatModel.Type.ValueString(),
	}

	// Handle optional boolean fields
	if !formatModel.IgnoreOtherFileFormats.IsNull() && !formatModel.IgnoreOtherFileFormats.IsUnknown() {
		value := formatModel.IgnoreOtherFileFormats.ValueBool()
		result.IgnoreOtherFileFormats = &value
	}

	if !formatModel.SkipFirstLine.IsNull() && !formatModel.SkipFirstLine.IsUnknown() {
		value := formatModel.SkipFirstLine.ValueBool()
		result.SkipFirstLine = &value
	}

	if !formatModel.ExtractHeader.IsNull() && !formatModel.ExtractHeader.IsUnknown() {
		value := formatModel.ExtractHeader.ValueBool()
		result.ExtractHeader = &value
	}

	if !formatModel.HasMergedCells.IsNull() && !formatModel.HasMergedCells.IsUnknown() {
		value := formatModel.HasMergedCells.ValueBool()
		result.HasMergedCells = &value
	}

	if !formatModel.AutoGenerateColumnNames.IsNull() && !formatModel.AutoGenerateColumnNames.IsUnknown() {
		value := formatModel.AutoGenerateColumnNames.ValueBool()
		result.AutoGenerateColumnNames = &value
	}

	if !formatModel.TrimHeader.IsNull() && !formatModel.TrimHeader.IsUnknown() {
		value := formatModel.TrimHeader.ValueBool()
		result.TrimHeader = &value
	}

	// Handle optional string fields
	if !formatModel.SheetName.IsNull() && !formatModel.SheetName.IsUnknown() {
		value := formatModel.SheetName.ValueString()
		result.SheetName = &value
	}

	if !formatModel.FieldDelimiter.IsNull() && !formatModel.FieldDelimiter.IsUnknown() {
		value := formatModel.FieldDelimiter.ValueString()
		result.FieldDelimiter = &value
	}

	if !formatModel.Quote.IsNull() && !formatModel.Quote.IsUnknown() {
		value := formatModel.Quote.ValueString()
		result.Quote = &value
	}

	if !formatModel.Comment.IsNull() && !formatModel.Comment.IsUnknown() {
		value := formatModel.Comment.ValueString()
		result.Comment = &value
	}

	if !formatModel.Escape.IsNull() && !formatModel.Escape.IsUnknown() {
		value := formatModel.Escape.ValueString()
		result.Escape = &value
	}

	if !formatModel.LineDelimiter.IsNull() && !formatModel.LineDelimiter.IsUnknown() {
		value := formatModel.LineDelimiter.ValueString()
		result.LineDelimiter = &value
	}

	return result, diags
}

// ConvertTableFormatToTerraformDatasource converts API TableFormatResponse to Terraform state for datasources.
// This function populates all fields including read-only ones from the API response.
//
// Parameters:
//   - ctx: Context for the operation
//   - apiFormat: The table format from the API response (can be nil)
//
// Returns:
//   - types.Object: The converted table format as a Terraform object with all fields
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertTableFormatToTerraformDatasource(
	ctx context.Context,
	apiFormat *models.TableFormatResponse,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := GetTableFormatDatasourceAttrTypes()

	if apiFormat == nil {
		return types.ObjectNull(attrTypes), diags
	}

	// Convert FullPath to types.List
	var fullPathList types.List
	if len(apiFormat.FullPath) > 0 {
		fullPathList, _ = types.ListValueFrom(ctx, types.StringType, apiFormat.FullPath)
	} else {
		fullPathList = types.ListNull(types.StringType)
	}

	// Convert Ctime to types.Int64
	var ctimeValue types.Int64
	if apiFormat.Ctime != nil {
		ctimeValue = types.Int64Value(int64(*apiFormat.Ctime))
	} else {
		ctimeValue = types.Int64Null()
	}

	// Populate all fields from API response using TableFormatDataSourceModel
	formatModel := models.TableFormatDataSourceModel{
		Type:                    types.StringValue(apiFormat.Type),
		IgnoreOtherFileFormats:  types.BoolPointerValue(apiFormat.IgnoreOtherFileFormats),
		SkipFirstLine:           types.BoolPointerValue(apiFormat.SkipFirstLine),
		ExtractHeader:           types.BoolPointerValue(apiFormat.ExtractHeader),
		HasMergedCells:          types.BoolPointerValue(apiFormat.HasMergedCells),
		SheetName:               types.StringPointerValue(apiFormat.SheetName),
		FieldDelimiter:          types.StringPointerValue(apiFormat.FieldDelimiter),
		Quote:                   types.StringPointerValue(apiFormat.Quote),
		Comment:                 types.StringPointerValue(apiFormat.Comment),
		Escape:                  types.StringPointerValue(apiFormat.Escape),
		LineDelimiter:           types.StringPointerValue(apiFormat.LineDelimiter),
		AutoGenerateColumnNames: types.BoolPointerValue(apiFormat.AutoGenerateColumnNames),
		TrimHeader:              types.BoolPointerValue(apiFormat.TrimHeader),
		// Read-only fields
		AutoCorrectCorruptDates: types.BoolPointerValue(apiFormat.AutoCorrectCorruptDates),
		Name:                    types.StringPointerValue(apiFormat.Name),
		FullPath:                fullPathList,
		Ctime:                   ctimeValue,
		IsFolder:                types.BoolPointerValue(apiFormat.IsFolder),
		Location:                types.StringPointerValue(apiFormat.Location),
	}

	formatObj, d := types.ObjectValueFrom(ctx, attrTypes, formatModel)
	diags.Append(d...)
	return formatObj, diags
}
