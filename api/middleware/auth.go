package middleware

import (
	"context"
	"net/http"
	"oauthive/api/helpers"
)

type InterceptorFunc = func(handler http.HandlerFunc) http.HandlerFunc

func BuildAuthGuard(cookieManager *helpers.CookieManager) InterceptorFunc {
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
