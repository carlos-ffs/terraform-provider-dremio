package models

// MySQLConfig represents MySQL source configuration
type MySQLConfig struct {
	AuthenticationType string                   `json:"authenticationType,omitempty"`
	Username           string                   `json:"username,omitempty"`
	Password           string                   `json:"password,omitempty"`
	Hostname           string                   `json:"hostname"`
	Port               string                   `json:"port"`
	NetWriteTimeout    int                      `json:"netWriteTimeout,omitempty"`
	FetchSize          int                      `json:"fetchSize,omitempty"`
	MaxIdleConns       int                      `json:"maxIdleConns,omitempty"`
	IdleTimeSec        int                      `json:"idleTimeSec,omitempty"`
	PropertyList       []map[string]interface{} `json:"propertyList,omitempty"`
}

// PostgreSQLConfig represents PostgreSQL source configuration
type PostgreSQLConfig struct {
	AuthenticationType       string                   `json:"authenticationType,omitempty"`
	Username                 string                   `json:"username,omitempty"`
	Password                 string                   `json:"password,omitempty"`
	SecretResourceUrl        string                   `json:"secretResourceUrl,omitempty"`
	Hostname                 string                   `json:"hostname"`
	Port                     string                   `json:"port"`
	DatabaseName             string                   `json:"databaseName,omitempty"`
	UseSsl                   bool                     `json:"useSsl,omitempty"`
	FetchSize                int                      `json:"fetchSize,omitempty"`
	MaxIdleConns             int                      `json:"maxIdleConns,omitempty"`
	IdleTimeSec              int                      `json:"idleTimeSec,omitempty"`
	EncryptionValidationMode string                   `json:"encryptionValidationMode,omitempty"`
	PropertyList             []map[string]interface{} `json:"propertyList,omitempty"`
}

// S3Config represents S3 source configuration
type S3Config struct {
	CredentialType              string                   `json:"credentialType"`
	AssumedRoleARN              string                   `json:"assumedRoleARN,omitempty"`
	AwsAccessKey                string                   `json:"awsAccessKey,omitempty"`
	AwsAccessSecret             string                   `json:"awsAccessSecret,omitempty"`
	ExternalBucketList          []string                 `json:"externalBucketList,omitempty"`
	Secure                      bool                     `json:"secure,omitempty"`
	EnableAsync                 bool                     `json:"enableAsync,omitempty"`
	CompatibilityMode           bool                     `json:"compatibilityMode,omitempty"`
	RequesterPays               bool                     `json:"requesterPays,omitempty"`
	EnableFileStatusCheck       bool                     `json:"enableFileStatusCheck,omitempty"`
	IsPartitionInferenceEnabled bool                     `json:"isPartitionInferenceEnabled,omitempty"`
	RootPath                    string                   `json:"rootPath,omitempty"`
	KmsKeyARN                   string                   `json:"kmsKeyARN,omitempty"`
	DefaultCtasFormat           string                   `json:"defaultCtasFormat,omitempty"`
	PropertyList                []map[string]interface{} `json:"propertyList,omitempty"`
	WhitelistedBuckets          []string                 `json:"whitelistedBuckets,omitempty"`
	IsCachingEnabled            bool                     `json:"isCachingEnabled,omitempty"`
	MaxCacheSpacePct            int                      `json:"maxCacheSpacePct,omitempty"`
}

// SnowflakeConfig represents Snowflake source configuration
type SnowflakeConfig struct {
	Hostname             string `json:"hostname"`
	Port                 string `json:"port"`
	Database             string `json:"database,omitempty"`
	Warehouse            string `json:"warehouse,omitempty"`
	Username             string `json:"username"`
	Password             string `json:"password,omitempty"`
	FetchSize            int    `json:"fetchSize,omitempty"`
	MaxIdleConns         int    `json:"maxIdleConns,omitempty"`
	IdleTimeSec          int    `json:"idleTimeSec,omitempty"`
	QueryTimeoutSec      int    `json:"queryTimeoutSec,omitempty"`
	AuthMode             string `json:"authMode,omitempty"`
	PrivateKey           string `json:"privateKey,omitempty"`
	PrivateKeyPassphrase string `json:"privateKeyPassphrase,omitempty"`
}

// BigQueryConfig represents BigQuery source configuration
type BigQueryConfig struct {
	Hostname        string                   `json:"hostname,omitempty"`
	Port            string                   `json:"port,omitempty"`
	ProjectId       string                   `json:"projectId"`
	AuthMode        string                   `json:"authMode,omitempty"`
	ClientEmail     string                   `json:"clientEmail,omitempty"`
	PrivateKey      interface{}              `json:"privateKey,omitempty"` // Can be string or JSON object
	FetchSize       int                      `json:"fetchSize,omitempty"`
	MaxIdleConns    int                      `json:"maxIdleConns,omitempty"`
	IdleTimeSec     int                      `json:"idleTimeSec,omitempty"`
	QueryTimeoutSec int                      `json:"queryTimeoutSec,omitempty"`
	PropertyList    []map[string]interface{} `json:"propertyList,omitempty"`
}

