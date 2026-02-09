package models

// SourceCreateRequest represents a request to create a source
// Reference: OpenAPI schema SourceCreateRequest
type SourceRequest struct {
	EntityType                       string                 `json:"entityType"`                                 // Always "source"
	Type                             string                 `json:"type"`                                       // Source type (ARCTIC, S3, SNOWFLAKE, etc.)
	Name                             string                 `json:"name"`                                       // User-defined name of the source
	Config                           map[string]interface{} `json:"config"`                                     // Configuration options specific to the source type
	MetadataPolicy                   *MetadataPolicy        `json:"metadataPolicy,omitempty"`                   // Metadata refresh policy
	AccelerationGracePeriodMs        int64                  `json:"accelerationGracePeriodMs,omitempty"`        // Time to keep Reflections before expiration (milliseconds)
	AccelerationRefreshPeriodMs      int64                  `json:"accelerationRefreshPeriodMs,omitempty"`      // Refresh frequency for Reflections (milliseconds)
	AccelerationActivePolicyType     string                 `json:"accelerationActivePolicyType,omitempty"`     // Policy for refreshing Reflections (NEVER, PERIOD, SCHEDULE)
	AccelerationRefreshSchedule      string                 `json:"accelerationRefreshSchedule,omitempty"`      // Cron expression for Reflection refresh schedule (UTC)
	AccelerationRefreshOnDataChanges bool                   `json:"accelerationRefreshOnDataChanges,omitempty"` // Refresh Reflections when Iceberg table snapshots change
	AccessControlList                *AccessControlList     `json:"accessControlList,omitempty"`
	Tag                              string                 `json:"tag,omitempty"` // Version tag for optimistic concurrency control
	ID                               string                 `json:"id,omitempty"`  // Unique identifier of the source
}

// SpaceCreateRequest represents a request to create a space
// Reference: https://docs.dremio.com/current/reference/api/catalog/container-space#create-a-space
type SpaceCreateRequest struct {
	EntityType        string             `json:"entityType"`                  // Always "space"
	Name              string             `json:"name"`                        // Name of the space
	AccessControlList *AccessControlList `json:"accessControlList,omitempty"` // Optional: User and role access settings
}

// SpaceUpdateRequest represents a request to update a space
// Reference: https://docs.dremio.com/current/reference/api/catalog/container-space#update-a-space
type SpaceUpdateRequest struct {
	EntityType        string             `json:"entityType"`                  // Always "space"
	ID                string             `json:"id"`                          // Unique identifier of the space
	Name              string             `json:"name"`                        // Name of the space
	Tag               string             `json:"tag"`                         // Version tag for optimistic concurrency control
	AccessControlList *AccessControlList `json:"accessControlList,omitempty"` // Optional: User and role access settings
}

// FolderCreateRequest represents a request to add a folder in an Arctic source
// Note: Adding folders is only supported in Arctic sources
// Reference: https://docs.dremio.com/cloud/reference/api/catalog/folder#adding-a-folder
type FolderCreateRequest struct {
	EntityType string   `json:"entityType"` // Always "folder"
	Path       []string `json:"path"`       // Path including the new folder name as the last item
}

// FolderUpdateRequest represents a request to update a folder in a non-Arctic source
// Note: Updating folders is only supported in non-Arctic sources
// Reference: https://docs.dremio.com/cloud/reference/api/catalog/folder#updating-a-folder
type FolderUpdateRequest struct {
	EntityType        string             `json:"entityType"`                  // Always "folder"
	ID                string             `json:"id"`                          // Unique identifier of the folder
	Path              []string           `json:"path"`                        // New path for the folder (for renaming/moving)
	Tag               string             `json:"tag,omitempty"`               // Optional: Version tag for optimistic concurrency control
	AccessControlList *AccessControlList `json:"accessControlList,omitempty"` // Optional: User and role access settings
}

