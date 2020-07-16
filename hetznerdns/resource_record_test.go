package hetznerdns

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRecordResources(t *testing.T) {
	// aZoneName must be a valid DNS domain name with an existing TLD
	aZoneName := fmt.Sprintf("%s.online", acctest.RandString(10))

	aValue := "192.168.1.1"
	aName := acctest.RandString(10)
	aType := "A"
	aTTL := 60

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccAPITokenPresent(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccRecordResourceConfigCreate(aZoneName, aName, aType, aValue, aTTL),
				PreventDiskCleanup: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hetznerdns_record.record1", "id"),
					resource.TestCheckResourceAttr(
						"hetznerdns_record.record1", "type", aType),
					resource.TestCheckResourceAttr(
						"hetznerdns_record.record1", "name", aName),
					resource.TestCheckResourceAttr(
						"hetznerdns_record.record1", "value", aValue),
					resource.TestCheckResourceAttr(
						"hetznerdns_record.record1", "ttl", strconv.Itoa(aTTL)),
				),
			},
		},
	})
}

func testAccRecordResourceConfigCreate(aZoneName string, aName string, aType string, aValue string, aTTL int) string {
	return fmt.Sprintf(`
resource "hetznerdns_zone" "zone1" {
	name = "%s"
	ttl = %d
}

resource "hetznerdns_record" "record1" {
	zone_id = "${hetznerdns_zone.zone1.id}"
	type = "%s"
	name = "%s"
	value = "%s"
	ttl = %d
}
`, aZoneName, aTTL, aType, aName, aValue, aTTL)
}

func TestAccTwoRecordResources(t *testing.T) {
	// aZoneName must be a valid DNS domain name with an existing TLD
	aZoneName := fmt.Sprintf("%s.online", acctest.RandString(10))

	aValue := "192.168.1.1"
	anotherValue := "192.168.1.2"
	aName := acctest.RandString(10)
	anotherName := acctest.RandString(10)
	aType := "A"
	aTTL := 60

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccAPITokenPresent(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccRecordResourceConfigCreateTwo(aZoneName, aName, anotherName, aType, aValue, anotherValue, aTTL),
				PreventDiskCleanup: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hetznerdns_record.record1", "id"),
					resource.TestCheckResourceAttrSet(
						"hetznerdns_record.record2", "id"),
				),
			},
		},
	})
}

func testAccRecordResourceConfigCreateTwo(aZoneName string, aName string, anotherName string, aType string, aValue string, anotherValue string, aTTL int) string {
	return fmt.Sprintf(`
resource "hetznerdns_zone" "zone1" {
	name = "%s"
	ttl = %d
}

resource "hetznerdns_record" "record1" {
	zone_id = "${hetznerdns_zone.zone1.id}"
	type = "%s"
	name = "%s"
	value = "%s"
	ttl = %d
}

resource "hetznerdns_record" "record2" {
	zone_id = "${hetznerdns_zone.zone1.id}"
	type = "%s"
	name = "%s"
	value = "%s"
	ttl = %d
}
`, aZoneName, aTTL, aType, aName, aValue, aTTL, aType, anotherName, anotherValue, aTTL)
}
