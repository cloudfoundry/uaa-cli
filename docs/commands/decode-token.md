# decode-token

[← Command Reference](../commands.md)

Decode a JWT token and display its claims as JSON. If no token is provided, the access token from the active context is used.

## Usage

```
uaa decode-token [TOKEN] [flags]
```

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `TOKEN` | No | A base64url-encoded JWT. Defaults to the access token in the active context. |

A second positional argument (token type, e.g. `bearer`) is accepted and ignored, for compatibility with `uaac token decode`.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--key` | | | PEM-encoded public key or certificate for signature verification |
| `--decode-times` | `-d` | `false` | Print human-readable timestamps for date fields (`iat`, `exp`, `nbf`, etc.) |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Output

By default the command prints the decoded JWT claims as indented JSON, suitable for piping to `jq` or other tools.

```json
{
  "aud": ["uaa", "admin"],
  "client_id": "admin",
  "exp": 1507671823,
  "iat": 1505079823,
  "iss": "https://uaa.example.com/oauth/token",
  "jti": "abc123",
  "scope": ["uaa.admin"],
  "sub": "admin",
  "zid": "uaa"
}
```

With `--decode-times`, a human-readable timestamp section is appended after the JSON:

```
--- Decoded timestamps ---
iat          (Issued At):      2017-09-10 18:03:43 UTC  (7 years ago)
exp          (Expires At):     2017-10-11 06:03:43 UTC  (7 years ago)
```

With `--key`, the signature is verified before printing claims. If verification succeeds, `Valid token signature.` is printed first. If it fails, the command exits with a non-zero status.

## Examples

### Decode the token from the active context

```bash
uaa target https://uaa.example.com
uaa get-client-credentials-token admin -s adminsecret
uaa decode-token
```

### Decode a token passed directly

```bash
uaa decode-token eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Decode with human-readable timestamps

```bash
uaa decode-token --decode-times
uaa decode-token eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9... -d
```

### Verify signature against a signing key

```bash
# Fetch the UAA's public signing key
uaa get-token-key

# Verify and decode
uaa decode-token --key "$(cat signing-key.pem)"
uaa decode-token eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9... --key "$(cat signing-key.pem)" --decode-times
```

## Troubleshooting

**`no token provided and no token found in active context`**
Run a token command first, e.g. `uaa get-client-credentials-token` or `uaa get-password-token`.

**`invalid token signature`**
The key passed to `--key` does not match the key used to sign the token. Fetch the correct signing key with `uaa get-token-key` or `uaa get-token-keys`.

**`unsupported algorithm`**
Only RSA (`RS256`, `RS384`, `RS512`) and EC (`ES256`, `ES384`, `ES512`) algorithms are supported for local signature verification.

## See Also

- [get-token-key](get-token-key.md) — view the UAA's current JWT signing key
- [get-token-keys](get-token-keys.md) — view all signing keys the UAA has used
- [context](context.md) — view the active context and its access token

---

[← Command Reference](../commands.md)
