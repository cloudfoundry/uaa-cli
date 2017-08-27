package help

func PasswordGrant() string {
  return `USAGE

  uaa target UAA_URL
  uaa get-password-token CLIENT_ID -s CLIENT_SECRET -u USERNAME -p PASSWORD

  After successfully running this command, the token is added to the CLI's
  current context. Access tokens saved in the context will be attached to subsequent
  requests when attempting to use CLI commands that hit UAA endpoints requiring
  Authorization.

BACKGROUND

  The "Resource Owner Password Credentials" grant type, often referred to by the
  shorter name "password grant", is one of the four authorization flows
  described in the RFC 6749. Like the other flows described in that spec, the
  goal of this flow is to obtain access tokens for a Client application that
  represent the delegated authorization of a User.  The Client may then present
  the token to a Resource Service in an Authorization header to perform
  priviledged action on the User's behalf.

          +----------+
          | Resource |
          |  Owner   |
          |  (User)  |
          +----------+
               v
               |    Resource Owner
              (A)   Password Credentials
               |    (i.e. username and password)
               v
          +----------+                                  +---------------+
          |          |>--(B)---- Resource Owner ------->|               |
          |          |         Password Credentials     | Authorization |
          |  Client  |                                  |     Server    |
          |(this CLI)|<--(C)---- Access Token ---------<|     (UAA)     |
          |          |    (w/ Optional Refresh Token)   |               |
          +----------+                                  +---------------+


  The get-password-token command allows you to test your client registration.
  By providing your client_id and client_secret to this CLI, it is able to
  play the part of your Client application in the Resource Owner Password
  Credentials flow.

WHEN SHOULD PASSWORD GRANT CLIENTS BE USED

  This flow is typically only used for Client applications that are considered
  highly trusted by the User, since they must provide their username and
  password directly to the Client who will use them on their behalf. CLIs or
  other native applications running on user-owned hardware are good candidates
  for the password grant flow.

  The Cloud Foundry CLI and the Credhub CLI are two examples of password grant
  clients within the CF ecosystem.

TROUBLESHOOTING FAQ

  Scenario: You are unable to get a token using get-password-token.

    - Ensure you are using valid client_id, client_secret, username, and password.
      If you are very confident in these values, run with --trace and make sure
      you have targeted the correct UAA for the credentials you are passing.

    - Ensure that "password" in the list of authorized_grant_types for your client.

  Scenario: You got a token but it does not have the expected scopes.

    - Verify "scope" field of your client registration includes the scopes you
      want to be included in the token.

    - Verify the user whose username and password you are using has the correct
      group memberships. The scopes in the issued token will be the intersection
      of the scopes your client requests (by listing them in the client
      registration) and the group memberships, or permissions, that the user
      has.
  `
}