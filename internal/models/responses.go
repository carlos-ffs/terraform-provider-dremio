package models

// SourceResponse represents a response for a source entity
// Reference: OpenAPI schema SourceResponse
type SourceResponse struct {
	ID                           *string                 `json:"id"`                                     // Unique identifier of the source
	Tag                          *string                 `json:"tag"`                                    // Version tag for optimistic concurrency control
	Type                         *string                 `json:"type"`                                   // Source type (ARCTIC, S3, SNOWFLAKE, etc.)
	Name                         *string                 `json:"name"`                                   // User-defined name of the source
	Config                       *map[string]interface{} `json:"config"`                                 // Configuration options specific to the source type
	MetadataPolicy               *MetadataPolicy         `json:"metadataPolicy,omitempty"`               // Metadata refresh policy
	AccelerationGracePeriodMs    *int64                  `json:"accelerationGracePeriodMs,omitempty"`    // Grace period before using Reflections
	AccelerationRefreshPeriodMs  *int64                  `json:"accelerationRefreshPeriodMs,omitempty"`  // Refresh period for Reflections
	AccelerationNeverExpire      *bool                   `json:"accelerationNeverExpire,omitempty"`      // Whether Reflections never expire
	AccelerationNeverRefresh     *bool                   `json:"accelerationNeverRefresh,omitempty"`     // Whether Reflections never refresh
	AccelerationActivePolicyType *string                 `json:"accelerationActivePolicyType,omitempty"` // Active policy type (PERIOD or NEVER)
	AccelerationRefreshSchedule  *string                 `json:"accelerationRefreshSchedule,omitempty"`  // Cron expression for refresh schedule
	Children                     *[]CatalogEntity        `json:"children,omitempty"`                     // Child entities
	AccessControlList            *AccessControlList      `json:"accessControlList,omitempty"`            // User and role access settings
	Permissions                  []string                `json:"permissions,omitempty"`                  // User's permissions on the source (as array of permission strings)
	Owner                        *Owner                  `json:"owner,omitempty"`                        // Owner information
}

// SpaceResponse represents a response for a space entity
// Reference: OpenAPI schema SpaceResponse
type SpaceResponse struct {
	EntityType        string             `json:"entityType,omitempty"`        // Type of catalog object (always "space")
	ID                string             `json:"id"`                          // Unique identifier of the space
	Name              string             `json:"name"`                        // Name of the space
	Tag               string             `json:"tag"`                         // Version tag for optimistic concurrency control
	CreatedAt         string             `json:"createdAt,omitempty"`         // Date and time the space was created (UTC)
	Children          []SpaceChild       `json:"children,omitempty"`          // Child entities in the space
	AccessControlList *AccessControlList `json:"accessControlList,omitempty"` // User and role access settings
	Permissions       []string           `json:"permissions,omitempty"`       // List of privileges on the space (only with include=permissions)
	Owner             *Owner             `json:"owner,omitempty"`             // Owner information
}

// SpaceChild represents a child entity in a space
type SpaceChild struct {
	ID            string   `json:"id"`                      // Unique identifier of the catalog object
	Path          []string `json:"path,omitempty"`          // Path of the catalog object within Dremio
	Tag           string   `json:"tag,omitempty"`           // Version tag
	Type          string   `json:"type,omitempty"`          // Type: CONTAINER, DATASET, FILE
	ContainerType string   `json:"containerType,omitempty"` // For CONTAINER: FOLDER, FUNCTION
	DatasetType   string   `json:"datasetType,omitempty"`   // For DATASET in space: always VIRTUAL
	CreatedAt     string   `json:"createdAt,omitempty"`     // Date and time the catalog object was created
}

// FolderResponse represents a response for a folder entity
// Reference: OpenAPI schema FolderResponse
type FolderResponse struct {
	ID                string             `json:"id"`                          // Unique identifier of the folder
	Path              []string           `json:"path"`                        // Full path to the folder
	Tag               string             `json:"tag"`                         // Version tag for optimistic concurrency control
	EntityType        string             `json:"entityType,omitempty"`        // Type of catalog object (always "folder")
	Children          []FolderChild      `json:"children,omitempty"`          // Child entities
	AccessControlList *AccessControlList `json:"accessControlList,omitempty"` // User and role access settings
	Permissions       []string           `json:"permissions,omitempty"`       // User's permissions on the folder
	Owner             *Owner             `json:"owner,omitempty"`             // Owner information
	NextPageToken     string             `json:"nextPageToken,omitempty"`     // Token for pagination
	StorageURI        string             `json:"storageUri,omitempty"`        // Storage URI for Open Catalog folders
}

// TableResponse represents a response for a table/dataset entity
// Reference: OpenAPI schema TableResponse
type TableResponse struct {
	ID                           string                     `json:"id"`                                     // Unique identifier of the table
	Type                         string                     `json:"type"`                                   // Dataset type (PHYSICAL_DATASET)
	Path                         []string                   `json:"path"`                                   // Full path to the table
	CreatedAt                    string                     `json:"createdAt,omitempty"`                    // Date and time the table was created (UTC)
	Tag                          string                     `json:"tag"`                                    // Version tag for optimistic concurrency control
	AccelerationRefreshPolicy    *AccelerationRefreshPolicy `json:"accelerationRefreshPolicy,omitempty"`    // Acceleration refresh policy
	Format                       *TableFormatResponse       `json:"format,omitempty"`                       // Table format information
	AccessControlList            *AccessControlList         `json:"accessControlList,omitempty"`            // User and role access settings
	Owner                        *Owner                     `json:"owner,omitempty"`                        // Owner information
	Fields                       []TableField               `json:"fields,omitempty"`                       // Table fields/columns
	ApproximateStatisticsAllowed bool                       `json:"approximateStatisticsAllowed,omitempty"` // Whether approximate statistics are allowed
}

