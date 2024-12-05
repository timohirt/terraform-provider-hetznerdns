package hetznerdns

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const validPublicIpAddress = "167.182.9.17"

func TestAccPrimaryServerResources(t *testing.T) {
	// aZoneName must be a valid DNS domain name with an existing TLD
	aZoneName := fmt.Sprintf("%s.online", acctest.RandString(10))
	aZoneTTL := 60

	psAddress := validPublicIpAddress
	psPort := 53

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccAPITokenPresent(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccPrimaryServerResourceConfigCreate(aZoneName, aZoneTTL, psAddress, psPort),
				PreventDiskCleanup: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hetznerdns_primary_server.ps1", "id"),
					resource.TestCheckResourceAttr(
						"hetznerdns_primary_server.ps1", "address", psAddress),
					resource.TestCheckResourceAttr(
						"hetznerdns_primary_server.ps1", "port", strconv.Itoa(psPort)),
				),
			},
		},
	})
}

func testAccPrimaryServerResourceConfigCreate(aZoneName string, aZoneTTL int, psAddress string, psPort int) string {
	return fmt.Sprintf(`
resource "hetznerdns_zone" "zone1" {
	name = "%s"
	ttl = %d
}

resource "hetznerdns_primary_server" "ps1" {
	zone_id = "${hetznerdns_zone.zone1.id}"
	address = "%s"
	port    = %d
}
`, aZoneName, aZoneTTL, psAddress, psPort)
}

func TestAccTwoPrimaryServersResources(t *testing.T) {
	// aZoneName must be a valid DNS domain name with an existing TLD
	aZoneName := fmt.Sprintf("%s.online", acctest.RandString(10))
	aZoneTTL := 60

	ps1Address := validPublicIpAddress
	ps1Port := 53

	ps2Address := "154.23.82.134"
	ps2Port := 53

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccAPITokenPresent(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccPrimaryServerResourceConfigCreateTwo(aZoneName, aZoneTTL, ps1Address, ps1Port, ps2Address, ps2Port),
				PreventDiskCleanup: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hetznerdns_primary_server.ps1", "id"),
					resource.TestCheckResourceAttrSet(
						"hetznerdns_primary_server.ps2", "id"),
				),
			},
		},
	})
}

func testAccPrimaryServerResourceConfigCreateTwo(aZoneName string, aTTL int, ps1Address string, ps1Port int, ps2Address string, ps2Port int) string {
	return fmt.Sprintf(`
resource "hetznerdns_zone" "zone1" {
	name = "%s"
	ttl = %d
}

resource "hetznerdns_primary_server" "ps1" {
	zone_id = "${hetznerdns_zone.zone1.id}"
	address = "%s"
	port    = %d
}

resource "hetznerdns_primary_server" "ps2" {
	zone_id = "${hetznerdns_zone.zone1.id}"
	address = "%s"
	port    = %d
}
`, aZoneName, aTTL, ps1Address, ps1Port, ps2Address, ps2Port)
}
