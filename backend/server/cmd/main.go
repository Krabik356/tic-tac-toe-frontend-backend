package main

import (
	"backend/server/internal/database"
	"backend/server/internal/loggers"
	"backend/server/internal/logic"
	"net/http"
)

func main() {
	dbManager := database.NewDBManager()
	defer dbManager.Pool.Close()

	manager := logic.NewManager(dbManager)
	go manager.Manage()
	handler := logic.NewHandler(manager)

	logger := loggers.NewLoggers()
	middleware := logic.NewMiddlewares(logger)
	authMiddleware := middleware.MiddlewareWithAuth(manager)
	http.Handle("/register", middleware.MiddlewareCORS(middleware.MiddlewareWithLoggs(http.HandlerFunc(handler.RegisterHandler))))
	http.Handle("/login", middleware.MiddlewareCORS(middleware.MiddlewareWithLoggs(http.HandlerFunc(handler.LoginHandler))))
	http.Handle("/getLeaderBoard", middleware.MiddlewareCORS(middleware.MiddlewareWithLoggs(authMiddleware(http.HandlerFunc(handler.GetLBHandler)))))
	http.Handle("/game", middleware.MiddlewareCORS(middleware.MiddlewareWithLoggs(authMiddleware(http.HandlerFunc(handler.GameWS)))))
	http.ListenAndServe(":8080", nil)
}
