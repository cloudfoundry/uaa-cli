package help

func Userinfo() string {
	return `  With a valid user token saved in context, use the me command to see
  information about the authenticated user.

      uaa target UAA_URL
      uaa get-password-token CLIENT_ID -s CLIENT_SECRET -u USERNAME -p PASSWORD
      uaa me

BACKGROUND

  The me command calls the UAA's /userinfo endpoint. This endpoint has been
  implemented in accordance with the OIDC 1.1 spec.

  The /userinfo endpoint is an OAuth 2.0 protected resource that returns claims
  about the authenticated user. To obtain the requested claims about the
  user, the client makes a request to the /userinfo endpoint using an access
  token obtained through OpenID Connect Authentication. 

TROUBLESHOOTING FAQ

  Scenario: You have obtained a token but are unable to see results from the
  "me" command.

  - Verify your token is still valid by using another command that requires
    authentication. If it is no longer valid, request another token and try
    again.

  - Check whether the output of "uaa context" shows "openid" in the list of
    scopes. This scope is required to access the /userinfo endpoint and it won't
    appear in your token unless requested by the client. You may need to update
    your client registration to include "openid" in the scope list.
`
}
