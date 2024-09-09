package api

import (
	"oauthive/api/handler"
	"oauthive/api/helpers"
	"oauthive/api/middleware"
	"oauthive/api/repository"
	"oauthive/domain/authenticator"

	"github.com/go-chi/chi/v5"
	mid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-rel/rel"
)

func NewMux(database rel.Repository, authenticator authenticator.Authenticator) *chi.Mux {
	mux := chi.NewMux()

	mux.Use(mid.StripSlashes)
	mux.Use(mid.Logger)
	mux.Use(mid.Recoverer)
	mux.Use(middleware.SetupSecureMiddleware())

	cookieManager := helpers.NewCookieManager(
		[]byte("KPMCSgGLbFsW"),
		[]byte("uvgJkz7wXraU"),
	)
	authGuard := middleware.BuildAuthGuard(cookieManager)

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

	mux.Mount("/auth", authHandler.SetupRoutes(authGuard))

	return mux
}
