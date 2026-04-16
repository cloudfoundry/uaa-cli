# delete-client

[← Command Reference](../commands.md)

Delete an OAuth client registration from the UAA.

## Usage

```
uaa delete-client CLIENT_ID [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--zone` | `-z` | | Identity zone subdomain in which to delete the client |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa delete-client my-app
uaa delete-client my-app --zone my-zone
```

## See Also

- [create-client](create-client.md)
- [list-clients](list-clients.md)

---

[← Command Reference](../commands.md)