// ViewResponse represents a response for a view/virtual dataset entity
// Reference: OpenAPI schema ViewResponse
type ViewResponse struct {
	ID                    string             `json:"id"`                              // Unique identifier of the view
	Type                  string             `json:"type"`                            // Dataset type (VIRTUAL_DATASET)
	Path                  []string           `json:"path"`                            // Full path to the view
	CreatedAt             string             `json:"createdAt,omitempty"`             // Date and time the view was created (UTC)
	Tag                   string             `json:"tag"`                             // Version tag for optimistic concurrency control
	SQL                   string             `json:"sql"`                             // SQL query defining the view
	SQLContext            []string           `json:"sqlContext,omitempty"`            // Context for SQL query execution
	Fields                []TableField       `json:"fields,omitempty"`                // View fields/columns
	IsMetadataExpired     bool               `json:"isMetadataExpired,omitempty"`     // Whether metadata is expired
	LastMetadataRefreshAt string             `json:"lastMetadataRefreshAt,omitempty"` // Last metadata refresh timestamp
	AccessControlList     *AccessControlList `json:"accessControlList,omitempty"`     // User and role access settings
	Permissions           []string           `json:"permissions,omitempty"`           // User's permissions on the view
	Owner                 *Owner             `json:"owner,omitempty"`                 // Owner information
}

// UDFResponse represents a response for a user-defined function entity
// Reference: OpenAPI schema UDFResponse
type UDFResponse struct {
	ID                *string            `json:"id"`                          // Unique identifier of the UDF
	Path              []string           `json:"path"`                        // Full path to the UDF
	Tag               *string            `json:"tag"`                         // Version tag for optimistic concurrency control
	CreatedAt         *string            `json:"createdAt,omitempty"`         // Date and time the UDF was created (UTC)
	LastModified      *string            `json:"lastModified,omitempty"`      // Date and time the UDF was last modified (UTC)
	IsScalar          *bool              `json:"isScalar"`                    // Whether the function is scalar
	FunctionArgList   *string            `json:"functionArgList,omitempty"`   // Function arguments as a string
	FunctionBody      *string            `json:"functionBody"`                // SQL body of the function
	ReturnType        *string            `json:"returnType,omitempty"`        // Return type as a string
	AccessControlList *AccessControlList `json:"accessControlList,omitempty"` // User and role access settings
	Permissions       []string           `json:"permissions,omitempty"`       // User's permissions on the UDF (as array of permission strings)
	Owner             *Owner             `json:"owner,omitempty"`             // Owner information
}

// LineageResponse represents a response for lineage information
// Reference: OpenAPI schema LineageResponse
type LineageResponse struct {
	Sources  []LineageSource  `json:"sources,omitempty"`  // Source entities in the lineage
	Parents  []LineageDataset `json:"parents,omitempty"`  // Parent datasets in the lineage
	Children []LineageDataset `json:"children,omitempty"` // Child datasets in the lineage
}

// LineageSource represents a source in lineage information
// Reference: OpenAPI schema LineageSource
type LineageSource struct {
	ID   string   `json:"id"`             // Unique identifier of the source
	Path []string `json:"path,omitempty"` // Path to the source
	Type string   `json:"type"`           // Source type
}

// LineageDataset represents a dataset in lineage information
// Reference: OpenAPI schema LineageDataset
type LineageDataset struct {
	ID        string   `json:"id"`                  // Unique identifier of the dataset
	Path      []string `json:"path,omitempty"`      // Path to the dataset
	Type      string   `json:"type"`                // Dataset type (PROMOTED or VIRTUAL)
	CreatedAt string   `json:"createdAt,omitempty"` // Date and time the dataset was created (UTC)
}

// TagResponse represents a response for tags on a catalog object
// Reference: OpenAPI schema TagResponse
type TagResponse struct {
	Tags    []string `json:"tags"`    // List of tags applied to the dataset
	Version string   `json:"version"` // Unique identifier of the set of tags
}

// WikiResponse represents a response for wiki information on a catalog object
// Reference: OpenAPI schema WikiResponse
type WikiResponse struct {
	Text    string `json:"text"`    // Text displayed in the wiki, formatted with GitHub-flavored Markdown
	Version int    `json:"version"` // Number for the most recent version of the wiki, starting with 0
}

// GrantsResponse represents a response for grants/permissions on a catalog object
type GrantsResponse struct {
	ID                  string             `json:"id"`                            // UUID of the catalog object.
	Grants              []GranteesResponse `json:"grants,omitempty"`              // Information about privileges available for each type of catalog object
	AvailablePrivileges []string           `json:"availablePrivileges,omitempty"` // List of available privileges on the catalog object. See Privileges for more information.
}

type GranteesResponse struct {
	ID          string   `json:"id"`                  // UUID of the user or role.
	Name        string   `json:"name"`                // Name of the user or role.
	FirstName   string   `json:"firstName,omitempty"` // The user's first name (not included if the object is a role)
	LastName    string   `json:"lastName,omitempty"`  // The user's last name (not included if the object is a role)
	Email       string   `json:"email,omitempty"`     // The user's email address (not included if the object is a role)
	GranteeType string   `json:"granteeType"`         // Type of catalog object
	Privileges  []string `json:"privileges"`          // List of available privileges on this type of catalog object
}

