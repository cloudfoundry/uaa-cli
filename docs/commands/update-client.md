# update-client

Update an existing OAuth client registration in the UAA.

## Usage

```
uaa update-client CLIENT_ID [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--client_secret` | `-s` | | Client secret |
| `--authorized_grant_types` | | | Comma-separated list of grant types allowed for this client |
| `--scope` | | | Comma-separated scopes for `authorization_code`, `implicit`, or `password` grants |
| `--authorities` | | | Comma-separated scopes for `client_credentials` grant |
| `--redirect_uri` | | | Comma-separated callback URLs for `authorization_code` and `implicit` grants |
| `--display_name` | | | A human-readable name for this client |
| `--access_token_validity` | | `0` | Seconds before issued access tokens expire |
| `--refresh_token_validity` | | `0` | Seconds before issued refresh tokens expire |
| `--zone` | `-z` | | Identity zone subdomain in which to update the client |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa update-client my-app --scope openid,profile,email
uaa update-client my-app --access_token_validity 3600
```

## See Also

- [create-client](create-client.md)
- [get-client](get-client.md)
- [set-client-secret](set-client-secret.md)
