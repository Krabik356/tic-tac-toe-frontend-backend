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
	handler := logic.NewHandler(manager)

	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/getLeaderBoard", handler.GetLBHandler)
	http.ListenAndServe(":8080", nil)
}
