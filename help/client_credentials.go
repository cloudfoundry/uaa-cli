package help

func ClientCredentials() string {
	return `USAGE

  uaa target UAA_URL
  uaa get-client-credentials-token CLIENT_ID -s CLIENT_SECRET

  After successfully running this command, the token is added to the CLI's
  current context. Access tokens saved in the context will be attached to subsequent
  requests when attempting to use CLI commands that hit UAA endpoints requiring
  Authorization.

BACKGROUND

  The Client Credentials grant type is one of the four authorization flows
  described in the RFC 6749. Like the other flows described in that spec, the
  goal of this flow is for Clients to obtain access tokens needed to perform
  priviledged actions against Resource Server APIs.  Unlike the other flows
  defined in the OAuth2 spec, the client_credentials grant type is used by
  Clients performing actions on their own behalf without the involvement of a
  User.

        +---------+                                  +---------------+
        |         |                                  |               |
        |         |>--(A)- Client Authentication --->| Authorization |
        | Client  |                                  |     Server    |
        |         |<--(B)---- Access Token ---------<|               |
        |         |                                  |               |
        +---------+                                  +---------------+

  The get-client-credentials-token command allows you to test your client
  registration.  By providing your client_id and client_secret to this CLI, it
  is able to play the part of your Client application in the flow.

WHEN SHOULD PASSWORD GRANT CLIENTS BE USED

  The client_credentials flow is typically used by Client applications wanting
  to perform service-to-service operations on its own behalf against a Resource
  Server. Many components within Cloud Foundry leaverage the client_credentials
  grant type to secure their service-to-service interactions.

  For example, if your client is a web application that need to send out a
  weekly email newsletter, it might make scheduled calls to an
  email/notifications service with a client_credentials token. In this example,
  the client_credentials grant type is appropriate because the client is calling
  the notification service on its own behalf and the resources it is accessing
  are not owned in any meaningful sense by the user.

TROUBLESHOOTING FAQ

  Scenario: You are unable to get a token using get-client-credentials-token.

    - Ensure you are using valid client_id and client_secret for the targeted
      UAA.  If you are very confident in these values, double check that you are
      targetting the correct UAA.

    - Ensure that "client_credentials" is included in the list of
      authorized_grant_types for your client. This flow may only be used with
      clients that have registered for the client_credentials grant type.

  Scenario: You got a token but it does not have the expected scopes.

    - Verify the "authorities" field of your client registration includes the
      scopes you want to be included in the token. The "scope" field of your
      client registration is not considered by the UAA when authoring tokens for
      the client_credentials grant type.
  `
}
