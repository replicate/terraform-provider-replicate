package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/replicate/replicate-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ModelVersionDataSource{}

func NewModelVersionDataSource() datasource.DataSource {
	return &ModelVersionDataSource{}
}

// ModelVersionDataSource defines the data source implementation.
type ModelVersionDataSource struct {
	client *replicate.Client
}

// ModelVersionDataSourceModel describes the data source data model.
type ModelVersionDataSourceModel struct {
	Model    types.String        `tfsdk:"model"`
	Versions []ModelVersionModel `tfsdk:"versions"`
	Id       types.String        `tfsdk:"id"`
}

type ModelVersionModel struct {
	ID         types.String `tfsdk:"id"`
	CreatedAt  types.String `tfsdk:"created_at"`
	CogVersion types.String `tfsdk:"cog_version"`
}

func (d *ModelVersionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model_version"
}

func (d *ModelVersionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves versions of a specific Replicate model",

		Attributes: map[string]schema.Attribute{
			"model": schema.StringAttribute{
				MarkdownDescription: "Model identifier ({model_owner}/{model_name})",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^/]+/[^/]+$`),
						"must match the format {model_owner}/{model_name}",
					),
				},
			},
			"versions": schema.ListNestedAttribute{
				MarkdownDescription: "List of model versions",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the model version",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "The creation time of the model version",
							Computed:            true,
						},
						"cog_version": schema.StringAttribute{
							MarkdownDescription: "The Cog version used for this model version",
							Computed:            true,
						},
					},
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier for this data source",
				Computed:            true,
			},
		},
	}
}

func (d *ModelVersionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*replicate.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *replicate.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ModelVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ModelVersionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Make API call to Replicate to get model versions
	components := strings.Split(data.Model.ValueString(), "/")
	if len(components) != 2 {
		resp.Diagnostics.AddError("Invalid model identifier", fmt.Sprintf("Expected {model_owner}/{model_name}, got: %s", data.Model.ValueString()))
		return
	}
	versions, err := d.client.ListModelVersions(ctx, components[0], components[1])
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read model versions, got error: %s", err))
		return
	}

	// Map the API response to our data model
	for _, version := range versions.Results {
		data.Versions = append(data.Versions, ModelVersionModel{
			ID:         types.StringValue(version.ID),
			CreatedAt:  types.StringValue(version.CreatedAt),
			CogVersion: types.StringValue(version.CogVersion),
		})
	}

	// Generate a unique ID for this data source
	data.Id = types.StringValue(fmt.Sprintf("%s/versions", data.Model.ValueString()))

	// Write logs using the tflog package
	tflog.Trace(ctx, "read model version data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
