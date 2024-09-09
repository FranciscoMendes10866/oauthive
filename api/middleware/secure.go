package middleware

import (
	"net/http"
	"oauthive/api/helpers"

	"github.com/unrolled/secure"
)

func SetupSecureMiddleware() func(next http.Handler) http.Handler {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		SSLRedirect:           helpers.IsProd,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		BrowserXssFilter:      true,
		AllowedHosts:          []string{helpers.FrontendURL, "https://discord.com"},
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; connect-src 'self' discord.com; object-src 'none'; frame-ancestors 'none';",
	})
	return secureMiddleware.Handler
}
