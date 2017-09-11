package help

func RefreshToken() string {
	return `USAGE

  uaa target UAA_URL
  uaa get-password-token CLIENT_ID -s CLIENT_SECRET -u USERNAME -p PASSWORD
  uaa context
  uaa refresh-token -s CLIENT_SECRET
  uaa context # the access_token should now be updated

  The refresh-token command is used by authorization_code and password clients
  to obtain a new, unexpired access_token token from the UAA. Refresh tokens are
  long-lived and should be kept confidential by clients.

TROUBLESHOOTING FAQ

  Scenario: You do not have refresh_token in your active context.

    - If your context shows grant_type implicit: Implicit Clients are unable to
      maintain secrets and for that reason are not issued refresh_tokens.

    - If your context shows grant_type authorization_code: Clients using the
      authorization_code grant type are only issued refresh tokens if they are
      also registered with the refresh_token in the authorized_grant_types list.
      Use an administrative client to view your client registration with "uaa
      get-client" and ensure refresh_token is in the authorized_grant_types
      list.

    - If your context shows grant_type password: Clients using the password
      grant type are only issued refresh tokens if they are also configured with
      the refresh_token in the authorized_grant_types list. Use an
      administrative client to view your client registration with "uaa
      get-client" and ensure refresh_token is in the authorized_grant_types
      list.

    - If your context shows grant_type client_credentials: Refresh tokens are
      never issued when using the client_credentials grant flow. Refresh tokens
      are used to get a renewed access_token representing a user's delegated
      authorization without the need for a user to reauthorize.  With
      client_credentials grant type, there is no user and so you may obtain
      another token at any time using your client_id and client_secret with the
      client_credentials grant flow.
`
}