// RedshiftConfig represents Redshift source configuration
type RedshiftConfig struct {
	ConnectionString   string                   `json:"connectionString"`
	AuthenticationType string                   `json:"authenticationType"`
	Username           string                   `json:"username,omitempty"`
	Password           string                   `json:"password,omitempty"`
	SecretResourceUrl  string                   `json:"secretResourceUrl,omitempty"`
	AwsProfile         string                   `json:"awsProfile,omitempty"`
	DbUser             string                   `json:"dbUser,omitempty"`
	FetchSize          int                      `json:"fetchSize,omitempty"`
	MaxIdleConns       int                      `json:"maxIdleConns,omitempty"`
	IdleTimeSec        int                      `json:"idleTimeSec,omitempty"`
	QueryTimeoutSec    int                      `json:"queryTimeoutSec,omitempty"`
	PropertyList       []map[string]interface{} `json:"propertyList,omitempty"`
}

// OracleConfig represents Oracle source configuration
type OracleConfig struct {
	AuthenticationType  string                   `json:"authenticationType,omitempty"`
	Username            string                   `json:"username,omitempty"`
	Password            string                   `json:"password,omitempty"`
	SecretResourceUrl   string                   `json:"secretResourceUrl,omitempty"`
	UseKerberos         bool                     `json:"useKerberos,omitempty"`
	Hostname            string                   `json:"hostname"`
	Port                string                   `json:"port"`
	Instance            string                   `json:"instance,omitempty"`
	UseSsl              bool                     `json:"useSsl,omitempty"`
	NativeEncryption    string                   `json:"nativeEncryption,omitempty"`
	UseTimezoneAsRegion bool                     `json:"useTimezoneAsRegion,omitempty"`
	IncludeSynonyms     bool                     `json:"includeSynonyms,omitempty"`
	MapDateToTimestamp  bool                     `json:"mapDateToTimestamp,omitempty"`
	FetchSize           int                      `json:"fetchSize,omitempty"`
	MaxIdleConns        int                      `json:"maxIdleConns,omitempty"`
	IdleTimeSec         int                      `json:"idleTimeSec,omitempty"`
	UseLdap             bool                     `json:"useLdap,omitempty"`
	BindDN              string                   `json:"bindDN,omitempty"`
	SslServerCertDN     string                   `json:"sslServerCertDN,omitempty"`
	PropertyList        []map[string]interface{} `json:"propertyList,omitempty"`
}

// MSSQLConfig represents MS SQL Server source configuration
type MSSQLConfig struct {
	AuthenticationType         string                   `json:"authenticationType,omitempty"`
	Username                   string                   `json:"username,omitempty"`
	Password                   string                   `json:"password,omitempty"`
	Hostname                   string                   `json:"hostname"`
	Port                       string                   `json:"port"`
	Database                   string                   `json:"database,omitempty"`
	UseSsl                     bool                     `json:"useSsl,omitempty"`
	FetchSize                  int                      `json:"fetchSize,omitempty"`
	MaxIdleConns               int                      `json:"maxIdleConns,omitempty"`
	IdleTimeSec                int                      `json:"idleTimeSec,omitempty"`
	ShowOnlyConnectionDatabase bool                     `json:"showOnlyConnectionDatabase,omitempty"`
	EncryptionValidationMode   string                   `json:"encryptionValidationMode,omitempty"`
	PropertyList               []map[string]interface{} `json:"propertyList,omitempty"`
}

// AzureStorageConfig represents Azure Storage source configuration
type AzureStorageConfig struct {
	AccountKind                 string                   `json:"accountKind"`
	AccountName                 string                   `json:"accountName"`
	AccessKey                   string                   `json:"accessKey,omitempty"`
	CredentialsType             string                   `json:"credentialsType,omitempty"`
	ClientId                    string                   `json:"clientId,omitempty"`
	ClientSecret                string                   `json:"clientSecret,omitempty"`
	TokenEndpoint               string                   `json:"tokenEndpoint,omitempty"`
	Containers                  []string                 `json:"containers,omitempty"`
	RootPath                    string                   `json:"rootPath,omitempty"`
	EnableAsync                 bool                     `json:"enableAsync,omitempty"`
	EnableSSL                   bool                     `json:"enableSSL,omitempty"`
	PropertyList                []map[string]interface{} `json:"propertyList,omitempty"`
	IsCachingEnabled            bool                     `json:"isCachingEnabled,omitempty"`
	MaxCacheSpacePct            int                      `json:"maxCacheSpacePct,omitempty"`
	DefaultCtasFormat           string                   `json:"defaultCtasFormat,omitempty"`
	IsPartitionInferenceEnabled bool                     `json:"isPartitionInferenceEnabled,omitempty"`
}

