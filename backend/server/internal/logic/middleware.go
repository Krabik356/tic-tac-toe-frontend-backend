package logic

import (
	"backend/server/internal/loggers"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Middlewares struct {
	logger *loggers.Loggers
}

type ResponseWriter struct {
	code int
	http.ResponseWriter
}

func (rW *ResponseWriter) WriteHeader(code int) {
	rW.ResponseWriter.WriteHeader(code)
	rW.code = code
}

func NewMiddlewares(logger *loggers.Loggers) *Middlewares {
	return &Middlewares{
		logger: logger,
	}
}

func (m *Middlewares) MiddlewareWithLoggs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rW := &ResponseWriter{
			code:           200,
			ResponseWriter: w,
		}

		b := make([]byte, 16)
		rand.Read(b)
		reqId := hex.EncodeToString(b)
		ctx := context.WithValue(r.Context(), "request_id", reqId)

		defer func() {
			code := rW.code
			if code >= 200 && code <= 299 {
				m.logger.HttpLogger.Info("log",
					zap.String("request_id", reqId),
					zap.Int("status_code", code),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.Duration("duration", time.Since(start)),
				)
			} else if code >= 400 && code <= 499 {
				m.logger.HttpLogger.Warn("log",
					zap.String("request_id", reqId),
					zap.Int("status_code", code),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.Duration("duration", time.Since(start)),
				)
			} else if code >= 500 && code <= 599 {
				m.logger.HttpLogger.Error("log",
					zap.String("request_id", reqId),
					zap.Int("status_code", code),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.Duration("duration", time.Since(start)),
				)
			}
		}()

		m.logger.HttpLogger.Info("log",
			zap.String("request_id", reqId),
			zap.String("status", "started"),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("start_time", start.Format(time.RFC3339)),
		)

		next.ServeHTTP(rW, r.WithContext(ctx))
	})
}

func (m *Middlewares) MiddlewareCORS(next http.Handler) http.Handler {
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

func (m *Middlewares) MiddlewareWithAuth(manager *Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := r.Header.Get("session_token")
			if token == "" {
				http.Error(w, "Unauthorized", 401)
				return
			}

			log.Println(token)
			isCorrect, err := manager.CheckToken(token)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					http.Error(w, "Unauthorized", 401)
					return
				}
				http.Error(w, "Servers error", 500)
				return
			}
			if !isCorrect {
				http.Error(w, "Unauthorized", 401)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
