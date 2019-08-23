package fixtures

const ExternalGroupsApiResponse = `{
  "resources": [
    {
      "displayName": "organizations.acme",
      "externalGroup": "cn=test_org,ou=people,o=springsource,o=org",
      "groupId": "59d28fbb-456c-4ab1-a1e3-0a5575ec4529",
      "origin": "ldap"
    }
  ],
  "startIndex": 1,
  "itemsPerPage": 1,
  "totalResults": 1,
  "schemas": [
    "urn:scim:schemas:core:1.0"
  ]
}`

const ExternalGroupsApiResponseInsufficientScope = `{
 "error": "insufficient_scope",
 "error_description": "Insufficient scope for this resource",
 "scope": "uaa.admin scim.read zones.uaa.admin"
}`
