# create-group

[← Command Reference](../commands.md)

Create a group (scope) in the UAA.

## Usage

```
uaa create-group GROUPNAME [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--description` | `-d` | | A human-readable description of the group |
| `--zone` | `-z` | | Identity zone subdomain in which to create the group |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa create-group cloud_controller.read
uaa create-group cloud_controller.read --description "Read access to Cloud Controller resources"
uaa create-group my-scope --zone my-zone
```

## See Also

- [get-group](get-group.md)
- [list-groups](list-groups.md)
- [add-member](add-member.md)

---

[← Command Reference](../commands.md)
