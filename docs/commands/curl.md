# curl

[← Command Reference](../commands.md)

Make an authenticated HTTP request to a UAA endpoint.

## Usage

```
uaa curl PATH [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--method` | `-X` | `GET` | HTTP method (GET, POST, PUT, DELETE, PATCH, etc.) |
| `--data` | `-d` | | HTTP request body |
| `--header` | `-H` | | Custom request header (flag may be specified multiple times) |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
# GET request
uaa curl /info

# POST with a JSON body
uaa curl /oauth/clients \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"client_id":"test","client_secret":"secret","authorized_grant_types":["client_credentials"]}'

# DELETE
uaa curl /oauth/clients/my-app -X DELETE

# Multiple custom headers
uaa curl /Users \
    -H "Accept: application/json" \
    -H "X-Identity-Zone-Id: my-zone"
```

---

[← Command Reference](../commands.md)
