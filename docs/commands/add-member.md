# add-member

Add a user to a group.

## Usage

```
uaa add-member GROUPNAME USERNAME
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa add-member cloud_controller.read bob
```

## See Also

- [remove-member](remove-member.md)
- [get-group](get-group.md)
- [create-group](create-group.md)
