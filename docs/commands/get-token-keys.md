# get-token-keys

View all keys the UAA has used to sign JWT tokens.

## Usage

```
uaa get-token-keys
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa target https://uaa.example.com
uaa get-token-keys
```

## See Also

- [get-token-key](get-token-key.md) — view the current (single) signing key
