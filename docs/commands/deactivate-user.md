# deactivate-user

Deactivate a user account by username. The account is suspended but not deleted.

## Usage

```
uaa deactivate-user USERNAME [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--zone` | `-z` | | Identity zone subdomain from which to deactivate the user |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa deactivate-user bob
uaa deactivate-user bob --zone my-zone
```

## See Also

- [activate-user](activate-user.md)
- [delete-user](delete-user.md)
