# targets

[← Command Reference](../commands.md)

List all registered targets. The active target is marked with `*`.

## Usage

```
uaa targets
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Output

Each target is printed on its own line with a 1-based index. The currently active target is prefixed with `*`:

```
* 1: https://uaa.example.com
  2: http://localhost:8080/uaa
```

If no targets have been registered, the command prints `No targets set.` and exits 0.

## Examples

```bash
# Register two targets and list them
uaa target https://uaa.example.com
uaa target http://localhost:8080/uaa --skip-ssl-validation
uaa targets
```

## See Also

- [target](target.md) — set or display the current target

---

[← Command Reference](../commands.md)
