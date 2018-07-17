package uaa

// MFAProvidersEndpoint is the path to the MFA providers resource.
const MFAProvidersEndpoint string = "/mfa-providers"

// MFAProviderConfig is configuration for an MFA provider
type MFAProviderConfig struct {
	Issuer              string `json:"issuer"`
	ProviderDescription string `json:"providerDescription"`
}

// MFAProvider is a UAA MFA provider
// http://docs.cloudfoundry.org/api/uaa/version/4.19.0/index.html#get-2
type MFAProvider struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	IdentityZoneID string            `json:"identityZoneId"`
	Config         MFAProviderConfig `json:"config"`
	Type           string            `json:"type"`
	Created        int               `json:"created"`
	LastModified   int               `json:"last_modified"`
}

// Identifier returns the field used to uniquely identify a MFAProvider.
func (m MFAProvider) Identifier() string {
	return m.ID
}
