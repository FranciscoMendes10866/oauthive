package middleware

import (
	"context"
	"errors"
	"net/http"
	"oauthive/api/helpers"
)

type AuthMiddlewareFunc = func(next http.Handler) http.Handler

func BuildAuthMiddleware(cookieManager *helpers.CookieManager) AuthMiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := cookieManager.GetCookie(r, helpers.AuthSessionCookie)
			if err != nil {
				helpers.Reply(w, errors.New("Not Authorized"), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), helpers.CtxSessionID, cookie.SessionID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
