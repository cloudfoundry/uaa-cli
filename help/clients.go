package help

func CreateClient() string {
	return `
  Clients wishing to obtain tokens representing a user's delegated authorization
  must perform a one-time registration with the UAA before initiating any OAuth
  flows. There are many configurable properties on clients. The most important
  ones are the selection of the authorized grant types for your client.
  Depending on the grant type you choose other properties may or may not be
  required.

  Most errors that occur during create-client can be easily diagnosed by using
  the --trace option to inspect the error response sent from the UAA.

EXAMPLE USAGE

  uaa create-client shinymail \
                    --client_secret secret \
                    --authorized_grant_types authorization_code \
                    --redirect_uri http://localhost:8080/*,http://email-sass.om/callback \
                    --scope mail.send,mail.read \
                    --display_name "Shinymail Web Mail Reader"

  uaa create-client background_emailer \
                    --client_secret secret \
                    --authorized_grant_types client_credentials \
                    --authorities notifications.write \
                    --display_name "Weekly newsletter email service"

  uaa create-client single_page_todo_app \
                    --authorized_grant_types implicit \
                    --redirect_uri http://localhost:8080/*,http://reactapp.com/callback \
                    --scope todo.read,todo.write \
                    --display_name "A Single-Page Todo App"

  uaa create-client trusted_cli \
                    --client_secret mumstheword \
                    --authorized_grant_types password \
                    --scope cloud_controller.admin,uaa.admin,network.write

  uaa create-client trusted_cli_copy \
                    --clone trusted_cli \
                    --client_secret donttellanyone
`
}
