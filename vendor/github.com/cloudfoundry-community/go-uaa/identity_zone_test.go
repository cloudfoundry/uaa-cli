package uaa_test

import uaa "github.com/cloudfoundry-community/go-uaa"

var identityzoneResponse string = `{
  "id" : "00000000-0000-0000-0000-000000000001",
  "subdomain" : "twiglet-get",
  "config" : {
    "clientSecretPolicy" : {
      "minLength" : -1,
      "maxLength" : -1,
      "requireUpperCaseCharacter" : -1,
      "requireLowerCaseCharacter" : -1,
      "requireDigit" : -1,
      "requireSpecialCharacter" : -1
    },
    "tokenPolicy" : {
      "accessTokenValidity" : 3600,
      "refreshTokenValidity" : 7200,
      "jwtRevocable" : false,
      "refreshTokenUnique" : false,
      "refreshTokenFormat" : "jwt",
      "activeKeyId" : "active-key-1"
    },
    "samlConfig" : {
      "assertionSigned" : true,
      "requestSigned" : true,
      "wantAssertionSigned" : true,
      "wantAuthnRequestSigned" : false,
      "assertionTimeToLiveSeconds" : 600,
      "activeKeyId" : "legacy-saml-key",
      "keys" : {
        "legacy-saml-key" : {
          "certificate" : "-----BEGIN CERTIFICATE-----\nMIICEjCCAXsCAg36MA0GCSqGSIb3DQEBBQUAMIGbMQswCQYDVQQGEwJKUDEOMAwG\nA1UECBMFVG9reW8xEDAOBgNVBAcTB0NodW8ta3UxETAPBgNVBAoTCEZyYW5rNERE\nMRgwFgYDVQQLEw9XZWJDZXJ0IFN1cHBvcnQxGDAWBgNVBAMTD0ZyYW5rNEREIFdl\nYiBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBmcmFuazRkZC5jb20wHhcNMTIw\nODIyMDUyNjU0WhcNMTcwODIxMDUyNjU0WjBKMQswCQYDVQQGEwJKUDEOMAwGA1UE\nCAwFVG9reW8xETAPBgNVBAoMCEZyYW5rNEREMRgwFgYDVQQDDA93d3cuZXhhbXBs\nZS5jb20wXDANBgkqhkiG9w0BAQEFAANLADBIAkEAm/xmkHmEQrurE/0re/jeFRLl\n8ZPjBop7uLHhnia7lQG/5zDtZIUC3RVpqDSwBuw/NTweGyuP+o8AG98HxqxTBwID\nAQABMA0GCSqGSIb3DQEBBQUAA4GBABS2TLuBeTPmcaTaUW/LCB2NYOy8GMdzR1mx\n8iBIu2H6/E2tiY3RIevV2OW61qY2/XRQg7YPxx3ffeUugX9F4J/iPnnu1zAxxyBy\n2VguKv4SWjRFoRkIfIlHX0qVviMhSlNy2ioFLy7JcPZb+v3ftDGywUqcBiVDoea0\nHn+GmxZA\n-----END CERTIFICATE-----\n"
        }
      },
      "entityID" : "cloudfoundry-saml-login",
      "disableInResponseToCheck" : false,
      "certificate" : "-----BEGIN CERTIFICATE-----\nMIICEjCCAXsCAg36MA0GCSqGSIb3DQEBBQUAMIGbMQswCQYDVQQGEwJKUDEOMAwG\nA1UECBMFVG9reW8xEDAOBgNVBAcTB0NodW8ta3UxETAPBgNVBAoTCEZyYW5rNERE\nMRgwFgYDVQQLEw9XZWJDZXJ0IFN1cHBvcnQxGDAWBgNVBAMTD0ZyYW5rNEREIFdl\nYiBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBmcmFuazRkZC5jb20wHhcNMTIw\nODIyMDUyNjU0WhcNMTcwODIxMDUyNjU0WjBKMQswCQYDVQQGEwJKUDEOMAwGA1UE\nCAwFVG9reW8xETAPBgNVBAoMCEZyYW5rNEREMRgwFgYDVQQDDA93d3cuZXhhbXBs\nZS5jb20wXDANBgkqhkiG9w0BAQEFAANLADBIAkEAm/xmkHmEQrurE/0re/jeFRLl\n8ZPjBop7uLHhnia7lQG/5zDtZIUC3RVpqDSwBuw/NTweGyuP+o8AG98HxqxTBwID\nAQABMA0GCSqGSIb3DQEBBQUAA4GBABS2TLuBeTPmcaTaUW/LCB2NYOy8GMdzR1mx\n8iBIu2H6/E2tiY3RIevV2OW61qY2/XRQg7YPxx3ffeUugX9F4J/iPnnu1zAxxyBy\n2VguKv4SWjRFoRkIfIlHX0qVviMhSlNy2ioFLy7JcPZb+v3ftDGywUqcBiVDoea0\nHn+GmxZA\n-----END CERTIFICATE-----\n"
    },
    "corsPolicy" : {
      "xhrConfiguration" : {
        "allowedOrigins" : [ ".*" ],
        "allowedOriginPatterns" : [ ],
        "allowedUris" : [ ".*" ],
        "allowedUriPatterns" : [ ],
        "allowedHeaders" : [ "Accept", "Authorization", "Content-Type" ],
        "allowedMethods" : [ "GET" ],
        "allowedCredentials" : false,
        "maxAge" : 1728000
      },
      "defaultConfiguration" : {
        "allowedOrigins" : [ ".*" ],
        "allowedOriginPatterns" : [ ],
        "allowedUris" : [ ".*" ],
        "allowedUriPatterns" : [ ],
        "allowedHeaders" : [ "Accept", "Authorization", "Content-Type" ],
        "allowedMethods" : [ "GET" ],
        "allowedCredentials" : false,
        "maxAge" : 1728000
      }
    },
    "links" : {
      "logout" : {
        "redirectUrl" : "/login",
        "redirectParameterName" : "redirect",
        "disableRedirectParameter" : false,
        "whitelist" : null
      },
      "homeRedirect" : "http://my.hosted.homepage.com/",
      "selfService" : {
        "selfServiceLinksEnabled" : true,
        "signup" : null,
        "passwd" : null
      }
    },
    "prompts" : [ {
      "name" : "username",
      "type" : "text",
      "text" : "Email"
    }, {
      "name" : "password",
      "type" : "password",
      "text" : "Password"
    }, {
      "name" : "passcode",
      "type" : "password",
      "text" : "Temporary Authentication Code (Get on at /passcode)"
    } ],
    "idpDiscoveryEnabled" : false,
    "branding" : {
      "companyName" : "Test Company",
      "productLogo" : "VGVzdFByb2R1Y3RMb2dv",
      "squareLogo" : "VGVzdFNxdWFyZUxvZ28=",
      "footerLegalText" : "Test footer legal text",
      "footerLinks" : {
        "Support" : "http://support.example.com"
      },
      "banner" : {
        "logo" : "VGVzdFByb2R1Y3RMb2dv",
        "text" : "Announcement",
        "textColor" : "#000000",
        "backgroundColor" : "#89cff0",
        "link" : "http://announce.example.com"
      },
      "consent" : {
        "text" : "Some Policy",
        "link" : "http://policy.example.com"
      }
    },
    "accountChooserEnabled" : false,
    "userConfig" : {
      "defaultGroups" : [ "openid", "password.write", "uaa.user", "approvals.me", "profile", "roles", "user_attributes", "uaa.offline_token" ]
    },
    "mfaConfig" : {
      "enabled" : false
    },
    "issuer" : "http://localhost:8080/uaa"
  },
  "name" : "The Twiglet Zone",
  "version" : 0,
  "created" : 1527032707746,
  "last_modified" : 1527032707746
}`

