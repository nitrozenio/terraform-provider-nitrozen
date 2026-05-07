package resources

import (
	"context"
	"encoding/json"

	"terraform-provider-nitrozen/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EntryResource struct {
	client *client.Client
}

func NewEntryResource() resource.Resource {
	return &EntryResource{}
}

func (r *EntryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "nitrozen_entry"
}

func (r *EntryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a changelog entry within a Nitrozen project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The unique ID of the entry.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the project this entry belongs to.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "The title of the changelog entry.",
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "The body content of the changelog entry (supports Markdown).",
			},
			"category": schema.StringAttribute{
				Required:    true,
				Description: "The category of the entry. One of: new, improvement, fix, announcement.",
			},
			"is_published": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the entry is publicly visible. Defaults to false.",
			},
		},
	}
}

func (r *EntryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

type EntryModel struct {
	ID          types.Int64  `tfsdk:"id"`
	ProjectID   types.Int64  `tfsdk:"project_id"`
	Title       types.String `tfsdk:"title"`
	Content     types.String `tfsdk:"content"`
	Category    types.String `tfsdk:"category"`
	IsPublished types.Bool   `tfsdk:"is_published"`
}

func (r *EntryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EntryModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, err := r.client.CreateEntry(
		plan.ProjectID.ValueInt64(),
		plan.Title.ValueString(),
		plan.Content.ValueString(),
		plan.Category.ValueString(),
		plan.IsPublished.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating entry", err.Error())
		return
	}

	id, err := client.ExtractID(body)
	if err != nil {
		resp.Diagnostics.AddError("Error reading entry ID", err.Error())
		return
	}

	plan.ID = types.Int64Value(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EntryModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, err := r.client.GetEntry(state.ProjectID.ValueInt64(), state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading entry", err.Error())
		return
	}

	var result struct {
		Data struct {
			ID          int64  `json:"id"`
			Title       string `json:"title"`
			Content     string `json:"content"`
			Category    string `json:"category"`
			IsPublished bool   `json:"is_published"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing entry response", err.Error())
		return
	}

	state.ID = types.Int64Value(result.Data.ID)
	state.Title = types.StringValue(result.Data.Title)
	state.Content = types.StringValue(result.Data.Content)
	state.Category = types.StringValue(result.Data.Category)
	state.IsPublished = types.BoolValue(result.Data.IsPublished)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EntryModel
	var state EntryModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, err := r.client.UpdateEntry(
		state.ProjectID.ValueInt64(),
		state.ID.ValueInt64(),
		plan.Title.ValueString(),
		plan.Content.ValueString(),
		plan.Category.ValueString(),
		plan.IsPublished.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating entry", err.Error())
		return
	}

	var result struct {
		Data struct {
			ID          int64  `json:"id"`
			Title       string `json:"title"`
			Content     string `json:"content"`
			Category    string `json:"category"`
			IsPublished bool   `json:"is_published"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing entry response", err.Error())
		return
	}

	state.Title = types.StringValue(result.Data.Title)
	state.Content = types.StringValue(result.Data.Content)
	state.Category = types.StringValue(result.Data.Category)
	state.IsPublished = types.BoolValue(result.Data.IsPublished)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EntryModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteEntry(state.ProjectID.ValueInt64(), state.ID.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Error deleting entry", err.Error())
	}
}
