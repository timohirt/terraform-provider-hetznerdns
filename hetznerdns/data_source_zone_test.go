package hetznerdns

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("data_source_zone", &resource.Sweeper{
		Name: "hetznerdns_zone_data_source",
	})
}

func TestAccHcloudDataSourceDatasources(t *testing.T) {
	// aName must be a valid DNS domain name with an existing TLD
	aName := fmt.Sprintf("%s.online", acctest.RandString(10))
	aTTL := 60
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccAPITokenPresent(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneDataSourceConfig(aName, aTTL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.hetznerdns_zone.zone1", "name", aName),
					resource.TestCheckResourceAttr(
						"data.hetznerdns_zone.zone1", "ttl", strconv.Itoa(aTTL)),
					resource.TestCheckResourceAttrSet("data.hetznerdns_zone.zone1", "id"),
				),
			},
		},
	})
}

func testAccZoneDataSourceConfig(name string, ttl int) string {
	return fmt.Sprintf(`
resource "hetznerdns_zone" "zone1" {
	name = "%s"
	ttl = "%d"
}

data "hetznerdns_zone" "zone1" {
	name = "${hetznerdns_zone.zone1.name}"
}
`, name, ttl)
}
