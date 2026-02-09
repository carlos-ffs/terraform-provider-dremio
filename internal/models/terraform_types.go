package models

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dremioSourceModel describes the resource data model.
type DremioSourceModel struct {
	ID                               types.String         `tfsdk:"id"`
	EntityType                       types.String         `tfsdk:"entity_type"`
	Type                             types.String         `tfsdk:"type"`
	Name                             types.String         `tfsdk:"name"`
	Config                           jsontypes.Normalized `tfsdk:"config"`
	MetadataPolicy                   types.Object         `tfsdk:"metadata_policy"`
	AccelerationGracePeriodMs        types.Int64          `tfsdk:"acceleration_grace_period_ms"`
	AccelerationRefreshPeriodMs      types.Int64          `tfsdk:"acceleration_refresh_period_ms"`
	AccelerationActivePolicyType     types.String         `tfsdk:"acceleration_active_policy_type"`
	AccelerationRefreshSchedule      types.String         `tfsdk:"acceleration_refresh_schedule"`
	AccelerationRefreshOnDataChanges types.Bool           `tfsdk:"acceleration_refresh_on_data_changes"`
	AccessControlList                types.Object         `tfsdk:"access_control_list"`
	Tag                              types.String         `tfsdk:"tag"`
}

// DremioSourceDataSourceModel describes the data source data model.
type DremioSourceDataSourceModel struct {
	ID                           types.String         `tfsdk:"id"`
	Name                         types.String         `tfsdk:"name"`
	Tag                          types.String         `tfsdk:"tag"`
	Type                         types.String         `tfsdk:"type"`
	Config                       jsontypes.Normalized `tfsdk:"config"`
	MetadataPolicy               types.Object         `tfsdk:"metadata_policy"`
	AccelerationGracePeriodMs    types.Int64          `tfsdk:"acceleration_grace_period_ms"`
	AccelerationRefreshPeriodMs  types.Int64          `tfsdk:"acceleration_refresh_period_ms"`
	AccelerationNeverExpire      types.Bool           `tfsdk:"acceleration_never_expire"`
	AccelerationNeverRefresh     types.Bool           `tfsdk:"acceleration_never_refresh"`
	AccelerationActivePolicyType types.String         `tfsdk:"acceleration_active_policy_type"`
	AccelerationRefreshSchedule  types.String         `tfsdk:"acceleration_refresh_schedule"`
	Children                     types.List           `tfsdk:"children"`
	AccessControlList            types.Object         `tfsdk:"access_control_list"`
	Permissions                  types.List           `tfsdk:"permissions"`
	Owner                        types.Object         `tfsdk:"owner"`
}

// MetadataPolicyModel represents the metadata policy nested object
type MetadataPolicyModel struct {
	AuthTTLMs                 types.Int64  `tfsdk:"auth_ttl_ms"`
	NamesRefreshMs            types.Int64  `tfsdk:"names_refresh_ms"`
	DatasetRefreshAfterMs     types.Int64  `tfsdk:"dataset_refresh_after_ms"`
	DatasetExpireAfterMs      types.Int64  `tfsdk:"dataset_expire_after_ms"`
	DatasetUpdateMode         types.String `tfsdk:"dataset_update_mode"`
	DeleteUnavailableDatasets types.Bool   `tfsdk:"delete_unavailable_datasets"`
	AutoPromoteDatasets       types.Bool   `tfsdk:"auto_promote_datasets"`
}

// ChildEntityModel represents a child entity in the source
type ChildEntityModel struct {
	ID            types.String `tfsdk:"id"`
	Path          types.List   `tfsdk:"path"`
	Tag           types.String `tfsdk:"tag"`
	Type          types.String `tfsdk:"type"`
	ContainerType types.String `tfsdk:"container_type"`
	DatasetType   types.String `tfsdk:"dataset_type"`
}

// AccessControlListModel represents the access control list nested object
type AccessControlListModel struct {
	Users types.List `tfsdk:"users"`
	Roles types.List `tfsdk:"roles"`
}

// UserPermissionModel represents a user permission entry
type UserPermissionModel struct {
	ID          types.String `tfsdk:"id"`
	Permissions types.List   `tfsdk:"permissions"`
}

