# context

Display information about the currently active CLI context, including the cached access token and its metadata.

## Usage

```
uaa context [flags]
```

A context represents a previously fetched access token and associated metadata such as the scopes that token contains. The uaa CLI caches these results in a local file so that they may be used when issuing requests that require an Authorization header.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--access_token` | | `false` | Display the context's raw access token |
| `--auth_header` | `-a` | `false` | Display the context's token type and access token (suitable for use as an Authorization header value) |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Print additional info on HTTP requests |

## Examples

```bash
# Show current context
uaa context

# Show the raw access token
uaa context --access_token

# Show the Authorization header value
uaa context --auth_header
```
