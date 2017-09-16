package fixtures

import (
	"encoding/json"
	"fmt"
)

const MarcusUserResponse = `{
	  "id" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
	  "externalId" : "marcus-user",
	  "meta" : {
		"version" : 1,
		"created" : "2017-01-15T16:54:15.677Z",
		"lastModified" : "2017-08-15T16:54:15.677Z"
	  },
	  "userName" : "marcus@stoicism.com",
	  "name" : {
		"familyName" : "Aurelius",
		"givenName" : "Marcus"
	  },
	  "emails" : [ {
		"value" : "marcus@stoicism.com",
		"primary" : false
	  } ],
	  "groups" : [ {
		"value" : "ac2ab20e-0a2d-4b68-82e4-817ee6b258b4",
		"display" : "philosophy.read",
		"type" : "DIRECT"
	  }, {
		"value" : "110b2434-4a30-439b-b5fc-f4cf47fc04f0",
		"display" : "philosophy.write",
		"type" : "DIRECT"
	  }],
	  "approvals" : [ {
		"userId" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
		"clientId" : "shinyclient",
		"scope" : "philosophy.read",
		"status" : "APPROVED",
		"lastUpdatedAt" : "2017-08-15T16:54:15.765Z",
		"expiresAt" : "2017-08-15T16:54:25.765Z"
	  }, {
		"userId" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
		"clientId" : "identity",
		"scope" : "uaa.user",
		"status" : "APPROVED",
		"lastUpdatedAt" : "2017-08-15T16:54:45.767Z",
		"expiresAt" : "2017-08-15T16:54:45.767Z"
	  } ],
	  "phoneNumbers" : [ {
		"value" : "5555555555"
	  } ],
	  "active" : true,
	  "verified" : true,
	  "origin" : "uaa",
	  "zoneId" : "uaa",
	  "passwordLastModified" : "2017-08-15T16:54:15.000Z",
	  "previousLogonTime" : 1502816055768,
	  "lastLogonTime" : 1502816055768,
	  "schemas" : [ "urn:scim:schemas:core:1.0" ]
	}`

const DrSeussUserResponse = `{
	  "id" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
	  "externalId" : "seuss-user",
	  "meta" : {
		"version" : 1,
		"created" : "2017-01-15T16:54:15.677Z",
		"lastModified" : "2017-08-15T16:54:15.677Z"
	  },
	  "userName" : "drseuss@whoville.com",
	  "name" : {
		"familyName" : "Theodore",
		"givenName" : "Giesel"
	  },
	  "emails" : [ {
		"value" : "drseuss@whoville.com",
		"primary" : true
	  } ],
	  "groups" : [ {
		"value" : "ac2ab20e-0a2d-4b68-82e4-817ee6b258b4",
		"display" : "cat_in_hat.read",
		"type" : "DIRECT"
	  }, {
		"value" : "110b2434-4a30-439b-b5fc-f4cf47fc04f0",
		"display" : "cat_in_hat.write",
		"type" : "DIRECT"
	  }],
	  "approvals" : [ {
		"userId" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
		"clientId" : "shinyclient",
		"scope" : "cat_in_hat.read",
		"status" : "APPROVED",
		"lastUpdatedAt" : "2017-08-15T16:54:15.765Z",
		"expiresAt" : "2017-08-15T16:54:25.765Z"
	  }, {
		"userId" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
		"clientId" : "identity",
		"scope" : "cat_in_hat.write",
		"status" : "APPROVED",
		"lastUpdatedAt" : "2017-08-15T16:54:45.767Z",
		"expiresAt" : "2017-08-15T16:54:45.767Z"
	  } ],
	  "phoneNumbers" : [ {
		"value" : "5555555555"
	  } ],
	  "active" : true,
	  "verified" : true,
	  "origin" : "uaa",
	  "zoneId" : "uaa",
	  "passwordLastModified" : "2017-08-15T16:54:15.000Z",
	  "previousLogonTime" : 1502816055768,
	  "lastLogonTime" : 1502816055768,
	  "schemas" : [ "urn:scim:schemas:core:1.0" ]
	}`

const PaginatedResponseTmpl = `{
		"resources": [%v,%v],
		"startIndex" : 1,
		"itemsPerPage" : 50,
		"totalResults" : 2,
		"schemas" : [ "urn:scim:schemas:core:1.0"]
	}`

func PaginatedResponse(resources ...interface{}) string {
	bytes, _ := json.Marshal(resources)

	return fmt.Sprintf(`{
		"resources": %v,
		"startIndex" : 1,
		"itemsPerPage" : 50,
		"totalResults" : %v,
		"schemas" : [ "urn:scim:schemas:core:1.0"]
	}`, string(bytes), len(resources))
}
