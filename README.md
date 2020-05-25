# Terraform Provider for Hetzner DNS

![CI Build](https://github.com/timohirt/terraform-provider-hetznerdns/workflows/CI%20Build/badge.svg?branch=master)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/) 1.14 (to build the provider plugin)

## Installing and Using this Plugin

This is a third-party provier, which is not maintained by Hashicorp. You have
to download and install it on your own. Pick the binary for the target operating
system from the [releases list](https://github.com/timohirt/terraform-provider-hetznerdns/releases)
and download the archive.

Terraform looks for providers in `~/.terraform.d/plugins`. Extract the archive and
copy the executable to this directory. See [Terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
for more details.

```bash
$ mkdir -p ~/.terraform.d/plugins
$ tar xzf terraform-provider-hetznerdns_1.0.0_linux_amd64.tar.gz
$ mv ./terraform-provider-hetznerdns ~/.terraform.d/plugins
```

Next add a hetznerdns resource or data source to your project and run 
terraform init.

```bash
$ cat << EOF > main.tf
resource "hetznerdns_zone" "z1" {
    name = "my.domain.tld"
    ttl = 60
}
EOF
$ terraform init

Initializing the backend...

Initializing provider plugins...

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.

$ terraform plan
provider.hetznerdns.apitoken
  The API access token to authenticate at Hetzner DNS API.

  Enter a value:
```

Enter your Hetzer DNS API key and hit enter. 

Set Hetzner DNS API key environment variable and you don't have to enter
it manually: `export HETZNER_DNS_API_TOKEN=<your api token>`

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
  Must be a valid domain with top level domain (eg. .de, .com, .io).

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
