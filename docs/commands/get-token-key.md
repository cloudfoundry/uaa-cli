# get-token-key

[← Command Reference](../commands.md)

View the key the UAA uses to sign JWT tokens.

## Usage

```
uaa get-token-key
```

Aliases: `token-key`

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
uaa target https://uaa.example.com
uaa get-token-key
```

## See Also

- [get-token-keys](get-token-keys.md) — view all signing keys the UAA has used

---

[← Command Reference](../commands.md)
