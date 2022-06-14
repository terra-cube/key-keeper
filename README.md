# key-keeper

Example config.yml

```yaml
vault:
  address: "address vault"
  token: "unexpired token"
  timeout: "timeout request: 1s, 3m, 5h"
certificates:
  root_path: "root ca path"
  cert_path: "path for generate certificate"
  valid_interval: "valid expire cert interval: 1s, 3m, 5h"
  ca:
    common_name: "common_name"
    vault_path: "path for store of ca in vault"
    path: "path for store of ca in host"
```

Build:

    go build -o key-keeper cmd/key-keeper/main.go

Run:

    key-keeper -config /path/to/config
