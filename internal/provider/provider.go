package provider

import (
	"context"
	"os"

	client "github.com/carlos-ffs/dremio-terraform-provider/internal/client"
	dremioDatasources "github.com/carlos-ffs/dremio-terraform-provider/internal/datasources"
	dremioResources "github.com/carlos-ffs/dremio-terraform-provider/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DremioProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type dremioProviderModel struct {
	Host                types.String `tfsdk:"host"`
	PersonalAccessToken types.String `tfsdk:"personal_access_token"`
	ProjectId           types.String `tfsdk:"project_id"`
	Ptype               types.String `tfsdk:"type"`
}

func (p *DremioProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dremio"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *DremioProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Dremio API Host. Defaults to https://api.dremio.cloud",
			},
			"personal_access_token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Dremio Personal Access Token",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Dremio Account Type. Defaults to cloud",
			},
			"project_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Dremio Project ID. Required for Dremio Cloud",
			},
		},
	}
}

func (p *DremioProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration

	personalAccessToken := os.Getenv("DREMIO_PAT")
	host := os.Getenv("DREMIO_HOST")
	ptype := os.Getenv("DREMIO_TYPE")
	projectId := os.Getenv("DREMIO_PROJECT_ID")

	var config dremioProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check configuration data, which should take precedence over
	// environment variable data, if found.
	if config.Host.ValueString() != "" {
		host = config.Host.ValueString()
	}
	if config.PersonalAccessToken.ValueString() != "" {
		personalAccessToken = config.PersonalAccessToken.ValueString()
	}
	if config.Ptype.ValueString() != "" {
		ptype = config.Ptype.ValueString()
	}
	if config.ProjectId.ValueString() != "" {
		projectId = config.ProjectId.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Dremio API Host",
			"The provider cannot create the Dremio API client as there is an unknown configuration value for the Dremio API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DREMIO_HOST environment variable.",
		)
	}

	if personalAccessToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_access_token"),
			"Unknown Dremio API Personal Access Token",
			"The provider cannot create the Dremio API client as there is an unknown configuration value for the Dremio API personal access token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DREMIO_PAT environment variable.",
		)
	}
	if ptype == "" {
		ptype = "cloud"
	}
	if ptype == "cloud" && projectId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("project_id"),
			"Unknown Dremio Project ID",
			"The provider cannot create the Dremio API client as there is an unknown configuration value for the Dremio Project ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DREMIO_PROJECT_ID environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Dremio client using the configuration values
	client, err := client.NewClient(&host, &personalAccessToken, &ptype, &projectId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Dremio API Client",
			"An unexpected error occurred when creating the Dremio API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Dremio Client Error: "+err.Error(),
		)
		return
	}

	// Make the Dremio client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *DremioProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		dremioDatasources.NewDremioSourceDataSource,
		dremioDatasources.NewDremioFolderDataSource,
		dremioDatasources.NewDremioFileDataSource,
		dremioDatasources.NewDremioTableDataSource,
		dremioDatasources.NewDremioUDFDataSource,
		dremioDatasources.NewDremioDatasetTagsDataSource,
		dremioDatasources.NewDremioDatasetWikiDataSource,
		dremioDatasources.NewDremioViewDataSource,
		dremioDatasources.NewDremioGrantsDataSource,
		dremioDatasources.NewDremioEngineDataSource,
		dremioDatasources.NewDremioEngineRuleSetDataSource,
		dremioDatasources.NewDremioDataMaintenanceTaskDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *DremioProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		dremioResources.NewDremioSourceResource,
		dremioResources.NewDremioFolderResource,
		dremioResources.NewDremioTableResource,
		dremioResources.NewDremioUDFResource,
		dremioResources.NewDremioDatasetTagsResource,
		dremioResources.NewDremioDatasetWikiResource,
		dremioResources.NewDremioViewResource,
		dremioResources.NewDremioGrantsResource,
		dremioResources.NewDremioEngineResource,
		dremioResources.NewDremioEngineRuleSetResource,
		dremioResources.NewDremioDataMaintenanceResource,
	}
}

func (p *DremioProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DremioProvider{
			version: version,
		}
	}
}
