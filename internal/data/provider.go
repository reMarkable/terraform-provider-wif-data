// Copyright (c) HashiCorp, Inc.

package data

import (
	// Documentation:
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema

	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider              = &WifDataProvider{}
	_ provider.ProviderWithFunctions = &WifDataProvider{}
)

// WifDataProvider defines the provider implementation.
type WifDataProvider struct {
	version string
}
type WifDataProviderModel struct {
	PoolId    types.String `tfsdk:"pool_id"`
	ProjectId types.Int64  `tfsdk:"project_id"`
}

func (p *WifDataProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "wif"
	resp.Version = p.version
}

func (p *WifDataProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `A utility data source for generating principal and principalset URLs.

    This is primarily used to provide IAM permissions directly to an external entity.
    `,
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				MarkdownDescription: "Project ID where the pool is located",
				Required:            true,
			},
			"pool_id": schema.StringAttribute{
				MarkdownDescription: "Pool ID",
				Required:            true,
			},
		},
	}
}

func (p *WifDataProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data WifDataProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = data
}

func (p *WifDataProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPrincipalDataSource,
		NewPrincipalSetDataSource,
	}
}

func (p *WifDataProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *WifDataProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WifDataProvider{
			version: version,
		}
	}
}
