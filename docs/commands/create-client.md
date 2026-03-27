# create-client

Create an OAuth client registration in the UAA.

## Usage

```
uaa create-client CLIENT_ID -s CLIENT_SECRET --authorized_grant_types GRANT_TYPES [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--client_secret` | `-s` | | Client secret |
| `--authorized_grant_types` | | | Comma-separated list of grant types allowed for this client (e.g. `authorization_code`, `client_credentials`, `implicit`, `password`) |
| `--scope` | | | Comma-separated scopes requested during `authorization_code`, `implicit`, or `password` grants |
| `--authorities` | | | Comma-separated scopes requested during `client_credentials` grant |
| `--redirect_uri` | | | Comma-separated callback URLs allowed for `authorization_code` and `implicit` grants |
| `--display_name` | | | A human-readable name for this client |
| `--access_token_validity` | | `0` | Seconds before issued access tokens expire |
| `--refresh_token_validity` | | `0` | Seconds before issued refresh tokens expire |
| `--clone` | | | Client ID of an existing client to clone configuration from |
| `--zone` | `-z` | | Identity zone subdomain in which to create the client |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
# Authorization code client
uaa create-client shinymail \
    --client_secret secret \
    --authorized_grant_types authorization_code \
    --redirect_uri http://localhost:9090/*,https://shinymail.example.com/callback \
    --scope mail.send,mail.read \
    --display_name "Shinymail Web Mail Reader"

# Client credentials (service-to-service) client
uaa create-client background-emailer \
    --client_secret secret \
    --authorized_grant_types client_credentials \
    --authorities notifications.write \
    --display_name "Weekly newsletter email service"

# Implicit (single-page app) client
uaa create-client my-spa \
    --authorized_grant_types implicit \
    --redirect_uri http://localhost:9090/*,https://myapp.example.com/callback \
    --scope openid,todo.read,todo.write \
    --display_name "My Single-Page App"

# Password grant client
uaa create-client trusted-cli \
    --client_secret mumstheword \
    --authorized_grant_types password \
    --scope cloud_controller.admin,uaa.admin

# Clone an existing client
uaa create-client trusted-cli-copy \
    --clone trusted-cli \
    --client_secret donttellanyone
```

## See Also

- [update-client](update-client.md)
- [delete-client](delete-client.md)
- [get-client](get-client.md)
- [list-clients](list-clients.md)
