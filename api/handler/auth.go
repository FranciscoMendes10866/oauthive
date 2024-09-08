package handler

import (
	"net/http"
	"oauthive/api/helpers"
	"oauthive/api/repository"
	"oauthive/db/entities"
	"oauthive/domain/authenticator"
	"time"

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
		ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logged out successfully"}`))
}

func (self *AuthHandler) RenewSession(w http.ResponseWriter, r *http.Request) {
	sessionStatus, err := helpers.CheckAuthSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	switch sessionStatus {
	case helpers.CookieValid:
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Current session is still valid"}`))
		return

	case helpers.CookieRenew:
		authCookie, err := self.cookieManager.GetCookie(r, helpers.AuthSessionCookie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
			ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
			UserID:    session.UserID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		self.cookieManager.SetCookie(w, helpers.AuthSessionCookie, &helpers.CookieContent{
			SessionID: newSession.ID,
		}, helpers.AuthSessionMaxAge)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Session renewed successfully"}`))
		return
	}
}