// JobResponse represents a response for a job
// Reference: OpenAPI schema JobResponse
type JobResponse struct {
	JobState     string           `json:"jobState"`               // The job's status
	RowCount     int              `json:"rowCount,omitempty"`     // Number of rows the job returned
	ErrorMessage string           `json:"errorMessage,omitempty"` // Error message for failed jobs
	StartedAt    string           `json:"startedAt,omitempty"`    // Date and time when the job started (UTC)
	EndedAt      string           `json:"endedAt,omitempty"`      // Date and time when the job ended (UTC)
	Acceleration *JobAcceleration `json:"acceleration,omitempty"` // Reflection information for the job
	QueryType    string           `json:"queryType,omitempty"`    // Job type
	User         string           `json:"user,omitempty"`         // User who ran the job
	QueryText    string           `json:"queryText,omitempty"`    // SQL query text
}

// JobAcceleration represents acceleration/reflection information for a job
type JobAcceleration struct {
	ReflectionRelationships []ReflectionRelationship `json:"reflectionRelationships,omitempty"` // Information about Reflections and their relationships to the job
}

// ReflectionRelationship represents the relationship between a Reflection and a job
type ReflectionRelationship struct {
	DatasetID    string `json:"datasetId"`    // Unique identifier for the dataset associated with the Reflection
	ReflectionID string `json:"reflectionId"` // Unique identifier for the Reflection
	Relationship string `json:"relationship"` // Type of relationship (CONSIDERED, MATCHED, CHOSEN)
}

// JobResultsResponse represents a response for job results
// Reference: OpenAPI schema JobResultsResponse
type JobResultsResponse struct {
	RowCount int                      `json:"rowCount"`       // Number of rows the job returned
	Schema   []JobResultSchemaField   `json:"schema"`         // Array of schema definitions for the data
	Rows     []map[string]interface{} `json:"rows,omitempty"` // Array of the data the job returned for each row
}

// JobResultSchemaField represents a schema field in job results
type JobResultSchemaField struct {
	Name string              `json:"name"` // Column name
	Type JobResultSchemaType `json:"type"` // Column type information
}

// JobResultSchemaType represents the type information for a schema field
type JobResultSchemaType struct {
	Name string `json:"name"` // Data type name
}

type EngineResponse struct {
	ID                        string `json:"id"`                                  // Unique identifier of the engine. API Create response only returns id.
	Name                      string `json:"name,omitempty"`                      // User-defined name for the engine
	Size                      string `json:"size,omitempty"`                      // Size of the engine
	ActiveReplicas            int    `json:"activeReplicas,omitempty"`            // Number of engine replicas currently active
	MinReplicas               int    `json:"minReplicas,omitempty"`               // Minimum number of engine replicas
	MaxReplicas               int    `json:"maxReplicas,omitempty"`               // Maximum number of engine replicas
	AutoStopDelaySeconds      int    `json:"autoStopDelaySeconds,omitempty"`      // Time that auto stop is delayed
	QueueTimeLimitSeconds     int    `json:"queueTimeLimitSeconds,omitempty"`     // Max time a query will wait in queue
	RuntimeLimitSeconds       int    `json:"runtimeLimitSeconds,omitempty"`       // Max time a query can run
	DrainTimeLimitSeconds     int    `json:"drainTimeLimitSeconds,omitempty"`     // Max time replica continues after resize/disable/delete
	MaxConcurrency            int    `json:"maxConcurrency,omitempty"`            // Max concurrent queries per replica
	State                     string `json:"state,omitempty"`                     // The current state of the engine
	QueriedAt                 string `json:"queriedAt,omitempty"`                 // The date and time that the engine was last used to execute a query
	StatusChangedAt           string `json:"statusChangedAt,omitempty"`           // The date and time (in UTC time) that the state of the engine changed
	Description               string `json:"description,omitempty"`               // Description for the engine
	InstanceFamily            string `json:"instanceFamily,omitempty"`            // Instance family (M5D, M6ID, M6GD, DDV4, DDV5)
	AdditionalEngineStateInfo string `json:"additionalEngineStateInfo,omitempty"` // Not used. Has the value NONE.
}

// ReflectionResponse represents a response for a Reflection
// Reference: OpenAPI schema ReflectionResponse
type ReflectionResponse struct {
	ID                            string                        `json:"id"`                                      // Unique identifier of the Reflection
	Type                          string                        `json:"type"`                                    // Type of Reflection (RAW or AGGREGATION)
	Name                          string                        `json:"name"`                                    // Name of the Reflection
	Tag                           string                        `json:"tag"`                                     // Version tag for optimistic concurrency control
	CreatedAt                     string                        `json:"createdAt,omitempty"`                     // Timestamp when the Reflection was created
	UpdatedAt                     string                        `json:"updatedAt,omitempty"`                     // Timestamp when the Reflection was last updated
	DatasetID                     string                        `json:"datasetId"`                               // ID of the dataset this Reflection is based on
	CurrentSizeBytes              int64                         `json:"currentSizeBytes,omitempty"`              // Current size of the Reflection in bytes
	TotalSizeBytes                int64                         `json:"totalSizeBytes,omitempty"`                // Total size of the Reflection in bytes
	Enabled                       bool                          `json:"enabled"`                                 // Whether the Reflection is enabled
	Status                        *ReflectionStatus             `json:"status,omitempty"`                        // Status of the Reflection
	DisplayFields                 []ReflectionDisplayField      `json:"displayFields,omitempty"`                 // Fields to display (for RAW Reflections)
	DimensionFields               []ReflectionDimensionField    `json:"dimensionFields,omitempty"`               // Dimension fields (for AGGREGATION Reflections)
	MeasureFields                 []ReflectionMeasureField      `json:"measureFields,omitempty"`                 // Measure fields (for AGGREGATION Reflections)
	DistributionFields            []ReflectionDistributionField `json:"distributionFields,omitempty"`            // Fields for data distribution across nodes
	PartitionFields               []ReflectionPartitionField    `json:"partitionFields,omitempty"`               // Fields for horizontal partitioning
	SortFields                    []ReflectionSortField         `json:"sortFields,omitempty"`                    // Fields for sorting data
	PartitionDistributionStrategy string                        `json:"partitionDistributionStrategy,omitempty"` // Strategy for partition distribution
	CanView                       bool                          `json:"canView,omitempty"`                       // Whether the user can view this Reflection
	CanAlter                      bool                          `json:"canAlter,omitempty"`                      // Whether the user can alter this Reflection
	EntityType                    string                        `json:"entityType,omitempty"`                    // Entity type (always "reflection")
}

