package middleware 

import (
	"fmt"
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Print log
	method := r.Method
	path := r.URL.Path
	protocol := r.Proto
	log.Println(fmt.Sprintf("%s %s %s", method, path, protocol))
	next.ServeHTTP(w, r)
	})
}