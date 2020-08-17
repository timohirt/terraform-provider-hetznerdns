# hetznerdns Provider

This providers helps you automate management of DNS zones
and records at Hetzner DNS.

## Example Usage

```hcl
data "hetznerdns_zone" "dns_zone" {
    name = "example.com"
}

data "hcloud_server" "web" {
    name = "web1"
}

resource "hetznerdns_record" "web" {
    zone_id = data.hetznerdns_zone.dns_zone.id
    name = "www"
    value = hcloud_server.web.ipv4_address
    type = "A"
    ttl= 60
}
```

## Argument Reference

The following arguments are supported:

- `apitoken` - (Required, string) The Hetzner DNS API token. You can 
  pass it using the env variable `HETZNER_DNS_API_TOKEN`as well.