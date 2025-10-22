package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func RecoveryWithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC RECOVERED:\n")
				log.Printf("  Error: %v\n", err)
				log.Printf("  Method: %s\n", r.Method)
				log.Printf("  URL: %s\n", r.URL.String())
				log.Printf("  Remote Address: %s\n", r.RemoteAddr)
				log.Printf("  Stack Trace:\n%s\n", debug.Stack())

				http.Error(
					w,
					fmt.Sprintf("Internal Server Error: %v", err),
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
