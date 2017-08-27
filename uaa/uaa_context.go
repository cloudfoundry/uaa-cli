package uaa

type Config struct {
	Trace bool
	SkipSSLValidation bool
	Context UaaContext
}

type UaaContext struct {
	BaseUrl string
	AccessToken string
}
