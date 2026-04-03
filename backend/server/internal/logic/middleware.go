package logic

import (
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func MiddlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:7010")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, session_token")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func MiddlewareWithAuth(manager *Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := r.Header.Get("session_token")
			if token == "" {
				http.Error(w, "Unauthorized", 401)
				log.Println("Unauthorized1")

				return
			}

			log.Println(token)
			isCorrect, err := manager.CheckToken(token)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					http.Error(w, "Unauthorized", 401)
					log.Println("Unauthorized2")

					return
				}
				http.Error(w, "Servers error", 500)
				log.Println("Servers error1xxx")
				log.Println(err)

				return
			}
			if !isCorrect {
				http.Error(w, "Unauthorized", 401)
				log.Println("Unauthorized3")

				return
			}
			log.Println("ok")
			next.ServeHTTP(w, r)
		})
	}
}
