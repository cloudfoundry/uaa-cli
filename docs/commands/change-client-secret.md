# uaa change-client-secret

## Overview

Change the secret for the currently authenticated client. This command allows a client to change its own secret by providing both the old secret and the new secret.

## Usage

```
uaa change-client-secret --old_secret OLD_SECRET --secret NEW_SECRET [flags]
```

## Required Authentication

This command requires an active client context obtained via the `client_credentials` grant type.

## Arguments

| Argument | Description |
|----------|-------------|
| `--old_secret` | The current secret for the client |
| `--secret`, `-s` | The new secret for the client |

## Options

| Option | Description |
|--------|-------------|
| `--zone`, `-z` | Identity zone subdomain where the client resides |
| `--verbose`, `-v` | Display verbose output including HTTP request/response details |

## Examples

### Change client secret with explicit values

```bash
uaa change-client-secret --old_secret currentsecret --secret newsecret
```

### Change client secret in a specific zone

```bash
uaa change-client-secret --old_secret currentsecret --secret newsecret --zone myzone
```

### Change client secret with verbose output

```bash
uaa change-client-secret --old_secret currentsecret --secret newsecret --verbose
```

## Prerequisites

1. You must have targeted a UAA server using `uaa target`
2. You must have an active client context with `client_credentials` grant type (obtained via `uaa get-client-credentials-token`)
3. The client must have the necessary permissions to change its own secret

## Notes

- This command is for self-service secret changes where a client changes its own secret
- Both the old and new secrets must be provided for security reasons
- After changing the secret, you will need to re-authenticate with the new secret
- Use `--verbose` to see the actual HTTP request being made to the UAA