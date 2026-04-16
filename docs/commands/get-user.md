# get-user

[← Command Reference](../commands.md)

Look up a user account by username.

## Usage

```
uaa get-user USERNAME [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--origin` | `-o` | | Identity provider to search (e.g. `uaa`, `ldap`). Searches all origins if omitted. |
| `--attributes` | `-a` | | Comma-separated list of user attributes to return (improves query performance) |
| `--zone` | `-z` | | Identity zone subdomain in which to find the user |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa get-user bob
uaa get-user bob --origin ldap
uaa get-user bob --attributes id,userName,emails
```

## See Also

- [list-users](list-users.md)
- [create-user](create-user.md)
- [delete-user](delete-user.md)

---

[← Command Reference](../commands.md)
