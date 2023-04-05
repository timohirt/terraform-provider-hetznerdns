package hetznerdns

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/timohirt/terraform-provider-hetznerdns/hetznerdns/api"
)

func dataSourceHetznerDNSRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHetznerDNSRecordsRead,
		Schema: map[string]*schema.Schema{
			"records": {
				Description: "list of records",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceHetznerDNSRecordsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*api.Client)
	zoneID, isNonZeroZoneID := d.GetOk("zone_id")
	if !isNonZeroZoneID {
		return fmt.Errorf("Data source records has no 'name' set")
	}

	records, err := client.GetRecordsByZoneID(zoneID.(string))
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error getting zone state. %s", err)
	}

	recordsList := []interface{}{}
	for _, record := range *records {
		values := map[string]interface{}{
			"zone_id": record.ZoneID,
			"type":    record.Type,
			"name":    record.Name,
			"value":   record.Value,
			"ttl":     record.TTL,
		}
		recordsList = append(recordsList, values)
	}

	d.Set("records", recordsList)
	d.SetId(zoneID.(string))

	return nil
}
