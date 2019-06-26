package uaa_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var suite spec.Suite

func init() {
	suite = spec.New("uaa", spec.Report(report.Terminal{}))
	suite("new", testNew)
	suite("clientExtra", testClientExtra)
	suite("curl", testCurl)
	suite("groupsExtra", testGroupsExtra)
	suite("isHealthy", testIsHealthy)
	suite("info", testInfo)
	suite("me", testMe)
	suite("tokenKey", testTokenKey)
	suite("tokenKeys", testTokenKeys)
	suite("buildSubdomainURL", testBuildSubdomainURL)
	suite("users", testUsers)

	// Generated
	suite("client", testClient)
	suite("group", testGroup)
	suite("identityZone", testIdentityZone)
	suite("mfaProvider", testMFAProvider)
	suite("user", testUser)
}

func TestUAA(t *testing.T) {
	suite.Run(t)
}
