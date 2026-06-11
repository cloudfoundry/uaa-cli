# set-password

[← Command Reference](../commands.md)

Set password for a user account by username (admin operation). This command allows administrators to change a user's password without knowing the current password.

## Usage

```
uaa set-password USERNAME [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--password` | `-p` | | New password for the user (will prompt if not provided) |
| `--origin` | `-o` | | Identity provider in which to search. Examples: uaa, ldap, etc. |
| `--zone` | `-z` | | Identity zone subdomain in which to set the password |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
# Set password interactively (will prompt)
uaa set-password testuser

# Set password via command line flag
uaa set-password testuser --password newpassword123

# Set password for user in specific origin
uaa set-password testuser --password newpass --origin ldap

# Set password for user in specific zone
uaa set-password testuser --password newpass --zone my-zone

# Use verbose mode to see request details
uaa set-password testuser --password newpass --verbose
```

## Authentication Requirements

This command requires administrative privileges. You need a token with `password.write` or `scim.write` scopes.

## See Also

- [create-user](create-user.md)
- [get-user](get-user.md)
- [activate-user](activate-user.md)
- [deactivate-user](deactivate-user.md)

---

[← Command Reference](../commands.md)