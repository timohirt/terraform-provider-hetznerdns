# hetznerdns_records Data Source

Provides details about all records of a Hetzner DNS Zone.

## Example Usage

```hcl
data "hetznerdns_records" "zone1" {
  zone_id = var.zone_id
}
```

## Argument Reference

- `zone_id` - (Required, string) The ID of the DNS zone.

## Read-Only

- `records` (List of Object)

### Nested Schema for `records`

Read-Only:

- `zone_id` (string)
- `name` (string)
- `value` (string)
- `type` (string)
- `ttl` (int)