// ReflectionListResponse represents a response for a list of Reflections
// Reference: OpenAPI schema ReflectionListResponse
type ReflectionListResponse struct {
	Data                []ReflectionResponse `json:"data"`                          // Array of Reflection objects
	CanAlterReflections bool                 `json:"canAlterReflections,omitempty"` // Whether the user can alter Reflections
}

// ReflectionStatus represents the status of a Reflection
// Reference: OpenAPI schema ReflectionStatus
type ReflectionStatus struct {
	Config         string `json:"config,omitempty"`         // Configuration status (OK, INVALID)
	Refresh        string `json:"refresh,omitempty"`        // Refresh status (GIVEN_UP, MANUAL, RUNNING, SCHEDULED)
	Availability   string `json:"availability,omitempty"`   // Availability status (NONE, EXPIRED, AVAILABLE)
	CombinedStatus string `json:"combinedStatus,omitempty"` // Combined status of the Reflection
	FailureCount   int    `json:"failureCount,omitempty"`   // Number of failures
	LastDataFetch  string `json:"lastDataFetch,omitempty"`  // Timestamp of last data fetch
	ExpiresAt      string `json:"expiresAt,omitempty"`      // Timestamp when the Reflection expires
}

// ReflectionDisplayField represents a display field in a Reflection
// Reference: OpenAPI schema ReflectionDisplayField
type ReflectionDisplayField struct {
	Name string `json:"name"` // Name of the field to display
}

// ReflectionDimensionField represents a dimension field in a Reflection
// Reference: OpenAPI schema ReflectionDimensionField
type ReflectionDimensionField struct {
	Name        string `json:"name"`                  // Name of the dimension field
	Granularity string `json:"granularity,omitempty"` // Granularity for the dimension field (DATE, NORMAL)
}

// ReflectionMeasureField represents a measure field in a Reflection
// Reference: OpenAPI schema ReflectionMeasureField
type ReflectionMeasureField struct {
	Name            string   `json:"name"`                      // Name of the measure field
	MeasureTypeList []string `json:"measureTypeList,omitempty"` // List of measure types (SUM, COUNT, MIN, MAX, AVG, APPROX_COUNT_DISTINCT)
}

// ReflectionDistributionField represents a distribution field in a Reflection
// Reference: OpenAPI schema ReflectionDistributionField
type ReflectionDistributionField struct {
	Name string `json:"name"` // Name of the distribution field
}

// ReflectionPartitionField represents a partition field in a Reflection
// Reference: OpenAPI schema ReflectionPartitionField
type ReflectionPartitionField struct {
	Name string `json:"name"` // Name of the partition field
}

// ReflectionSortField represents a sort field in a Reflection
// Reference: OpenAPI schema ReflectionSortField
type ReflectionSortField struct {
	Name string `json:"name"` // Name of the sort field
}

// ReflectionRecommendationsResponse represents a response for Reflection recommendations
// Reference: OpenAPI schema ReflectionRecommendationsResponse
type ReflectionRecommendationsResponse struct {
	Data []ReflectionRecommendation `json:"data,omitempty"` // List of recommended Reflection objects
}

// ReflectionRecommendation represents a recommended Reflection
// Reference: OpenAPI schema ReflectionRecommendation
type ReflectionRecommendation struct {
	Type            string                              `json:"type"`                      // Reflection type (RAW or AGGREGATION)
	Enabled         bool                                `json:"enabled"`                   // If the Reflection is available for accelerating queries
	DisplayFields   []ReflectionRecommendationField     `json:"displayFields,omitempty"`   // Fields displayed (for raw Reflections)
	DimensionFields []ReflectionRecommendationDimension `json:"dimensionFields,omitempty"` // Dimension fields (for aggregation Reflections)
	MeasureFields   []ReflectionRecommendationMeasure   `json:"measureFields,omitempty"`   // Measure fields (for aggregation Reflections)
	PartitionFields []ReflectionRecommendationField     `json:"partitionFields,omitempty"` // Partition fields
	EntityType      string                              `json:"entityType"`                // Entity type
}

// ReflectionRecommendationField represents a field in a Reflection recommendation
type ReflectionRecommendationField struct {
	Name string `json:"name"` // Field name
}

