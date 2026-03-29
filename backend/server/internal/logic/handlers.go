package logic

import (
	"encoding/json"
	"net/http"
	"time"
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

	token, err := h.manager.GenerateToken(regData.Name)
	if err != nil {
		http.Error(w, "Servers error, problem with generating tocken", 500)
		return
	}
	regData.Token = token
	http.SetCookie(w, &http.Cookie{
		Name:     "session_tocken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(3 * 24 * time.Hour),
	})
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
		token, err := h.manager.GenerateToken(regData.Name)
		if err != nil {
			http.Error(w, "Servers error, problem with generating tocken", 500)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_tocken",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			Expires:  time.Now().Add(3 * 24 * time.Hour),
		})
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
	if r.Method != http.MethodGet {
		http.Error(w, "Only get method", 405)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unsended cookies", 400)
		return
	}

	token := cookie.Value
	isCorrect, err := h.manager.CheckToken(token)
	if err != nil {
		http.Error(w, "Servers error", 500)
		return
	}
	if !isCorrect {
		http.Error(w, "There is not user with this token", 401)
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
