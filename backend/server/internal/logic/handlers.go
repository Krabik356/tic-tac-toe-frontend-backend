package logic

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	manager *Manager
}

func NewHandler(m *Manager) *Handler {
	return &Handler{
		manager: m,
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only post method", 405)
		return
	}

	var regData RegisterData
	err := json.NewDecoder(r.Body).Decode(&regData)
	if err != nil {
		http.Error(w, "Invalid data", 400)
		return
	}

	err = h.manager.RegisterUser(regData)
	if err != nil {
		http.Error(w, "Servers error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(SuccessfulRegistartion{
		Status: "success",
		Name:   regData.Name,
	})
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method", 405)
		return
	}

	var regData RegisterData
	err := json.NewDecoder(r.Body).Decode(&regData)
	if err != nil {
		http.Error(w, "Invalid data", 400)
		return
	}

	isLogined, rank, err := h.manager.LoginUser(regData)
	if err != nil {
		http.Error(w, "Servers error", 500)
		return
	}

	if isLogined {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(SuccessfulLogin{
			Status: "success",
			Name:   regData.Name,
			Rank:   rank,
		})
	} else {
		http.Error(w, "Invalid username or password", 401)
		return
	}
}

func (h *Handler) GetLBHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method", 405)
		return
	}

	leaderBoard, err := h.manager.GetLeaderBoard()
	if err != nil {
		http.Error(w, "Servers error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(SuccessfulLeaderBoard{
		Status: "success",
		LB:     leaderBoard,
	})

}
