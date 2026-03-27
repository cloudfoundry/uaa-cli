# list-clients

List all client registrations in the targeted UAA.

## Usage

```
uaa list-clients [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--zone` | `-z` | | Identity zone subdomain in which to list clients |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa list-clients
uaa list-clients --zone my-zone
```

## See Also

- [get-client](get-client.md)
- [create-client](create-client.md)
