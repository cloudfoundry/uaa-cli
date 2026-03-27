# get-client-credentials-token

[← Command Reference](../commands.md)

Obtain an access token using the `client_credentials` grant type and save it to the active context.

## Usage

```
uaa get-client-credentials-token CLIENT_ID -s CLIENT_SECRET [flags]
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

The client credentials grant type is used by clients performing actions on their own behalf, without the involvement of a user. It is one of the four authorization flows described in RFC 6749.

After successfully running this command, the token is added to the CLI's current context and will be attached to subsequent requests to UAA endpoints requiring authorization.

This flow is typically used for service-to-service operations. The `authorities` field of the client registration (not `scope`) determines what permissions are included in the token.

## Examples

```bash
uaa target https://uaa.example.com
uaa get-client-credentials-token my-client -s my-secret
uaa context
```

## Troubleshooting

**Unable to get a token:**
- Verify the `client_id` and `client_secret` are correct and you are targeting the right UAA.
- Ensure `client_credentials` is listed in the client's `authorized_grant_types`.

**Token does not have expected scopes:**
- Check the `authorities` field of the client registration — this controls scopes in `client_credentials` tokens, not the `scope` field.

---

[← Command Reference](../commands.md)
