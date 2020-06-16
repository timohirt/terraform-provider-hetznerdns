# Terraform Provider for Hetzner DNS

![CI Build](https://github.com/timohirt/terraform-provider-hetznerdns/workflows/CI%20Build/badge.svg?branch=master)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/timohirt/terraform-provider-hetznerdns)
![GitHub](https://img.shields.io/github/license/timohirt/terraform-provider-hetznerdns)

Read about what I learnt while [implementing this Terraform Provider](http://www.timohirt.de/blog/implementing-a-terraform-provider/).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/) 1.14 (to build the provider plugin)

## Installing and Using this Plugin

See [INSTALL.md](./INSTALL.md).

## Resources

### hetznerdns_zone

Provides a Hetzner DNS Zone resource to create, update and delete DNS Zones.

#### Example Usage

```
resource "hetznerdns_zone" "zone1" {
    name = "zone1.online"
    ttl = 3600
}
```

#### Argument Reference

The following arguments are supported:

- `name` - (Required, string) Name of the DNS zone to create. Must be a valid
  domain with top level domain. Meaning `<domain>.de` or `<domain>.io`. Don't
  include sub domains on this level. So, no `sub.<domain>.io`. The Hetzner API
  rejects attempts to create a zone with a sub domain name. Use a record to
  create the sub domain.

- `ttl` - (Required, int) Time to live of this zone.

#### Import

A Zone can be imported using its `id`. Log in to the Hetzner DNS web frontend,
navigate to the zone you want to import, and copy the id from the URL in your
browser.

```
terraform import hetznerdns_zone.zone1 <id>
```

### hetznerdns_record

Provides a Hetzner DNS Recrods resource to create, update and delete DNS Records.

#### Example Usage

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

#### Argument Reference

The following arguments are supported:

- `zone_id` - (Required, string) Id of the DNS zone to create
  the record in. 

- `name` - (Required, string) Name of the DNS record to create. 
  Must be a valid domain with top level domain (eg. .de, .com, .io).

- `value` - (Required, string) The value of the record (eg. 192.168.1.1).

- `type` - (Required, string) The type of the record.

- `ttl` - (Required, int) Time to live of this record.

#### Import

tbd

## Data Sources

### hetznerdns_zone

Provides details about a Hetzner DNS Zone.

#### Example Usage

```
data "hetznerdns_zone" "zone1" {
	name = "zone1.online"
}
```

#### Argument Reference

- `id` - (Required, string) The ID of the DNS zone.

- `name` - (Required, string) Name of the DNS zone to create. 
  Must be a valid domain with top level domain (eg. .de, .com, .io)

- `ttl` - (Required, int) Time to live of this zone.
