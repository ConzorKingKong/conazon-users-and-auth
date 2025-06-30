package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOauthConfig = oauth2.Config{
		Scopes:   []string{"email", "profile", "openid"},
		Endpoint: google.Endpoint,
	}
	SecureCookie   bool
	DatabaseURLEnv string
	SECRETKEY      []byte
)