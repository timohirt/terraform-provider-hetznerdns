# hetznerdns_zone Resource

Provides a Hetzner DNS Zone resource to create, update and delete DNS Zones.

## Example Usage

```hcl
resource "hetznerdns_zone" "zone1" {
    name = "zone1.online"
    ttl = 3600
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) Name of the DNS zone to create. Must be a valid
  domain with top level domain. Meaning `<domain>.de` or `<domain>.io`. Don't
  include sub domains on this level. So, no `sub.<domain>.io`. The Hetzner API
  rejects attempts to create a zone with a sub domain name. Use a record to
  create the sub domain.

- `ttl` - (Required, int) Time to live of this zone.

## Import

A Zone can be imported using its `id`. Log in to the Hetzner DNS web frontend,
navigate to the zone you want to import, and copy the id from the URL in your
browser.

```
terraform import hetznerdns_zone.zone1 rMu2waTJPbHr4
```
