# get-group

[← Command Reference](../commands.md)

Look up a group by group name.

## Usage

```
uaa get-group GROUPNAME [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--attributes` | `-a` | | Comma-separated attributes to include in the response (improves query performance) |
| `--zone` | `-z` | | Identity zone subdomain in which to look up the group |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa get-group cloud_controller.read
uaa get-group cloud_controller.read --attributes id,displayName,members
```

## See Also

- [list-groups](list-groups.md)
- [create-group](create-group.md)
- [add-member](add-member.md)

---

[← Command Reference](../commands.md)
