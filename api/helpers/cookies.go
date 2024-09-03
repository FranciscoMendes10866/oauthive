package helpers

import (
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
)

type CookieContent struct {
	SessionID int
}

type CookieManager struct {
	secureCookie *securecookie.SecureCookie
}

func NewCookieManager(hashKey, blockKey []byte) *CookieManager {
	return &CookieManager{
		secureCookie: securecookie.New(hashKey, blockKey),
	}
}

func (self *CookieManager) SetCookie(w http.ResponseWriter, name string, value *CookieContent, maxAge int) {
	encoded, err := self.secureCookie.Encode(name, *value)
	if err != nil {
		log.Println("Error encoding cookie:", err)
		return
	}
	cookie := &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   IsProd,
		MaxAge:   maxAge,
	}
	http.SetCookie(w, cookie)
}

func (self *CookieManager) GetCookie(r *http.Request, name string) (*CookieContent, error) {
	if cookie, err := r.Cookie(name); err == nil {
		value := &CookieContent{}
		if err = self.secureCookie.Decode(name, cookie.Value, &value); err == nil {
			return value, nil
		}
	}
	return nil, http.ErrNoCookie
}

func (self *CookieManager) ClearCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   IsProd,
	}
	http.SetCookie(w, cookie)
}
