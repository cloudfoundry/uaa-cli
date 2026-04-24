# target

[← Command Reference](../commands.md)

Set the URL of the UAA you'd like to target, or display the current target.

## Usage

```
uaa target UAA_URL [flags]
uaa target
```

Aliases: `api`

When called with no arguments, displays the currently targeted UAA URL and its status.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--skip-ssl-validation` | `-k` | `false` | Disable SSL certificate validation for requests to this target |
| `--force` | `-f` | `false` | Save the target without verifying connectivity (skip the `/info` check) |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
# Set a target
uaa target https://uaa.example.com

# Set a target, skipping SSL validation
uaa target https://uaa.example.com --skip-ssl-validation

# Set a target without checking connectivity (useful for unreachable or not-yet-running UAAs)
uaa target https://uaa.example.com --force

# Display the current target
uaa target
```

---

[← Command Reference](../commands.md)
