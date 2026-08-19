package main

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = (*VscaleProvider)(nil)

type VscaleProvider struct {
	version string
}

type providerModel struct {
	Token  types.String `tfsdk:"token"`
	APIURL types.String `tfsdk:"api_url"`
}

func NewProvider() provider.Provider {
	return &VscaleProvider{version: version}
}

func (p *VscaleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vscale"
	resp.Version = p.version
}

func (p *VscaleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Vscale API token. Defaults to VSCALE_API_TOKEN.",
			},
			"api_url": schema.StringAttribute{
				Optional:    true,
				Description: "Vscale API base URL. Defaults to https://api.vscale.io/v1/.",
			},
		},
	}
}

func (p *VscaleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerModel
	if req.Config.Raw.IsNull() {
		data = providerModel{}
	} else {
		resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	token := data.Token.ValueString()
	if token == "" {
		token = os.Getenv("VSCALE_API_TOKEN")
	}
	if strings.TrimSpace(token) == "" {
		resp.Diagnostics.AddError("Missing Vscale API token", "Set token in the provider configuration or VSCALE_API_TOKEN.")
		return
	}

	baseURL := data.APIURL.ValueString()
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.vscale.io/v1/"
	}

	client, err := NewAPIClient(token, baseURL)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Vscale API URL", err.Error())
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *VscaleProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServerResource,
	}
}

func (p *VscaleProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
