package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/lidarr"
)

const notificationsDataSourceName = "notifications"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NotificationsDataSource{}

func NewNotificationsDataSource() datasource.DataSource {
	return &NotificationsDataSource{}
}

// NotificationsDataSource defines the notifications implementation.
type NotificationsDataSource struct {
	client *lidarr.Lidarr
}

// Notifications describes the notifications data model.
type Notifications struct {
	Notifications types.Set    `tfsdk:"notifications"`
	ID            types.String `tfsdk:"id"`
}

func (d *NotificationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationsDataSourceName
}

func (d *NotificationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Notifications -->List all available [Notifications](../resources/notification).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"notifications": schema.SetNestedAttribute{
				MarkdownDescription: "Notification list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"on_grab": schema.BoolAttribute{
							MarkdownDescription: "On grab flag.",
							Computed:            true,
						},
						"on_release_import": schema.BoolAttribute{
							MarkdownDescription: "On release import flag.",
							Computed:            true,
						},
						"on_upgrade": schema.BoolAttribute{
							MarkdownDescription: "On upgrade flag.",
							Computed:            true,
						},
						"on_rename": schema.BoolAttribute{
							MarkdownDescription: "On rename flag.",
							Computed:            true,
						},
						"on_download_failure": schema.BoolAttribute{
							MarkdownDescription: "On download failure flag.",
							Computed:            true,
						},
						"on_import_failure": schema.BoolAttribute{
							MarkdownDescription: "On import failure flag.",
							Computed:            true,
						},
						"on_track_retag": schema.BoolAttribute{
							MarkdownDescription: "On track retag.",
							Computed:            true,
						},
						"on_health_issue": schema.BoolAttribute{
							MarkdownDescription: "On health issue flag.",
							Computed:            true,
						},
						"on_application_update": schema.BoolAttribute{
							MarkdownDescription: "On application update flag.",
							Computed:            true,
						},
						"include_health_warnings": schema.BoolAttribute{
							MarkdownDescription: "Include health warnings.",
							Computed:            true,
						},
						"config_contract": schema.StringAttribute{
							MarkdownDescription: "Notification configuration template.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "Notification implementation name.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Notification name.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Notification ID.",
							Computed:            true,
						},
						// Field values
						"always_update": schema.BoolAttribute{
							MarkdownDescription: "Always update flag.",
							Computed:            true,
						},
						"clean_library": schema.BoolAttribute{
							MarkdownDescription: "Clean library flag.",
							Computed:            true,
						},
						"direct_message": schema.BoolAttribute{
							MarkdownDescription: "Direct message flag.",
							Computed:            true,
						},
						"notify": schema.BoolAttribute{
							MarkdownDescription: "Notify flag.",
							Computed:            true,
						},
						"require_encryption": schema.BoolAttribute{
							MarkdownDescription: "Require encryption flag.",
							Computed:            true,
						},
						"send_silently": schema.BoolAttribute{
							MarkdownDescription: "Add silently flag.",
							Computed:            true,
						},
						"update_library": schema.BoolAttribute{
							MarkdownDescription: "Update library flag.",
							Computed:            true,
						},
						"use_eu_endpoint": schema.BoolAttribute{
							MarkdownDescription: "Use EU endpoint flag.",
							Computed:            true,
						},
						"use_ssl": schema.BoolAttribute{
							MarkdownDescription: "Use SSL flag.",
							Computed:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port.",
							Computed:            true,
						},
						"grab_fields": schema.Int64Attribute{
							MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart.",
							Computed:            true,
						},
						"import_fields": schema.Int64Attribute{
							MarkdownDescription: "Import fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Codecs, `5` Group, `6` Size, `7` Languages, `8` Subtitles, `9` Links, `10` Release, `11` Poster, `12` Fanart.",
							Computed:            true,
						},
						"method": schema.Int64Attribute{
							MarkdownDescription: "Method. `1` POST, `2` PUT.",
							Computed:            true,
						},
						"priority": schema.Int64Attribute{
							MarkdownDescription: "Priority.", // TODO: add values in description
							Computed:            true,
						},
						"access_token": schema.StringAttribute{
							MarkdownDescription: "Access token.",
							Computed:            true,
						},
						"access_token_secret": schema.StringAttribute{
							MarkdownDescription: "Access token secret.",
							Computed:            true,
						},
						"api_key": schema.StringAttribute{
							MarkdownDescription: "API key.",
							Computed:            true,
						},
						"app_token": schema.StringAttribute{
							MarkdownDescription: "App token.",
							Computed:            true,
						},
						"arguments": schema.StringAttribute{
							MarkdownDescription: "Arguments.",
							Computed:            true,
						},
						"author": schema.StringAttribute{
							MarkdownDescription: "Author.",
							Computed:            true,
						},
						"auth_token": schema.StringAttribute{
							MarkdownDescription: "Auth token.",
							Computed:            true,
						},
						"auth_user": schema.StringAttribute{
							MarkdownDescription: "Auth user.",
							Computed:            true,
						},
						"avatar": schema.StringAttribute{
							MarkdownDescription: "Avatar.",
							Computed:            true,
						},
						"bcc": schema.StringAttribute{
							MarkdownDescription: "BCC.",
							Computed:            true,
						},
						"bot_token": schema.StringAttribute{
							MarkdownDescription: "Bot token.",
							Computed:            true,
						},
						"cc": schema.StringAttribute{
							MarkdownDescription: "CC.",
							Computed:            true,
						},
						"channel": schema.StringAttribute{
							MarkdownDescription: "Channel.",
							Computed:            true,
						},
						"chat_id": schema.StringAttribute{
							MarkdownDescription: "Chat ID.",
							Computed:            true,
						},
						"consumer_key": schema.StringAttribute{
							MarkdownDescription: "Consumer key.",
							Computed:            true,
						},
						"consumer_secret": schema.StringAttribute{
							MarkdownDescription: "Consumer secret.",
							Computed:            true,
						},
						"device_names": schema.StringAttribute{
							MarkdownDescription: "Device names.",
							Computed:            true,
						},
						"display_time": schema.StringAttribute{
							MarkdownDescription: "Display time.",
							Computed:            true,
						},
						"expire": schema.StringAttribute{
							MarkdownDescription: "Expire.",
							Computed:            true,
						},
						"expires": schema.StringAttribute{
							MarkdownDescription: "Expires.",
							Computed:            true,
						},
						"from": schema.StringAttribute{
							MarkdownDescription: "From.",
							Computed:            true,
						},
						"host": schema.StringAttribute{
							MarkdownDescription: "Host.",
							Computed:            true,
						},
						"icon": schema.StringAttribute{
							MarkdownDescription: "Icon.",
							Computed:            true,
						},
						"mention": schema.StringAttribute{
							MarkdownDescription: "Mention.",
							Computed:            true,
						},
						"password": schema.StringAttribute{
							MarkdownDescription: "password.",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							MarkdownDescription: "Path.",
							Computed:            true,
						},
						"refresh_token": schema.StringAttribute{
							MarkdownDescription: "Refresh token.",
							Computed:            true,
						},
						"retry": schema.StringAttribute{
							MarkdownDescription: "Retry.",
							Computed:            true,
						},
						"sender_domain": schema.StringAttribute{
							MarkdownDescription: "Sender domain.",
							Computed:            true,
						},
						"sender_id": schema.StringAttribute{
							MarkdownDescription: "Sender ID.",
							Computed:            true,
						},
						"server": schema.StringAttribute{
							MarkdownDescription: "server.",
							Computed:            true,
						},
						"sign_in": schema.StringAttribute{
							MarkdownDescription: "Sign in.",
							Computed:            true,
						},
						"sound": schema.StringAttribute{
							MarkdownDescription: "Sound.",
							Computed:            true,
						},
						"to": schema.StringAttribute{
							MarkdownDescription: "To.",
							Computed:            true,
						},
						"token": schema.StringAttribute{
							MarkdownDescription: "Token.",
							Computed:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "URL.",
							Computed:            true,
						},
						"user_key": schema.StringAttribute{
							MarkdownDescription: "User key.",
							Computed:            true,
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "Username.",
							Computed:            true,
						},
						"web_hook_url": schema.StringAttribute{
							MarkdownDescription: "Web hook url.",
							Computed:            true,
						},
						"channel_tags": schema.SetAttribute{
							MarkdownDescription: "Channel tags.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"device_ids": schema.SetAttribute{
							MarkdownDescription: "Device IDs.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"devices": schema.SetAttribute{
							MarkdownDescription: "Devices.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"recipients": schema.SetAttribute{
							MarkdownDescription: "Recipients.",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *NotificationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*lidarr.Lidarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *lidarr.Lidarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *NotificationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Notifications

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get notifications current value
	response, err := d.client.GetNotificationsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationsDataSourceName)
	// Map response body to resource schema attribute
	profiles := make([]Notification, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.Notifications.Type(context.Background()), &data.Notifications)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
