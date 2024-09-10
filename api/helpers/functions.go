package helpers

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func GetSessionID(ctx context.Context) int {
	return ctx.Value(CtxSessionID).(int)
}

type CookieStatus = string

const (
	CookieValid CookieStatus = "valid"
	CookieRenew CookieStatus = "renew"
)

func CheckAuthSession(r *http.Request, issuedAt, expiresAt int64) (CookieStatus, error) {
	elapsedTime := time.Now().Unix() - issuedAt
	remainingAge := expiresAt - elapsedTime

	if remainingAge > 0 && remainingAge <= AuthRenewThreshold {
		return CookieRenew, nil
	} else if remainingAge > 0 {
		return CookieValid, nil
	}

	return "", errors.New("Cookie expired")
}
