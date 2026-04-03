package main

import (
	"backend/server/internal/database"
	"backend/server/internal/logic"
	"net/http"
)

func main() {
	dbManager := database.NewDBManager()
	defer dbManager.Pool.Close()

	manager := logic.NewManager(dbManager)
	go manager.Manage()
	handler := logic.NewHandler(manager)

	authMiddleware := logic.MiddlewareWithAuth(manager)
	http.Handle("/register", logic.MiddlewareCORS(http.HandlerFunc(handler.RegisterHandler)))
	http.Handle("/login", logic.MiddlewareCORS(http.HandlerFunc(handler.LoginHandler)))
	http.Handle("/getLeaderBoard", logic.MiddlewareCORS(authMiddleware(http.HandlerFunc(handler.GetLBHandler))))
	http.Handle("/game", logic.MiddlewareCORS(authMiddleware(http.HandlerFunc(handler.GameWS))))
	http.ListenAndServe(":8080", nil)
}