// ReflectionRecommendationDimension represents a dimension field in a Reflection recommendation
type ReflectionRecommendationDimension struct {
	Name        string `json:"name"`                  // Field name
	Granularity string `json:"granularity,omitempty"` // Granularity for the dimension field
}

// ReflectionRecommendationMeasure represents a measure field in a Reflection recommendation
type ReflectionRecommendationMeasure struct {
	Name            string   `json:"name"`                      // Field name
	MeasureTypeList []string `json:"measureTypeList,omitempty"` // List of measure types
}

// ReflectionSummaryResponse represents a response for Reflection summaries
// Reference: OpenAPI schema ReflectionSummaryResponse
type ReflectionSummaryResponse struct {
	Data                  []ReflectionSummaryItem `json:"data,omitempty"`                  // List of Reflection summary objects
	NextPageToken         string                  `json:"nextPageToken,omitempty"`         // Token to retrieve the next page of results
	IsCanAlterReflections bool                    `json:"isCanAlterReflections,omitempty"` // Whether the current user has project-level privileges to alter Reflections
}

// ReflectionSummaryItem represents a summary item for a Reflection
// Reference: OpenAPI schema ReflectionSummaryItem
type ReflectionSummaryItem struct {
	CreatedAt          string                   `json:"createdAt,omitempty"`          // Date and time the Reflection was created (UTC)
	UpdatedAt          string                   `json:"updatedAt,omitempty"`          // Date and time the Reflection was last updated (UTC)
	ID                 string                   `json:"id"`                           // Unique identifier for the Reflection
	ReflectionType     string                   `json:"reflectionType"`               // Type of Reflection (RAW or AGGREGATION)
	ReflectionMode     string                   `json:"reflectionMode,omitempty"`     // How the Reflection was created (Manual or Autonomous)
	Name               string                   `json:"name"`                         // User-provided name for the Reflection
	CurrentSizeBytes   int64                    `json:"currentSizeBytes,omitempty"`   // Data size of the latest Reflection job in bytes
	OutputRecords      int64                    `json:"outputRecords,omitempty"`      // Number of records returned for the latest Reflection
	TotalSizeBytes     int64                    `json:"totalSizeBytes,omitempty"`     // Data size of all Reflection jobs that have not been pruned in bytes
	DatasetID          string                   `json:"datasetId"`                    // Unique identifier for the anchor dataset
	DatasetType        string                   `json:"datasetType"`                  // Type of anchor dataset (PHYSICAL_DATASET or VIRTUAL_DATASET)
	DatasetPath        []string                 `json:"datasetPath,omitempty"`        // Path to the anchor dataset
	Status             *ReflectionSummaryStatus `json:"status,omitempty"`             // Status of the Reflection
	ConsideredCount    int                      `json:"consideredCount,omitempty"`    // Number of jobs that considered the Reflection during planning
	MatchedCount       int                      `json:"matchedCount,omitempty"`       // Number of jobs that matched the Reflection during planning
	ChosenCount        int                      `json:"chosenCount,omitempty"`        // Number of jobs accelerated by the Reflection
	ConsideredJobsLink string                   `json:"consideredJobsLink,omitempty"` // Link to list of considered jobs for the Reflection
	MatchedJobsLink    string                   `json:"matchedJobsLink,omitempty"`    // Link to list of matched jobs for the Reflection
	ChosenJobsLink     string                   `json:"chosenJobsLink,omitempty"`     // Link to list of chosen jobs for the Reflection
	CanView            bool                     `json:"canView,omitempty"`            // Whether the user can view this Reflection
	CanAlter           bool                     `json:"canAlter,omitempty"`           // Whether the user can alter this Reflection
}

// ReflectionSummaryStatus represents the status of a Reflection in a summary
// Reference: OpenAPI schema ReflectionSummaryStatus
type ReflectionSummaryStatus struct {
	ConfigStatus              string `json:"configStatus,omitempty"`              // Status of the Reflection configuration (OK, INVALID)
	RefreshStatus             string `json:"refreshStatus,omitempty"`             // Status of the Reflection refresh
	AvailabilityStatus        string `json:"availabilityStatus,omitempty"`        // Status of the Reflection's availability for accelerating queries
	CombinedStatus            string `json:"combinedStatus,omitempty"`            // Combined status based on configStatus, refreshStatus, and availabilityStatus
	RefreshMethod             string `json:"refreshMethod,omitempty"`             // Method used for the most recent refresh (NONE, FULL, INCREMENTAL)
	FailureCount              int    `json:"failureCount,omitempty"`              // Number of times refresh attempts failed
	LastFailureMessage        string `json:"lastFailureMessage,omitempty"`        // Error message from the last failed refresh
	LastDataFetchAt           string `json:"lastDataFetchAt,omitempty"`           // Date and time the Reflection data was last refreshed (UTC)
	ExpiresAt                 string `json:"expiresAt,omitempty"`                 // Date and time the Reflection will expire (UTC)
	LastRefreshDurationMillis int64  `json:"lastRefreshDurationMillis,omitempty"` // Duration of the most recent refresh in milliseconds
}

// ScriptResponse represents a response for a script
// Reference: OpenAPI schema ScriptResponse
type ScriptResponse struct {
	ID         string   `json:"id"`                   // Unique identifier of the script
	Name       string   `json:"name"`                 // User-provided name of the script
	Content    string   `json:"content"`              // The script's SQL
	Context    []string `json:"context,omitempty"`    // Path where the SQL query runs
	Owner      string   `json:"owner,omitempty"`      // User ID who owns the script
	CreatedAt  string   `json:"createdAt,omitempty"`  // Date and time the script was created (UTC)
	CreatedBy  string   `json:"createdBy,omitempty"`  // User ID who created the script
	ModifiedAt string   `json:"modifiedAt,omitempty"` // Date and time the script was last modified (UTC)
	ModifiedBy string   `json:"modifiedBy,omitempty"` // User ID who last modified the script
}