// RolePermissionModel represents a role permission entry
type RolePermissionModel struct {
	ID          types.String `tfsdk:"id"`
	Permissions types.List   `tfsdk:"permissions"`
}

// OwnerModel represents the owner information
type OwnerModel struct {
	OwnerID   types.String `tfsdk:"owner_id"`
	OwnerType types.String `tfsdk:"owner_type"`
}

type DremioFolderModel struct {
	ID                types.String `tfsdk:"id"`
	EntityType        types.String `tfsdk:"entity_type"`
	Path              types.List   `tfsdk:"path"`
	AccessControlList types.Object `tfsdk:"access_control_list"`
	Tag               types.String `tfsdk:"tag"`
}

// DremioUDFModel describes the UDF resource data model.
type DremioUDFModel struct {
	ID                types.String `tfsdk:"id"`
	EntityType        types.String `tfsdk:"entity_type"`
	Path              types.List   `tfsdk:"path"`
	IsScalar          types.Bool   `tfsdk:"is_scalar"`
	FunctionArgList   types.String `tfsdk:"function_arg_list"`
	FunctionBody      types.String `tfsdk:"function_body"`
	ReturnType        types.String `tfsdk:"return_type"`
	AccessControlList types.Object `tfsdk:"access_control_list"`
	Tag               types.String `tfsdk:"tag"`
}

// DremioFolderDataSourceModel describes the folder data source data model.
type DremioFolderDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Path              types.List   `tfsdk:"path"`
	EntityType        types.String `tfsdk:"entity_type"`
	MaxChildren       types.Int64  `tfsdk:"max_children"`
	Tag               types.String `tfsdk:"tag"`
	Children          types.List   `tfsdk:"children"`
	AccessControlList types.Object `tfsdk:"access_control_list"`
	Permissions       types.List   `tfsdk:"permissions"`
	Owner             types.Object `tfsdk:"owner"`
	StorageURI        types.String `tfsdk:"storage_uri"`
}

// FolderChildModel represents a child object within a folder
type FolderChildModel struct {
	ID            types.String `tfsdk:"id"`
	Path          types.List   `tfsdk:"path"`
	Tag           types.String `tfsdk:"tag"`
	Type          types.String `tfsdk:"type"`
	ContainerType types.String `tfsdk:"container_type"`
	DatasetType   types.String `tfsdk:"dataset_type"`
	CreatedAt     types.String `tfsdk:"created_at"`
}

// DremioUDFDataSourceModel describes the UDF data source data model.
type DremioUDFDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Path              types.List   `tfsdk:"path"`
	Tag               types.String `tfsdk:"tag"`
	CreatedAt         types.String `tfsdk:"created_at"`
	LastModified      types.String `tfsdk:"last_modified"`
	IsScalar          types.Bool   `tfsdk:"is_scalar"`
	FunctionArgList   types.String `tfsdk:"function_arg_list"`
	FunctionBody      types.String `tfsdk:"function_body"`
	ReturnType        types.String `tfsdk:"return_type"`
	AccessControlList types.Object `tfsdk:"access_control_list"`
	Permissions       types.List   `tfsdk:"permissions"`
	Owner             types.Object `tfsdk:"owner"`
}

// DremioFileDataSourceModel describes the file data source data model.
type DremioFileDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Path       types.List   `tfsdk:"path"`
	EntityType types.String `tfsdk:"entity_type"`
}

// DremioTableDataSourceModel describes the table data source data model.
type DremioTableDataSourceModel struct {
	ID                           types.String `tfsdk:"id"`
	Path                         types.List   `tfsdk:"path"`
	Type                         types.String `tfsdk:"type"`
	CreatedAt                    types.String `tfsdk:"created_at"`
	Tag                          types.String `tfsdk:"tag"`
	AccelerationRefreshPolicy    types.Object `tfsdk:"acceleration_refresh_policy"`
	Format                       types.Object `tfsdk:"format"`
	AccessControlList            types.Object `tfsdk:"access_control_list"`
	Owner                        types.Object `tfsdk:"owner"`
	Fields                       types.String `tfsdk:"fields"` // JSON string representation of table fields
	ApproximateStatisticsAllowed types.Bool   `tfsdk:"approximate_statistics_allowed"`
}

