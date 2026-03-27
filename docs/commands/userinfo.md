# userinfo

[← Command Reference](../commands.md)

Display claims about the currently authenticated user by calling the UAA's `/userinfo` endpoint.

## Usage

```
uaa userinfo
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Background

The `userinfo` command calls the UAA's `/userinfo` endpoint, which is implemented in accordance with the OIDC 1.1 spec. It returns claims about the authenticated user using the access token saved in the current context.

## Examples

```bash
uaa target https://uaa.example.com
uaa get-password-token cf -s "" -u admin -p admin-secret
uaa userinfo
```

## Troubleshooting

**Got a token but `userinfo` returns an error:**
- Verify your token is still valid by running another authenticated command.
- Check that `openid` appears in the scopes listed by `uaa context`. The `/userinfo` endpoint requires the `openid` scope. Update your client registration to include `openid` in its `scope` list if needed.

---

[← Command Reference](../commands.md)
