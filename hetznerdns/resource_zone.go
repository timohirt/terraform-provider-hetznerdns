package hetznerdns

import (
	"log"

	"github.com/timohirt/terraform-provider-hetznerdns/v2/hetznerdns/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneCreate,
		Read:   resourceZoneRead,
		Update: resourceZoneUpdate,
		Delete: resourceZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceZoneCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Creating resource zone")

	client := m.(*api.Client)

	var opts api.CreateZoneOpts
	if name, isOk := d.GetOk("name"); isOk {
		opts.Name = name.(string)
	}

	if ttl, isOk := d.GetOk("ttl"); isOk {
		opts.TTL = ttl.(int)
	}

	resp, err := client.CreateZone(opts)
	if err != nil {
		log.Printf("[ERROR] Creating resource zone failed: %s", err)
		d.SetId("")
		return err
	}
	d.SetId(resp.ID)

	return resourceZoneRead(d, m)
}

func resourceZoneRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Reading resource zone")
	client := m.(*api.Client)
	zoneID := d.Id()
	zone, err := client.GetZone(zoneID)
	if err != nil {
		log.Printf("[ERROR] Reading resource zone failed: %s", err)
		return err
	}

	if zone == nil {
		log.Printf("[WARN] DNS zone with id %s doesn't exist, removing it from state", zoneID)
		d.SetId("")
		return nil
	}

	d.Set("name", zone.Name)
	d.Set("ttl", zone.TTL)

	return nil
}

func resourceZoneUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Updating resource zone")
	client := m.(*api.Client)
	zoneID := d.Id()
	zone, err := client.GetZone(zoneID)
	if err != nil {
		return err
	}

	if zone == nil {
		log.Printf("[WARN] DNS zone with id %s doesn't exist, removing it from state", zoneID)
		d.SetId("")
		return nil
	}

	d.Partial(true)
	if d.HasChange("ttl") {
		zone.TTL = d.Get("ttl").(int)
		zone, err = client.UpdateZone(*zone)
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceZoneRead(d, m)
}

func resourceZoneDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Deleting resource zone")

	client := m.(*api.Client)
	zoneID := d.Id()

	err := client.DeleteZone(zoneID)
	if err != nil {
		log.Printf("[ERROR] Error deleting zone %s: %s", zoneID, err)
		return err
	}

	return nil
}
