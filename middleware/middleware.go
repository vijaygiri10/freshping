package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Logger(next http.HandlerFunc, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))

			if err := recover(); err != nil {
				fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start)), err)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Logged connection from %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}
