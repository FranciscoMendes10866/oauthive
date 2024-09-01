package handler

import (
	"net/http"
	"oauthive/api/repository"
	"oauthive/db/entities"
	"oauthive/domain/authenticator"

	"github.com/araddon/dateparse"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authenticator authenticator.Authenticator
	userRepo      repository.UserRepository
	accountRepo   repository.AccountRepository
}

func NewAuthHandler(
	authenticator authenticator.Authenticator,
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository,
) *AuthHandler {
	return &AuthHandler{
		authenticator,
		userRepo,
		accountRepo,
	}
}

func (self *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		http.Error(w, "Provider URL param is required", http.StatusBadRequest)
		return
	}
	self.authenticator.InitializeLogin(provider, w, r)
}

func (self *AuthHandler) LoginCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		http.Error(w, "Provider URL param is required", http.StatusBadRequest)
		return
	}

	user, err := self.authenticator.CompleteLogin(provider, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = self.userRepo.UpsertUser(r.Context(), &entities.User{
		Name:  user.Name,
		Email: user.Email,
		Image: user.AvatarURL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userRecord, err := self.userRepo.FindUserByEmail(r.Context(), user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parsedTime, err := dateparse.ParseAny(user.ExpiresAt.String())
	if err != nil {
		http.Error(w, "Error parsing date", http.StatusInternalServerError)
		return
	}

	err = self.accountRepo.UpsertAccount(r.Context(), &entities.Account{
		Provider:          user.Provider,
		ProviderAccountID: user.UserID,
		RefreshToken:      user.RefreshToken,
		AccessToken:       user.AccessToken,
		IDToken:           user.IDToken,
		TokenType:         "Bearer",
		UserID:            userRecord.ID,
		ExpiresAt:         parsedTime.Unix(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: create session -> sessionRepo.CreateSession
	// TODO: create secure cookie -> github.com/gorilla/securecookie
	// TODO: redirect user -> http.Redirect(w, r, "FRONTEND_URL", http.StatusFound)
}
