# list-users

Search and list user accounts using SCIM filters.

## Usage

```
uaa list-users [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--filter` | | | A SCIM filter expression, e.g. `'userName eq "bob@example.com"'` |
| `--attributes` | `-a` | | Comma-separated attributes to include in results (improves query performance) |
| `--sortBy` | `-b` | | Attribute to sort results by (e.g. `created`, `userName`) |
| `--sortOrder` | `-o` | | Sort direction: `ascending` or `descending` |
| `--zone` | `-z` | | Identity zone subdomain in which to list users |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Background

SCIM (System for Cross-Domain Identity Management) filters provide a limited query interface for searching users. Fetching group membership data requires SQL joins on the UAA side, so use `--attributes` to limit fields when querying large user sets.

## Examples

```bash
# List all users
uaa list-users

# Find users with a gmail.com address
uaa list-users --filter 'userName co "gmail.com"' --attributes id,emails

# Find users from a specific identity provider
uaa list-users --filter 'userName eq "bob@example.com" and origin eq "ldap"'

# Find unverified users
uaa list-users --filter 'verified eq false' --attributes id,userName,name,emails

# Find users whose username starts with "z"
uaa list-users --filter 'userName sw "z"'

# See client approvals for a specific user
uaa list-users --filter 'userName eq "bob@example.com"' --attributes approvals

# See full details including group memberships
uaa list-users --filter 'userName eq "bob@example.com"'
```

## See Also

- [get-user](get-user.md)
- [create-user](create-user.md)
