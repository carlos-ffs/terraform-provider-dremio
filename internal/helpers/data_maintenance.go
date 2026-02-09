package helpers

import (
	"context"
	"fmt"

	"github.com/carlos-ffs/dremio-terraform-provider/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ConvertMaintenanceTaskToTerraform converts API MaintenanceTaskResponse to Terraform state.
//
// Parameters:
//   - ctx: Context for the operation
//   - taskResp: The maintenance task response from the API
//   - state: The Terraform state model to update
//
// Returns:
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertMaintenanceTaskToTerraform(
	ctx context.Context,
	taskResp *models.MaintenanceTaskResponse,
	state *models.DremioDataMaintenanceModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	state.ID = types.StringValue(taskResp.ID)
	state.TaskType = types.StringValue(taskResp.TaskType)
	state.Level = types.StringValue(taskResp.Level)
	state.SourceName = types.StringValue(taskResp.SourceName)
	state.IsEnabled = types.BoolValue(taskResp.IsEnabled)

	if taskResp.TaskConfig != nil {
		state.TableID = types.StringValue(taskResp.TaskConfig.TableID)
	}

	tflog.Debug(ctx, fmt.Sprintf("ConvertMaintenanceTaskToTerraform: ID=%s, Type=%s, Level=%s, SourceName=%s, IsEnabled=%t",
		taskResp.ID, taskResp.TaskType, taskResp.Level, taskResp.SourceName, taskResp.IsEnabled))

	return diags
}

// ConvertMaintenanceTaskFromTerraform converts Terraform state/plan to API request body.
//
// Parameters:
//   - ctx: Context for the operation
//   - data: The Terraform state/plan model
//
// Returns:
//   - *models.MaintenanceTaskRequest: The API request body
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertMaintenanceTaskFromTerraform(
	ctx context.Context,
	data *models.DremioDataMaintenanceModel,
) (*models.MaintenanceTaskRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	reqBody := &models.MaintenanceTaskRequest{
		TaskType:  data.TaskType.ValueString(),
		IsEnabled: data.IsEnabled.ValueBool(),
		TaskConfig: &models.MaintenanceTaskConfig{
			TableID: data.TableID.ValueString(),
		},
	}

	return reqBody, diags
}

// ConvertMaintenanceTaskToDataSource converts API MaintenanceTaskResponse to Terraform data source state.
//
// Parameters:
//   - ctx: Context for the operation
//   - taskResp: The maintenance task response from the API
//   - state: The Terraform data source state model to update
//
// Returns:
//   - diag.Diagnostics: Any diagnostics encountered during conversion
func ConvertMaintenanceTaskToDataSource(
	ctx context.Context,
	taskResp *models.MaintenanceTaskResponse,
	state *models.DremioDataMaintenanceDataSourceModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	state.ID = types.StringValue(taskResp.ID)
	state.TaskType = types.StringValue(taskResp.TaskType)
	state.Level = types.StringValue(taskResp.Level)
	state.SourceName = types.StringValue(taskResp.SourceName)
	state.IsEnabled = types.BoolValue(taskResp.IsEnabled)

	if taskResp.TaskConfig != nil {
		state.TableID = types.StringValue(taskResp.TaskConfig.TableID)
	}

	tflog.Debug(ctx, fmt.Sprintf("ConvertMaintenanceTaskToDataSource: ID=%s, Type=%s, Level=%s, SourceName=%s, IsEnabled=%t",
		taskResp.ID, taskResp.TaskType, taskResp.Level, taskResp.SourceName, taskResp.IsEnabled))

	return diags
}
