package logic

import "github.com/gorilla/websocket"

type RegisterData struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Token     string `json:"-"`
	TokenTime string `json:"-"`
}

type UnRegisterData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SuccessfulRegistartion struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	Token  string `json:"token"`
}

type SuccessfulLogin struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	Rank   int    `json:"rank"`
	Token  string `json:"token"`
}

type Leaders struct {
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

type LeaderBoard struct {
	Data []Leaders `json:"data"`
}

type SuccessfulLeaderBoard struct {
	Status string      `json:"status"`
	LB     LeaderBoard `json:"lb"`
}

type Client struct {
	Name     string
	Marker   string
	Conn     *websocket.Conn
	Send     chan interface{}
	RoomChan chan interface{}
}

// RoomsMovement
func NewClient(name string, conn *websocket.Conn) *Client {
	return &Client{
		Name:   name,
		Marker: "",
		Conn:   conn,
		Send:   make(chan interface{}, 10),
	}
}

type GameStart struct {
	Status string `json:"status"`
	Name1  string `json:"name1"`
	Name2  string `json:"name2"`
}

type RoomsError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type Movement struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type RoomsMovement struct {
	Client *Client
	X      int
	Y      int
}

type MoveMade struct {
	Who string `json:"who"`
	X   int    `json:"x"`
	Y   int    `json:"y"`
}

type ClientDisconect struct {
	Status string
	Client *Client
}

type AfterGame struct {
	Status  string `json:"status"`
	Name1   string `json:"name1"`
	Name2   string `json:"name2"`
	Amount1 int    `json:"amount1"`
	Amount2 int    `json:"amount2"`
}
