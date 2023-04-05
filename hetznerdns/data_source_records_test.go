package hetznerdns

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("data_source_records", &resource.Sweeper{
		Name: "hetznerdns_records_data_source",
	})
}

func TestAccHcloudDataSourceRecords(t *testing.T) {
	// aName must be a valid DNS domain name with an existing TLD
	aName := fmt.Sprintf("%s.online", acctest.RandString(10))
	aZoneTTL := 60
	aRecordValue := "192.168.1.1"
	aRecordName := acctest.RandString(10)
	aRecordType := "A"
	aRecordTTL := aZoneTTL * 2
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccAPITokenPresent(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordsDataSourceConfig(aName, aZoneTTL, aRecordValue, aRecordName, aRecordType, aRecordTTL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hetznerdns_zone.zone1", "name", aName),
					resource.TestCheckResourceAttr(
						"data.hetznerdns_records.zone1", "records.0.value", aRecordValue),
					resource.TestCheckResourceAttr(
						"data.hetznerdns_records.zone1", "records.0.name", aRecordName),
					resource.TestCheckResourceAttr(
						"data.hetznerdns_records.zone1", "records.0.type", aRecordType),
					resource.TestCheckResourceAttr(
						"data.hetznerdns_records.zone1", "records.0.ttl", strconv.Itoa(aRecordTTL)),
				),
			},
		},
	})
}

func testAccRecordsDataSourceConfig(name string, ttl int, recordValue string, recordName string, recordType string, recordTTL int) string {
	return fmt.Sprintf(`
resource "hetznerdns_zone" "zone1" {
	name = "%s"
	ttl  = %d
}

resource "hetznerdns_record" "record1" {
	zone_id = hetznerdns_zone.zone1.id
	name    = "%s"
	type    = "%s"
	value   = "%s"
	ttl     = %d
}

data "hetznerdns_records" "zone1" {
	depends_on = [hetznerdns_record.record1]

	zone_id = hetznerdns_zone.zone1.id
}

`, name, ttl, recordName, recordType, recordValue, recordTTL)
}
