# update-user

[← Command Reference](../commands.md)

Update an existing user account in the UAA.

## Usage

```
uaa update-user USERNAME [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--given_name` | | | Given (first) name |
| `--family_name` | | | Family (last) name |
| `--emails` | | | Email addresses (flag may be specified multiple times) |
| `--phones` | | | Phone numbers (flag may be specified multiple times) |
| `--origin` | `-o` | | Identity provider origin to search for user (e.g. `uaa`, `ldap`) |
| `--del_attrs` | | | Attributes to remove (e.g. `phoneNumbers`, `name`) |
| `--zone` | `-z` | | Identity zone subdomain in which to update the user |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
# Update a user's name
uaa update-user bob --given_name Robert --family_name Smith

# Update a user's email addresses
uaa update-user alice --emails alice@newdomain.com --emails alice.jones@work.com

# Update a user from a specific origin
uaa update-user carol --origin ldap --given_name Caroline

# Remove phone numbers from a user
uaa update-user bob --del_attrs phoneNumbers

# Update multiple attributes at once
uaa update-user alice \
    --given_name Alice \
    --family_name Johnson \
    --emails alice.johnson@example.com \
    --phones 555-1234

# Update user in a specific zone with verbose output
uaa update-user bob \
    --given_name Robert \
    --zone my-zone \
    --verbose
```

## Notes

- The command first retrieves the existing user, then merges the specified updates
- At least one update flag must be specified
- When using `--del_attrs`, be careful not to remove required attributes
- The `--emails` attribute cannot be deleted as it may make the user unusable
- Use `--verbose` to see the HTTP PUT request details

## See Also

- [create-user](create-user.md)
- [get-user](get-user.md)
- [list-users](list-users.md)
- [delete-user](delete-user.md)

---

[← Command Reference](../commands.md)