package config

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"path"
	"time"
)

const (
	REFRESH_TOKEN      = GrantType("refresh_token")
	AUTHCODE           = GrantType("authorization_code")
	IMPLICIT           = GrantType("implicit")
	PASSWORD           = GrantType("password")
	CLIENT_CREDENTIALS = GrantType("client_credentials")
)

type Config struct {
	Verbose          bool
	ZoneSubdomain    string
	Targets          map[string]Target
	ActiveTargetName string
}

type Target struct {
	BaseUrl           string
	SkipSSLValidation bool
	Contexts          map[string]UaaContext
	ActiveContextName string
}

type GrantType string

type UaaContext struct {
	ClientId  string    `json:"client_id"`
	GrantType GrantType `json:"grant_type"`
	Username  string    `json:"username"`
	Token     oauth2.Token
}

func NewConfig() Config {
	c := Config{}
	c.Targets = map[string]Target{}
	return c
}

func NewConfigWithServerURL(url string) Config {
	c := NewConfig()
	t := NewTarget()
	t.BaseUrl = url
	c.AddTarget(t)
	return c
}

func NewContextWithToken(accessToken string) UaaContext {
	ctx := UaaContext{
		Token: oauth2.Token{
			AccessToken: accessToken,
			Expiry:      time.Now().Add(1 * time.Hour),
		},
	}
	return ctx
}

func NewTarget() Target {
	t := Target{}
	t.Contexts = map[string]UaaContext{}
	return t
}

func (c Config) GetActiveTarget() Target {
	return c.Targets[c.ActiveTargetName]
}

func (c *Config) AddTarget(newTarget Target) {
	c.Targets[newTarget.name()] = newTarget
	c.ActiveTargetName = newTarget.name()
}

func (c *Config) AddContext(newContext UaaContext) {
	if c.Targets == nil {
		c.Targets = map[string]Target{}
	}
	t := c.Targets[c.ActiveTargetName]
	if t.Contexts == nil {
		t.Contexts = map[string]UaaContext{}
	}
	t.Contexts[newContext.name()] = newContext
	t.ActiveContextName = newContext.name()
	c.AddTarget(t)
}

func (c Config) GetActiveContext() UaaContext {
	return c.GetActiveTarget().GetActiveContext()
}

func (t Target) GetActiveContext() UaaContext {
	return t.Contexts[t.ActiveContextName]
}

func (t Target) name() string {
	return "url:" + t.BaseUrl
}

func (uc UaaContext) name() string {
	return fmt.Sprintf("client:%v user:%v grant_type:%v", uc.ClientId, uc.Username, uc.GrantType)
}

func ConfigDir() string {
	return path.Join(userHomeDir(), ".uaa")
}

func ConfigPath() string {
	return path.Join(ConfigDir(), "config.json")
}

func ReadConfig() Config {
	c := NewConfig()

	data, err := ioutil.ReadFile(ConfigPath())
	if err != nil {
		return c
	}

	json.Unmarshal(data, &c)

	return c
}

func WriteConfig(c Config) error {
	err := makeDirectory()
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	configPath := ConfigPath()
	return ioutil.WriteFile(configPath, data, 0600)
}

func RemoveConfig() error {
	return os.Remove(ConfigPath())
}
