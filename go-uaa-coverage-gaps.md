# go-uaa / UAA coverage gaps

Notes from auditing `github.com/cloudfoundry-community/go-uaa` (client library
used by uaa-cli) against the current UAA server source, triggered by
`uaa list-clients` breaking on `allowpublic` being serialized as a JSON
string instead of a boolean. Two distinct bug classes were found. Fixing
these is not yet scheduled — capturing findings here for future triage.

## Bug class 1: strict-typed fields on loosely-typed UAA data (fixed)

`Client` is the only UAA resource whose config is stored in a generic
`Map<String,Object> additionalInformation` and flattened into JSON via
`@JsonAnyGetter`/`@JsonAnySetter` with zero type coercion — so a field's JSON
type depends on how/when the client was created, not a fixed schema. Fixed
in go-uaa branch `fix/allowpublic-string-value` (2 commits) by converting
these to `FooRaw interface{}` + `Foo()` accessor pairs:

- `AllowPublic` (the original bug)
- `ApprovalsDeleted`, `LastModified`, `AllowedProviders`, `RequiredUserGroups`

`CreatedWith`/`TokenSalt` go through the same untyped path but are always
used as plain strings server-side, so left as-is.

The new `.github/workflows/integration-test.yml` / `integration/uaa_test.go`
suite will fail against uaa-cli's currently-pinned go-uaa v0.4.2 until this
fix is released and go.mod is bumped — `scripts/boot/uaa.yml` (the default
profile the CI job boots) seeds `login`, `client_federated_jwt_trust`,
`client_with_allowpublic_and_jwks_uri_trust`, and
`oauth_showcase_saml2_bearer` with `allowpublic`/`autoapprove` as literal
YAML strings, which reproduces the exact original bug. That failure is
expected and is the point of the new CI job.

## Bug class 2: fields/resources UAA has that go-uaa doesn't model at all

Unlike bug class 1, these are silent gaps, not crashes — Go's
`json.Unmarshal` ignores unknown fields, so uaa-cli simply can't read, set,
or expose them.

### IdentityProvider — no support at all (biggest gap)

go-uaa has **zero** `IdentityProvider` type. UAA has a full resource here
([`model/src/main/java/org/cloudfoundry/identity/uaa/provider/`](https://github.com/cloudfoundry/uaa/tree/develop/model/src/main/java/org/cloudfoundry/identity/uaa/provider))
with LDAP/SAML/OIDC/Keystone/UAA-internal config classes, including real
CA-cert/TLS trust fields:
- SAML (`SamlIdentityProviderDefinition`) & OIDC/OAuth
  (`AbstractExternalOAuthIdentityProviderDefinition`): `skipSslValidation`
  (bool), `caCertificates` (`List<String>`)
- LDAP (`LdapIdentityProviderDefinition`): `tlsConfiguration`
  (none/simple/ldaps), `caCertificates`

This means no `create-idp`/`get-idp`/`update-idp`/`list-idps` in uaa-cli, and
no way to configure or inspect TLS trust for federated auth at all.

### Client — additional gaps beyond bug class 1

- `refreshTokenUnique` (`ClientConstants.REFRESH_TOKEN_UNIQUE`) — per-client
  override of concurrent-session/refresh-token uniqueness. Security-relevant,
  actively read in `UaaTokenServices`.
- `jwt_creds` (`ClientJwtConfiguration.JWT_CREDS`) — client JWT-bearer
  federation credentials, validated in `ClientAdminEndpointsValidator`.
  Security-relevant.
- `signup_redirect_url` / `change_email_redirect_url` — post-signup/
  post-email-change redirect targets. Open-redirect surface if unvalidated.
- (`client_jwt_config` is intentionally stripped server-side before GET
  responses — not a real gap; go-uaa's `ChangeClientJWT` already covers the
  write path.)

### IdentityZone — several gaps, ranked

- `issuer` — per-zone custom token issuer override. **Security-critical**:
  affects JWT `iss` validation (`KeyInfoService`, `UaaTokenUtils`). Entirely
  missing.
- `defaultIdentityProvider` — default IdP routing per zone
  (`IdentityZoneConfigurationBootstrap`).
- `UserConfig`: go-uaa only has `defaultGroups`; missing `allowedGroups`,
  `maxUsers`, `checkOriginEnabled`, `allowOriginLoop`.
- `SamlConfig.entityID` — custom per-zone SAML SP entity ID.
- `BrandingInformation`: go-uaa has 3 of 8 fields; missing
  `footerLegalText`, `footerLinks`, `banner`, `consent`, `loginConsent`
  (mostly cosmetic; `consent`/`loginConsent` have mild compliance relevance).

### User — real but lower-priority gaps

- `aliasZid` / `aliasId` — cross-zone shadow-user/alias federation
  (`ScimUser.patch`). Worth prioritizing over the rest of this list.
- `salt` — password hashing salt (has a plain getter/setter, not
  `@JsonIgnore`).
- Cosmetic SCIM attributes: `displayName`, `nickName`, `profileUrl`, `title`,
  `userType`, `preferredLanguage`, `locale`, `timezone`.
- (`passwordLastModified`, `previousLogonTime`, `lastLogonTime` already
  present and correct. No MFA/TOTP or `passwordChangeRequired` fields exist
  on `ScimUser` itself.)

### Group — complete, no gaps found.

## Bombshell: MfaProvider is dead code

UAA removed the entire MFA-provider feature server-side in February 2024
(migration `V4_106__Drop_MFA_Tables.sql` drops `mfa_providers` and
`user_google_mfa_credentials`; no `MfaProvider`/`GoogleMfaProviderConfig`
Java class exists anywhere in the current `develop` branch — only a stray
`AuditEventType` enum remnant). go-uaa's `MFAProvider`/`MFAProviderConfig`
types and `/mfa-providers` endpoint model a resource that **no longer exists
on current UAA**. Any code relying on this is pointing at a dead endpoint.
Worth confirming what, if anything, in uaa-cli actually surfaces this before
deciding whether to remove it from go-uaa.
