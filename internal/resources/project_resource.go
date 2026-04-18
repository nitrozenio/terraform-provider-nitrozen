package resources

import (
	"context"
	"encoding/json"

	"terraform-provider-nitrozen/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProjectResource struct {
	client *client.Client
}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "nitrozen_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

// ✅ MODEL (IMPORTANT)
type ProjectModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// ✅ FIXED CREATE
func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ✅ Handle optional description safely
	description := ""
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		description = plan.Description.ValueString()
	}

	body, err := r.client.DoRequest("POST", "/projects", map[string]string{
		"name":        plan.Name.ValueString(),
		"description": description,
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}

	id, err := client.ExtractID(body)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", err.Error())
		return
	}

	// ✅ Build state PROPERLY
	state := ProjectModel{
		ID:          types.Int64Value(id),
		Name:        plan.Name,
		Description: plan.Description,
	}

	// ✅ VERY IMPORTANT: capture diagnostics
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProjectModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, err := r.client.GetProject(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}

	var result struct {
		Data struct {
			ID          int64  `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", err.Error())
		return
	}

	state.ID = types.Int64Value(result.Data.ID)
	state.Name = types.StringValue(result.Data.Name)
	state.Description = types.StringValue(result.Data.Description)

	resp.State.Set(ctx, &state)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProjectModel
	var state ProjectModel

	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	description := ""
	if !plan.Description.IsNull() {
		description = plan.Description.ValueString()
	}

	body, err := r.client.UpdateProject(
		state.ID.ValueInt64(),
		plan.Name.ValueString(),
		description,
	)

	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}

	var result struct {
		Data struct {
			ID          int64  `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"data"`
	}

	json.Unmarshal(body, &result)

	state.ID = types.Int64Value(result.Data.ID)
	state.Name = types.StringValue(result.Data.Name)
	state.Description = types.StringValue(result.Data.Description)

	resp.State.Set(ctx, &state)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProjectModel

	req.State.Get(ctx, &state)

	err := r.client.DeleteProject(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}
}
