package main

import (
	"log"
	"net/http"
)

func (app *application) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.authenticateToken(r)
		if err != nil {
			log.Println("err: ", err.Error())
			app.invalidCredentials(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
