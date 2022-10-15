package hetznerdns

import (
	"fmt"

	"github.com/timohirt/terraform-provider-hetznerdns/v2/hetznerdns/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHetznerDNSZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHetznerDNSZoneRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceHetznerDNSZoneRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*api.Client)
	name, isNonZeroName := d.GetOk("name")
	if !isNonZeroName {
		return fmt.Errorf("Data source zone has no 'name' set")
	}

	zone, err := client.GetZoneByName(name.(string))
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error getting zone state. %s", err)
	}

	if zone == nil {
		return fmt.Errorf("DNS zone '%s' doesn't exist", name.(string))
	}

	d.Set("name", zone.Name)
	d.Set("ttl", zone.TTL)
	d.SetId(zone.ID)

	return nil
}
