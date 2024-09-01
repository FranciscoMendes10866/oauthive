package authenticator

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"golang.org/x/oauth2"
)

type Authenticator interface {
	InitializeLogin(provider string, w http.ResponseWriter, r *http.Request)
	CompleteLogin(provider string, w http.ResponseWriter, r *http.Request) (*goth.User, error)
	RefreshToken(provider, refreshToken string) (*oauth2.Token, error)
}

type authenticator struct {
	discordProvider *discord.Provider
}

func NewAuthenticator() Authenticator {
	discordProvider := discord.New(
		os.Getenv("DISCORD_CLIENT_ID"),
		os.Getenv("DISCORD_CLIENT_SECRET"),
		os.Getenv("DISCORD_CALLBACK_URL"),
		[]string{discord.ScopeIdentify, discord.ScopeEmail}...,
	)

	goth.UseProviders([]goth.Provider{discordProvider}...)

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(60 * 4) // 4 minutes
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = os.Getenv("API_ENV") == "prod"

	gothic.Store = store

	return &authenticator{
		discordProvider: discordProvider,
	}
}

func (self *authenticator) InitializeLogin(provider string, w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func (self *authenticator) CompleteLogin(provider string, w http.ResponseWriter, r *http.Request) (*goth.User, error) {
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (self *authenticator) RefreshToken(provider, refreshToken string) (*oauth2.Token, error) {
	switch provider {
	case "discord":
		{
			newTokens, err := self.discordProvider.RefreshToken(refreshToken)
			if err != nil {
				return nil, err
			}
			return newTokens, nil
		}
	default:
		return nil, errors.New("Oauth provider not supported")
	}
}
