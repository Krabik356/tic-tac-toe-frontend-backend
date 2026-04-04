package logic

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Manager struct {
	rooms      map[int]*Room
	idCounter  int
	JoinChan   chan *Client
	DeleteChan chan int
	DataBase   interface {
		Register(RegisterData) error
		Login(RegisterData) (bool, int, error)
		GetLeaderBoard() (LeaderBoard, error)
		SetToken(string, string, string) error
		CheckToken(string) (bool, string, error)
		GetName(string) (string, error)
		Retribution(string, string, int, int) error
		CheckName(string) error
	}
}

func NewManager(dbM interface {
	Register(RegisterData) error
	Login(RegisterData) (bool, int, error)
	GetLeaderBoard() (LeaderBoard, error)
	SetToken(string, string, string) error
	CheckToken(string) (bool, string, error)
	GetName(string) (string, error)
	Retribution(string, string, int, int) error
	CheckName(string) error
}) *Manager {
	return &Manager{
		rooms:      make(map[int]*Room),
		idCounter:  0,
		JoinChan:   make(chan *Client, 100),
		DataBase:   dbM,
		DeleteChan: make(chan int, 100),
	}
}

func (m *Manager) RegisterUser(data RegisterData) error {
	err := m.DataBase.Register(data)
	return err
}

func (m *Manager) IsThereThisName(name string) (bool, error) {
	err := m.DataBase.CheckName(name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *Manager) LoginUser(data RegisterData) (bool, int, error) {
	isLogined, rank, err := m.DataBase.Login(data)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, 0, nil
	}
	return isLogined, rank, err
}

func (m *Manager) GetLeaderBoard() (LeaderBoard, error) {
	leaderBoard, err := m.DataBase.GetLeaderBoard()
	return leaderBoard, err
}

func (m *Manager) GenerateToken(name string, withSave bool) (string, string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", "", err
	}
	token := hex.EncodeToString(b)
	tokenTime := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	if withSave {
		err = m.DataBase.SetToken(name, token, tokenTime)
		return token, tokenTime, err
	}
	return token, tokenTime, nil
}

func (m *Manager) CheckToken(token string) (bool, error) {
	isCorrect, tokenTime, err := m.DataBase.CheckToken(token)
	if err != nil {
		return false, err
	}
	if isCorrect {
		timeNow := time.Now()
		tokenTimeTime, err := time.Parse(time.RFC3339, tokenTime)
		if err != nil {
			return false, err
		}
		if timeNow.After(tokenTimeTime) {
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

func (m *Manager) Authorize(token string) (string, error) {
	name, err := m.DataBase.GetName(token)
	return name, err
}

func (m *Manager) roomFindOrAdd(c *Client) {
	for _, room := range m.rooms {
		if room.Status == "waiting" {
			room.Clients[1] = c
			room.Status = "game"
			c.Marker = "X"
			go room.Run()
			return
		}
	}
	id := m.idCounter
	m.rooms[id] = NewRoom(id, m.DataBase, m.DeleteChan)
	m.rooms[id].Clients[0] = c
	m.idCounter++
	c.Marker = "0"
}

func (m *Manager) Manage() {
	for {
		select {
		case c := <-m.JoinChan:
			m.roomFindOrAdd(c)
		case id := <-m.DeleteChan:
			room, ok := m.rooms[id]
			if ok {
				room.CloseRoom()
				delete(m.rooms, room.Id)
			}
		}
	}
}
