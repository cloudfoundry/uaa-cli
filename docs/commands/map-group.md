# map-group

[← Command Reference](../commands.md)

Map a UAA group (scope) to an external group defined within an external identity provider.

## Usage

```
uaa map-group EXTERNAL_GROUPNAME GROUPNAME [flags]
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
# Map an LDAP group to a UAA scope (defaults to ldap origin)
uaa map-group "cn=developers,ou=groups,dc=example,dc=com" cloud_controller.read

# Map to a specific origin
uaa map-group "developers" cloud_controller.read --origin okta
```

## See Also

- [unmap-group](unmap-group.md)
- [list-group-mappings](list-group-mappings.md)

---

[← Command Reference](../commands.md)
