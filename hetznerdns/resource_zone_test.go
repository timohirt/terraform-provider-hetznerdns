package hetznerdns

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("resource_source_zone", &resource.Sweeper{
		Name: "hetznerdns_zone_resource",
	})
}

func TestAccZoneResources(t *testing.T) {
	// aName must be a valid DNS domain name with an existing TLD
	aName := fmt.Sprintf("%s.online", acctest.RandString(10))
	aTTL := 60

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccAPITokenPresent(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccZoneResourceConfigCreate(aName, aTTL),
				PreventDiskCleanup: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hetznerdns_zone.zone1", "name", aName),
					resource.TestCheckResourceAttr(
						"hetznerdns_zone.zone1", "ttl", strconv.Itoa(aTTL)),
				),
			},
		},
	})
}

func testAccZoneResourceConfigCreate(name string, ttl int) string {
	return fmt.Sprintf(`
resource "hetznerdns_zone" "zone1" {
	name = "%s"
	ttl = %d
}
`, name, ttl)
}
