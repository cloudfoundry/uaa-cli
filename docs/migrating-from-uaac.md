# Migrating from uaac

This guide helps users of the Ruby-based [cf-uaac](https://github.com/cloudfoundry/cf-uaac) (`uaac`) CLI transition to the Go-based `uaa` CLI.
See the [UAA-CLI Command Reference](./commands.md) for the full list of `uaa` commands and their options.

## Key Differences

|                    | uaac (Ruby)                                      | uaa (Go)                                                                     |
|--------------------|--------------------------------------------------|------------------------------------------------------------------------------|
| **Command style**  | Hierarchical topic/verb: `uaac token client get` | Flat VERB-NOUN with hyphens: `uaa get-client-credentials-token`              |
| **Output format**  | Human-readable text                              | Defaults to Machine-parseable JSON, with some human-readable status messages |
| **SSL validation** | `uaac target URL --skip-ssl-validation`          | `uaa target URL --skip-ssl-validation`                                       |
| **Identity zones** | `-z` / `--zone` global flag                      | `-z` / `--zone` flag on most commands                                        |
| **Verbose/trace**  | `-t` / `--trace` for verbose debug output        | `-v` / `--verbose` on all commands                                           |

### Output Format Notes

The `uaa` CLI outputs a combination of human-readable status messages and JSON data to stdout. Many commands (like `create-client`, `update-client`, etc.) print a success message followed by JSON output. For reliable machine parsing:

- Extract the full JSON block starting at the first line that begins with `{` or `[` rather than using only the last line, for example: `uaa create-client myclient -s secret --authorized_grant_types client_credentials | sed -n '/^[[:space:]]*[{\[]/,$p' | jq`
- Be aware that status messages and JSON are both written to stdout, not separated by stderr
- Commands that only retrieve data (like `get-client`, `list-clients`) typically output only JSON

---

## Command Reference

### Targeting & Context

| uaac | uaa | Notes |
|------|-----|-------|
| `uaac target [uaa_url]` | [`uaa target UAA_URL`](commands/target.md) | |
| `uaac targets` | *(no equivalent)* | |
| `uaac context [name]` | [`uaa context`](commands/context.md) | |
| `uaac contexts` | `uaa contexts` | |
| `uaac version` | [`uaa version`](commands/version.md) | |

### System Information

| uaac | uaa | Notes |
|------|-----|-------|
| `uaac info` | [`uaa info`](commands/info.md) | |
| `uaac me` | [`uaa userinfo`](commands/userinfo.md) | |
| `uaac signing key` | [`uaa get-token-key`](commands/get-token-key.md) | uaa-cli splits this into two commands |
| `uaac signing key` | [`uaa get-token-keys`](commands/get-token-keys.md) | Returns all signing keys |
| `uaac prompts` | *(no equivalent)* | Use `uaa curl /info` to inspect prompts |
| `uaac password strength [password]` | *(no equivalent)* | Use `uaa curl /password/score -X POST -d 'password=...'` |

### Tokens

| uaac | uaa | Notes |
|------|-----|-------|
| `uaac token client get [id]` | [`uaa get-client-credentials-token CLIENT_ID -s SECRET`](commands/get-client-credentials-token.md) | |
| `uaac token owner get [client] [user]` | [`uaa get-password-token CLIENT_ID -s SECRET -u USER -p PASS`](commands/get-password-token.md) | |
| `uaac token authcode get` | [`uaa get-authcode-token CLIENT_ID -s SECRET --port PORT`](commands/get-authcode-token.md) | |
| `uaac token implicit get` | [`uaa get-implicit-token CLIENT_ID --port PORT`](commands/get-implicit-token.md) | |
| `uaac token refresh [refreshtoken]` | [`uaa refresh-token -s SECRET`](commands/refresh-token.md) | |
| `uaac token get [credentials...]` | *(no equivalent)* | Use `uaa get-password-token` for password grant |
| `uaac token sso get [client]` | *(no equivalent)* | Passcode/SSO grant not yet supported |
| `uaac token decode [token] [tokentype]` | *(no equivalent)* | Decode manually: `uaa context --access_token \| cut -d. -f2 \| base64 -d` |
| `uaac token delete [contexts...]` | *(no equivalent)* | Delete `~/.uaa/config.json` to clear all contexts |

### Users

| uaac | uaa | Notes |
|------|-----|-------|
| `uaac user add [name]` | [`uaa create-user USERNAME --givenName G --familyName F -p PASS --email EMAIL`](commands/create-user.md) | |
| `uaac user get [name]` | [`uaa get-user USERNAME`](commands/get-user.md) | |
| `uaac users [filter]` | [`uaa list-users --filter '...'`](commands/list-users.md) | |
| `uaac user delete [name]` | [`uaa delete-user USERNAME`](commands/delete-user.md) | |
| `uaac user activate [name]` | [`uaa activate-user USERNAME`](commands/activate-user.md) | |
| `uaac user deactivate [name]` | [`uaa deactivate-user USERNAME`](commands/deactivate-user.md) | |
| `uaac user update [name]` | *(no equivalent)* | Use `uaa curl /Users/USER_ID -X PUT -d '{...}'` |
| `uaac user ids [username\|id...]` | *(no equivalent)* | Use `uaa get-user USERNAME` for individual lookups |
| `uaac user unlock [name]` | *(no equivalent)* | Use `uaa curl /Users/USER_ID/status -X PATCH -d '{"locked":false}'` |
| `uaac password set [name]` | *(no equivalent)* | Use `uaa curl /Users/USER_ID/password -X PUT -d '{"password":"NEW"}'` |
| `uaac password change` | *(no equivalent)* | Use `uaa curl /Users/USER_ID/password -X PUT -d '{"oldPassword":"OLD","password":"NEW"}'` |

### Clients

| uaac | uaa | Notes |
|------|-----|-------|
| `uaac client add [id]` | [`uaa create-client CLIENT_ID -s SECRET ...`](commands/create-client.md) | uaac supports `--interactive` / `-i` to prompt for each field; uaa-cli does not |
| `uaac client get [id]` | [`uaa get-client CLIENT_ID`](commands/get-client.md) | |
| `uaac clients [filter]` | [`uaa list-clients`](commands/list-clients.md) | |
| `uaac client update [id]` | [`uaa update-client CLIENT_ID ...`](commands/update-client.md) | uaac supports `--interactive` / `-i`; uaa-cli does not |
| `uaac client delete [id]` | [`uaa delete-client CLIENT_ID`](commands/delete-client.md) | |
| `uaac secret set [id]` | [`uaa set-client-secret CLIENT_ID -s SECRET`](commands/set-client-secret.md) | |
| `uaac secret change` | *(no equivalent)* | Use `uaa curl /oauth/clients/CLIENT_ID/secret -X PUT -d '{"oldSecret":"OLD","secret":"NEW"}'` |
| `uaac client jwt add [id]` | *(no equivalent)* | Use `uaa curl /oauth/clients/CLIENT_ID/clientjwt -X PUT -d '{...}'` |
| `uaac client jwt update [id]` | *(no equivalent)* | Use `uaa curl /oauth/clients/CLIENT_ID/clientjwt -X PUT -d '{...}'` |
| `uaac client jwt delete [id]` | *(no equivalent)* | Use `uaa curl /oauth/clients/CLIENT_ID/clientjwt -X DELETE` |

### Groups

| uaac | uaa | Notes |
|------|-----|-------|
| `uaac group add [name]` | [`uaa create-group GROUPNAME`](commands/create-group.md) | |
| `uaac group get [name]` | [`uaa get-group GROUPNAME`](commands/get-group.md) | |
| `uaac groups [filter]` | [`uaa list-groups --filter '...'`](commands/list-groups.md) | |
| `uaac group delete [name]` | *(no equivalent)* | Use `uaa curl /Groups/GROUP_ID -X DELETE` |
| `uaac member add [name] [users...]` | [`uaa add-member GROUPNAME USERNAME`](commands/add-member.md) | uaa-cli adds one user at a time |
| `uaac member delete [name] [users...]` | [`uaa remove-member GROUPNAME USERNAME`](commands/remove-member.md) | uaa-cli removes one user at a time |
| `uaac group map [external_group]` | [`uaa map-group EXTERNAL_GROUPNAME GROUPNAME`](commands/map-group.md) | See argument order note below |
| `uaac group unmap [group_name] [external_group]` | [`uaa unmap-group EXTERNAL_GROUPNAME GROUPNAME`](commands/unmap-group.md) | See argument order note below |
| `uaac group mappings` | [`uaa list-group-mappings`](commands/list-group-mappings.md) | |

### Miscellaneous

| uaac | uaa | Notes |
|------|-----|-------|
| `uaac curl [path]` | [`uaa curl PATH`](commands/curl.md) | |

---

## Notable Argument Differences

### Group mapping argument order

The `map-group` and `unmap-group` commands take arguments in a different order than uaac:

```bash
# uaac: UAA group name first (via --name flag), external group as positional arg
uaac group map --name cloud_controller.read "cn=devs,ou=groups,dc=example,dc=com"

# uaa: external group name first, UAA group name second
uaa map-group "cn=devs,ou=groups,dc=example,dc=com" cloud_controller.read
```

### `member add` / `member delete` — one user at a time

uaac accepts multiple usernames in a single call; uaa-cli accepts one at a time:

```bash
# uaac: add multiple users at once
uaac member add my-group alice bob carol

# uaa: one user per call
uaa add-member my-group alice
uaa add-member my-group bob
uaa add-member my-group carol
```

---

## Using `uaa curl` as a Fallback

For uaac commands that have no direct equivalent, `uaa curl` provides authenticated access to any UAA API endpoint. First obtain a token, then use `uaa curl` with the active context's credentials:

```bash
uaa target https://uaa.example.com
uaa get-client-credentials-token admin -s admin-secret

# Example: delete a group by ID
GROUP_ID=$(uaa get-group my-group | jq -r .id)
uaa curl /Groups/$GROUP_ID -X DELETE
```
