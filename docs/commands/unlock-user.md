# unlock-user

[← Command Reference](../commands.md)

Unlock a user account by username. This removes lockouts caused by failed login attempts.

## Usage

```
uaa unlock-user USERNAME [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--origin` | `-o` | | Identity provider to search for the user (e.g., uaa, ldap) |
| `--zone` | `-z` | | Identity zone subdomain from which to unlock the user |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa unlock-user bob
uaa unlock-user bob --origin ldap
uaa unlock-user bob --zone my-zone
```

## See Also

- [activate-user](activate-user.md)
- [deactivate-user](deactivate-user.md)
- [get-user](get-user.md)

---

[← Command Reference](../commands.md)