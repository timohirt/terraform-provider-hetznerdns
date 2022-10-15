package hetznerdns

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

	"github.com/timohirt/terraform-provider-hetznerdns/v2/hetznerdns/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordCreate,
		ReadContext:   resourceRecordRead,
		UpdateContext: resourceRecordUpdate,
		DeleteContext: resourceRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceRecordCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating resource record")
	client := m.(*api.Client)

	zoneID, zoneIDNonEmpty := d.GetOk("zone_id")
	if !zoneIDNonEmpty {
		return diag.Errorf("Zone ID of record not set")
	}

	name, nameNonEmpty := d.GetOk("name")
	if !nameNonEmpty {
		return diag.Errorf("Name of record not set")
	}

	recordType, typeNonEmpty := d.GetOk("type")
	if !typeNonEmpty {
		return diag.Errorf("Type of record not set")
	}

	value, valueNonEmpty := d.GetOk("value")
	if !valueNonEmpty {
		return diag.Errorf("Value of record not set")
	}

	opts := api.CreateRecordOpts{
		ZoneID: zoneID.(string),
		Name:   name.(string),
		Type:   recordType.(string),
		Value:  value.(string),
	}

	tTL, tTLNonEmpty := d.GetOk("ttl")
	if tTLNonEmpty {
		nonEmptyTTL := tTL.(int)
		opts.TTL = &nonEmptyTTL
	}

	record, err := client.CreateRecord(opts)
	if err != nil {
		log.Printf("[ERROR] Error creating DNS record %s: %s", opts.Name, err)
		return diag.Errorf("Error creating DNS record %s: %s", opts.Name, err)
	}

	d.SetId(record.ID)
	return resourceRecordRead(c, d, m)
}

func resourceRecordRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading resource record")
	client := m.(*api.Client)

	id := d.Id()
	record, err := client.GetRecord(id)
	if err != nil {
		return diag.Errorf("Error getting record with id %s: %s", id, err)
	}

	if record == nil {
		log.Printf("[WARN] DNS record with id %s doesn't exist, removing it from state", id)
		d.SetId("")
		return nil
	}

	d.SetId(record.ID)
	d.Set("name", record.Name)
	d.Set("zone_id", record.ZoneID)
	d.Set("type", record.Type)

	d.Set("ttl", nil)
	if record.HasTTL() {
		d.Set("ttl", record.TTL)
	}
	d.Set("value", record.Value)

	return nil
}

func resourceRecordUpdate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating resource record")
	client := m.(*api.Client)

	id := d.Id()
	record, err := client.GetRecord(id)
	if err != nil {
		return diag.Errorf("Error getting record with id %s: %s", id, err)
	}

	if record == nil {
		log.Printf("[WARN] DNS record with id %s doesn't exist, removing it from state", id)
		d.SetId("")
		return nil
	}

	if d.HasChanges("name", "ttl", "type", "value") {
		record.Name = d.Get("name").(string)

		record.TTL = nil
		ttl, ttlNonEmpty := d.GetOk("ttl")
		if ttlNonEmpty {
			ttl := ttl.(int)
			record.TTL = &ttl
		}
		record.Type = d.Get("type").(string)
		record.Value = d.Get("value").(string)

		record, err = client.UpdateRecord(*record)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRecordRead(c, d, m)
}

func resourceRecordDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting resource record")

	client := m.(*api.Client)
	recordID := d.Id()

	err := client.DeleteRecord(recordID)
	if err != nil {
		log.Printf("[ERROR] Error deleting record %s: %s", recordID, err)
		return diag.FromErr(err)
	}

	return nil
}
