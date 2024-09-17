// Copyright (c) HashiCorp, Inc.

package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PrincipalDataSource{}

func NewPrincipalDataSource() datasource.DataSource {
	return &PrincipalDataSource{}
}

// PrincipalDataSource defines the data source implementation.
type PrincipalDataSource struct {
	providerData WifDataProviderModel
}

// PrincipalDataSourceModel describes the data source data model.
type PrincipalDataSourceModel struct {
	Url     types.String `tfsdk:"url"`
	Subject types.String `tfsdk:"subject"`
	Id      types.String `tfsdk:"id"`
}

func (d *PrincipalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_principal"
}

func (d *PrincipalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data Source for generating principal urls.",

		Attributes: map[string]schema.Attribute{
			"subject": schema.StringAttribute{
				MarkdownDescription: "Subject identifier",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Principal URL",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "datasource identifier",
				Computed:            true,
			},
		},
	}
}

func (d *PrincipalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(WifDataProviderModel)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected WifDataProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.providerData = providerData
}

func (d *PrincipalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PrincipalDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	data.Id = types.StringValue("wif_principal:" + data.Subject.ValueString())
	url := fmt.Sprintf("principal://iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/subject/%s", d.providerData.ProjectId, d.providerData.PoolId.ValueString(), data.Subject.ValueString())
	data.Url = types.StringValue(url)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
