# list-groups

[← Command Reference](../commands.md)

Search and list groups using SCIM filters.

## Usage

```
uaa list-groups [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--filter` | | | A SCIM filter expression, e.g. `'id eq "a5e3f9fb-65a0-4033-a86c-11f4712e1fed"'` |
| `--attributes` | `-a` | | Comma-separated attributes to include in results (improves query performance) |
| `--sortBy` | `-b` | | Attribute to sort results by (e.g. `created`, `displayName`) |
| `--sortOrder` | `-o` | | Sort direction: `ascending` or `descending` |
| `--zone` | `-z` | | Identity zone subdomain from which to list groups |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
# List all groups
uaa list-groups

# Find a group by display name
uaa list-groups --filter 'displayName eq "cloud_controller.read"'

# List groups sorted by name
uaa list-groups --sortBy displayName --sortOrder ascending

# Return only id and displayName
uaa list-groups --attributes id,displayName
```

## See Also

- [get-group](get-group.md)
- [create-group](create-group.md)

---

[← Command Reference](../commands.md)
