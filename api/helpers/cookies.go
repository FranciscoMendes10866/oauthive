package helpers

import (
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
)

type CookieContent struct {
	SessionID        int
	IssuedAt         int64
	RenewalTimeframe int64
}

type CookieManager struct {
	secureCookie *securecookie.SecureCookie
	IsProd       bool
}

func NewCookieManager(hashKey, blockKey []byte) *CookieManager {
	return &CookieManager{
		secureCookie: securecookie.New(hashKey, blockKey),
		IsProd:       os.Getenv("API_ENV") == "prod",
	}
}

func (self *CookieManager) SetCookie(w http.ResponseWriter, name string, value CookieContent, maxAge int) error {
	encoded, err := self.secureCookie.Encode(name, value)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   self.IsProd,
		MaxAge:   maxAge,
	})
	return nil
}

func (self *CookieManager) GetCookie(r *http.Request, name string) (CookieContent, error) {
	value := CookieContent{}
	if cookie, err := r.Cookie(name); err == nil {
		if err = self.secureCookie.Decode(name, cookie.Value, &value); err == nil {
			return value, nil
		}
	}
	return value, http.ErrNoCookie
}

func (self *CookieManager) ClearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   self.IsProd,
	})
}
