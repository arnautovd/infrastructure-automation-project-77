package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

type ServerResource struct {
	client *APIClient
}

func (r *ServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *ServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"make_from": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rplan": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"location": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ssh_keys": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"public_address": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

type serverModel struct {
	Name          types.String `tfsdk:"name"`
	MakeFrom      types.String `tfsdk:"make_from"`
	Rplan         types.String `tfsdk:"rplan"`
	Location      types.String `tfsdk:"location"`
	SSHKeys       types.Set    `tfsdk:"ssh_keys"`
	PublicAddress types.String `tfsdk:"public_address"`
}

func (r *ServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*APIClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "The Vscale provider did not provide an API client.")
		return
	}
	r.client = client
}

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "The Vscale provider must be configured before creating a server.")
		return
	}
	var plan serverModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	keys, err := stringSetToIDs(ctx, plan.SSHKeys)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("ssh_keys"), "Invalid SSH key ID", err.Error())
		return
	}
	server, err := r.client.CreateServer(ctx, createServerRequest{
		Name: plan.Name.ValueString(), MakeFrom: plan.MakeFrom.ValueString(),
		Rplan: plan.Rplan.ValueString(), Location: plan.Location.ValueString(),
		DoStart: true, Keys: keys,
	})
	if err != nil {
		resp.Diagnostics.AddError("Create Vscale server failed", err.Error())
		return
	}
	resp.State.SetAttribute(ctx, path.Root("name"), server.Name)
	resp.State.SetAttribute(ctx, path.Root("make_from"), server.MadeFrom)
	resp.State.SetAttribute(ctx, path.Root("rplan"), server.Rplan)
	resp.State.SetAttribute(ctx, path.Root("location"), server.Location)
	resp.State.SetAttribute(ctx, path.Root("public_address"), server.PublicAddresses.Address)
	resp.State.SetAttribute(ctx, path.Root("ssh_keys"), plan.SSHKeys)
	resp.State.SetAttribute(ctx, path.Root("id"), strconv.FormatInt(server.CTID, 10))
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateID types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &stateID)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := strconv.ParseInt(stateID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Vscale server ID", err.Error())
		return
	}
	server, err := r.client.GetServer(ctx, id)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Vscale server failed", err.Error())
		return
	}
	resp.State.SetAttribute(ctx, path.Root("name"), server.Name)
	resp.State.SetAttribute(ctx, path.Root("make_from"), server.MadeFrom)
	resp.State.SetAttribute(ctx, path.Root("rplan"), server.Rplan)
	resp.State.SetAttribute(ctx, path.Root("location"), server.Location)
	resp.State.SetAttribute(ctx, path.Root("public_address"), server.PublicAddresses.Address)
}

func (r *ServerResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Server replacement required", "All Vscale server arguments require replacement; Terraform should schedule a new resource.")
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateID types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &stateID)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := strconv.ParseInt(stateID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Vscale server ID", err.Error())
		return
	}
	if err := r.client.DeleteServer(ctx, id); err != nil {
		resp.Diagnostics.AddError("Delete Vscale server failed", err.Error())
	}
}

func stringSetToIDs(ctx context.Context, value types.Set) ([]int64, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}
	var values []string
	if diags := value.ElementsAs(ctx, &values, false); diags.HasError() {
		return nil, fmt.Errorf("SSH keys must be a set of numeric IDs")
	}
	result := make([]int64, 0, len(values))
	for _, value := range values {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("SSH key ID %q is not numeric", value)
		}
		result = append(result, id)
	}
	return result, nil
}
