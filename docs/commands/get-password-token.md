# get-password-token

Obtain an access token using the resource owner password credentials (`password`) grant type and save it to the active context.

## Usage

```
uaa get-password-token CLIENT_ID -s CLIENT_SECRET -u USERNAME -p PASSWORD [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--client_secret` | `-s` | | Client secret |
| `--username` | `-u` | | Username |
| `--password` | `-p` | | User password |
| `--format` | | `jwt` | Token format. Available values: `jwt`, `opaque` |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Background

The password grant type is used by highly trusted client applications (such as CLIs) where the user provides their credentials directly to the client. The client then exchanges those credentials for a token at the UAA.

After successfully running this command, the token is added to the CLI's current context and will be attached to subsequent requests to UAA endpoints requiring authorization.

## Examples

```bash
uaa target https://uaa.example.com
uaa get-password-token cf -s "" -u admin -p admin-secret
uaa context
```

## Troubleshooting

**Unable to get a token:**
- Verify `client_id`, `client_secret`, `username`, and `password` are all correct.
- Ensure `password` is listed in the client's `authorized_grant_types`.

**Token does not have expected scopes:**
- Check the `scope` field of the client registration.
- Verify the user has the correct group memberships. Token scopes are the intersection of the client's `scope` and the user's group memberships.
