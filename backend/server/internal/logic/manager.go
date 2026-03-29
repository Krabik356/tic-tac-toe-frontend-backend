package logic

import (
	"crypto/rand"
	"encoding/hex"
)

type Manager struct {
	DataBase interface {
		Register(RegisterData) error
		Login(RegisterData) (bool, int, error)
		GetLeaderBoard() (LeaderBoard, error)
		SetInfo(string, []string, []string) error
		CheckToken(string) (bool, error)
	}
}

func NewManager(dbM interface {
	Register(RegisterData) error
	Login(RegisterData) (bool, int, error)
	GetLeaderBoard() (LeaderBoard, error)
	SetInfo(string, []string, []string) error
	CheckToken(string) (bool, error)
}) *Manager {
	return &Manager{
		DataBase: dbM,
	}
}

func (m *Manager) RegisterUser(data RegisterData) error {
	err := m.DataBase.Register(data)
	return err
}

func (m *Manager) LoginUser(data RegisterData) (bool, int, error) {
	isLogined, rank, err := m.DataBase.Login(data)
	return isLogined, rank, err
}

func (m *Manager) GetLeaderBoard() (LeaderBoard, error) {
	leaderBoard, err := m.DataBase.GetLeaderBoard()
	return leaderBoard, err
}

func (m *Manager) GenerateToken(name string) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	err = m.DataBase.SetInfo(name, []string{"session_tocken"}, []string{token})
	return token, err
}

func (m *Manager) CheckToken(token string) (bool, error) {
	return m.DataBase.CheckToken(token)
}
