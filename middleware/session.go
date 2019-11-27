package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// TODO: use Go's crypto/rand or securecookie.GenerateRandomKey(32) and persist the result
var (
	key = []byte("super-secret-key")
	Store = sessions.NewCookieStore(key)
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	session, _ := Store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	next.ServeHTTP(w, r)
	})
}
