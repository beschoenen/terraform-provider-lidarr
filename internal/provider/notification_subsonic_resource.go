package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	notificationSubsonicResourceName   = "notification_subsonic"
	notificationSubsonicImplementation = "Xbmc"
	notificationSubsonicConfigContract = "XbmcSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationSubsonicResource{}
	_ resource.ResourceWithImportState = &NotificationSubsonicResource{}
)

func NewNotificationSubsonicResource() resource.Resource {
	return &NotificationSubsonicResource{}
}

// NotificationSubsonicResource defines the notification implementation.
type NotificationSubsonicResource struct {
	client *lidarr.APIClient
}

// NotificationSubsonic describes the notification data model.
type NotificationSubsonic struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Host                  types.String `tfsdk:"host"`
	Name                  types.String `tfsdk:"name"`
	Username              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	URLBase               types.String `tfsdk:"url_base"`
	Port                  types.Int64  `tfsdk:"port"`
	ID                    types.Int64  `tfsdk:"id"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	UseSSL                types.Bool   `tfsdk:"use_ssl"`
	Notify                types.Bool   `tfsdk:"notify"`
	UpdateLibrary         types.Bool   `tfsdk:"update_library"`
	OnReleaseImport       types.Bool   `tfsdk:"on_release_import"`
	OnTrackRetag          types.Bool   `tfsdk:"on_track_retag"`
	OnRename              types.Bool   `tfsdk:"on_rename"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnUpgrade             types.Bool   `tfsdk:"on_upgrade"`
}

func (n NotificationSubsonic) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		Port:                  n.Port,
		Host:                  n.Host,
		URLBase:               n.URLBase,
		Password:              n.Password,
		Username:              n.Username,
		Name:                  n.Name,
		ID:                    n.ID,
		UseSSL:                n.UseSSL,
		Notify:                n.Notify,
		UpdateLibrary:         n.UpdateLibrary,
		OnGrab:                n.OnGrab,
		OnReleaseImport:       n.OnReleaseImport,
		OnRename:              n.OnRename,
		OnTrackRetag:          n.OnTrackRetag,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		OnHealthIssue:         n.OnHealthIssue,
		OnUpgrade:             n.OnUpgrade,
	}
}

func (n *NotificationSubsonic) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Port = notification.Port
	n.URLBase = notification.URLBase
	n.Host = notification.Host
	n.Password = notification.Password
	n.Username = notification.Username
	n.Name = notification.Name
	n.ID = notification.ID
	n.UseSSL = notification.UseSSL
	n.Notify = notification.Notify
	n.UpdateLibrary = notification.UpdateLibrary
	n.OnGrab = notification.OnGrab
	n.OnReleaseImport = notification.OnReleaseImport
	n.OnTrackRetag = notification.OnTrackRetag
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnRename = notification.OnRename
	n.OnUpgrade = notification.OnUpgrade
}

func (r *NotificationSubsonicResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationSubsonicResourceName
}

func (r *NotificationSubsonicResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Subsonic resource.\nFor more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Subsonic](https://wiki.servarr.com/lidarr/supported#xbmc).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Required:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Required:            true,
			},
			"on_rename": schema.BoolAttribute{
				MarkdownDescription: "On rename flag.",
				Required:            true,
			},
			"on_track_retag": schema.BoolAttribute{
				MarkdownDescription: "On movie file delete flag.",
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
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationSubsonic name.",
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
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"notify": schema.BoolAttribute{
				MarkdownDescription: "Notification flag.",
				Optional:            true,
				Computed:            true,
			},
			"update_library": schema.BoolAttribute{
				MarkdownDescription: "Update library flag.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "URL base.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Required:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Host.",
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
		},
	}
}

func (r *NotificationSubsonicResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationSubsonicResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationSubsonic

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationSubsonic
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationSubsonicResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationSubsonicResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSubsonicResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationSubsonic

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationSubsonic current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(int(notification.ID.ValueInt64()))).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationSubsonicResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationSubsonicResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSubsonicResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationSubsonic

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationSubsonic
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationSubsonicResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationSubsonicResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSubsonicResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationSubsonic

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationSubsonic current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationSubsonicResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationSubsonicResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationSubsonicResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationSubsonicResourceName+": "+req.ID)
}

func (n *NotificationSubsonic) write(ctx context.Context, notification *lidarr.NotificationResource) {
	genericNotification := Notification{
		OnGrab:                types.BoolValue(notification.GetOnGrab()),
		OnUpgrade:             types.BoolValue(notification.GetOnUpgrade()),
		OnRename:              types.BoolValue(notification.GetOnRename()),
		OnTrackRetag:          types.BoolValue(notification.GetOnTrackRetag()),
		OnReleaseImport:       types.BoolValue(notification.GetOnReleaseImport()),
		OnHealthIssue:         types.BoolValue(notification.GetOnHealthIssue()),
		IncludeHealthWarnings: types.BoolValue(notification.GetIncludeHealthWarnings()),
		ID:                    types.Int64Value(int64(notification.GetId())),
		Name:                  types.StringValue(notification.GetName()),
	}
	genericNotification.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	genericNotification.writeFields(ctx, notification.Fields)
	n.fromNotification(&genericNotification)
}

func (n *NotificationSubsonic) read(ctx context.Context) *lidarr.NotificationResource {
	tags := make([]*int32, len(n.Tags.Elements()))
	tfsdk.ValueAs(ctx, n.Tags, &tags)

	notification := lidarr.NewNotificationResource()
	notification.SetOnGrab(n.OnGrab.ValueBool())
	notification.SetOnUpgrade(n.OnUpgrade.ValueBool())
	notification.SetOnRename(n.OnRename.ValueBool())
	notification.SetOnTrackRetag(n.OnTrackRetag.ValueBool())
	notification.SetOnReleaseImport(n.OnReleaseImport.ValueBool())
	notification.SetOnHealthIssue(n.OnHealthIssue.ValueBool())
	notification.SetIncludeHealthWarnings(n.IncludeHealthWarnings.ValueBool())
	notification.SetConfigContract(notificationSubsonicConfigContract)
	notification.SetImplementation(notificationSubsonicImplementation)
	notification.SetId(int32(n.ID.ValueInt64()))
	notification.SetName(n.Name.ValueString())
	notification.SetTags(tags)
	notification.SetFields(n.toNotification().readFields(ctx))

	return notification
}
