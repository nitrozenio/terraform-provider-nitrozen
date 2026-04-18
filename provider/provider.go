package provider

import (
	"context"

	"terraform-provider-nitrozen/internal/client"
	"terraform-provider-nitrozen/internal/resources"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NitrozenProvider struct{}

func New() provider.Provider {
	return &NitrozenProvider{}
}

func (p *NitrozenProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nitrozen"
}

func (p *NitrozenProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

// ✅ MODEL (IMPORTANT)
type providerModel struct {
	Token types.String `tfsdk:"token"`
}

// ✅ FIXED Configure
func (p *NitrozenProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Token.IsNull() || data.Token.IsUnknown() {
		resp.Diagnostics.AddError("Missing Token", "Token is required")
		return
	}

	// 🔥 THIS WAS THE ISSUE
	c := client.NewClient(data.Token.ValueString())

	resp.ResourceData = c
}

func (p *NitrozenProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewProjectResource,
	}
}

func (p *NitrozenProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
