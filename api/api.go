package api

import (
	"oauthive/api/handler"
	"oauthive/api/helpers"
	"oauthive/api/repository"
	"oauthive/domain/authenticator"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-rel/rel"
)

func NewMux(database rel.Repository) *chi.Mux {
	mux := chi.NewMux()

	mux.Use(middleware.StripSlashes)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	authenticator := authenticator.NewAuthenticator()
	cookieManager := helpers.NewCookieManager(
		[]byte("KPMCSgGLbFsW"),
		[]byte("uvgJkz7wXraU"),
	)

	userRepo := repository.NewUserRepository(database)
	accountRepo := repository.NewAccountRepository(database)
	sessionRepo := repository.NewSessionRepository(database)

	authHandler := handler.NewAuthHandler(
		authenticator,
		userRepo,
		accountRepo,
		sessionRepo,
		cookieManager,
	)

	mux.Get("/login/{provider}", authHandler.Login)
	mux.Get("/{provider}/callback", authHandler.LoginCallback)
	mux.Get("/refresh", authHandler.RenewSession)
	mux.Get("/logout", authHandler.Logout)

	return mux
}
