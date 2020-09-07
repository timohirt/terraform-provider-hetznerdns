# hetznerdns_zone Data Source

Provides details about a Hetzner DNS Zone.

## Example Usage

```hcl
data "hetznerdns_zone" "zone1" {
	name = "zone1.online"
}
```

## Argument Reference

- `id` - (Required, string) The ID of the DNS zone.

- `name` - (Required, string) Name of the DNS zone to get data from. 

- `ttl` - (Required, int) Time to live of this zone.