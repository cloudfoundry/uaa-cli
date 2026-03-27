# refresh-token

Obtain a new access token using the `refresh_token` grant type and update the active context.

## Usage

```
uaa refresh-token -s CLIENT_SECRET [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--client_secret` | `-s` | | Client secret |
| `--format` | | `jwt` | Token format. Available values: `jwt`, `opaque` |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Background

Refresh tokens are long-lived credentials issued alongside access tokens for `authorization_code` and `password` grant clients. Use this command to exchange a refresh token for a new, unexpired access token without requiring the user to re-authenticate.

## Examples

```bash
uaa target https://uaa.example.com
uaa get-password-token cf -s "" -u admin -p admin-secret
uaa context                          # note the refresh_token
uaa refresh-token -s ""
uaa context                          # access_token should now be updated
```

## Troubleshooting

**No refresh_token in active context:**
- **Implicit grant:** Implicit clients cannot maintain secrets and are never issued refresh tokens.
- **Authorization code / password grant:** The client must include `refresh_token` in its `authorized_grant_types`. Run `uaa get-client CLIENT_ID` to verify.
- **Client credentials grant:** Refresh tokens are never issued for `client_credentials`. Re-authenticate using `get-client-credentials-token` at any time.
