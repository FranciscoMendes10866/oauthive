package handler

import (
	"fmt"
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
	sessionRepo   repository.SessionRepository
}

func NewAuthHandler(
	authenticator authenticator.Authenticator,
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository,
	sessionRepo repository.SessionRepository,
) *AuthHandler {
	return &AuthHandler{
		authenticator,
		userRepo,
		accountRepo,
		sessionRepo,
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

	err = self.accountRepo.UpsertAccount(r.Context(), &entities.Account{
		Provider:          user.Provider,
		ProviderAccountID: user.UserID,
		UserID:            userRecord.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parsedTime, err := dateparse.ParseAny(user.ExpiresAt.String())
	if err != nil {
		http.Error(w, "Error parsing date", http.StatusInternalServerError)
		return
	}

	newSession, err := self.sessionRepo.CreateSession(r.Context(), &entities.Session{
		ExpiresAt: parsedTime.Unix(),
		UserID:    userRecord.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(newSession.ID)

	// TODO: create secure cookie -> github.com/gorilla/securecookie
	// TODO: redirect user -> http.Redirect(w, r, "FRONTEND_URL", http.StatusFound)
}
