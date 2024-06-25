package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/replicate/replicate-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DeploymentResource{}
var _ resource.ResourceWithImportState = &DeploymentResource{}

func NewDeploymentResource() resource.Resource {
	return &DeploymentResource{}
}

// DeploymentResource defines the resource implementation.
type DeploymentResource struct {
	client *replicate.Client
}

// DeploymentResourceModel describes the resource data model.
type DeploymentResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Owner        types.String `tfsdk:"owner"`
	Model        types.String `tfsdk:"model"`
	Version      types.String `tfsdk:"version"`
	Hardware     types.String `tfsdk:"hardware"`
	MinInstances types.Int64  `tfsdk:"min_instances"`
	MaxInstances types.Int64  `tfsdk:"max_instances"`
	Id           types.String `tfsdk:"id"`
}

func (r *DeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment"
}

func (r *DeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Deployment resource",

		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				MarkdownDescription: "Owner of the deployment",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the deployment",
				Required:            true,
			},
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
			"version": schema.StringAttribute{
				MarkdownDescription: "Model version ID",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-fA-F0-9]+$`), "must be a valid version ID"),
				},
			},
			"hardware": schema.StringAttribute{
				MarkdownDescription: "Hardware SKU for the deployment",
				Required:            true,
			},
			"min_instances": schema.Int64Attribute{
				MarkdownDescription: "Minimum number of instances",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"max_instances": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of instances",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.AtLeastSumOf(path.MatchRoot("min_instances")),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *DeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeploymentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create deployment with API
	deployment, err := r.client.CreateDeployment(ctx, replicate.CreateDeploymentOptions{
		Name:         data.Name.ValueString(),
		Model:        data.Model.ValueString(),
		Version:      data.Version.ValueString(),
		Hardware:     data.Hardware.ValueString(),
		MinInstances: int(data.MinInstances.ValueInt64()),
		MaxInstances: int(data.MaxInstances.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create deployment, got error: %s", err))
		return
	}

	// Update the model with the latest data
	data.Owner = types.StringValue(deployment.Owner)
	data.Name = types.StringValue(deployment.Name)
	data.Model = types.StringValue(deployment.CurrentRelease.Model)
	data.Version = types.StringValue(deployment.CurrentRelease.Version)
	data.Hardware = types.StringValue(deployment.CurrentRelease.Configuration.Hardware)
	data.MinInstances = types.Int64Value(int64(deployment.CurrentRelease.Configuration.MinInstances))
	data.MaxInstances = types.Int64Value(int64(deployment.CurrentRelease.Configuration.MaxInstances))
	data.Id = types.StringValue(fmt.Sprintf("%s/%s", deployment.Owner, deployment.Name))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeploymentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts := strings.Split(data.Id.ValueString(), "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Expected ID in format owner/name, got: %s", data.Id.ValueString()))
		return
	}

	// Get deployment from API
	deployment, err := r.client.GetDeployment(ctx, parts[0], parts[1])
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read deployment, got error: %s", err))
		return
	}

	// Update the model with the latest data
	data.Owner = types.StringValue(deployment.Owner)
	data.Name = types.StringValue(deployment.Name)
	data.Model = types.StringValue(deployment.CurrentRelease.Model)
	data.Version = types.StringValue(deployment.CurrentRelease.Version)
	data.Hardware = types.StringValue(deployment.CurrentRelease.Configuration.Hardware)
	data.MinInstances = types.Int64Value(int64(deployment.CurrentRelease.Configuration.MinInstances))
	data.MaxInstances = types.Int64Value(int64(deployment.CurrentRelease.Configuration.MaxInstances))
	data.Id = types.StringValue(fmt.Sprintf("%s/%s", deployment.Owner, deployment.Name))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeploymentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update deployment with API
	opts := replicate.UpdateDeploymentOptions{}
	if !data.Version.IsNull() {
		opts.Version = data.Version.ValueStringPointer()
	}
	if !data.Hardware.IsNull() {
		opts.Hardware = data.Hardware.ValueStringPointer()
	}
	if !data.MinInstances.IsNull() {
		minInstances := int(*data.MinInstances.ValueInt64Pointer())
		opts.MinInstances = &minInstances
	}
	if !data.MaxInstances.IsNull() {
		maxInstances := int(*data.MaxInstances.ValueInt64Pointer())
		opts.MaxInstances = &maxInstances
	}
	_, err := r.client.UpdateDeployment(ctx, data.Owner.ValueString(), data.Name.ValueString(), opts)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update deployment, got error: %s", err))
		return
	}
	data.Id = types.StringValue(fmt.Sprintf("%s/%s", data.Owner.ValueString(), data.Name.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeploymentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDeployment(ctx, data.Owner.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete deployment, got error: %s", err))
		return
	}
}

func (r *DeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*replicate.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *replicate.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}
