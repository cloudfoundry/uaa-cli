package help

func ImplicitGrant() string {
	return `USAGE

    uaa target UAA_URL
    uaa get-implicit-token CLIENT_ID --port REDIRECT_URI_PORT

  This command will launch a browser window where the user will be prompted to
  login and authorize the implicit grant client. After authorizing, the user is
  redirected to a local server process the CLI has started to recieve the
  callback from UAA.
  
  After successfully running this command, the access token is added to the
  CLI's current context. Access tokens saved in the context will be attached to
  subsequent requests when attempting to use CLI commands that hit UAA endpoints
  requiring Authorization headers.

  The value of the --port argument must correspond to a localhost redirect uri
  in the client registration for the implicit client. A fuller example,
  including the creation of the implicit client, would be

    uaa target UAA_URL
    uaa get-client-credentials-token admin -s ADMIN_CLIENT_SECRET
    uaa create-client my_implicit_client \
            --authorized_grant_types \
            --scope cloud_controller.read,cloud_controller.write \
            --redirect_uri https://myprodsite.com/**,http://localhost:8080/**
    uaa get-implicit-token my_implicit_client --port 8080
    uaa context

  The important thing to notice in the above example is that the port given in
  the arguments to get-implicit-token command must match the redirect_uri port
  in the create-client step.

BACKGROUND

  The implicit grant type is one of the four authorization flows described in
  the RFC 6749. Like the other flows described in that spec, the goal of this
  flow is to obtain access tokens for a Client application that represent the
  delegated authorization of a User.  The Client may then present the token to a
  Resource Service in an Authorization header to perform priviledged action on
  the User's behalf.

  Implicit grant type is most commonly used for single-page web applications.
  These are said to be "public clients" because they have no way of maintaining
  secrecy. After the user authorizes, an access token is transmistted to the
  client as part of the URI fragment, which means it may be available to
  unauthorized parties or malicious JavaScripts running in the user's browser.
  Consequently, implicit clients should be configured with small values for
  access_token_validity and the minimal set required permissions.
 
     +----------+
     | Resource |
     |  Owner   |
     | (User)   |
     +----------+
          ^
          |
         (B)
     +----|-----+          Client Identifier     +---------------+
     |         -+----(A)-- & Redirection URI --->|               |
     |  User-   |                                | Authorization |
     |  Agent  -|----(B)-- User authenticates -->|     Server    |
     | (Browser)|                                |    (UAA)      |
     |          |<---(C)--- Redirection URI ----<|               |
     |          |          with Access Token     +---------------+
     |          |            in Fragment
     |          |                                +-------------------+
     |          |----(D)--- Redirection URI ---->|     Web-Hosted    |
     |          |          without Fragment      |       Client      |
     |          |                                |      Resource     |
     |     (F)  |<---(E)------- Script ---------<| (Single-Page App) |
     |          |                                +-------------------+
     +-|--------+
       |    |
      (A)  (G) Access Token
       |    |
       ^    v
     +---------+
     |         |
     |  Client |
     |         |
     +---------+


TROUBLESHOOTING FAQ

  Scenario: After signing in, I see unknown redirect_uri error.

    - Try running "uaa get-client <your-implicit-client-id>" to view the client
      registration for your client. To get tokens when doing development on your
      local machine, including using this CLI, you must have an entry in the
      redirect_uri list for http://localhost:PORT where PORT is the value passed
      when invoking "uaa get-implicit-token CLIENT_ID --port PORT"

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
