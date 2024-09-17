// Copyright (c) HashiCorp, Inc.

package data

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PrincipalSetDataSource{}

func NewPrincipalSetDataSource() datasource.DataSource {
	return &PrincipalSetDataSource{}
}

// PrincipalSetDataSource defines the data source implementation.
type PrincipalSetDataSource struct {
	providerData WifDataProviderModel
}

// PrincipalSetDataSourceModel describes the data source data model.
type PrincipalSetDataSourceModel struct {
	Url              types.String `tfsdk:"url"`
	Target           types.String `tfsdk:"target"`
	SourceExpression types.String `tfsdk:"source_expression"`
	Id               types.String `tfsdk:"id"`
}

func (d *PrincipalSetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_principal_set"
}

func (d *PrincipalSetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data Source for generating principalSet urls.",

		Attributes: map[string]schema.Attribute{
			"target": schema.StringAttribute{
				MarkdownDescription: "Target identifier",
				Required:            true,
			},
			"source_expression": schema.StringAttribute{
				MarkdownDescription: "Source Epression (valid CEL expression)",
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

func (d *PrincipalSetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PrincipalSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PrincipalSetDataSourceModel

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

	se := data.SourceExpression.ValueString()
	t := data.Target.ValueString()
	env, err := cel.NewEnv(cel.Variable(t, cel.StringType))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to setup cel parser, got error: %s", err))
		return
	}
	_, issues := env.Parse(se)
	if issues != nil && issues.Err() != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse cel expression, got error: %v", issues))
		return
	}

	data.Id = types.StringValue("wif_principal_set:" + data.Target.ValueString() + ":" + data.SourceExpression.ValueString())
	url := fmt.Sprintf("principalSet://iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/%s/%s", d.providerData.ProjectId, d.providerData.PoolId.ValueString(), t, se)
	data.Url = types.StringValue(url)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