// TableFieldModel represents a field/column in a table or view
type TableFieldModel struct {
	Name types.String `tfsdk:"name"`
	Type types.Object `tfsdk:"type"`
}

// FieldTypeModel represents the type information for a field
type FieldTypeModel struct {
	Name        types.String `tfsdk:"name"`
	Precision   types.Int64  `tfsdk:"precision"`
	Scale       types.Int64  `tfsdk:"scale"`
	SubSchema   types.List   `tfsdk:"sub_schema"`
	ElementType types.Object `tfsdk:"element_type"`
}

// AccelerationRefreshPolicyModel represents the acceleration refresh policy for a dataset
type AccelerationRefreshPolicyModel struct {
	ActivePolicyType types.String `tfsdk:"active_policy_type"` // Policy for refreshing Reflections (NEVER, PERIOD, SCHEDULE, REFRESH_ON_DATA_CHANGES)
	RefreshPeriodMs  types.Int64  `tfsdk:"refresh_period_ms"`  // Refresh period in milliseconds (minimum 3600000, default 3600000)
	RefreshSchedule  types.String `tfsdk:"refresh_schedule"`   // Cron expression for refresh schedule (UTC), e.g., "0 0 8 * * ?"
	GracePeriodMs    types.Int64  `tfsdk:"grace_period_ms"`    // Maximum age for Reflection data in milliseconds
	Method           types.String `tfsdk:"method"`             // Method for refreshing Reflections (AUTO, FULL, INCREMENTAL)
	RefreshField     types.String `tfsdk:"refresh_field"`      // Field to use for incremental refresh
	NeverExpire      types.Bool   `tfsdk:"never_expire"`       // Whether Reflections never expire
}

// TableFormatModel represents the format information for a table (resource)
// This includes only the writable fields that can be sent in requests (TableFormatRequest).
type TableFormatModel struct {
	Type                    types.String `tfsdk:"type"`
	IgnoreOtherFileFormats  types.Bool   `tfsdk:"ignore_other_file_formats"`
	SkipFirstLine           types.Bool   `tfsdk:"skip_first_line"`
	ExtractHeader           types.Bool   `tfsdk:"extract_header"`
	HasMergedCells          types.Bool   `tfsdk:"has_merged_cells"`
	SheetName               types.String `tfsdk:"sheet_name"`
	FieldDelimiter          types.String `tfsdk:"field_delimiter"`
	Quote                   types.String `tfsdk:"quote"`
	Comment                 types.String `tfsdk:"comment"`
	Escape                  types.String `tfsdk:"escape"`
	LineDelimiter           types.String `tfsdk:"line_delimiter"`
	AutoGenerateColumnNames types.Bool   `tfsdk:"auto_generate_column_names"`
	TrimHeader              types.Bool   `tfsdk:"trim_header"`
}

// TableFormatDataSourceModel represents the format information for a table (datasource)
// This includes all fields from the API response, including read-only fields.
type TableFormatDataSourceModel struct {
	Type                    types.String `tfsdk:"type"`
	IgnoreOtherFileFormats  types.Bool   `tfsdk:"ignore_other_file_formats"`
	SkipFirstLine           types.Bool   `tfsdk:"skip_first_line"`
	ExtractHeader           types.Bool   `tfsdk:"extract_header"`
	HasMergedCells          types.Bool   `tfsdk:"has_merged_cells"`
	SheetName               types.String `tfsdk:"sheet_name"`
	FieldDelimiter          types.String `tfsdk:"field_delimiter"`
	Quote                   types.String `tfsdk:"quote"`
	Comment                 types.String `tfsdk:"comment"`
	Escape                  types.String `tfsdk:"escape"`
	LineDelimiter           types.String `tfsdk:"line_delimiter"`
	AutoGenerateColumnNames types.Bool   `tfsdk:"auto_generate_column_names"`
	TrimHeader              types.Bool   `tfsdk:"trim_header"`
	// Read-only fields (only in datasources)
	AutoCorrectCorruptDates types.Bool   `tfsdk:"auto_correct_corrupt_dates"`
	Name                    types.String `tfsdk:"name"`
	FullPath                types.List   `tfsdk:"full_path"` // List of strings
	Ctime                   types.Int64  `tfsdk:"ctime"`
	IsFolder                types.Bool   `tfsdk:"is_folder"`
	Location                types.String `tfsdk:"location"`
}