// ArcticConfig represents Arctic catalog source configuration
type ArcticConfig struct {
	StorageProvider         string                   `json:"storageProvider"` // AWS or AZURE
	CredentialType          string                   `json:"credentialType,omitempty"`
	AwsAccessKey            string                   `json:"awsAccessKey,omitempty"`
	AwsAccessSecret         string                   `json:"awsAccessSecret,omitempty"`
	AwsRootPath             string                   `json:"awsRootPath,omitempty"`
	AssumedRoleARN          string                   `json:"assumedRoleARN,omitempty"`
	AzureStorageAccount     string                   `json:"azureStorageAccount,omitempty"`
	AzureRootPath           string                   `json:"azureRootPath,omitempty"`
	AzureAuthenticationType string                   `json:"azureAuthenticationType,omitempty"`
	AzureAccessKey          string                   `json:"azureAccessKey,omitempty"`
	AzureApplicationId      string                   `json:"azureApplicationId,omitempty"`
	AzureClientSecret       string                   `json:"azureClientSecret,omitempty"`
	AzureOAuthTokenEndpoint string                   `json:"azureOAuthTokenEndpoint,omitempty"`
	PropertyList            []map[string]interface{} `json:"propertyList,omitempty"`
	AsyncEnabled            bool                     `json:"asyncEnabled,omitempty"`
	IsCachingEnabled        bool                     `json:"isCachingEnabled,omitempty"`
	MaxCacheSpacePct        int                      `json:"maxCacheSpacePct,omitempty"`
	DefaultCtasFormat       string                   `json:"defaultCtasFormat,omitempty"`
	ArcticCatalogId         string                   `json:"arcticCatalogId"`
}

// AWSGlueConfig represents AWS Glue Data Catalog source configuration
type AWSGlueConfig struct {
	CredentialType                       string                   `json:"credentialType"`
	AssumedRoleARN                       string                   `json:"assumedRoleARN,omitempty"`
	AwsAccessKey                         string                   `json:"awsAccessKey,omitempty"`
	AwsAccessSecret                      string                   `json:"awsAccessSecret,omitempty"`
	RegionNameSelection                  string                   `json:"regionNameSelection"`
	Secure                               bool                     `json:"secure,omitempty"`
	EnableAsync                          bool                     `json:"enableAsync,omitempty"`
	LakeFormationEnableAccessPermissions bool                     `json:"lakeFormationEnableAccessPermissions,omitempty"`
	PropertyList                         []map[string]interface{} `json:"propertyList,omitempty"`
	AllowedDatabases                     []string                 `json:"allowedDatabases,omitempty"`
	IsCachingEnabled                     bool                     `json:"isCachingEnabled,omitempty"`
	MaxCacheSpacePct                     int                      `json:"maxCacheSpacePct,omitempty"`
}

// Db2Config represents IBM Db2 source configuration
type Db2Config struct {
	Database        string                   `json:"database"`
	Hostname        string                   `json:"hostname"`
	Username        string                   `json:"username"`
	Password        string                   `json:"password,omitempty"`
	Port            string                   `json:"port"`
	FetchSize       int                      `json:"fetchSize,omitempty"`
	MaxIdleConns    int                      `json:"maxIdleConns,omitempty"`
	IdleTimeSec     int                      `json:"idleTimeSec,omitempty"`
	QueryTimeoutSec int                      `json:"queryTimeoutSec,omitempty"`
	PropertyList    []map[string]interface{} `json:"propertyList,omitempty"`
}

// IcebergRESTCatalogConfig represents Iceberg REST Catalog source configuration
type IcebergRESTCatalogConfig struct {
	PropertyList                 []map[string]interface{} `json:"propertyList,omitempty"`
	SecretPropertyList           []map[string]interface{} `json:"secretPropertyList,omitempty"`
	EnableAsync                  bool                     `json:"enableAsync,omitempty"`
	IsCachingEnabled             bool                     `json:"isCachingEnabled,omitempty"`
	MaxCacheSpacePct             int                      `json:"maxCacheSpacePct,omitempty"`
	RestEndpointUri              string                   `json:"restEndpointUri"`
	AllowedNamespaces            []string                 `json:"allowedNamespaces,omitempty"`
	IsUsingVendedCredentials     bool                     `json:"isUsingVendedCredentials,omitempty"`
	IsRecursiveAllowedNamespaces bool                     `json:"isRecursiveAllowedNamespaces,omitempty"`
}