// TableCreateRequest represents a request to format a file or folder as a table
// Reference: https://docs.dremio.com/cloud/reference/api/catalog/table#formatting-a-file-or-folder-as-a-table
type TableRequest struct {
	SourceOrFolderID          string                     `json:"id"`                                  // Unique identifier of the source or folder
	EntityType                string                     `json:"entityType"`                          // Always "dataset"
	Path                      []string                   `json:"path"`                                // Path to the file or folder
	Type                      string                     `json:"type"`                                // Always "PHYSICAL_DATASET"
	Tag                       string                     `json:"tag,omitempty"`                       // Optional: Version tag for optimistic concurrency control
	AccelerationRefreshPolicy *AccelerationRefreshPolicy `json:"accelerationRefreshPolicy,omitempty"` // Acceleration refresh policy
	Format                    *TableFormatRequest        `json:"format"`                              // Format parameters
	AccessControlList         *AccessControlList         `json:"accessControlList,omitempty"`         // Optional: User and role access settings
}

// UDFCreateRequest represents a request to create a user-defined function
// Reference: https://docs.dremio.com/cloud/reference/api/catalog/user-defined-function
type UDFRequest struct {
	EntityType        string             `json:"entityType"`                  // Always "function"
	Path              []string           `json:"path"`                        // Path where the UDF should be created
	IsScalar          bool               `json:"isScalar"`                    // true for scalar function, false for tabular function
	FunctionArgList   string             `json:"functionArgList"`             // Arguments and their data types
	FunctionBody      string             `json:"functionBody"`                // Statement that the UDF executes
	ReturnType        string             `json:"returnType"`                  // Data type(s) that the UDF returns
	AccessControlList *AccessControlList `json:"accessControlList,omitempty"` // Optional: User and role access settings
	// Used for update operations
	Tag string `json:"tag,omitempty"` // Version tag for optimistic concurrency control
	ID  string `json:"id,omitempty"`  // Unique identifier of the UDF
}

// ViewCreateRequest represents a request to create a view
// Reference: https://docs.dremio.com/cloud/reference/api/catalog/view
type ViewRequest struct {
	EntityType        string             `json:"entityType"`                  // Always "dataset"
	Path              []string           `json:"path"`                        // Path to the location where the view should be created
	SQL               string             `json:"sql"`                         // SQL query to use to create the view
	Type              string             `json:"type"`                        // Always "VIRTUAL_DATASET"
	SQLContext        []string           `json:"sqlContext"`                  // Context for the SQL query
	AccessControlList *AccessControlList `json:"accessControlList,omitempty"` // Optional: User and role access settings
	ID                string             `json:"id,omitempty"`                // Unique identifier for the view
	Tag               string             `json:"tag,omitempty"`               // Version tag for optimistic concurrency control
}

// TagRequest represents a request to set tags
type TagRequest struct {
	Tags []string `json:"tags"`
	// UUID of the current set of tags. Dremio uses the version value to ensure that you are updating the most recent version of the tags.
	// Required for updates and deletes, omit for initial creation.
	Version string `json:"version,omitempty"`
}

// WikiRequest represents a request to set wiki content
type WikiRequest struct {
	ID string `json:"id,omitempty"` // UUID of the source, folder, or dataset whose wiki you want to delete.

	// Text to display in the wiki. Use GitHub-flavored Markdown for wiki formatting and \n for new lines and blank lines.
	// Each wiki may have a maximum of 100,000 characters. Send an empty string to delete the wiki.
	Text string `json:"text"`

	// Number specified as the version value for the most recent existing wiki.
	// Dremio uses the version value to ensure that you are deleting the most recent version of the wiki.
	Version *int `json:"version,omitempty"`
}

// GrantsRequest represents a request to set grants
type GrantsRequest struct {
	ID     string           `json:"id,omitempty"` // UUID of the Dremio catalog object.
	Grants []GranteeRequest `json:"grants"`
	Tag    string           `json:"tag,omitempty"` // For Arctic catalog sources only
}

// GranteeRequest represents a grantee in a grants request
type GranteeRequest struct {
	Privileges  []string `json:"privileges"`
	GranteeType string   `json:"granteeType"` // USER or ROLE
	ID          string   `json:"id"`
}

