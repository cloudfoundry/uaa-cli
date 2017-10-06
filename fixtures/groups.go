package fixtures

const AdminGroupResponse = `{
	"id" : "05a0c169-3592-4a45-b109-a16d9246e0ab",
	"meta" : {
		"version" : 1,
		"created" : "2017-01-15T16:54:15.677Z",
		"lastModified" : "2017-08-15T16:54:15.677Z"
	},
	"displayName" : "admin",
	"description" : "admin",
	"members" : [ {
		"origin" : "uaa",
		"type" : "USER",
		"value" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"
	} ],
	"zoneId" : "uaa",
	"schemas" : [ "urn:scim:schemas:core:1.0" ]
}`

const ReadGroupResponse = `{
	"id" : "ea777017-883e-48ba-800a-637c71409b5e",
	"meta" : {
		"version" : 1,
		"created" : "2017-01-15T16:54:15.677Z",
		"lastModified" : "2017-08-15T16:54:15.677Z"
	},
	"displayName" : "read",
	"members" : [ {
		"origin" : "uaa",
		"type" : "USER",
		"value" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"
	} ],
	"zoneId" : "uaa",
	"schemas" : [ "urn:scim:schemas:core:1.0" ]
}`
