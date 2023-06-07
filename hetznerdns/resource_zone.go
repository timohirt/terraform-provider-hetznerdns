package hetznerdns

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/timohirt/terraform-provider-hetznerdns/hetznerdns/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZoneCreate,
		ReadContext:   resourceZoneRead,
		UpdateContext: resourceZoneUpdate,
		DeleteContext: resourceZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceZoneCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating resource zone")

	client := m.(*api.Client)

	var opts api.CreateZoneOpts
	if name, isOk := d.GetOk("name"); isOk {
		opts.Name = name.(string)
	}

	if ttl, isOk := d.GetOk("ttl"); isOk {
		opts.TTL = ttl.(int)
	}

	zoneExists, err := client.ZoneExistsByName(opts.Name)
	if err != nil {
		log.Printf("[ERROR] Checking if resource zone exists failed: %s", err)
		d.SetId("")
		return diag.FromErr(err)
	}

	if zoneExists {
		errMsg := fmt.Sprintf("DNS zone with name %s already exists", opts.Name)
		log.Printf("[ERROR] %s", errMsg)
		return diag.Errorf(errMsg)
	}

	resp, err := client.CreateZone(opts)
	if err != nil {
		log.Printf("[ERROR] Creating resource zone failed: %s", err)
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(resp.ID)

	return resourceZoneRead(c, d, m)
}

func resourceZoneRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading resource zone")
	client := m.(*api.Client)
	zoneID := d.Id()
	zone, err := client.GetZone(zoneID)
	if err != nil {
		log.Printf("[ERROR] Reading resource zone failed: %s", err)
		return diag.FromErr(err)
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

func resourceZoneUpdate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating resource zone")
	client := m.(*api.Client)
	zoneID := d.Id()
	zone, err := client.GetZone(zoneID)
	if err != nil {
		return diag.FromErr(err)
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
			return diag.FromErr(err)
		}
	}

	d.Partial(false)

	return resourceZoneRead(c, d, m)
}

func resourceZoneDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting resource zone")

	client := m.(*api.Client)
	zoneID := d.Id()

	err := client.DeleteZone(zoneID)
	if err != nil {
		log.Printf("[ERROR] Error deleting zone %s: %s", zoneID, err)
		return diag.FromErr(err)
	}

	return nil
}