// CloudCreateRequest represents a request to create a cloud
type CloudCreateRequest struct {
	Name       string           `json:"name"`                // User-defined name for the cloud
	Attributes *CloudAttributes `json:"attributes"`          // Cloud attributes (AWS or Azure specific)
	RequestID  string           `json:"requestId,omitempty"` // User-defined idempotency key for retries
}

// CloudUpdateRequest represents a request to update a cloud
type CloudUpdateRequest struct {
	Name       string           `json:"name"`       // User-defined name for the cloud
	Attributes *CloudAttributes `json:"attributes"` // Cloud attributes (AWS or Azure specific)
}

// CloudAttributes represents cloud-specific attributes
// This is a placeholder - actual structure varies by cloud provider (AWS/Azure)
type CloudAttributes map[string]interface{}

// EngineRulesRequest represents a request to update engine routing rules
type EngineRulesRequest struct {
	RuleSet *RuleSet `json:"ruleSet"` // The rule set containing all routing rules
}

// RuleSet represents a set of engine routing rules
type RuleSet struct {
	RuleInfos       []*RuleInfo `json:"ruleInfos"`       // List of all the rules in the project
	RuleInfoDefault *RuleInfo   `json:"ruleInfoDefault"` // The default rule (cannot be deleted)
	Tag             string      `json:"tag"`             // UUID of a tag that routes JDBC queries to a particular session. When the JDBC connection property ROUTING_TAG is set, the specified tag value is associated with all queries executed within that connection's session.
}

// RuleInfo represents a single engine routing rule
type RuleInfo struct {
	Name          string `json:"name"`          // User-defined name for the rule
	Condition     string `json:"condition"`     // The routing condition for the rule. You can use SQL syntax to create this condition. see Workload Management https://docs.dremio.com/dremio-cloud/admin/engines/workload-management/ for more information.
	EngineName    string `json:"engineName"`    // The name of the engine to which jobs will be routed. When action is REJECT, leave this parameter empty.
	Action        string `json:"action"`        // The rule type. When a query is routed to a particular engine, the value is ROUTE. When a query is rejected, the value is REJECT.
	RejectMessage string `json:"rejectMessage"` // The message displayed to the user if the rule rejects jobs.
}

// CloudTag represents a cloud tag (AWS or Azure)
type CloudTag struct {
	Key   string `json:"key"`   // The key identifier for the tag
	Value string `json:"value"` // The value of the tag
}

// EngineRequest represents a request to create/update an engine
type EngineRequest struct {
	Name                  string `json:"name,omitempty"`                  // User-defined name for the engine
	Size                  string `json:"size,omitempty"`                  // Size of the engine (XX_SMALL_V1, X_SMALL_V1, etc.). Mandatory for create, optional for update
	MinReplicas           int    `json:"minReplicas,omitempty"`           // Minimum number of engine replicas. Mandatory for create, optional for update
	MaxReplicas           int    `json:"maxReplicas,omitempty"`           // Maximum number of engine replicas. Mandatory for create, optional for update
	AutoStopDelaySeconds  int    `json:"autoStopDelaySeconds,omitempty"`  // Time that auto stop is delayed. Mandatory for create, optional for update
	QueueTimeLimitSeconds int    `json:"queueTimeLimitSeconds,omitempty"` // Max time a query will wait in queue. Mandatory for create, optional for update
	RuntimeLimitSeconds   int    `json:"runtimeLimitSeconds,omitempty"`   // Max time a query can run. Mandatory for create, optional for update
	DrainTimeLimitSeconds int    `json:"drainTimeLimitSeconds,omitempty"` // Max time replica continues after resize/disable/delete. Mandatory for create, optional for update
	MaxConcurrency        int    `json:"maxConcurrency,omitempty"`        // Max concurrent queries per replica. Mandatory for create, optional for update
	Description           string `json:"description"`                     // Description for the engine. Required (use empty string if not provided)
	RequestID             string `json:"requestId,omitempty"`             // User-defined idempotency key for retries. Mandatory for create, not used for update
}

