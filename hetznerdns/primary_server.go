package hetznerdns

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/timohirt/terraform-provider-hetznerdns/hetznerdns/api"
	"log"
)

func resourcePrimaryServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrimaryServerCreate,
		ReadContext:   resourcePrimaryServerRead,
		UpdateContext: resourcePrimaryServerUpdate,
		DeleteContext: resourcePrimaryServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePrimaryServerCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating primary server")
	client := m.(*api.Client)

	zoneID, zoneIDNonEmpty := d.GetOk("zone_id")
	if !zoneIDNonEmpty {
		return diag.Errorf("Zone ID of primary server not set")
	}

	address, addressNonEmpty := d.GetOk("address")
	if !addressNonEmpty {
		return diag.Errorf("Address of primaryServer not set")
	}

	port, portNonEmpty := d.GetOk("port")
	if !portNonEmpty {
		return diag.Errorf("Port of primary server not set")
	}
	portInt := port.(int)

	opts := api.CreatePrimaryServerRequest{
		ZoneID:  zoneID.(string),
		Address: address.(string),
		Port:    &portInt,
	}

	record, err := client.CreatePrimaryServer(opts)
	if err != nil {
		log.Printf("[ERROR] Error creating primary server %s: %s", opts.Address, err)
		return diag.Errorf("Error creating primary server %s: %s", opts.Address, err)
	}

	d.SetId(record.ID)
	return resourcePrimaryServerRead(c, d, m)
}

func resourcePrimaryServerRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading primary server")
	client := m.(*api.Client)

	id := d.Id()
	record, err := client.GetPrimaryServer(id)
	if err != nil {
		return diag.Errorf("Error getting primary server with id %s: %s", id, err)
	}

	if record == nil {
		log.Printf("[WARN] Primary server with id %s doesn't exist, removing it from state", id)
		d.SetId("")
		return nil
	}

	d.SetId(record.ID)
	d.Set("address", record.Address)
	d.Set("zone_id", record.ZoneID)
	d.Set("port", record.Port)

	return nil
}

func resourcePrimaryServerUpdate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating primary server")
	client := m.(*api.Client)

	id := d.Id()
	record, err := client.GetPrimaryServer(id)
	if err != nil {
		return diag.Errorf("Error getting primary server with id %s: %s", id, err)
	}

	if record == nil {
		log.Printf("[WARN] Primary server with id %s doesn't exist, removing it from state", id)
		d.SetId("")
		return nil
	}

	if d.HasChanges("address", "port") {
		record.Address = d.Get("address").(string)
		port := d.Get("port").(int)
		record.Port = &port

		record, err = client.UpdatePrimaryServer(*record)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRecordRead(c, d, m)
}

func resourcePrimaryServerDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting resource record")

	client := m.(*api.Client)
	recordID := d.Id()

	err := client.DeletePrimaryServer(recordID)
	if err != nil {
		log.Printf("[ERROR] Error deleting primary server %s: %s", recordID, err)
		return diag.FromErr(err)
	}

	return nil
}
