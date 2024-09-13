package api

import (
	"oauthive/api/handler"
	"oauthive/api/helpers"
	"oauthive/api/middleware"
	"oauthive/api/repository"
	"oauthive/domain/authenticator"
	"os"

	"github.com/go-chi/chi/v5"
	mid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-rel/rel"
)

func NewMux(database rel.Repository, authenticator authenticator.Authenticator) *chi.Mux {
	mux := chi.NewMux()

	mux.Use(mid.StripSlashes)
	mux.Use(mid.Logger)
	mux.Use(mid.Recoverer)
	mux.Use(middleware.SetupCors())

	cookieManager := helpers.NewCookieManager(
		[]byte(os.Getenv("COOKIE_HASH_KEY")),
		[]byte(os.Getenv("COOKIE_BLOCK_KEY")),
	)
	authMiddleware := middleware.BuildAuthMiddleware(cookieManager)

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

	mux.Mount("/auth", authHandler.SetupRoutes(authMiddleware))

	return mux
}