// ExternalTokenProviderCreateRequest represents a request to create an external token provider
type ExternalTokenProviderCreateRequest struct {
	Name      string   `json:"name"`              // Name for the external token provider
	Audience  []string `json:"audience"`          // Intended recipients of the JWT
	UserClaim string   `json:"userClaim"`         // Key name for the target claim in the JWT
	IssuerURL string   `json:"issuerUrl"`         // URL that identifies the principal that issued the JWT
	JwksURL   string   `json:"jwksUrl,omitempty"` // Endpoint that hosts the JWK Set
	Enabled   bool     `json:"enabled,omitempty"` // If the provider is available
}

// ExternalTokenProviderUpdateRequest represents a request to update an external token provider
type ExternalTokenProviderUpdateRequest struct {
	Name      string   `json:"name"`              // Name for the external token provider
	Audience  []string `json:"audience"`          // Intended recipients of the JWT
	UserClaim string   `json:"userClaim"`         // Key name for the target claim in the JWT
	IssuerURL string   `json:"issuerUrl"`         // URL that identifies the principal that issued the JWT
	JwksURL   string   `json:"jwksUrl,omitempty"` // Endpoint that hosts the JWK Set
	Enabled   bool     `json:"enabled,omitempty"` // If the provider is available
}

// IdentityProviderCreateRequest represents a request to create an identity provider
type IdentityProviderCreateRequest struct {
	Type         string `json:"type"`                // Type of identity provider (GENERIC_OIDC, AZURE_AD, OKTA)
	IsActive     bool   `json:"isActive,omitempty"`  // Enable the provider as a login option
	IssuerURL    string `json:"issuerUrl,omitempty"` // Issuer URL for generic OIDC
	Domain       string `json:"domain,omitempty"`    // Publisher domain for Microsoft Entra ID
	OktaURL      string `json:"oktaUrl,omitempty"`   // URL for Okta
	ClientID     string `json:"clientID"`            // Client or application ID
	ClientSecret string `json:"clientSecret"`        // Client secret
}

// PipeLoadFilesRequest represents a request to load files with a pipe
type PipeLoadFilesRequest struct {
	Files []PipeFile `json:"files"` // Paths and sizes of files to load
}

// PipeFile represents a file to load via pipe
type PipeFile struct {
	Path string `json:"path"` // Path to the file
	Size string `json:"size"` // Estimated size (e.g., "80 MB")
}

// ProjectCreateRequest represents a request to create a project
type ProjectCreateRequest struct {
	Name         string              `json:"name"`                  // User-defined name for the project
	RequestID    string              `json:"requestId,omitempty"`   // User-defined idempotency key
	CloudID      string              `json:"cloudId"`               // ID of the cloud where compute resources will be created
	ProjectStore string              `json:"projectStore"`          // S3 bucket or Azure storage container
	Credentials  *ProjectCredentials `json:"credentials"`           // Storage credentials
	Type         string              `json:"type,omitempty"`        // Type of the project
	CatalogName  string              `json:"catalogName,omitempty"` // Name for the project's primary Arctic catalog (AWS only)
}

// ProjectUpdateRequest represents a request to update a project
type ProjectUpdateRequest struct {
	Name        string              `json:"name,omitempty"`        // User-defined name for the project
	Credentials *ProjectCredentials `json:"credentials,omitempty"` // Updated storage credentials
}

