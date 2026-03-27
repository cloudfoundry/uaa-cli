# set-client-secret

Update the secret for an OAuth client registration.

## Usage

```
uaa set-client-secret CLIENT_ID -s CLIENT_SECRET [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--client_secret` | `-s` | | New client secret |
| `--zone` | `-z` | | Identity zone subdomain where the client resides |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa set-client-secret my-app -s new-secret
uaa set-client-secret my-app -s new-secret --zone my-zone
```

## See Also

- [update-client](update-client.md)
