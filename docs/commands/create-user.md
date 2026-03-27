# create-user

[← Command Reference](../commands.md)

Create a user account in the UAA.

## Usage

```
uaa create-user USERNAME [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--givenName` | | | Given (first) name (required) |
| `--familyName` | | | Family (last) name (required) |
| `--email` | | | Email address (required; flag may be specified multiple times) |
| `--password` | `-p` | | User password (required for `uaa` origin) |
| `--origin` | `-o` | `uaa` | Identity provider origin (e.g. `uaa`, `ldap`) |
| `--phone` | | | Phone number (optional; flag may be specified multiple times) |
| `--zone` | `-z` | | Identity zone subdomain in which to create the user |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
# Create a UAA-origin user
uaa create-user bob \
    --givenName Bob \
    --familyName Smith \
    --email bob@example.com \
    --password s3cr3t

# Create a user from an external identity provider
uaa create-user alice \
    --givenName Alice \
    --familyName Jones \
    --email alice@example.com \
    --origin ldap

# Create a user with multiple email addresses
uaa create-user carol \
    --givenName Carol \
    --familyName White \
    --email carol@example.com \
    --email carol.white@work.example.com \
    --password s3cr3t
```

## See Also

- [get-user](get-user.md)
- [list-users](list-users.md)
- [delete-user](delete-user.md)

---

[← Command Reference](../commands.md)
