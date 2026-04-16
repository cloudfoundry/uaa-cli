# delete-user

[← Command Reference](../commands.md)

Delete a user account by username.

## Usage

```
uaa delete-user USERNAME [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--origin` | `-o` | `uaa` | Identity provider origin of the user to delete |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa delete-user bob
uaa delete-user bob --origin ldap
```

## See Also

- [get-user](get-user.md)
- [deactivate-user](deactivate-user.md)

---

[← Command Reference](../commands.md)
