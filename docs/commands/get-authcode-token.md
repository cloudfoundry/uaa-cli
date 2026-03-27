# get-authcode-token

Obtain an access token using the `authorization_code` grant type and save it to the active context. Opens a browser window for the user to authenticate.

## Usage

```
uaa get-authcode-token CLIENT_ID -s CLIENT_SECRET --port REDIRECT_URI_PORT [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--client_secret` | `-s` | | Client secret |
| `--port` | | | Port on which to run the local callback server. Must match a `localhost` redirect URI in the client registration. |
| `--scope` | | `openid` | Comma-separated scopes to request in the token |
| `--format` | | `jwt` | Token format. Available values: `jwt`, `opaque` |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa target https://uaa.example.com
uaa get-client-credentials-token admin -s admin-secret
uaa create-client my-app \
    --authorized_grant_types authorization_code \
    --redirect_uri http://localhost:9090/** \
    --scope openid,cloud_controller.read
uaa get-authcode-token my-app -s my-secret --port 9090
uaa context
```

## Troubleshooting

**Unknown redirect_uri error after signing in:**
- The `--port` value must match a `http://localhost:PORT` entry in the client's `redirect_uri` list. Run `uaa get-client CLIENT_ID` to inspect the registration.

**Token does not have expected scopes:**
- Check the `scope` field of the client registration.
- Verify the user has the correct group memberships.