// PermissionsModel represents user permissions on a catalog object (response-only)
type PermissionsModel struct {
	CanView                  types.Bool `tfsdk:"can_view"`
	CanAlter                 types.Bool `tfsdk:"can_alter"`
	CanDelete                types.Bool `tfsdk:"can_delete"`
	CanManageGrants          types.Bool `tfsdk:"can_manage_grants"`
	CanEditAccessControlList types.Bool `tfsdk:"can_edit_access_control_list"`
	CanCreateChildren        types.Bool `tfsdk:"can_create_children"`
	CanRead                  types.Bool `tfsdk:"can_read"`
	CanEditFormatSettings    types.Bool `tfsdk:"can_edit_format_settings"`
	CanSelect                types.Bool `tfsdk:"can_select"`
	CanViewReflections       types.Bool `tfsdk:"can_view_reflections"`
	CanAlterReflections      types.Bool `tfsdk:"can_alter_reflections"`
	CanCreateReflections     types.Bool `tfsdk:"can_create_reflections"`
	CanDropReflections       types.Bool `tfsdk:"can_drop_reflections"`
}

// FunctionArgModel represents a function argument (response-only)
type FunctionArgModel struct {
	Name types.String `tfsdk:"name"`
	Type types.Object `tfsdk:"type"`
}

// SpaceChildModel represents a child entity in a space (response-only)
type SpaceChildModel struct {
	ID            types.String `tfsdk:"id"`
	Path          types.List   `tfsdk:"path"`
	Tag           types.String `tfsdk:"tag"`
	Type          types.String `tfsdk:"type"`
	ContainerType types.String `tfsdk:"container_type"`
	DatasetType   types.String `tfsdk:"dataset_type"`
	CreatedAt     types.String `tfsdk:"created_at"`
}

// GranteeRequestModel represents a grantee in a grants request (request-only)
type GranteeRequestModel struct {
	Privileges  types.Set    `tfsdk:"privileges"`
	GranteeType types.String `tfsdk:"grantee_type"`
	ID          types.String `tfsdk:"id"`
}

// ScriptGranteeModel represents a user or role with privileges on a script (request-only)
type ScriptGranteeModel struct {
	GranteeID  types.String `tfsdk:"grantee_id"`
	Privileges types.List   `tfsdk:"privileges"`
}

// CloudTagModel represents a cloud tag (AWS or Azure) (request-only)
type CloudTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// RuleInfoModel represents a single engine routing rule (request-only)
type RuleInfoModel struct {
	Name          types.String `tfsdk:"name"`
	Condition     types.String `tfsdk:"condition"`
	EngineName    types.String `tfsdk:"engine_name"`
	Action        types.String `tfsdk:"action"`
	RejectMessage types.String `tfsdk:"reject_message"`
}

// RuleSetModel represents a set of engine routing rules (request-only)
type RuleSetModel struct {
	RuleInfos       types.List   `tfsdk:"rule_infos"`
	RuleInfoDefault types.Object `tfsdk:"rule_info_default"`
	Tag             types.String `tfsdk:"tag"`
}

// ProjectCredentialsModel represents project storage credentials (request-only)
type ProjectCredentialsModel struct {
	Type                types.String `tfsdk:"type"`
	AccessKeyID         types.String `tfsdk:"access_key_id"`
	SecretAccessKey     types.String `tfsdk:"secret_access_key"`
	RoleArn             types.String `tfsdk:"role_arn"`
	InstanceProfileArn  types.String `tfsdk:"instance_profile_arn"`
	ExternalID          types.String `tfsdk:"external_id"`
	ExternalIDSignature types.String `tfsdk:"external_id_signature"`
	TenantID            types.String `tfsdk:"tenant_id"`
	ClientID            types.String `tfsdk:"client_id"`
	ClientSecret        types.String `tfsdk:"client_secret"`
	AccountName         types.String `tfsdk:"account_name"`
}

