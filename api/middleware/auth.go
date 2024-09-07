package middleware

import (
	"context"
	"net/http"
	"oauthive/api/helpers"
)

func BuildAuthGuard(cookieManager *helpers.CookieManager) func(handler http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := cookieManager.GetCookie(r, helpers.AuthSessionCookie)
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
			r = r.WithContext(context.WithValue(r.Context(), helpers.CtxSessionID, cookie.SessionID))
			handler(w, r)
		})
	}
}
