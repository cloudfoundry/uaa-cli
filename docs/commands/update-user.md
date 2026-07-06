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
| `--givenName` | | | Given (first) name |
| `--familyName` | | | Family (last) name |
| `--email` | | | Email address (flag may be specified multiple times) |
| `--phone` | | | Phone number (flag may be specified multiple times) |
| `--origin` | `-o` | | Identity provider origin to search for user (e.g. `uaa`, `ldap`) |
| `--delAttrs` | | | Attributes to remove (e.g. `phoneNumbers`, `name`) |
| `--zone` | `-z` | | Identity zone subdomain in which to update the user |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
# Update a user's name
uaa update-user bob --givenName Robert --familyName Smith

# Update a user's email addresses
uaa update-user alice --email alice@newdomain.com --email alice.jones@work.com

# Update a user from a specific origin
uaa update-user carol --origin ldap --givenName Caroline

# Remove phone numbers from a user
uaa update-user bob --delAttrs phoneNumbers

# Update multiple attributes at once
uaa update-user alice \
    --givenName Alice \
    --familyName Johnson \
    --email alice.johnson@example.com \
    --phone 555-1234

# Update user in a specific zone with verbose output
uaa update-user bob \
    --givenName Robert \
    --zone my-zone \
    --verbose
```

## Notes

- The command first retrieves the existing user, then merges the specified updates
- At least one update flag (`--givenName`, `--familyName`, `--email`, `--phone`, or `--delAttrs`) must be specified
- When using `--delAttrs`, be careful not to remove required attributes
- The `--email` attribute cannot be deleted as it may make the user unusable
- Use `--verbose` to see the HTTP PUT request details

## See Also

- [create-user](create-user.md)
- [get-user](get-user.md)
- [list-users](list-users.md)
- [delete-user](delete-user.md)

---

[← Command Reference](../commands.md)