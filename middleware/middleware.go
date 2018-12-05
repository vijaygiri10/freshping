package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func Logger(next http.HandlerFunc, Name string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			fmt.Println(r.Method, "  ", Name, "  ", r.RequestURI, "  ", time.Since(start))
			if err := recover(); err != nil {
				fmt.Println("execption ocurried and recovred ")
			}
		}()
	})
}