// ScriptsListResponse represents a response for a list of scripts
// Reference: OpenAPI schema ScriptsListResponse
type ScriptsListResponse struct {
	Total int              `json:"total"` // Total number of scripts in the project
	Data  []ScriptResponse `json:"data"`  // List of scripts in the project
}

// ScriptBatchDeleteResponse represents a response for batch deleting scripts
// Reference: OpenAPI schema ScriptBatchDeleteResponse
type ScriptBatchDeleteResponse struct {
	UnauthorizedIDs []string `json:"unauthorizedIds,omitempty"` // IDs of scripts that could not be deleted due to lack of authorization
	NotFoundIDs     []string `json:"notFoundIds,omitempty"`     // IDs of scripts that were not found
	OtherErrorIDs   []string `json:"otherErrorIds,omitempty"`   // IDs of scripts that could not be deleted due to other errors
}

// ScriptGrantsResponse represents a response for script grants/permissions
// Reference: OpenAPI schema ScriptGrantsResponse
type ScriptGrantsResponse struct {
	Users []ScriptGrantee `json:"users,omitempty"` // Array of user privilege grants
	Roles []ScriptGrantee `json:"roles,omitempty"` // Array of role privilege grants
}

// SearchResponse represents a response for search results
// Reference: OpenAPI schema SearchResponse
type SearchResponse struct {
	SessionID     string               `json:"sessionId,omitempty"`     // Session identifier to correlate API calls during feedback collection
	NextPageToken string               `json:"nextPageToken,omitempty"` // Token of the next page of results to fetch in a paginated response
	Results       []SearchResultObject `json:"results,omitempty"`       // Array of search results
}

// SearchResultObject represents a single search result
// Reference: OpenAPI schema SearchResultObject
type SearchResultObject struct {
	Category         string                  `json:"category"`                   // The type of the result object
	CatalogObject    *SearchCatalogObject    `json:"catalogObject,omitempty"`    // If the result is a catalog object
	JobObject        *SearchJobObject        `json:"jobObject,omitempty"`        // If the result is a job
	ReflectionObject *SearchReflectionObject `json:"reflectionObject,omitempty"` // If the result is a reflection
	ScriptObject     *SearchScriptObject     `json:"scriptObject,omitempty"`     // If the result is a script
}

// SearchCatalogObject represents attributes for catalog objects in search results
// Reference: OpenAPI schema SearchCatalogObject
type SearchCatalogObject struct {
	Path        []string          `json:"path,omitempty"`        // Namespace path to the object
	Branch      string            `json:"branch,omitempty"`      // Versioned branch name
	Type        string            `json:"type,omitempty"`        // Type of catalog object
	Labels      []string          `json:"labels,omitempty"`      // User-defined labels
	Wiki        string            `json:"wiki,omitempty"`        // Markdown-formatted documentation or notes
	Owner       *SearchUserOrRole `json:"owner,omitempty"`       // User or role object
	CreatedAt   string            `json:"createdAt,omitempty"`   // Creation timestamp (RFC 3339)
	ModifiedAt  string            `json:"modifiedAt,omitempty"`  // Last modification timestamp (RFC 3339)
	Columns     []string          `json:"columns,omitempty"`     // Column names
	FunctionSQL string            `json:"functionSql,omitempty"` // SQL definition for functions
}

// SearchJobObject represents attributes for jobs in search results
// Reference: OpenAPI schema SearchJobObject
type SearchJobObject struct {
	ID              string                 `json:"id"`                        // Unique identifier for the job
	QueriedDatasets []SearchQueriedDataset `json:"queriedDatasets,omitempty"` // Datasets queried in the job
	SQL             string                 `json:"sql,omitempty"`             // Executed SQL statement
	JobType         string                 `json:"jobType,omitempty"`         // Type of job
	User            *SearchUserOrRole      `json:"user,omitempty"`            // User or role who ran the job
	StartTime       string                 `json:"startTime,omitempty"`       // Job start timestamp
	FinishTime      string                 `json:"finishTime,omitempty"`      // Job completion timestamp
	JobState        string                 `json:"jobState,omitempty"`        // Job status
	Error           string                 `json:"error,omitempty"`           // Error message if the job failed
}

// SearchQueriedDataset represents a dataset queried in a job
type SearchQueriedDataset struct {
	DatasetType string   `json:"datasetType,omitempty"` // Dataset type (TABLE or VIEW)
	DatasetPath []string `json:"datasetPath,omitempty"` // Path to the dataset
}

// SearchScriptObject represents attributes for scripts in search results
// Reference: OpenAPI schema SearchScriptObject
type SearchScriptObject struct {
	ID         string            `json:"id"`                   // Script identifier
	Name       string            `json:"name,omitempty"`       // Name of the script
	Owner      *SearchUserOrRole `json:"owner,omitempty"`      // User or role object
	Content    string            `json:"content,omitempty"`    // SQL content
	CreatedAt  string            `json:"createdAt,omitempty"`  // Creation timestamp
	ModifiedAt string            `json:"modifiedAt,omitempty"` // Last modified timestamp
}