// ProjectCredentials represents project storage credentials
type ProjectCredentials struct {
	Type                string `json:"type"`                          // ACCESS_KEY, IAM_ROLE, or AZURE_STORAGE_CLIENT_CREDENTIALS
	AccessKeyID         string `json:"accessKeyId,omitempty"`         // AWS access key (for ACCESS_KEY)
	SecretAccessKey     string `json:"secretAccessKey,omitempty"`     // AWS secret key (for ACCESS_KEY)
	RoleArn             string `json:"roleArn,omitempty"`             // AWS cross-account role (for IAM_ROLE)
	InstanceProfileArn  string `json:"instanceProfileArn,omitempty"`  // AWS instance profile (for IAM_ROLE)
	ExternalID          string `json:"externalId,omitempty"`          // AWS external ID (for IAM_ROLE)
	ExternalIDSignature string `json:"externalIdSignature,omitempty"` // AWS external ID signature (for IAM_ROLE)
	TenantID            string `json:"tenantId,omitempty"`            // Azure tenant ID (for AZURE_STORAGE_CLIENT_CREDENTIALS)
	ClientID            string `json:"clientId,omitempty"`            // Azure client ID (for AZURE_STORAGE_CLIENT_CREDENTIALS)
	ClientSecret        string `json:"clientSecret,omitempty"`        // Azure client secret (for AZURE_STORAGE_CLIENT_CREDENTIALS)
	AccountName         string `json:"accountName,omitempty"`         // Azure storage account name (for AZURE_STORAGE_CLIENT_CREDENTIALS)
}

// ReflectionCreateRequest represents a request to create a Reflection
type ReflectionCreateRequest struct {
	Type                          string                   `json:"type"`                                    // RAW or AGGREGATION
	Name                          string                   `json:"name"`                                    // Name of the Reflection
	DatasetID                     string                   `json:"datasetId"`                               // ID of the dataset
	Enabled                       bool                     `json:"enabled"`                                 // Whether the Reflection is enabled
	EntityType                    string                   `json:"entityType"`                              // Always "reflection"
	DisplayFields                 []map[string]interface{} `json:"displayFields,omitempty"`                 // Fields to display (for RAW)
	DimensionFields               []map[string]interface{} `json:"dimensionFields,omitempty"`               // Dimension fields (for AGGREGATION)
	MeasureFields                 []map[string]interface{} `json:"measureFields,omitempty"`                 // Measure fields (for AGGREGATION)
	DistributionFields            []map[string]interface{} `json:"distributionFields,omitempty"`            // Fields for data distribution
	PartitionFields               []map[string]interface{} `json:"partitionFields,omitempty"`               // Fields for partitioning
	SortFields                    []map[string]interface{} `json:"sortFields,omitempty"`                    // Fields for sorting
	PartitionDistributionStrategy string                   `json:"partitionDistributionStrategy,omitempty"` // CONSOLIDATED or STRIPED
	CanView                       bool                     `json:"canView,omitempty"`                       // Whether user can view
	CanAlter                      bool                     `json:"canAlter,omitempty"`                      // Whether user can alter
}

// ReflectionUpdateRequest represents a request to update a Reflection
type ReflectionUpdateRequest struct {
	ID                            string                   `json:"id"`                                      // Unique identifier
	Type                          string                   `json:"type"`                                    // RAW or AGGREGATION
	Name                          string                   `json:"name"`                                    // Name of the Reflection
	Tag                           string                   `json:"tag"`                                     // Version tag for concurrency control
	DatasetID                     string                   `json:"datasetId"`                               // ID of the dataset
	Enabled                       bool                     `json:"enabled"`                                 // Whether the Reflection is enabled
	EntityType                    string                   `json:"entityType"`                              // Always "reflection"
	DisplayFields                 []map[string]interface{} `json:"displayFields,omitempty"`                 // Fields to display (for RAW)
	DimensionFields               []map[string]interface{} `json:"dimensionFields,omitempty"`               // Dimension fields (for AGGREGATION)
	MeasureFields                 []map[string]interface{} `json:"measureFields,omitempty"`                 // Measure fields (for AGGREGATION)
	DistributionFields            []map[string]interface{} `json:"distributionFields,omitempty"`            // Fields for data distribution
	PartitionFields               []map[string]interface{} `json:"partitionFields,omitempty"`               // Fields for partitioning
	SortFields                    []map[string]interface{} `json:"sortFields,omitempty"`                    // Fields for sorting
	PartitionDistributionStrategy string                   `json:"partitionDistributionStrategy,omitempty"` // CONSOLIDATED or STRIPED
	CanView                       bool                     `json:"canView,omitempty"`                       // Whether user can view
	CanAlter                      bool                     `json:"canAlter,omitempty"`                      // Whether user can alter
}

