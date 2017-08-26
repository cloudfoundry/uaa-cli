package uaa

type Config struct {
	Trace bool
	Context UaaContext
}

type UaaContext struct {
	BaseUrl string
	AccessToken string
}
