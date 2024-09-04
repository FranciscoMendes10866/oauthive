package main

import (
	"oauthive/api/handler"
	"oauthive/api/helpers"
	m "oauthive/api/middleware"
	"oauthive/api/repository"
	"oauthive/db"
	"oauthive/domain/authenticator"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	database := db.Init("database.db")

	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authenticator := authenticator.NewAuthenticator()
	cookieManager := helpers.NewCookieManager([]byte("KPMCSgGLbFsW"), []byte("uvgJkz7wXraU"))
	authGuard := m.BuildAuthGuard(cookieManager)

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
}
