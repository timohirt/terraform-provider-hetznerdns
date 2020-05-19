package hetznerdns

import (
	"fmt"
	"log"

	"github.com/timohirt/terraform-provider-hetznerdns/hetznerdns/api"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceRecordCreate,
		Read:   resourceRecordRead,
		Update: resourceRecordUpdate,
		Delete: resourceRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Required: true,
			},
		},
	}
}

func resourceRecordCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Updaing resource record")
	client := m.(*api.Client)

	zoneID, zoneIDNonEmpty := d.GetOk("zone_id")
	if !zoneIDNonEmpty {
		return fmt.Errorf("Zone ID of record not set")
	}

	name, nameNonEmpty := d.GetOk("name")
	if !nameNonEmpty {
		return fmt.Errorf("Name of record not set")
	}

	recordType, typeNonEmpty := d.GetOk("type")
	if !typeNonEmpty {
		return fmt.Errorf("Type of record not set")
	}

	tTL, tTLNonEmpty := d.GetOk("ttl")
	if !tTLNonEmpty {
		return fmt.Errorf("TTL of record not set")
	}

	value, valueNonEmpty := d.GetOk("value")
	if !valueNonEmpty {
		return fmt.Errorf("Value of record not set")
	}

	opts := api.CreateRecordOpts{
		ZoneID: zoneID.(string),
		Name:   name.(string),
		Type:   recordType.(string),
		Value:  value.(string),
		TTL:    tTL.(int),
	}

	record, err := client.CreateRecord(opts)
	if err != nil {
		log.Printf("[ERROR] Error Creating DNs record %s: %s", opts.Name, err)
		return fmt.Errorf("Error creating DNS record %s: %s", opts.Name, err)
	}

	d.SetId(record.ID)
	return resourceRecordRead(d, m)
}

func resourceRecordRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Reading resource record")
	client := m.(*api.Client)

	id := d.Id()
	record, err := client.GetRecord(id)
	if err != nil {
		return fmt.Errorf("Error getting record with id %s: %s", id, err)
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
	d.Set("ttl", record.TTL)
	d.Set("value", record.Value)

	return nil
}

func resourceRecordUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Updating resource record")
	client := m.(*api.Client)

	id := d.Id()
	record, err := client.GetRecord(id)
	if err != nil {
		return fmt.Errorf("Error getting record with id %s: %s", id, err)
	}

	if record == nil {
		log.Printf("[WARN] DNS record with id %s doesn't exist, removing it from state", id)
		d.SetId("")
		return nil
	}

	if d.HasChanges("name", "ttl", "type", "value") {
		record.Name = d.Get("name").(string)
		record.TTL = d.Get("ttl").(int)
		record.Type = d.Get("type").(string)
		record.Value = d.Get("value").(string)

		record, err = client.UpdateRecord(*record)
		if err != nil {
			return err
		}
	}

	return resourceRecordRead(d, m)
}

func resourceRecordDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Deleting resource record")

	client := m.(*api.Client)
	recordID := d.Id()

	err := client.DeleteRecord(recordID)
	if err != nil {
		log.Printf("[ERROR] Error deleting record %s: %s", recordID, err)
		return err
	}

	return nil
}