// DremioTableModel describes the table resource data model.
type DremioTableModel struct {
	ID                        types.String `tfsdk:"id"`
	EntityType                types.String `tfsdk:"entity_type"`
	Type                      types.String `tfsdk:"type"`
	Path                      types.List   `tfsdk:"path"`
	FileOrFolderID            types.String `tfsdk:"file_or_folder_id"`
	AccelerationRefreshPolicy types.Object `tfsdk:"acceleration_refresh_policy"`
	Format                    types.Object `tfsdk:"format"`
	AccessControlList         types.Object `tfsdk:"access_control_list"`
	Tag                       types.String `tfsdk:"tag"`
}

// DremioDatasetTagsModel describes the dataset tags resource data model.
type DremioDatasetTagsModel struct {
	DatasetID types.String `tfsdk:"dataset_id"`
	Tags      types.List   `tfsdk:"tags"`
	Version   types.String `tfsdk:"version"`
}

// DremioDatasetTagsDataSourceModel describes the dataset tags datasource data model.
type DremioDatasetTagsDataSourceModel struct {
	DatasetID types.String `tfsdk:"dataset_id"`
	Tags      types.List   `tfsdk:"tags"`
	Version   types.String `tfsdk:"version"`
}

// DremioDatasetWikiModel describes the dataset wiki resource data model.
type DremioDatasetWikiModel struct {
	DatasetID types.String `tfsdk:"dataset_id"`
	Text      types.String `tfsdk:"text"`
	Version   types.Int64  `tfsdk:"version"`
}

// DremioDatasetWikiDataSourceModel describes the dataset wiki datasource data model.
type DremioDatasetWikiDataSourceModel struct {
	DatasetID types.String `tfsdk:"dataset_id"`
	Text      types.String `tfsdk:"text"`
	Version   types.Int64  `tfsdk:"version"`
}

// DremioGrantsModel describes the grants resource data model.
type DremioGrantsModel struct {
	CatalogObjectID     types.String `tfsdk:"catalog_object_id"`
	Grants              types.Set    `tfsdk:"grants"` // Set of GranteeRequestModel
	AvailablePrivileges types.List   `tfsdk:"available_privileges"`
}

// DremioGrantsDataSourceModel describes the grants data source data model.
type DremioGrantsDataSourceModel struct {
	CatalogObjectID     types.String `tfsdk:"catalog_object_id"`
	Grants              types.Set    `tfsdk:"grants"` // Set of GranteesResponse as objects
	AvailablePrivileges types.List   `tfsdk:"available_privileges"`
}

// DremioViewModel describes the view resource data model.
type DremioViewModel struct {
	ID                types.String `tfsdk:"id"`
	EntityType        types.String `tfsdk:"entity_type"`
	Type              types.String `tfsdk:"type"`
	Path              types.List   `tfsdk:"path"`
	SQL               types.String `tfsdk:"sql"`
	SQLContext        types.List   `tfsdk:"sql_context"`
	AccessControlList types.Object `tfsdk:"access_control_list"`
	Tag               types.String `tfsdk:"tag"`
	Fields            types.String `tfsdk:"fields"` // JSON string representation of view fields
}

// DremioViewDataSourceModel describes the view data source data model.
type DremioViewDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Path              types.List   `tfsdk:"path"`
	Type              types.String `tfsdk:"type"`
	CreatedAt         types.String `tfsdk:"created_at"`
	Tag               types.String `tfsdk:"tag"`
	SQL               types.String `tfsdk:"sql"`
	SQLContext        types.List   `tfsdk:"sql_context"`
	Fields            types.String `tfsdk:"fields"` // JSON string representation of view fields
	AccessControlList types.Object `tfsdk:"access_control_list"`
	Permissions       types.List   `tfsdk:"permissions"`
	Owner             types.Object `tfsdk:"owner"`
}