// AzureSynapseConfig represents Microsoft Azure Synapse Analytics source configuration
type AzureSynapseConfig struct {
	Hostname                 string                   `json:"hostname"`
	Port                     string                   `json:"port,omitempty"`
	Username                 string                   `json:"username,omitempty"`
	Password                 string                   `json:"password,omitempty"`
	AuthenticationType       string                   `json:"authenticationType"`
	FetchSize                int                      `json:"fetchSize,omitempty"`
	UseSsl                   bool                     `json:"useSsl,omitempty"`
	EnableServerVerification bool                     `json:"enableServerVerification,omitempty"`
	MaxIdleConns             int                      `json:"maxIdleConns,omitempty"`
	IdleTimeSec              int                      `json:"idleTimeSec,omitempty"`
	QueryTimeoutSec          int                      `json:"queryTimeoutSec,omitempty"`
	Database                 string                   `json:"database,omitempty"`
	PropertyList             []map[string]interface{} `json:"propertyList,omitempty"`
}

// SAPHANAConfig represents SAP HANA source configuration
type SAPHANAConfig struct {
	Hostname        string                   `json:"hostname"`
	Port            string                   `json:"port"`
	Schema          string                   `json:"schema,omitempty"`
	Username        string                   `json:"username"`
	Password        string                   `json:"password,omitempty"`
	MaxIdleConns    int                      `json:"maxIdleConns,omitempty"`
	IdleTimeSec     int                      `json:"idleTimeSec,omitempty"`
	QueryTimeoutSec int                      `json:"queryTimeoutSec,omitempty"`
	FetchSize       int                      `json:"fetchSize,omitempty"`
	PropertyList    []map[string]interface{} `json:"propertyList,omitempty"`
}

// SnowflakeOpenCatalogConfig represents Snowflake Open Catalog source configuration
type SnowflakeOpenCatalogConfig struct {
	PropertyList                     []map[string]interface{} `json:"propertyList,omitempty"`
	SecretPropertyList               []map[string]interface{} `json:"secretPropertyList,omitempty"`
	EnableAsync                      bool                     `json:"enableAsync,omitempty"`
	IsCachingEnabled                 bool                     `json:"isCachingEnabled,omitempty"`
	MaxCacheSpacePct                 int                      `json:"maxCacheSpacePct,omitempty"`
	RestEndpointUri                  string                   `json:"restEndpointUri"`
	AllowedNamespaces                []string                 `json:"allowedNamespaces,omitempty"`
	IsUsingVendedCredentials         bool                     `json:"isUsingVendedCredentials,omitempty"`
	IsRecursiveAllowedNamespaces     bool                     `json:"isRecursiveAllowedNamespaces,omitempty"`
	SnowflakeOpenCatalogWarehouse    string                   `json:"snowflakeOpenCatalogWarehouse"`
	SnowflakeOpenCatalogClientId     string                   `json:"snowflakeOpenCatalogClientId"`
	SnowflakeOpenCatalogClientSecret string                   `json:"snowflakeOpenCatalogClientSecret,omitempty"`
	SnowflakeOpenCatalogScope        string                   `json:"snowflakeOpenCatalogScope,omitempty"`
}

// UnityCatalogConfig represents Unity Catalog source configuration
type UnityCatalogConfig struct {
	PropertyList                 []map[string]interface{} `json:"propertyList,omitempty"`
	SecretPropertyList           []map[string]interface{} `json:"secretPropertyList,omitempty"`
	EnableAsync                  bool                     `json:"enableAsync,omitempty"`
	IsCachingEnabled             bool                     `json:"isCachingEnabled,omitempty"`
	MaxCacheSpacePct             int                      `json:"maxCacheSpacePct,omitempty"`
	RestEndpointUri              string                   `json:"restEndpointUri"`
	AllowedNamespaces            []string                 `json:"allowedNamespaces,omitempty"`
	IsUsingVendedCredentials     bool                     `json:"isUsingVendedCredentials,omitempty"`
	IsRecursiveAllowedNamespaces bool                     `json:"isRecursiveAllowedNamespaces,omitempty"`
	UnityAuthToken               string                   `json:"unityAuthToken,omitempty"`
	UnityCatalog                 string                   `json:"unityCatalog"`
}

// VerticaConfig represents Vertica source configuration
type VerticaConfig struct {
	Database        string                   `json:"database"`
	Hostname        string                   `json:"hostname"`
	Username        string                   `json:"username"`
	Password        string                   `json:"password,omitempty"`
	Port            string                   `json:"port"`
	FetchSize       int                      `json:"fetchSize,omitempty"`
	MaxIdleConns    int                      `json:"maxIdleConns,omitempty"`
	IdleTimeSec     int                      `json:"idleTimeSec,omitempty"`
	QueryTimeoutSec int                      `json:"queryTimeoutSec,omitempty"`
	PropertyList    []map[string]interface{} `json:"propertyList,omitempty"`
}
