package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientTorrentDownloadStationResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientTorrentDownloadStationResourceConfig("resourceTorrentDownloadStationTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_download_client_torrent_download_station.test", "use_ssl", "false"),
					resource.TestCheckResourceAttrSet("lidarr_download_client_torrent_download_station.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientTorrentDownloadStationResourceConfig("resourceTorrentDownloadStationTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_download_client_torrent_download_station.test", "use_ssl", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_download_client_torrent_download_station.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientTorrentDownloadStationResourceConfig(name, ssl string) string {
	return fmt.Sprintf(`
	resource "lidarr_download_client_torrent_download_station" "test" {
		enable = false
		use_ssl = %s
		priority = 1
		name = "%s"
		host = "torrent-download-station"
		port = 9091
	}`, ssl, name)
}
