# Releasing and Publishing Provider

A release of this provider includes publishing the new version
at the Terraform registry. Therefore, a file containing hashes
of the binaries must be singed with my private GPG key.

First create a new tag and push it to GitHub. This starts a
workflow which drafts a release.

```bash
$ git tag v1.0.17
$ git push --tags
```

Then navigate your browser to the GitHub [releases page of this repository](https://github.com/timohirt/terraform-provider-hetznerdns/releases). 

Download `terraform-provider-hetznerdns_1.0.17_SHA256SUMS` from the
list of assets associated with this release. The version might by
different.

Now create a signature file using gpg.

```bash
$ gpg --detach-sign ./terraform-provider-hetznerdns_1.0.17_SHA256SUMS
```

This create another file `terraform-provider-hetznerdns_1.0.17_SHA256SUMS.sig`.
Add it to the assets of the current draft release.

Then publish the release. Terraform registry adds the new release automatically.