// SearchReflectionObject represents attributes for Reflections in search results
// Reference: OpenAPI schema SearchReflectionObject
type SearchReflectionObject struct {
	ID            string   `json:"id"`                      // Reflection ID
	Name          string   `json:"name,omitempty"`          // Name of the Reflection
	DatasetType   string   `json:"datasetType,omitempty"`   // Type of dataset (TABLE or VIEW)
	DatasetPath   []string `json:"datasetPath,omitempty"`   // Path to the dataset
	DatasetBranch string   `json:"datasetBranch,omitempty"` // Dataset branch
	CreatedAt     string   `json:"createdAt,omitempty"`     // Creation timestamp
	ModifiedAt    string   `json:"modifiedAt,omitempty"`    // Last modified timestamp
}

// SearchUserOrRole represents a user or role object used in owner or user fields
// Reference: OpenAPI schema SearchUserOrRole
type SearchUserOrRole struct {
	ID       string `json:"id"`                 // Unique ID of the user or role
	Type     string `json:"type"`               // Type of entity (USER or ROLE)
	Username string `json:"username,omitempty"` // Present for USER type
	RoleName string `json:"roleName,omitempty"` // Present for ROLE type
}

// TokenListResponse represents a response for a list of Personal Access Tokens
// Reference: OpenAPI schema TokenListResponse
type TokenListResponse struct {
	Data []Token `json:"data,omitempty"` // Array of token objects
}

// Token represents a Personal Access Token
// Reference: OpenAPI schema Token
type Token struct {
	TID       string `json:"tid"`       // Unique identifier of the PAT
	UID       string `json:"uid"`       // Unique identifier of the user
	Label     string `json:"label"`     // User-provided name of the PAT
	CreatedAt string `json:"createdAt"` // Date and time that the PAT was created (UTC)
	ExpiresAt string `json:"expiresAt"` // Date and time that the PAT will expire (UTC)
}

// UsageResponse represents a response for usage information
// Reference: OpenAPI schema UsageResponse
type UsageResponse struct {
	Data              []UsageObject `json:"data,omitempty"`              // List of usage objects for the project or engine
	PreviousPageToken *string       `json:"previousPageToken,omitempty"` // Token for retrieving the previous page of usage objects
	NextPageToken     *string       `json:"nextPageToken,omitempty"`     // Token for retrieving the next page of usage objects
}

// UsageObject represents usage information for a project or engine
// Reference: OpenAPI schema UsageObject
type UsageObject struct {
	ID        string  `json:"id"`        // The ID for the project or engine
	Type      string  `json:"type"`      // Specifies whether the usage is reported for a project or for an engine (PROJECT or ENGINE)
	StartTime string  `json:"startTime"` // The starting date and time of the period for which the usage is reported
	EndTime   string  `json:"endTime"`   // The ending date and time of the period for which the usage is reported
	Usage     float64 `json:"usage"`     // The usage for the object in Dremio Consumption Units (DCUs)
}

// SQLResponse represents a response for a SQL query submission
// Reference: OpenAPI schema SQLResponse
type SQLResponse struct {
	ID string `json:"id"` // Job ID associated with the SQL query
}

// PipeLoadFilesResponse represents a response for pipe load files operation
// Reference: OpenAPI schema PipeLoadFilesResponse
type PipeLoadFilesResponse struct {
	RequestID string `json:"requestId"` // Request ID associated with the batch of files submitted
}

// Permissions represents user permissions on a catalog object
type Permissions struct {
	CanView                  bool `json:"canView,omitempty"`                  // Whether the user can view the object
	CanAlter                 bool `json:"canAlter,omitempty"`                 // Whether the user can alter the object
	CanDelete                bool `json:"canDelete,omitempty"`                // Whether the user can delete the object
	CanManageGrants          bool `json:"canManageGrants,omitempty"`          // Whether the user can manage grants on the object
	CanEditAccessControlList bool `json:"canEditAccessControlList,omitempty"` // Whether the user can edit access control list on the object
	CanCreateChildren        bool `json:"canCreateChildren,omitempty"`        // Whether the user can create children on the object
	CanRead                  bool `json:"canRead,omitempty"`                  // Whether the user can read the object
	CanEditFormatSettings    bool `json:"canEditFormatSettings,omitempty"`    // Whether the user can edit format settings on the object
	CanSelect                bool `json:"canSelect,omitempty"`                // Whether the user can select data from the object
	CanViewReflections       bool `json:"canViewReflections,omitempty"`       // Whether the user can view reflections on the object
	CanAlterReflections      bool `json:"canAlterReflections,omitempty"`      // Whether the user can alter reflections on the object
	CanCreateReflections     bool `json:"canCreateReflections,omitempty"`     // Whether the user can create reflections on the object
	CanDropReflections       bool `json:"canDropReflections,omitempty"`       // Whether the user can drop reflections on the object
}

// FunctionArg represents a function argument
type FunctionArg struct {
	Name string     `json:"name"`           // Argument name
	Type *FieldType `json:"type,omitempty"` // Argument type
}

// CatalogEntity represents a child entity in a catalog container
// Reference: OpenAPI schema CatalogEntity
type CatalogEntity struct {
	ID            string   `json:"id"`                      // Unique identifier of the entity
	Path          []string `json:"path"`                    // Full path to the entity
	Tag           string   `json:"tag,omitempty"`           // Version tag
	Type          string   `json:"type"`                    // Entity type (CONTAINER or DATASET)
	ContainerType string   `json:"containerType,omitempty"` // Container type (SPACE, SOURCE, FOLDER, HOME)
	DatasetType   string   `json:"datasetType,omitempty"`   // Dataset type (VIRTUAL_DATASET or PHYSICAL_DATASET)
}