// JobBasedRecommendationsRequest represents a request to submit job IDs for Reflection recommendations
type JobBasedRecommendationsRequest struct {
	JobIDs []string `json:"jobIds"` // Job IDs of queries for which to request recommendations
}

// UsageBasedReflectionCreateRequest represents a request to create a Reflection from usage-based recommendation
type UsageBasedReflectionCreateRequest struct {
	ReflectionRequestBody map[string]interface{} `json:"reflectionRequestBody"` // Reflection body from recommendation
	RecommendationID      string                 `json:"recommendationId"`      // ID of the usage-based recommendation
}

// ScriptCreateRequest represents a request to create a script
type ScriptCreateRequest struct {
	Name    string   `json:"name"`              // Name for the script
	Content string   `json:"content"`           // SQL for the script
	Context []string `json:"context,omitempty"` // Path where the SQL query should run
	Owner   string   `json:"owner,omitempty"`   // User ID who should own the script
}

// ScriptUpdateRequest represents a request to update a script
type ScriptUpdateRequest struct {
	Name    string   `json:"name,omitempty"`    // Updated name for the script
	Content string   `json:"content,omitempty"` // Updated SQL for the script
	Context []string `json:"context,omitempty"` // Updated path where the SQL query should run
	Owner   string   `json:"owner,omitempty"`   // Updated owner user ID
}

// ScriptBatchDeleteRequest represents a request to batch delete scripts
type ScriptBatchDeleteRequest struct {
	IDs []string `json:"ids"` // Array of script IDs to delete
}

// ScriptGrantsUpdateRequest represents a request to update script grants
type ScriptGrantsUpdateRequest struct {
	Users []ScriptGrantee `json:"users,omitempty"` // Array of user privilege grants
	Roles []ScriptGrantee `json:"roles,omitempty"` // Array of role privilege grants
}

// ScriptGrantee represents a user or role with privileges on a script
type ScriptGrantee struct {
	GranteeID  string   `json:"granteeId"`  // User or role ID
	Privileges []string `json:"privileges"` // Array of privileges (VIEW, MODIFY, DELETE, MANAGE_GRANTS)
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query      string `json:"query"`                // Search string
	Filter     string `json:"filter,omitempty"`     // Optional CEL filter expression
	PageToken  string `json:"pageToken,omitempty"`  // Token to retrieve next page
	MaxResults int    `json:"maxResults,omitempty"` // Maximum number of results per page
}

// TokenCreateRequest represents a request to create a personal access token
type TokenCreateRequest struct {
	Label                string `json:"label"`                // User-defined description for the token
	MillisecondsToExpire int64  `json:"millisecondsToExpire"` // Lifespan of the token in milliseconds
}

// SQLRequest represents a SQL query request
type SQLRequest struct {
	SQL        string                 `json:"sql"`                  // SQL query to run
	Context    []string               `json:"context,omitempty"`    // Path to the container where the query should run
	References map[string]interface{} `json:"references,omitempty"` // Additional references
}

// FolderChild represents a child object within a folder
type FolderChild struct {
	ID            string   `json:"id"`
	Path          []string `json:"path"`
	Tag           string   `json:"tag"`
	Type          string   `json:"type"`                    // CONTAINER or DATASET
	ContainerType string   `json:"containerType,omitempty"` // FOLDER (if type is CONTAINER)
	DatasetType   string   `json:"datasetType,omitempty"`   // VIRTUAL or PROMOTED (if type is DATASET)
	CreatedAt     string   `json:"createdAt,omitempty"`
}

