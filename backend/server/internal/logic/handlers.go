package logic

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	upgrade *websocket.Upgrader
	manager *Manager
}

func NewHandler(m *Manager) *Handler {
	return &Handler{
		upgrade: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return r.Header.Get("Origin") == "http://localhost:7010"
			},
		},
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

	token, tokenTime, err := h.manager.GenerateToken(regData.Name, false)
	if err != nil {
		http.Error(w, "Servers error, problem with generating tocken", 500)
		return
	}
	regData.Token = token
	regData.TokenTime = tokenTime
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
		Token:  token,
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
		token, _, err := h.manager.GenerateToken(regData.Name, true)
		if err != nil {
			http.Error(w, "Servers error, problem with generating tocken", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(SuccessfulLogin{
			Status: "success",
			Name:   regData.Name,
			Rank:   rank,
			Token:  token,
		})
	} else {
		http.Error(w, "Invalid username or password", 401)
		return
	}
}

func (h *Handler) GetLBHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only post method", 405)
		return
	}

	token := r.Header.Get("session_token")
	if token == "" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	isCorrect, err := h.manager.CheckToken(token)
	if err != nil {
		http.Error(w, "Servers error", 500)
		return
	}
	if !isCorrect {
		http.Error(w, "Unauthorized", 401)
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

func (h *Handler) GameWS(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("session_token")
	if token == "" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	isCorrect, err := h.manager.CheckToken(token)
	if err != nil {
		http.Error(w, "Servers error", 500)
		return
	}
	if !isCorrect {
		http.Error(w, "Unauthorized", 401)
		return
	}

	name, err := h.manager.Authorize(token)
	if err != nil {
		http.Error(w, "Servers problem", 500)
		return
	}

	conn, err := h.upgrade.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Servers problem", 500)
		return
	}

	client := NewClient(name, conn)
	h.manager.JoinChan <- client

	go func() {
		defer client.Conn.Close()

		for msg := range client.Send {
			err := client.Conn.WriteJSON(msg)
			if err != nil {
				if ch := client.RoomChan; ch != nil {
					select {
					case ch <- ClientDisconect{
						Status: "disconected",
						Client: client,
					}:
					default:
					}
				}
				return
			}
		}
	}()

	go func() {
		defer client.Conn.Close()
		for {
			var move Movement
			err := client.Conn.ReadJSON(&move)
			if err != nil {
				if ch := client.RoomChan; ch != nil {
					select {
					case ch <- ClientDisconect{
						Status: "disconected",
						Client: client,
					}:
					default:
					}
				}
				return
			}
			if ch := client.RoomChan; ch != nil {
				select {
				case ch <- RoomsMovement{
					Client: client,
					X:      move.X,
					Y:      move.Y,
				}:
				default:
					ch <- ClientDisconect{
						Status: "disconected",
						Client: client,
					}
				}
			} else {
				return
			}
		}
	}()
}
