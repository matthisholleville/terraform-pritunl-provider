---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pritunl Provider"
subcategory: ""
description: |-
  
---

# pritunl Provider



## Example Usage

```terraform
provider "pritunl" {
  url      = "https://vpn.server.com"
  token    = "api-token"
  secret   = "api-secret"
  insecure = false
}

data "pritunl_server" "server" {
  server_name = "Devops_test"
}

resource "pritunl_route" "test" {
  server_id = data.pritunl_server.server.id
  nat       = true
  comment   = "my custom route"
  network   = "34.56.43.32/32"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `connection_check` (Boolean)
- `insecure` (Boolean)
- `secret` (String)
- `token` (String)
- `url` (String)