package helpers

import (
	"context"
	"errors"
	"net/http"
)

func GetSessionID(ctx context.Context) int {
	return ctx.Value(CtxSessionID).(int)
}

type CookieStatus = string

const (
	CookieValid   CookieStatus = "valid"
	CookieRenew   CookieStatus = "renew"
	CookieExpired CookieStatus = "expired"
)

func CheckAuthSession(r *http.Request) (CookieStatus, error) {
	cookie, err := r.Cookie(AuthSessionCookie)
	if err != nil {
		return "", errors.New("Cookie doesn't exist")
	}

	remainingAge := cookie.MaxAge
	if remainingAge < 0 {
		return CookieExpired, nil
	}

	if remainingAge > 0 && remainingAge <= AuthRenewThreshold {
		return CookieRenew, nil
	}

	return CookieValid, nil
}
