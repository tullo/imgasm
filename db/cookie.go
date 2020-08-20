package db

import (
	"net/http"

	"github.com/packago/config"
	"github.com/tullo/cookie"
)

func init() {
	NewCookie()
}

func NewCookie() {
	var opts cookie.CookieOptions
	switch config.File().GetString("environment") {
	case "development":
		opts = cookie.CookieOptions{
			AuthenticationKey: []byte(config.File().GetString("session.development.authenticationKey")),
			EncryptionKey:     []byte(config.File().GetString("session.development.encryptionKey")),
			Domain:            config.File().GetString("session.development.domain"),
			Path:              "/",
			MaxAge:            config.File().GetInt("session.development.maxAge"),
			Secure:            false,
			SameSite:          http.SameSiteStrictMode,
		}
	case "production":
		opts = cookie.CookieOptions{
			AuthenticationKey: []byte(config.File().GetString("session.production.authenticationKey")),
			EncryptionKey:     []byte(config.File().GetString("session.production.encryptionKey")),
			Domain:            config.File().GetString("session.production.domain"),
			Path:              "/",
			MaxAge:            config.File().GetInt("session.production.maxAge"),
			Secure:            true,
			SameSite:          http.SameSiteStrictMode,
		}
	default:
		panic("./config.json environment not set to development or production")
	}
	cookie.NewCookieStore(config.File().GetString("session.key"), opts)
}
