package uaa_test

import uaa "github.com/cloudfoundry-community/go-uaa"

var mfaproviderResponse string = `{
	"id": "00000000-0000-0000-0000-000000000001",
	"name": "sampleGoogleMfaProvider8ZFKcx",
	"identityZoneId": "uaa",
	"config": {
		"issuer": "uaa",
		"providerDescription": "Google MFA for default zone"
	},
	"type": "google-authenticator",
	"created": 1529690500934,
	"last_modified": 1529690500934
}`

var mfaproviderListResponse string = `[{
	"id": "00000000-0000-0000-0000-000000000001",
	"name": "sampleGoogleMfaProviderCJTjGb",
	"identityZoneId": "uaa",
	"config" : {
		"issuer": "uaa",
		"providerDescription": "Google MFA for default zone"
	},
	"type": "google-authenticator",
	"created": 1529690500558,
	"last_modified": 1529690500558
}, {
	"id": "00000000-0000-0000-0000-000000000002",
	"name": "sampleGoogleMfaProviderUKaW73",
	"identityZoneId": "uaa",
	"config": {
		"issuer": "uaa",
		"providerDescription": "Google MFA for default zone"
	},
	"type": "google-authenticator",
	"created" : 1529690500430,
	"last_modified" : 1529690500430
}]`

var testMFAProviderValue uaa.MFAProvider = uaa.MFAProvider{
	ID:             "00000000-0000-0000-0000-000000000001",
	Name:           "sampleGoogleMfaProvider8ZFKcx",
	IdentityZoneID: "uaa",
	Config: uaa.MFAProviderConfig{
		Issuer:              "uaa",
		ProviderDescription: "Google MFA for default zone",
	},
	Type:         "google-authenticator",
	Created:      1529690500934,
	LastModified: 1529690500934,
}

var testMFAProviderJSON string = `{
	"id": "00000000-0000-0000-0000-000000000001",
	"name": "sampleGoogleMfaProvider8ZFKcx",
	"identityZoneId": "uaa",
	"config": {
		"issuer": "uaa",
		"providerDescription": "Google MFA for default zone"
	},
	"type": "google-authenticator",
	"created": 1529690500934,
	"last_modified": 1529690500934
}`
