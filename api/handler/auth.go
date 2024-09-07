package handler

import (
	"net/http"
	"oauthive/api/helpers"
	"oauthive/api/repository"
	"oauthive/db/entities"
	"oauthive/domain/authenticator"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authenticator authenticator.Authenticator
	userRepo      *repository.UserRepository
	accountRepo   *repository.AccountRepository
	sessionRepo   *repository.SessionRepository
	cookieManager *helpers.CookieManager
}

func NewAuthHandler(
	authenticator authenticator.Authenticator,
	userRepo *repository.UserRepository,
	accountRepo *repository.AccountRepository,
	sessionRepo *repository.SessionRepository,
	cookieManager *helpers.CookieManager,
) *AuthHandler {
	return &AuthHandler{
		authenticator,
		userRepo,
		accountRepo,
		sessionRepo,
		cookieManager,
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

	newSession, err := self.sessionRepo.CreateSession(r.Context(), &entities.Session{
		ExpiresAt: helpers.AuthSessionMaxAge,
		UserID:    userRecord.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	self.cookieManager.SetCookie(w, helpers.AuthSessionCookie, &helpers.CookieContent{
		SessionID: newSession.ID,
	}, helpers.AuthSessionMaxAge)

	http.Redirect(w, r, helpers.FrontendURL, http.StatusFound)
}

func (self *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := helpers.GetSessionID(r.Context())

	err := self.sessionRepo.DeleteSessionByID(r.Context(), sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	self.cookieManager.ClearCookie(w, helpers.AuthSessionCookie)

	w.Write([]byte("Logged out successfully"))
}

func (self *AuthHandler) RenewSession(w http.ResponseWriter, r *http.Request) {
	sessionStatus, err := helpers.CheckAuthSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	authCookie, err := self.cookieManager.GetCookie(r, helpers.AuthSessionCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch sessionStatus {
	case helpers.CookieExpired:
		err := self.sessionRepo.DeleteSessionByID(r.Context(), authCookie.SessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		self.cookieManager.ClearCookie(w, helpers.AuthSessionCookie)
		http.Error(w, "Expired session", http.StatusForbidden)
		return

	case helpers.CookieRenew:
		session, err := self.sessionRepo.FindSessionByID(r.Context(), authCookie.SessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = self.sessionRepo.DeleteSessionByID(r.Context(), session.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newSession, err := self.sessionRepo.CreateSession(r.Context(), &entities.Session{
			ExpiresAt: helpers.AuthSessionMaxAge,
			UserID:    session.UserID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		self.cookieManager.SetCookie(w, helpers.AuthSessionCookie, &helpers.CookieContent{
			SessionID: newSession.ID,
		}, helpers.AuthSessionMaxAge)
		return

	case helpers.CookieValid:
		// TODO -> should I do anything? maybe return a 200 (ok)
		return
	}
}
