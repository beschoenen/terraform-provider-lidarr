package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/lidarr"
)

const (
	notificationEmailResourceName   = "notification_email"
	notificationEmailImplementation = "Email"
	notificationEmailConfigContract = "EmailSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationEmailResource{}
	_ resource.ResourceWithImportState = &NotificationEmailResource{}
)

func NewNotificationEmailResource() resource.Resource {
	return &NotificationEmailResource{}
}

// NotificationEmailResource defines the notification implementation.
type NotificationEmailResource struct {
	client *lidarr.Lidarr
}

// NotificationEmail describes the notification data model.
type NotificationEmail struct {
	Tags                  types.Set    `tfsdk:"tags"`
	To                    types.Set    `tfsdk:"to"`
	Cc                    types.Set    `tfsdk:"cc"`
	Bcc                   types.Set    `tfsdk:"bcc"`
	From                  types.String `tfsdk:"from"`
	Server                types.String `tfsdk:"server"`
	Name                  types.String `tfsdk:"name"`
	Username              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	ID                    types.Int64  `tfsdk:"id"`
	Port                  types.Int64  `tfsdk:"port"`
	RequireEncryption     types.Bool   `tfsdk:"require_encryption"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	OnReleaseImport       types.Bool   `tfsdk:"on_release_import"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnDownloadFailure     types.Bool   `tfsdk:"on_download_failure"`
	OnUpgrade             types.Bool   `tfsdk:"on_upgrade"`
	OnImportFailure       types.Bool   `tfsdk:"on_import_failure"`
}

func (n NotificationEmail) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		From:                  n.From,
		To:                    n.To,
		Cc:                    n.Cc,
		Bcc:                   n.Bcc,
		Server:                n.Server,
		Port:                  n.Port,
		Username:              n.Username,
		Password:              n.Password,
		Name:                  n.Name,
		ID:                    n.ID,
		RequireEncryption:     n.RequireEncryption,
		OnGrab:                n.OnGrab,
		OnReleaseImport:       n.OnReleaseImport,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		OnApplicationUpdate:   n.OnApplicationUpdate,
		OnHealthIssue:         n.OnHealthIssue,
		OnDownloadFailure:     n.OnDownloadFailure,
		OnUpgrade:             n.OnUpgrade,
		OnImportFailure:       n.OnImportFailure,
	}
}

func (n *NotificationEmail) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.From = notification.From
	n.To = notification.To
	n.Cc = notification.Cc
	n.Bcc = notification.Bcc
	n.Server = notification.Server
	n.Port = notification.Port
	n.Username = notification.Username
	n.Password = notification.Password
	n.Name = notification.Name
	n.ID = notification.ID
	n.RequireEncryption = notification.RequireEncryption
	n.OnGrab = notification.OnGrab
	n.OnReleaseImport = notification.OnReleaseImport
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnDownloadFailure = notification.OnDownloadFailure
	n.OnUpgrade = notification.OnUpgrade
	n.OnImportFailure = notification.OnImportFailure
}

func (r *NotificationEmailResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationEmailResourceName
}

func (r *NotificationEmailResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Email resource.\nFor more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Email](https://wiki.servarr.com/lidarr/supported#email).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Required:            true,
			},
			"on_import_failure": schema.BoolAttribute{
				MarkdownDescription: "On download flag.",
				Required:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Required:            true,
			},
			"on_download_failure": schema.BoolAttribute{
				MarkdownDescription: "On movie delete flag.",
				Required:            true,
			},
			"on_release_import": schema.BoolAttribute{
				MarkdownDescription: "On movie file delete for upgrade flag.",
				Required:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Required:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Required:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationEmail name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"require_encryption": schema.BoolAttribute{
				MarkdownDescription: "Require encryption flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"server": schema.StringAttribute{
				MarkdownDescription: "Server.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"from": schema.StringAttribute{
				MarkdownDescription: "From.",
				Required:            true,
			},
			"to": schema.SetAttribute{
				MarkdownDescription: "To.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"cc": schema.SetAttribute{
				MarkdownDescription: "Cc.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"bcc": schema.SetAttribute{
				MarkdownDescription: "Bcc.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationEmailResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*lidarr.Lidarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *lidarr.Lidarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NotificationEmailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationEmail

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationEmail
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationEmailResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationEmailResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationEmailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationEmail

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationEmail current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationEmailResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationEmailResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationEmailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationEmail

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationEmail
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationEmailResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationEmailResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationEmailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationEmail

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationEmail current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationEmailResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationEmailResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationEmailResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationEmailResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationEmail) write(ctx context.Context, notification *lidarr.NotificationOutput) {
	genericNotification := Notification{
		OnGrab:                types.BoolValue(notification.OnGrab),
		OnImportFailure:       types.BoolValue(notification.OnImportFailure),
		OnUpgrade:             types.BoolValue(notification.OnUpgrade),
		OnDownloadFailure:     types.BoolValue(notification.OnDownloadFailure),
		OnReleaseImport:       types.BoolValue(notification.OnReleaseImport),
		OnHealthIssue:         types.BoolValue(notification.OnHealthIssue),
		OnApplicationUpdate:   types.BoolValue(notification.OnApplicationUpdate),
		IncludeHealthWarnings: types.BoolValue(notification.IncludeHealthWarnings),
		ID:                    types.Int64Value(notification.ID),
		Name:                  types.StringValue(notification.Name),
	}
	genericNotification.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	genericNotification.writeFields(ctx, notification.Fields)
	n.fromNotification(&genericNotification)
}

func (n *NotificationEmail) read(ctx context.Context) *lidarr.NotificationInput {
	var tags []int

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	return &lidarr.NotificationInput{
		OnGrab:                n.OnGrab.ValueBool(),
		OnImportFailure:       n.OnImportFailure.ValueBool(),
		OnUpgrade:             n.OnUpgrade.ValueBool(),
		OnDownloadFailure:     n.OnDownloadFailure.ValueBool(),
		OnReleaseImport:       n.OnReleaseImport.ValueBool(),
		OnHealthIssue:         n.OnHealthIssue.ValueBool(),
		OnApplicationUpdate:   n.OnApplicationUpdate.ValueBool(),
		IncludeHealthWarnings: n.IncludeHealthWarnings.ValueBool(),
		ConfigContract:        notificationEmailConfigContract,
		Implementation:        notificationEmailImplementation,
		ID:                    n.ID.ValueInt64(),
		Name:                  n.Name.ValueString(),
		Tags:                  tags,
		Fields:                n.toNotification().readFields(ctx),
	}
}
