package helpers

import (
	"context"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetProjectCredentialsAttrTypes returns the attribute type definitions for ProjectCredentials structures.
func GetProjectCredentialsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                   types.StringType,
		"access_key_id":          types.StringType,
		"secret_access_key":      types.StringType,
		"role_arn":               types.StringType,
		"instance_profile_arn":   types.StringType,
		"external_id":            types.StringType,
		"external_id_signature":  types.StringType,
		"tenant_id":              types.StringType,
		"client_id":              types.StringType,
		"client_secret":          types.StringType,
		"account_name":           types.StringType,
	}
}

// ConvertProjectCredentialsFromTerraform converts Terraform ProjectCredentials state to API request format.
// This is a request-only structure, so there is no ToTerraform conversion.
//
// Parameters:
//   - ctx: Context for the operation
//   - credentialsObj: The project credentials from Terraform state/plan
//
// Returns:
//   - *models.ProjectCredentials: The converted project credentials for API requests (nil if input is null/unknown)
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertProjectCredentialsFromTerraform(
	ctx context.Context,
	credentialsObj types.Object,
) (*models.ProjectCredentials, diag.Diagnostics) {
	var diags diag.Diagnostics

	if credentialsObj.IsNull() || credentialsObj.IsUnknown() {
		return nil, diags
	}

	var credentialsModel models.ProjectCredentialsModel
	diagsL := credentialsObj.As(ctx, &credentialsModel, basetypes.ObjectAsOptions{})
	if diagsL.HasError() {
		diags.Append(diagsL...)
		return nil, diags
	}

	result := &models.ProjectCredentials{
		Type: credentialsModel.Type.ValueString(),
	}

	if !credentialsModel.AccessKeyID.IsNull() && !credentialsModel.AccessKeyID.IsUnknown() {
		result.AccessKeyID = credentialsModel.AccessKeyID.ValueString()
	}

	if !credentialsModel.SecretAccessKey.IsNull() && !credentialsModel.SecretAccessKey.IsUnknown() {
		result.SecretAccessKey = credentialsModel.SecretAccessKey.ValueString()
	}

	if !credentialsModel.RoleArn.IsNull() && !credentialsModel.RoleArn.IsUnknown() {
		result.RoleArn = credentialsModel.RoleArn.ValueString()
	}

	if !credentialsModel.InstanceProfileArn.IsNull() && !credentialsModel.InstanceProfileArn.IsUnknown() {
		result.InstanceProfileArn = credentialsModel.InstanceProfileArn.ValueString()
	}

	if !credentialsModel.ExternalID.IsNull() && !credentialsModel.ExternalID.IsUnknown() {
		result.ExternalID = credentialsModel.ExternalID.ValueString()
	}

	if !credentialsModel.ExternalIDSignature.IsNull() && !credentialsModel.ExternalIDSignature.IsUnknown() {
		result.ExternalIDSignature = credentialsModel.ExternalIDSignature.ValueString()
	}

	if !credentialsModel.TenantID.IsNull() && !credentialsModel.TenantID.IsUnknown() {
		result.TenantID = credentialsModel.TenantID.ValueString()
	}

	if !credentialsModel.ClientID.IsNull() && !credentialsModel.ClientID.IsUnknown() {
		result.ClientID = credentialsModel.ClientID.ValueString()
	}

	if !credentialsModel.ClientSecret.IsNull() && !credentialsModel.ClientSecret.IsUnknown() {
		result.ClientSecret = credentialsModel.ClientSecret.ValueString()
	}

	if !credentialsModel.AccountName.IsNull() && !credentialsModel.AccountName.IsUnknown() {
		result.AccountName = credentialsModel.AccountName.ValueString()
	}

	return result, diags
}