// MetadataPolicy represents metadata refresh policy settings
type MetadataPolicy struct {
	AuthTTLMs                 *int64  `json:"authTTLMs,omitempty"`                 // Length of time, in milliseconds, that source permissions are cached
	NamesRefreshMs            *int64  `json:"namesRefreshMs,omitempty"`            // When to run a refresh of a source, in milliseconds
	DatasetRefreshAfterMs     *int64  `json:"datasetRefreshAfterMs,omitempty"`     // How often the metadata in the dataset is refreshed, in milliseconds
	DatasetExpireAfterMs      *int64  `json:"datasetExpireAfterMs,omitempty"`      // Amount of time, in milliseconds, to keep the metadata before it expires
	DatasetUpdateMode         *string `json:"datasetUpdateMode,omitempty"`         // Metadata policy for when a dataset is updated (e.g., PREFETCH_QUERIED)
	DeleteUnavailableDatasets *bool   `json:"deleteUnavailableDatasets,omitempty"` // Remove dataset definitions if underlying data is unavailable
	AutoPromoteDatasets       *bool   `json:"autoPromoteDatasets,omitempty"`       // Automatically format files into tables when a query is issued
}

// AccessControlList represents access control settings
type AccessControlList struct {
	Users []UserPermission `json:"users,omitempty"`
	Roles []RolePermission `json:"roles,omitempty"`
}

// UserPermission represents a user's permissions
type UserPermission struct {
	ID          string   `json:"id"`
	Permissions []string `json:"permissions"`
}

// RolePermission represents a role's permissions
type RolePermission struct {
	ID          string   `json:"id"`
	Permissions []string `json:"permissions"`
}

// Owner represents the owner of a catalog object
type Owner struct {
	OwnerID   string `json:"ownerId"`
	OwnerType string `json:"ownerType"` // USER or ROLE
}

// ========================================
// User Management
// ========================================

// UserCreateRequest represents a request to create a user
// Reference: https://docs.dremio.com/cloud/reference/api/user/
type UserCreateRequest struct {
	UserName  string `json:"userName"`            // Username for the user
	FirstName string `json:"firstName,omitempty"` // First name of the user
	LastName  string `json:"lastName,omitempty"`  // Last name of the user
	Email     string `json:"email"`               // Email address of the user
	Tag       string `json:"tag,omitempty"`       // Version tag for optimistic concurrency control
}

// UserUpdateRequest represents a request to update a user
// Reference: https://docs.dremio.com/cloud/reference/api/user/
type UserUpdateRequest struct {
	UserName  string `json:"userName,omitempty"`  // Updated username
	FirstName string `json:"firstName,omitempty"` // Updated first name
	LastName  string `json:"lastName,omitempty"`  // Updated last name
	Email     string `json:"email,omitempty"`     // Updated email address
	Tag       string `json:"tag,omitempty"`       // Version tag for optimistic concurrency control
}

// ========================================
// Role Management
// ========================================

// RoleCreateRequest represents a request to create a role
// Reference: https://docs.dremio.com/cloud/reference/api/role/
type RoleCreateRequest struct {
	Name        string   `json:"name"`                  // Name of the role
	Description string   `json:"description,omitempty"` // Description of the role
	Members     []string `json:"members,omitempty"`     // Array of user IDs who are members of the role
}

// RoleUpdateRequest represents a request to update a role
// Reference: https://docs.dremio.com/cloud/reference/api/role/
type RoleUpdateRequest struct {
	Name        string   `json:"name,omitempty"`        // Updated name of the role
	Description string   `json:"description,omitempty"` // Updated description
	Members     []string `json:"members,omitempty"`     // Updated array of user IDs
	Tag         string   `json:"tag,omitempty"`         // Version tag for optimistic concurrency control
}

// RoleMembersUpdateRequest represents a request to update role members
type RoleMembersUpdateRequest struct {
	Add    []string `json:"add,omitempty"`    // User IDs to add to the role
	Remove []string `json:"remove,omitempty"` // User IDs to remove from the role
}

// ========================================
// Data Maintenance
// ========================================

// MaintenanceTaskCreateRequest represents a request to create a maintenance task
// Reference: https://docs.dremio.com/cloud/reference/api/catalog/maintenance/
type MaintenanceTaskRequest struct {
	TaskType   string                 `json:"type"`                // Type of maintenance task (OPTIMIZE, EXPIRE_SNAPSHOTS)
	IsEnabled  bool                   `json:"isEnabled,omitempty"` // Whether the task is enabled
	TaskConfig *MaintenanceTaskConfig `json:"config,omitempty"`    // An object that contains a fully qualified object name in the indicated catalog as the target for the maintenance task.
}
type MaintenanceTaskConfig struct {
	TableID string `json:"tableId,omitempty"` // Unique identifier of the table
}

