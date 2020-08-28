# hetznerdns_record Resource

Provides a Hetzner DNS Recrods resource to create, update and delete DNS Records.

## Example Usage

```
data "hetznerdns_zone" "zone1" {
    name = "zone1.online"
}

resource "hetznerdns_record" "www" {
    zone_id = hetznerdns_zone.z1.id
    name = "www"
    value = "192.168.1.1"
    type = "A"
    ttl= 60
}
```

## Argument Reference

The following arguments are supported:

- `zone_id` - (Required, string) Id of the DNS zone to create
  the record in. 

- `name` - (Required, string) Name of the DNS record to create. 

- `value` - (Required, string) The value of the record (eg. 192.168.1.1).

- `type` - (Required, string) The type of the record.

- `ttl` - (Optional, int) Time to live of this record.
