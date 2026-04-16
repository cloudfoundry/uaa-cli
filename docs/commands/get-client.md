# get-client

[← Command Reference](../commands.md)

View the registration details for an OAuth client.

## Usage

```
uaa get-client CLIENT_ID [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--zone` | `-z` | | Identity zone subdomain in which to look up the client |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa get-client my-app
uaa get-client my-app --zone my-zone
```

## See Also

- [list-clients](list-clients.md)
- [create-client](create-client.md)
- [update-client](update-client.md)

---

[← Command Reference](../commands.md)
