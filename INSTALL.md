# Installing and Using this Plugin

This is a third-party provier, which is not maintained by Hashicorp. You have
to download and install it on your own. Pick the binary for the target operating
system from the [releases list](https://github.com/timohirt/terraform-provider-hetznerdns/releases)
and download the archive.

## Installing the Provider 

Terraform looks for providers in `~/.terraform.d/plugins`. Extract the archive and
copy the executable to this directory. See [Terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
for more details.

```bash
$ mkdir -p ~/.terraform.d/plugins
$ tar xzf terraform-provider-hetznerdns_1.0.0_linux_amd64.tar.gz
$ mv ./terraform-provider-hetznerdns ~/.terraform.d/plugins
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
