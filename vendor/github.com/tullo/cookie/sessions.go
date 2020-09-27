package cookie

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	cookies = make(map[string]*sessions.CookieStore)

	ErrCookieNotFound = errors.New("Cookie not found")
)

type CookieOptions struct {
	AuthenticationKey []byte
	EncryptionKey     []byte
	Domain            string
	Path              string
	MaxAge            int
	Secure            bool
	SameSite          http.SameSite
}

func NewCookieStore(name string, c CookieOptions) {
	if cookies[name] == nil {
		cookies[name] = sessions.NewCookieStore(c.AuthenticationKey, c.EncryptionKey)
		cookies[name].Options = &sessions.Options{
			Domain:   c.Domain,
			Path:     c.Path,
			MaxAge:   c.MaxAge,
			Secure:   c.Secure,
			SameSite: c.SameSite,
		}
	}
}

func GetSession(r *http.Request, name string) (*sessions.Session, error) {
	if cookies[name] == nil {
		return nil, ErrCookieNotFound
	}
	return cookies[name].Get(r, name)
}

func ClearSession(sess *sessions.Session) {
	for key := range sess.Values {
		delete(sess.Values, key)
	}
}
