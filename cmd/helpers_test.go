package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	"fmt"
	"github.com/cloudfoundry-community/go-uaa"
	"net/http"

	. "github.com/onsi/gomega/ghttp"
)

var buildConfig = func(target string) config.Config {
	cfg := config.NewConfigWithServerURL(target)
	ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
	cfg.AddContext(ctx)

	return cfg
}

var mockGroupLookup = func(id, groupname string) {
	server.RouteToHandler("GET", "/Groups", CombineHandlers(
		VerifyRequest("GET", "/Groups", fmt.Sprintf("filter=displayName+eq+%%22%s%%22&startIndex=1&count=100", groupname)),
		RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.Group{ID: id, DisplayName: groupname})),
	))
}

var mockExternalGroupMapping = func(externalGroupname, internalGroupId, internalGroupname, origin string) {
	server.RouteToHandler("POST", "/Groups/External", CombineHandlers(
		VerifyRequest("POST", "/Groups/External"),
		VerifyJSONRepresenting(map[string]interface{}{
			"groupId":       internalGroupId,
			"externalGroup": externalGroupname,
			"origin":        origin,
		}),
		RespondWith(http.StatusCreated, fixtures.EntityResponse(
			uaa.GroupMapping{
				GroupID:       internalGroupId,
				ExternalGroup: externalGroupname,
				DisplayName:   internalGroupname,
				Origin:        origin,
				Schemas:       []string{"urn:scim:schemas:core:1.0"},
			})),
	))
}

var mockExternalGroupUnmapping = func(externalGroupname, internalGroupId, internalGroupname, origin string) {
	path := fmt.Sprintf("/Groups/External/groupId/%s/externalGroup/%s/origin/%s", internalGroupId, externalGroupname, origin)
	server.RouteToHandler("DELETE", path, CombineHandlers(
		VerifyRequest("DELETE", path),
		RespondWith(http.StatusOK, fixtures.EntityResponse(
			uaa.GroupMapping{
				GroupID:       internalGroupId,
				ExternalGroup: externalGroupname,
				DisplayName:   internalGroupname,
				Origin:        origin,
				Schemas:       []string{"urn:scim:schemas:core:1.0"},
			})),
	))
}
