# Terraform Provider for Hetzner DNS

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/) 1.14 (to build the provider plugin)

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

- `name` - (Required, string) Name of the DNS zone to create. 
  Must be a valid domain with top level domain (eg. .de, .com, .io)

- `ttl` - (Required, int) Time to live of this zone.

#### Import

A Zone can be imported using its `id`. Log in to the Hetzner DNS web frontend,
navigate to the zone you want to import, and copy the id from the URL in your
browser.

```
terraform import hetznerdns_zone.zone1 <id>
```

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