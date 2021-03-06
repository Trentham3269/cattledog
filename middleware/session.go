package middleware

import (
	"log"
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
		session, err := Store.Get(r, "cookie-name")
		if err != nil {
			log.Println(err)
		}

		val := session.Values["authenticated"]

		// Check if user is authenticated
		if auth, ok := val.(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