var identityzoneListResponse string = `[{
  "id" : "00000000-0000-0000-0000-000000000001",
  "subdomain" : "twiglet-get",
  "config" : {
    "clientSecretPolicy" : {
      "minLength" : -1,
      "maxLength" : -1,
      "requireUpperCaseCharacter" : -1,
      "requireLowerCaseCharacter" : -1,
      "requireDigit" : -1,
      "requireSpecialCharacter" : -1
    },
    "tokenPolicy" : {
      "accessTokenValidity" : 3600,
      "refreshTokenValidity" : 7200,
      "jwtRevocable" : false,
      "refreshTokenUnique" : false,
      "refreshTokenFormat" : "jwt",
      "activeKeyId" : "active-key-1"
    },
    "samlConfig" : {
      "assertionSigned" : true,
      "requestSigned" : true,
      "wantAssertionSigned" : true,
      "wantAuthnRequestSigned" : false,
      "assertionTimeToLiveSeconds" : 600,
      "activeKeyId" : "legacy-saml-key",
      "keys" : {
        "legacy-saml-key" : {
          "certificate" : "-----BEGIN CERTIFICATE-----\nMIICEjCCAXsCAg36MA0GCSqGSIb3DQEBBQUAMIGbMQswCQYDVQQGEwJKUDEOMAwG\nA1UECBMFVG9reW8xEDAOBgNVBAcTB0NodW8ta3UxETAPBgNVBAoTCEZyYW5rNERE\nMRgwFgYDVQQLEw9XZWJDZXJ0IFN1cHBvcnQxGDAWBgNVBAMTD0ZyYW5rNEREIFdl\nYiBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBmcmFuazRkZC5jb20wHhcNMTIw\nODIyMDUyNjU0WhcNMTcwODIxMDUyNjU0WjBKMQswCQYDVQQGEwJKUDEOMAwGA1UE\nCAwFVG9reW8xETAPBgNVBAoMCEZyYW5rNEREMRgwFgYDVQQDDA93d3cuZXhhbXBs\nZS5jb20wXDANBgkqhkiG9w0BAQEFAANLADBIAkEAm/xmkHmEQrurE/0re/jeFRLl\n8ZPjBop7uLHhnia7lQG/5zDtZIUC3RVpqDSwBuw/NTweGyuP+o8AG98HxqxTBwID\nAQABMA0GCSqGSIb3DQEBBQUAA4GBABS2TLuBeTPmcaTaUW/LCB2NYOy8GMdzR1mx\n8iBIu2H6/E2tiY3RIevV2OW61qY2/XRQg7YPxx3ffeUugX9F4J/iPnnu1zAxxyBy\n2VguKv4SWjRFoRkIfIlHX0qVviMhSlNy2ioFLy7JcPZb+v3ftDGywUqcBiVDoea0\nHn+GmxZA\n-----END CERTIFICATE-----\n"
        }
      },
      "entityID" : "cloudfoundry-saml-login",
      "disableInResponseToCheck" : false,
      "certificate" : "-----BEGIN CERTIFICATE-----\nMIICEjCCAXsCAg36MA0GCSqGSIb3DQEBBQUAMIGbMQswCQYDVQQGEwJKUDEOMAwG\nA1UECBMFVG9reW8xEDAOBgNVBAcTB0NodW8ta3UxETAPBgNVBAoTCEZyYW5rNERE\nMRgwFgYDVQQLEw9XZWJDZXJ0IFN1cHBvcnQxGDAWBgNVBAMTD0ZyYW5rNEREIFdl\nYiBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBmcmFuazRkZC5jb20wHhcNMTIw\nODIyMDUyNjU0WhcNMTcwODIxMDUyNjU0WjBKMQswCQYDVQQGEwJKUDEOMAwGA1UE\nCAwFVG9reW8xETAPBgNVBAoMCEZyYW5rNEREMRgwFgYDVQQDDA93d3cuZXhhbXBs\nZS5jb20wXDANBgkqhkiG9w0BAQEFAANLADBIAkEAm/xmkHmEQrurE/0re/jeFRLl\n8ZPjBop7uLHhnia7lQG/5zDtZIUC3RVpqDSwBuw/NTweGyuP+o8AG98HxqxTBwID\nAQABMA0GCSqGSIb3DQEBBQUAA4GBABS2TLuBeTPmcaTaUW/LCB2NYOy8GMdzR1mx\n8iBIu2H6/E2tiY3RIevV2OW61qY2/XRQg7YPxx3ffeUugX9F4J/iPnnu1zAxxyBy\n2VguKv4SWjRFoRkIfIlHX0qVviMhSlNy2ioFLy7JcPZb+v3ftDGywUqcBiVDoea0\nHn+GmxZA\n-----END CERTIFICATE-----\n"
    },
    "corsPolicy" : {
      "xhrConfiguration" : {
        "allowedOrigins" : [ ".*" ],
        "allowedOriginPatterns" : [ ],
        "allowedUris" : [ ".*" ],
        "allowedUriPatterns" : [ ],
        "allowedHeaders" : [ "Accept", "Authorization", "Content-Type" ],
        "allowedMethods" : [ "GET" ],
        "allowedCredentials" : false,
        "maxAge" : 1728000
      },
      "defaultConfiguration" : {
        "allowedOrigins" : [ ".*" ],
        "allowedOriginPatterns" : [ ],
        "allowedUris" : [ ".*" ],
        "allowedUriPatterns" : [ ],
        "allowedHeaders" : [ "Accept", "Authorization", "Content-Type" ],
        "allowedMethods" : [ "GET" ],
        "allowedCredentials" : false,
        "maxAge" : 1728000
      }
    },
    "links" : {
      "logout" : {
        "redirectUrl" : "/login",
        "redirectParameterName" : "redirect",
        "disableRedirectParameter" : false,
        "whitelist" : null
      },
      "homeRedirect" : "http://my.hosted.homepage.com/",
      "selfService" : {
        "selfServiceLinksEnabled" : true,
        "signup" : null,
        "passwd" : null
      }
    },
    "prompts" : [ {
      "name" : "username",
      "type" : "text",
      "text" : "Email"
    }, {
      "name" : "password",
      "type" : "password",
      "text" : "Password"
    }, {
      "name" : "passcode",
      "type" : "password",
      "text" : "Temporary Authentication Code (Get on at /passcode)"
    } ],
    "idpDiscoveryEnabled" : false,
    "branding" : {
      "companyName" : "Test Company",
      "productLogo" : "VGVzdFByb2R1Y3RMb2dv",
      "squareLogo" : "VGVzdFNxdWFyZUxvZ28=",
      "footerLegalText" : "Test footer legal text",
      "footerLinks" : {
        "Support" : "http://support.example.com"
      },
      "banner" : {
        "logo" : "VGVzdFByb2R1Y3RMb2dv",
        "text" : "Announcement",
        "textColor" : "#000000",
        "backgroundColor" : "#89cff0",
        "link" : "http://announce.example.com"
      },
      "consent" : {
        "text" : "Some Policy",
        "link" : "http://policy.example.com"
      }
    },
    "accountChooserEnabled" : false,
    "userConfig" : {
      "defaultGroups" : [ "openid", "password.write", "uaa.user", "approvals.me", "profile", "roles", "user_attributes", "uaa.offline_token" ]
    },
    "mfaConfig" : {
      "enabled" : false
    },
    "issuer" : "http://localhost:8080/uaa"
  },
  "name" : "The Twiglet Zone",
  "version" : 0,
  "created" : 1527032707746,
  "last_modified" : 1527032707746
},{
  "id" : "00000000-0000-0000-0000-000000000002",
  "subdomain" : "twiglet-get",
  "config" : {
    "clientSecretPolicy" : {
      "minLength" : -1,
      "maxLength" : -1,
      "requireUpperCaseCharacter" : -1,
      "requireLowerCaseCharacter" : -1,
      "requireDigit" : -1,
      "requireSpecialCharacter" : -1
    },
    "tokenPolicy" : {
      "accessTokenValidity" : 3600,
      "refreshTokenValidity" : 7200,
      "jwtRevocable" : false,
      "refreshTokenUnique" : false,
      "refreshTokenFormat" : "jwt",
      "activeKeyId" : "active-key-1"
    },
    "samlConfig" : {
      "assertionSigned" : true,
      "requestSigned" : true,
      "wantAssertionSigned" : true,
      "wantAuthnRequestSigned" : false,
      "assertionTimeToLiveSeconds" : 600,
      "activeKeyId" : "legacy-saml-key",
      "keys" : {
        "legacy-saml-key" : {
          "certificate" : "-----BEGIN CERTIFICATE-----\nMIICEjCCAXsCAg36MA0GCSqGSIb3DQEBBQUAMIGbMQswCQYDVQQGEwJKUDEOMAwG\nA1UECBMFVG9reW8xEDAOBgNVBAcTB0NodW8ta3UxETAPBgNVBAoTCEZyYW5rNERE\nMRgwFgYDVQQLEw9XZWJDZXJ0IFN1cHBvcnQxGDAWBgNVBAMTD0ZyYW5rNEREIFdl\nYiBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBmcmFuazRkZC5jb20wHhcNMTIw\nODIyMDUyNjU0WhcNMTcwODIxMDUyNjU0WjBKMQswCQYDVQQGEwJKUDEOMAwGA1UE\nCAwFVG9reW8xETAPBgNVBAoMCEZyYW5rNEREMRgwFgYDVQQDDA93d3cuZXhhbXBs\nZS5jb20wXDANBgkqhkiG9w0BAQEFAANLADBIAkEAm/xmkHmEQrurE/0re/jeFRLl\n8ZPjBop7uLHhnia7lQG/5zDtZIUC3RVpqDSwBuw/NTweGyuP+o8AG98HxqxTBwID\nAQABMA0GCSqGSIb3DQEBBQUAA4GBABS2TLuBeTPmcaTaUW/LCB2NYOy8GMdzR1mx\n8iBIu2H6/E2tiY3RIevV2OW61qY2/XRQg7YPxx3ffeUugX9F4J/iPnnu1zAxxyBy\n2VguKv4SWjRFoRkIfIlHX0qVviMhSlNy2ioFLy7JcPZb+v3ftDGywUqcBiVDoea0\nHn+GmxZA\n-----END CERTIFICATE-----\n"
        }
      },
      "entityID" : "cloudfoundry-saml-login",
      "disableInResponseToCheck" : false,
      "certificate" : "-----BEGIN CERTIFICATE-----\nMIICEjCCAXsCAg36MA0GCSqGSIb3DQEBBQUAMIGbMQswCQYDVQQGEwJKUDEOMAwG\nA1UECBMFVG9reW8xEDAOBgNVBAcTB0NodW8ta3UxETAPBgNVBAoTCEZyYW5rNERE\nMRgwFgYDVQQLEw9XZWJDZXJ0IFN1cHBvcnQxGDAWBgNVBAMTD0ZyYW5rNEREIFdl\nYiBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBmcmFuazRkZC5jb20wHhcNMTIw\nODIyMDUyNjU0WhcNMTcwODIxMDUyNjU0WjBKMQswCQYDVQQGEwJKUDEOMAwGA1UE\nCAwFVG9reW8xETAPBgNVBAoMCEZyYW5rNEREMRgwFgYDVQQDDA93d3cuZXhhbXBs\nZS5jb20wXDANBgkqhkiG9w0BAQEFAANLADBIAkEAm/xmkHmEQrurE/0re/jeFRLl\n8ZPjBop7uLHhnia7lQG/5zDtZIUC3RVpqDSwBuw/NTweGyuP+o8AG98HxqxTBwID\nAQABMA0GCSqGSIb3DQEBBQUAA4GBABS2TLuBeTPmcaTaUW/LCB2NYOy8GMdzR1mx\n8iBIu2H6/E2tiY3RIevV2OW61qY2/XRQg7YPxx3ffeUugX9F4J/iPnnu1zAxxyBy\n2VguKv4SWjRFoRkIfIlHX0qVviMhSlNy2ioFLy7JcPZb+v3ftDGywUqcBiVDoea0\nHn+GmxZA\n-----END CERTIFICATE-----\n"
    },
    "corsPolicy" : {
      "xhrConfiguration" : {
        "allowedOrigins" : [ ".*" ],
        "allowedOriginPatterns" : [ ],
        "allowedUris" : [ ".*" ],
        "allowedUriPatterns" : [ ],
        "allowedHeaders" : [ "Accept", "Authorization", "Content-Type" ],
        "allowedMethods" : [ "GET" ],
        "allowedCredentials" : false,
        "maxAge" : 1728000
      },
      "defaultConfiguration" : {
        "allowedOrigins" : [ ".*" ],
        "allowedOriginPatterns" : [ ],
        "allowedUris" : [ ".*" ],
        "allowedUriPatterns" : [ ],
        "allowedHeaders" : [ "Accept", "Authorization", "Content-Type" ],
        "allowedMethods" : [ "GET" ],
        "allowedCredentials" : false,
        "maxAge" : 1728000
      }
    },
    "links" : {
      "logout" : {
        "redirectUrl" : "/login",
        "redirectParameterName" : "redirect",
        "disableRedirectParameter" : false,
        "whitelist" : null
      },
      "homeRedirect" : "http://my.hosted.homepage.com/",
      "selfService" : {
        "selfServiceLinksEnabled" : true,
        "signup" : null,
        "passwd" : null
      }
    },
    "prompts" : [ {
      "name" : "username",
      "type" : "text",
      "text" : "Email"
    }, {
      "name" : "password",
      "type" : "password",
      "text" : "Password"
    }, {
      "name" : "passcode",
      "type" : "password",
      "text" : "Temporary Authentication Code (Get on at /passcode)"
    } ],
    "idpDiscoveryEnabled" : false,
    "branding" : {
      "companyName" : "Test Company",
      "productLogo" : "VGVzdFByb2R1Y3RMb2dv",
      "squareLogo" : "VGVzdFNxdWFyZUxvZ28=",
      "footerLegalText" : "Test footer legal text",
      "footerLinks" : {
        "Support" : "http://support.example.com"
      },
      "banner" : {
        "logo" : "VGVzdFByb2R1Y3RMb2dv",
        "text" : "Announcement",
        "textColor" : "#000000",
        "backgroundColor" : "#89cff0",
        "link" : "http://announce.example.com"
      },
      "consent" : {
        "text" : "Some Policy",
        "link" : "http://policy.example.com"
      }
    },
    "accountChooserEnabled" : false,
    "userConfig" : {
      "defaultGroups" : [ "openid", "password.write", "uaa.user", "approvals.me", "profile", "roles", "user_attributes", "uaa.offline_token" ]
    },
    "mfaConfig" : {
      "enabled" : false
    },
    "issuer" : "http://localhost:8080/uaa"
  },
  "name" : "The Twiglet Zone",
  "version" : 0,
  "created" : 1527032707746,
  "last_modified" : 1527032707746
}]`

var testIdentityZoneValue uaa.IdentityZone = uaa.IdentityZone{
	ID:        "00000000-0000-0000-0000-000000000001",
	Subdomain: "twiglet-get",
	Name:      "The Twiglet Zone",
}

var testIdentityZoneJSON string = `{
  "id" : "00000000-0000-0000-0000-000000000001",
  "subdomain" : "twiglet-get",
  "config": {},
  "name" : "The Twiglet Zone"
}`
