# Command Reference

All commands support `-v` / `--verbose` to print detailed HTTP request/response information.

Each command name below links to a page with a full description, including all available flags and usage examples.

## Getting Started

| Command | Description |
|---------|-------------|
| [`target`](commands/target.md) | Set the URL of the UAA you'd like to target |
| [`context`](commands/context.md) | See information about the currently active CLI context |
| [`info`](commands/info.md) | See version and global configurations for the targeted UAA |
| [`version`](commands/version.md) | Print CLI version |

## Getting Tokens

| Command | Description |
|---------|-------------|
| [`get-client-credentials-token`](commands/get-client-credentials-token.md) | Obtain an access token using the `client_credentials` grant type |
| [`get-password-token`](commands/get-password-token.md) | Obtain an access token using the `password` grant type |
| [`get-authcode-token`](commands/get-authcode-token.md) | Obtain an access token using the `authorization_code` grant type |
| [`get-implicit-token`](commands/get-implicit-token.md) | Obtain an access token using the `implicit` grant type |
| [`refresh-token`](commands/refresh-token.md) | Obtain a new access token using a refresh token |
| [`get-token-key`](commands/get-token-key.md) | View the key for validating UAA's JWT token signatures |
| [`get-token-keys`](commands/get-token-keys.md) | View all keys the UAA has used to sign JWT tokens |

## Managing Clients

| Command | Description |
|---------|-------------|
| [`create-client`](commands/create-client.md) | Create an OAuth client registration in the UAA |
| [`update-client`](commands/update-client.md) | Update an OAuth client registration in the UAA |
| [`delete-client`](commands/delete-client.md) | Delete a client registration |
| [`get-client`](commands/get-client.md) | View a client registration |
| [`list-clients`](commands/list-clients.md) | See all clients in the targeted UAA |
| [`set-client-secret`](commands/set-client-secret.md) | Update the secret for a client |

## Managing Users

| Command | Description |
|---------|-------------|
| [`create-user`](commands/create-user.md) | Create a user |
| [`get-user`](commands/get-user.md) | Look up a user by username |
| [`list-users`](commands/list-users.md) | Search and list users with SCIM filters |
| [`delete-user`](commands/delete-user.md) | Delete a user by username |
| [`activate-user`](commands/activate-user.md) | Activate a user by username |
| [`deactivate-user`](commands/deactivate-user.md) | Deactivate a user by username |

## Managing Groups

| Command | Description |
|---------|-------------|
| [`create-group`](commands/create-group.md) | Create a group |
| [`get-group`](commands/get-group.md) | Look up a group by group name |
| [`list-groups`](commands/list-groups.md) | Search and list groups with SCIM filters |
| [`add-member`](commands/add-member.md) | Add a user to a group |
| [`remove-member`](commands/remove-member.md) | Remove a user from a group |
| [`map-group`](commands/map-group.md) | Map a UAA group to an external group from an identity provider |
| [`unmap-group`](commands/unmap-group.md) | Remove a mapping between a UAA group and an external group |
| [`list-group-mappings`](commands/list-group-mappings.md) | List all mappings between UAA groups and external groups |

## Miscellaneous

| Command | Description |
|---------|-------------|
| [`curl`](commands/curl.md) | Make an authenticated HTTP request to a UAA endpoint |
| [`userinfo`](commands/userinfo.md) | See claims about the authenticated user |
