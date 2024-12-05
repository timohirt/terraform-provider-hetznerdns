package hetznerdns

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/timohirt/terraform-provider-hetznerdns/hetznerdns/api"
)

// Provider creates and return a Terraform resource provider
// for Hetzern DNS
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apitoken": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HETZNER_DNS_API_TOKEN", nil),
				Description: "The API access token to authenticate at Hetzner DNS API.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hetznerdns_zone":           resourceZone(),
			"hetznerdns_record":         resourceRecord(),
			"hetznerdns_primary_server": resourcePrimaryServer(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hetznerdns_zone":    dataSourceHetznerDNSZone(),
			"hetznerdns_records": dataSourceHetznerDNSRecords(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(c context.Context, r *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return api.NewClient(r.Get("apitoken").(string))
}
