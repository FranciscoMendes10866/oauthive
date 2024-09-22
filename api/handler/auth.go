package handler

import (
	"net/http"
	"oauthive/api/helpers"
	"oauthive/api/middleware"
	"oauthive/api/repository"
	"oauthive/db/entities"
	"oauthive/domain/authenticator"
	"os"

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

func (self *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		helpers.Reply(w, "Provider URL param is required", http.StatusBadRequest)
		return
	}
	self.authenticator.InitializeLogin(provider, w, r)
}

func (self *AuthHandler) loginCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		helpers.Reply(w, "Provider URL param is required", http.StatusBadRequest)
		return
	}

	user, err := self.authenticator.CompleteLogin(provider, w, r)
	if err != nil {
		helpers.Reply(w, err, http.StatusInternalServerError)
		return
	}

	if err := self.userRepo.UpsertUser(r.Context(), entities.User{
		Name:  user.Name,
		Email: user.Email,
	}); err != nil {
		helpers.Reply(w, err, http.StatusInternalServerError)
		return
	}

	userRecord, err := self.userRepo.FindUserByEmail(r.Context(), user.Email)
	if err != nil {
		helpers.Reply(w, err, http.StatusInternalServerError)
		return
	}

	if err := self.accountRepo.UpsertAccount(r.Context(), &entities.Account{
		Provider:          user.Provider,
		ProviderAccountID: user.UserID,
		UserID:            userRecord.ID,
	}); err != nil {
		helpers.Reply(w, err, http.StatusInternalServerError)
		return
	}

	timeFactory := &helpers.TimeframeFactory{}

	newSession, err := self.sessionRepo.CreateSession(r.Context(), &entities.Session{
		ExpiresAt: timeFactory.GenerateExpiresAt(),
		UserID:    userRecord.ID,
	})
	if err != nil {
		helpers.Reply(w, err, http.StatusInternalServerError)
		return
	}

	if err := self.cookieManager.SetCookie(w, helpers.AuthSessionCookie, helpers.CookieContent{
		SessionID:        newSession.ID,
		IssuedAt:         timeFactory.GenerateIssuedAt(),
		RenewalTimeframe: timeFactory.GenerateRenewalTimeframe(),
	}, helpers.AuthSessionMaxAge); err != nil {
		helpers.Reply(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, os.Getenv("FRONTEND_URL"), http.StatusFound)
}

func (self *AuthHandler) renewSession(w http.ResponseWriter, r *http.Request) {
	authCookie, err := self.cookieManager.GetCookie(r, helpers.AuthSessionCookie)
	if err != nil {
		helpers.Reply(w, err, http.StatusUnauthorized)
		return
	}

	session, err := self.sessionRepo.FindSessionByID(r.Context(), authCookie.SessionID)
	if err != nil {
		helpers.Reply(w, err, http.StatusForbidden)
		return
	}

	timeFactory := &helpers.TimeframeFactory{}

	sessionStatus := timeFactory.Verify(
		authCookie.IssuedAt,
		session.ExpiresAt,
		authCookie.RenewalTimeframe,
	)

	switch sessionStatus {
	case helpers.CookieValid:
	default:
		helpers.Reply(w, "Current session is still valid", http.StatusOK)
		return

	case helpers.CookieRenew:
		if err := self.sessionRepo.DeleteSessionByID(r.Context(), session.ID); err != nil {
			helpers.Reply(w, err, http.StatusInternalServerError)
			return
		}

		newSession, err := self.sessionRepo.CreateSession(r.Context(), &entities.Session{
			ExpiresAt: timeFactory.GenerateExpiresAt(),
			UserID:    session.UserID,
		})
		if err != nil {
			helpers.Reply(w, err, http.StatusInternalServerError)
			return
		}

		if err := self.cookieManager.SetCookie(w, helpers.AuthSessionCookie, helpers.CookieContent{
			SessionID:        newSession.ID,
			IssuedAt:         timeFactory.GenerateIssuedAt(),
			RenewalTimeframe: timeFactory.GenerateRenewalTimeframe(),
		}, helpers.AuthSessionMaxAge); err != nil {
			helpers.Reply(w, err, http.StatusInternalServerError)
			return
		}

		helpers.Reply(w, "Session renewed successfully", http.StatusCreated)
		return
	}
}

func (self *AuthHandler) getUser(w http.ResponseWriter, r *http.Request) {
	sessionID, err := helpers.GetSessionID(r.Context())
	if err != nil {
		helpers.Reply(w, err, http.StatusForbidden)
		return
	}

	session, err := self.sessionRepo.FindSessionByID(r.Context(), sessionID)
	if err != nil {
		helpers.Reply(w, err, http.StatusForbidden)
		return
	}

	user, err := self.userRepo.FindUserByID(r.Context(), session.UserID)
	if err != nil {
		helpers.Reply(w, err, http.StatusInternalServerError)
		return
	}

	helpers.Reply(w, user, http.StatusOK)
}

func (self *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := helpers.GetSessionID(r.Context())
	if err != nil {
		helpers.Reply(w, err, http.StatusForbidden)
		return
	}

	if err := self.sessionRepo.DeleteSessionByID(r.Context(), sessionID); err != nil {
		helpers.Reply(w, err, http.StatusInternalServerError)
		return
	}

	self.cookieManager.ClearCookie(w, helpers.AuthSessionCookie)

	helpers.Reply(w, "Logged out successfully", http.StatusOK)
}

func (self *AuthHandler) SetupRoutes(authMiddleware middleware.AuthMiddlewareFunc) *chi.Mux {
	mux := chi.NewMux()

	mux.Get("/login/{provider}", self.login)
	mux.Get("/{provider}/callback", self.loginCallback)
	mux.Get("/refresh", self.renewSession)
	mux.With(authMiddleware).Get("/user", self.getUser)
	mux.With(authMiddleware).Get("/logout", self.logout)

	return mux
}