// ========================================
// Billing (Cloud-only)
// ========================================

// BillingAccountUpdateRequest represents a request to update a billing account
// Reference: https://docs.dremio.com/cloud/reference/api/billing/
type BillingAccountUpdateRequest struct {
	Name        string                 `json:"name,omitempty"`        // Updated name
	Description string                 `json:"description,omitempty"` // Updated description
	Metadata    map[string]interface{} `json:"metadata,omitempty"`    // Additional metadata
}

// ========================================
// Arctic (Cloud-only)
// ========================================

// ArcticCatalogCreateRequest represents a request to create an Arctic catalog
// Reference: https://docs.dremio.com/cloud/reference/api/arctic/catalogs/
type ArcticCatalogCreateRequest struct {
	Name      string `json:"name"`                // User-defined name for the catalog
	RequestID string `json:"requestId,omitempty"` // User-defined idempotency key (UUID)
}

// ArcticScheduleCreateRequest represents a request to create an Arctic catalog schedule
// Reference: https://docs.dremio.com/cloud/reference/api/arctic/schedules/
type ArcticScheduleCreateRequest struct {
	Name           string                 `json:"name"`                 // Name of the schedule
	CronExpression string                 `json:"cronExpression"`       // Cron expression for the schedule
	Enabled        bool                   `json:"enabled,omitempty"`    // Whether the schedule is enabled
	TaskType       string                 `json:"taskType"`             // Type of task (OPTIMIZE, VACUUM, etc.)
	Parameters     map[string]interface{} `json:"parameters,omitempty"` // Task-specific parameters
	RequestID      string                 `json:"requestId,omitempty"`  // User-defined idempotency key
}

// ArcticScheduleUpdateRequest represents a request to update an Arctic catalog schedule
type ArcticScheduleUpdateRequest struct {
	Name           string                 `json:"name,omitempty"`           // Updated name
	CronExpression string                 `json:"cronExpression,omitempty"` // Updated cron expression
	Enabled        *bool                  `json:"enabled,omitempty"`        // Updated enabled status
	Parameters     map[string]interface{} `json:"parameters,omitempty"`     // Updated parameters
}

type TableFormatRequest struct {
	Type                    string  `json:"type"`                              // Type of data in the table (Delta, Excel, Iceberg, JSON, Parquet, Text, Unknown, XLS)
	IgnoreOtherFileFormats  *bool   `json:"ignoreOtherFileFormats,omitempty"`  // For Parquet folders, ignore non-Parquet files
	SkipFirstLine           *bool   `json:"skipFirstLine,omitempty"`           // Skip first line when creating table (Excel/Text)
	ExtractHeader           *bool   `json:"extractHeader,omitempty"`           // Extract column names from first line (Excel/Text)
	HasMergedCells          *bool   `json:"hasMergedCells,omitempty"`          // Expand merged cells (Excel)
	SheetName               *string `json:"sheetName,omitempty"`               // Sheet name for Excel files with multiple sheets
	FieldDelimiter          *string `json:"fieldDelimiter,omitempty"`          // Field delimiter character (Text), default: ","
	Quote                   *string `json:"quote,omitempty"`                   // Quote character (Text), default: "\""
	Comment                 *string `json:"comment,omitempty"`                 // Comment character (Text), default: "#"
	Escape                  *string `json:"escape,omitempty"`                  // Escape character (Text), default: "\""
	LineDelimiter           *string `json:"lineDelimiter,omitempty"`           // Line delimiter (Text), default: "\r\n"
	AutoGenerateColumnNames *bool   `json:"autoGenerateColumnNames,omitempty"` // Use existing column names (Text)
	TrimHeader              *bool   `json:"trimHeader,omitempty"`              // Trim column names (Text)
}
