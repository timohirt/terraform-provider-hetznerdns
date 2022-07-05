package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	plugin "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/timohirt/terraform-provider-hetznerdns/v2/hetznerdns"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return hetznerdns.Provider()
		},
	})
}
