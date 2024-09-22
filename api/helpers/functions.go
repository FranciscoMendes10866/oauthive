package helpers

import (
	"context"
	"errors"
	"time"
)

func GetSessionID(ctx context.Context) (int, error) {
	sessionID, ok := ctx.Value(CtxSessionID).(int)
	if !ok {
		return 0, errors.New("session ID not found in context")
	}
	return sessionID, nil
}

type CookieStatus string

const (
	CookieValid CookieStatus = "valid"
	CookieRenew CookieStatus = "renew"
)

type TimeframeFactory struct{}

func (tf *TimeframeFactory) GenerateIssuedAt() int64 {
	return time.Now().Unix()
}

func (tf *TimeframeFactory) GenerateExpiresAt() int64 {
	return time.Now().Add(7 * 24 * time.Hour).Unix()
}

func (tf *TimeframeFactory) GenerateRenewalTimeframe() int64 {
	return time.Now().Add(4 * 24 * time.Hour).Unix()
}

func (tf *TimeframeFactory) Verify(issuedAt, expiresAt, renewalTimeframe int64) CookieStatus {
	remainingTime := expiresAt - time.Now().Unix()
	if remainingTime > renewalTimeframe && remainingTime <= expiresAt {
		return CookieRenew
	}
	return CookieValid
}
