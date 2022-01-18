# Installing and Using this Plugin

This is a third-party provier, which is not maintained by Hashicorp. You have
to download and install it on your own. Pick the binary for the target operating
system from the [releases list](https://github.com/timohirt/terraform-provider-hetznerdns/releases)
and download the archive.

## Installing the Provider 

### Terraform 0.13

Terraform 0.13 changed the [location of custom providers](https://www.terraform.io/upgrade-guides/0-13.html#new-filesystem-layout-for-local-copies-of-providers), in order to work
nicly with the Terraform Registry. This means that you have to create or migrate
to a slighly different directory structure where you put the provide binaries.

```bash
$ mkdir -p ~/.terraform.d/plugins/github.com/timohirt/hetznerdns/1.0.5/linux_amd64
$ tar xzf terraform-provider-hetznerdns_1.0.5_linux_amd64.tar.gz
$ mv ./terraform-provider-hetznerdns ~/.terraform.d/plugins/github.com/timohirt/hetznerdns/1.0.5/linux_amd64/terraform-provider-hetznerdns_v1.0.5
```

As you can see above, the version as well as the operating system is now included in
the path and filename. Make sure you pick the right binaries for the os and use `darwin_amd64` 
or `win_amd64` instead of `linux_amd64` if necessary.

### Using the Provider on your Local Machine

Add the following to your `terraform.tf`.

```terraform
terraform {
  required_providers {
    hetznerdns = {
      source = "github.com/timohirt/hetznerdns"
    }
  }
  required_version = ">= 1.0"
}
```

## Testing 

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
