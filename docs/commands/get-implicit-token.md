# get-implicit-token

[← Command Reference](../commands.md)

Obtain an access token using the `implicit` grant type and save it to the active context. Opens a browser window for the user to authenticate.

## Usage

```
uaa get-implicit-token CLIENT_ID --port REDIRECT_URI_PORT [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | | | Port on which to run the local callback server. Must match a `localhost` redirect URI in the client registration. |
| `--scope` | | `openid` | Comma-separated scopes to request in the token |
| `--format` | | `jwt` | Token format. Available values: `jwt`, `opaque` |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Background

The implicit grant type is most commonly used for single-page web applications ("public clients") that cannot maintain a client secret. After the user authorizes, the access token is transmitted as part of the URI fragment. Because implicit clients have no secret, they are not issued refresh tokens.

## Examples

```bash
uaa target https://uaa.example.com
uaa get-client-credentials-token admin -s admin-secret
uaa create-client my-spa \
    --authorized_grant_types implicit \
    --redirect_uri http://localhost:9090/**,https://myapp.example.com/** \
    --scope openid,todo.read
uaa get-implicit-token my-spa --port 9090
uaa context
```

## Troubleshooting

**Unknown redirect_uri error after signing in:**
- The `--port` value must match a `http://localhost:PORT` entry in the client's `redirect_uri` list.

**Token does not have expected scopes:**
- Check the `scope` field of the client registration.
- Verify the user has the correct group memberships.

---

[← Command Reference](../commands.md)
