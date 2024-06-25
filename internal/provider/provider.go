package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/replicate/replicate-go"
)

// Ensure ReplicateProvider satisfies various provider interfaces.
var _ provider.Provider = &ReplicateProvider{}
var _ provider.ProviderWithFunctions = &ReplicateProvider{}

const (
	EnvAccApiToken = "REPLICATE_API_TOKEN"
)

var (
	UserAgent = "terraform-provider-replicate"
)

// ReplicateProvider defines the provider implementation.
type ReplicateProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ReplicateProviderModel describes the provider data model.
type ReplicateProviderModel struct {
	ApiToken types.String `tfsdk:"api_token"`
	BaseURL  types.String `tfsdk:"base_url"`
}

func (p *ReplicateProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "replicate"
	resp.Version = p.version
}

func (p *ReplicateProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "Replicate API token for authentication",
				Required:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Replicate API base URL",
				Optional:            true,
			},
		},
	}
}

func (p *ReplicateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ReplicateProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ApiToken.IsNull() {
		resp.Diagnostics.AddError(
			"Missing API Token",
			"The api_token attribute is required for the Replicate provider",
		)
		return
	}

	opts := []replicate.ClientOption{
		replicate.WithUserAgent(UserAgent + "/" + p.version),
		replicate.WithToken(data.ApiToken.ValueString()),
	}

	if !data.BaseURL.IsNull() {
		opts = append(opts, replicate.WithBaseURL(data.BaseURL.ValueString()))
	}

	client, err := replicate.NewClient(opts...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Replicate client",
			"An error occurred while creating the Replicate client: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ReplicateProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeploymentResource,
	}
}

func (p *ReplicateProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewHardwareDataSource,
		NewModelVersionDataSource,
	}
}

func (p *ReplicateProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ReplicateProvider{
			version: version,
		}
	}
}