type FileResponse struct {
	ID         string   `json:"id"`         // Unique identifier of the file
	EntityType string   `json:"entityType"` // Always "file" for files
	Path       []string `json:"path"`       // Path to the file
}

// TableField represents a field/column in a table or view
// Reference: OpenAPI schema TableField
type TableField struct {
	Name string     `json:"name"`           // Field name
	Type *FieldType `json:"type,omitempty"` // Field type information
}

// FieldType represents the type information for a field
// Reference: OpenAPI schema FieldType
type FieldType struct {
	Name        string       `json:"name,omitempty"`        // Type name (e.g., INTEGER, VARCHAR, STRUCT, LIST)
	Precision   int          `json:"precision,omitempty"`   // Precision for numeric types
	Scale       int          `json:"scale,omitempty"`       // Scale for numeric types
	SubSchema   []TableField `json:"subSchema,omitempty"`   // Sub-schema for STRUCT types
	ElementType *FieldType   `json:"elementType,omitempty"` // Element type for LIST types
}

// AccelerationRefreshPolicy represents the acceleration refresh policy for a dataset
// Reference: OpenAPI schema AccelerationRefreshPolicy
type AccelerationRefreshPolicy struct {
	ActivePolicyType *string `json:"activePolicyType,omitempty"` // Policy for refreshing Reflections (NEVER, PERIOD, SCHEDULE, REFRESH_ON_DATA_CHANGES)
	RefreshPeriodMs  *int64  `json:"refreshPeriodMs,omitempty"`  // Refresh period in milliseconds (minimum 3600000, default 3600000)
	RefreshSchedule  *string `json:"refreshSchedule,omitempty"`  // Cron expression for refresh schedule (UTC), e.g., "0 0 8 * * ?"
	GracePeriodMs    *int64  `json:"gracePeriodMs,omitempty"`    // Maximum age for Reflection data in milliseconds
	Method           *string `json:"method,omitempty"`           // Method for refreshing Reflections (AUTO, FULL, INCREMENTAL)
	RefreshField     *string `json:"refreshField,omitempty"`     // Field to use for incremental refresh
	NeverExpire      *bool   `json:"neverExpire,omitempty"`      // Whether Reflections never expire
}

// TableFormat represents the format information for a table
// Reference: OpenAPI schema TableFormat
type TableFormatResponse struct {
	Type                    string   `json:"type"`                              // Type of data in the table (Delta, Excel, Iceberg, JSON, Parquet, Text, Unknown, XLS)
	Name                    *string  `json:"name,omitempty"`                    // Table name
	FullPath                []string `json:"fullPath,omitempty"`                // Full path to the table
	Ctime                   *int     `json:"ctime,omitempty"`                   // Not used (always 0)
	IsFolder                *bool    `json:"isFolder,omitempty"`                // Whether the table was created from a folder
	Location                *string  `json:"location,omitempty"`                // Location where table metadata is stored
	IgnoreOtherFileFormats  *bool    `json:"ignoreOtherFileFormats,omitempty"`  // For Parquet folders, ignore non-Parquet files
	SkipFirstLine           *bool    `json:"skipFirstLine,omitempty"`           // Skip first line when creating table (Excel/Text)
	ExtractHeader           *bool    `json:"extractHeader,omitempty"`           // Extract column names from first line (Excel/Text)
	HasMergedCells          *bool    `json:"hasMergedCells,omitempty"`          // Expand merged cells (Excel)
	SheetName               *string  `json:"sheetName,omitempty"`               // Sheet name for Excel files with multiple sheets
	FieldDelimiter          *string  `json:"fieldDelimiter,omitempty"`          // Field delimiter character (Text), default: ","
	Quote                   *string  `json:"quote,omitempty"`                   // Quote character (Text), default: "\""
	Comment                 *string  `json:"comment,omitempty"`                 // Comment character (Text), default: "#"
	Escape                  *string  `json:"escape,omitempty"`                  // Escape character (Text), default: "\""
	LineDelimiter           *string  `json:"lineDelimiter,omitempty"`           // Line delimiter (Text), default: "\r\n"
	AutoGenerateColumnNames *bool    `json:"autoGenerateColumnNames,omitempty"` // Use existing column names (Text)
	TrimHeader              *bool    `json:"trimHeader,omitempty"`              // Trim column names (Text)
	AutoCorrectCorruptDates *bool    `json:"autoCorrectCorruptDates,omitempty"` // Auto-correct corrupted date fields (read-only)
}

// EngineRulesResponse represents a response for engine routing rules
// Reference: https://docs.dremio.com/dremio-cloud/api/engine-rules
type EngineRulesResponse struct {
	RuleSet *RuleSet `json:"ruleSet"` // The rule set containing all routing rules
}

type MaintenanceTaskResponse struct {
	ID         string                 `json:"id"`                  // Unique identifier of the maintenance task
	TaskType   string                 `json:"type"`                // Type of maintenance task (OPTIMIZE, EXPIRE_SNAPSHOTS)
	Level      string                 `json:"level"`               // Level of the maintenance task (TABLE)
	SourceName string                 `json:"sourceName"`          // Name of the source
	IsEnabled  bool                   `json:"isEnabled,omitempty"` // Whether the task is enabled
	TaskConfig *MaintenanceTaskConfig `json:"config,omitempty"`    // An object that contains a fully qualified object name in the indicated catalog as the target for the maintenance task.
}
