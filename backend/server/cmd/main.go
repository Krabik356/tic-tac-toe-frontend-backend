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

	http.Handle("/register", logic.Middleware(http.HandlerFunc(handler.RegisterHandler)))
	http.Handle("/login", logic.Middleware(http.HandlerFunc(handler.LoginHandler)))
	http.Handle("/getLeaderBoard", logic.Middleware(http.HandlerFunc(handler.GetLBHandler)))
	http.HandleFunc("/game", handler.GameWS)
	http.ListenAndServe(":8080", nil)
}
