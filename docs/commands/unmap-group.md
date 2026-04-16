# unmap-group

[← Command Reference](../commands.md)

Remove the mapping between an external group and a UAA group (scope).

## Usage

```
uaa unmap-group EXTERNAL_GROUPNAME GROUPNAME [flags]
```

Note the argument order: the external group name comes first, followed by the UAA group name.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--origin` | | `ldap` | Identity provider origin for the external group mapping |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
# Unmap an LDAP group from a UAA scope
uaa unmap-group "cn=developers,ou=groups,dc=example,dc=com" cloud_controller.read

# Unmap from a specific origin
uaa unmap-group "developers" cloud_controller.read --origin okta
```

## See Also

- [map-group](map-group.md)
- [list-group-mappings](list-group-mappings.md)

---

[← Command Reference](../commands.md)
