package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/replicate/replicate-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &HardwareDataSource{}

func NewHardwareDataSource() datasource.DataSource {
	return &HardwareDataSource{}
}

// HardwareDataSource defines the data source implementation.
type HardwareDataSource struct {
	client *replicate.Client
}

// HardwareDataSourceModel describes the data source data model.
type HardwareDataSourceModel struct {
	Hardware []HardwareModel `tfsdk:"hardware"`
	Id       types.String    `tfsdk:"id"`
}

type HardwareModel struct {
	Name types.String `tfsdk:"name"`
	SKU  types.String `tfsdk:"sku"`
}

func (d *HardwareDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hardware"
}

func (d *HardwareDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves available hardware for Replicate models",

		Attributes: map[string]schema.Attribute{
			"hardware": schema.ListNestedAttribute{
				MarkdownDescription: "List of available hardware options",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the hardware option",
							Computed:            true,
						},
						"sku": schema.StringAttribute{
							MarkdownDescription: "SKU of the hardware option",
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

func (d *HardwareDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *HardwareDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HardwareDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Make API call to Replicate to get hardware options
	hardwareOptions, err := d.client.ListHardware(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read hardware options, got error: %s", err))
		return
	}

	// Map the API response to our data model
	for _, hw := range *hardwareOptions {
		data.Hardware = append(data.Hardware, HardwareModel{
			Name: types.StringValue(hw.Name),
			SKU:  types.StringValue(hw.SKU),
		})
	}

	// Generate a unique ID for this data source
	data.Id = types.StringValue("replicate_hardware")

	// Write logs using the tflog package
	tflog.Trace(ctx, "read hardware data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