// DremioEngineModel describes the engine resource data model.
type DremioEngineModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Size                  types.String `tfsdk:"size"`
	MinReplicas           types.Int64  `tfsdk:"min_replicas"`
	MaxReplicas           types.Int64  `tfsdk:"max_replicas"`
	AutoStopDelaySeconds  types.Int64  `tfsdk:"auto_stop_delay_seconds"`
	QueueTimeLimitSeconds types.Int64  `tfsdk:"queue_time_limit_seconds"`
	RuntimeLimitSeconds   types.Int64  `tfsdk:"runtime_limit_seconds"`
	DrainTimeLimitSeconds types.Int64  `tfsdk:"drain_time_limit_seconds"`
	MaxConcurrency        types.Int64  `tfsdk:"max_concurrency"`
	Description           types.String `tfsdk:"description"`
	Enable                types.Bool   `tfsdk:"enable"`
	// Computed fields
	State                     types.String `tfsdk:"state"`
	ActiveReplicas            types.Int64  `tfsdk:"active_replicas"`
	QueriedAt                 types.String `tfsdk:"queried_at"`
	StatusChangedAt           types.String `tfsdk:"status_changed_at"`
	InstanceFamily            types.String `tfsdk:"instance_family"`
	AdditionalEngineStateInfo types.String `tfsdk:"additional_engine_state_info"`
}

// DremioEngineDataSourceModel describes the engine data source model.
type DremioEngineDataSourceModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Size                      types.String `tfsdk:"size"`
	MinReplicas               types.Int64  `tfsdk:"min_replicas"`
	MaxReplicas               types.Int64  `tfsdk:"max_replicas"`
	AutoStopDelaySeconds      types.Int64  `tfsdk:"auto_stop_delay_seconds"`
	QueueTimeLimitSeconds     types.Int64  `tfsdk:"queue_time_limit_seconds"`
	RuntimeLimitSeconds       types.Int64  `tfsdk:"runtime_limit_seconds"`
	DrainTimeLimitSeconds     types.Int64  `tfsdk:"drain_time_limit_seconds"`
	MaxConcurrency            types.Int64  `tfsdk:"max_concurrency"`
	Description               types.String `tfsdk:"description"`
	State                     types.String `tfsdk:"state"`
	ActiveReplicas            types.Int64  `tfsdk:"active_replicas"`
	QueriedAt                 types.String `tfsdk:"queried_at"`
	StatusChangedAt           types.String `tfsdk:"status_changed_at"`
	AdditionalEngineStateInfo types.String `tfsdk:"additional_engine_state_info"`
}

// DremioEngineRuleSetModel describes the engine rule set resource data model.
type DremioEngineRuleSetModel struct {
	RuleInfos       types.List   `tfsdk:"rule_infos"`        // List of RuleInfoModel
	RuleInfoDefault types.Object `tfsdk:"rule_info_default"` // The default rule (cannot be deleted)
	Tag             types.String `tfsdk:"tag"`               // UUID for routing JDBC queries
}

// DremioDataMaintenanceModel describes the data maintenance task resource data model.
type DremioDataMaintenanceModel struct {
	ID         types.String `tfsdk:"id"`          // Unique identifier of the maintenance task
	TaskType   types.String `tfsdk:"type"`        // Type of maintenance task (OPTIMIZE, EXPIRE_SNAPSHOTS)
	Level      types.String `tfsdk:"level"`       // Level of the maintenance task (TABLE) - computed
	SourceName types.String `tfsdk:"source_name"` // Name of the source - computed
	IsEnabled  types.Bool   `tfsdk:"is_enabled"`  // Whether the task is enabled
	TableID    types.String `tfsdk:"table_id"`    // Fully qualified table name (e.g., "folder1.folder2.table1")
}

// DremioDataMaintenanceDataSourceModel describes the data maintenance task data source model.
type DremioDataMaintenanceDataSourceModel struct {
	ID         types.String `tfsdk:"id"`          // Unique identifier of the maintenance task (required)
	TaskType   types.String `tfsdk:"type"`        // Type of maintenance task (OPTIMIZE, EXPIRE_SNAPSHOTS)
	Level      types.String `tfsdk:"level"`       // Level of the maintenance task (TABLE)
	SourceName types.String `tfsdk:"source_name"` // Name of the source
	IsEnabled  types.Bool   `tfsdk:"is_enabled"`  // Whether the task is enabled
	TableID    types.String `tfsdk:"table_id"`    // Fully qualified table name (e.g., "folder1.folder2.table1")
}
