package logic

import (
	"errors"
	"os"
	"strconv"
)

type Room struct {
	Id             int
	Status         string
	Clients        [2]*Client
	Turn           int
	Field          [3][3]string
	Input          chan interface{}
	Win            int
	Draw           int
	Withdraw       int
	DeleteRoomChan chan int
	Database       interface {
		Retribution(string, string, int, int) error
	}
	IsClosed bool
}

func NewRoom(id int, db interface {
	Retribution(string, string, int, int) error
}, deleteChan chan int) *Room {
	withdraw, _ := strconv.Atoi(os.Getenv("SERVER_WITHDRAW"))
	win, _ := strconv.Atoi(os.Getenv("SERVER_WIN"))
	draw, _ := strconv.Atoi(os.Getenv("SERVER_DRAW"))

	return &Room{
		Id:      id,
		Status:  "waiting",
		Clients: [2]*Client{},
		Turn:    0,
		Field: [3][3]string{
			{"*", "*", "*"},
			{"*", "*", "*"},
			{"*", "*", "*"},
		},
		Input:          make(chan interface{}, 100),
		Win:            win,
		Draw:           draw,
		Withdraw:       withdraw,
		Database:       db,
		DeleteRoomChan: deleteChan,
	}
}

func (r *Room) isCorrectMove(row, col int, who *Client) error {
	if r.Field[row][col] == "*" {
		if who == r.Clients[r.Turn] {
			return nil
		} else {
			return errors.New("Not your turn")
		}
	}
	return errors.New("Incorrect move")
}

func winnerId(el string) int {
	if el == "X" {
		return 0
	}
	return 1
}

func (r *Room) isWin() (int, bool) {
	field := r.Field
	var toCompare string
	for _, row := range field {
		isWin := true
		toCompare = row[0]
		for _, el := range row {
			if el != toCompare {
				isWin = false
				break
			}
		}
		if isWin == true && toCompare != "*" {
			return winnerId(toCompare), true
		}
	}

	for c := range field {
		isWin := true
		toCompare = field[0][c]
		for r := range field {
			if field[r][c] != toCompare {
				isWin = false
				break
			}
		}
		if isWin == true && toCompare != "*" {
			return winnerId(toCompare), true
		}
	}

	isWin := true
	toCompare = field[0][0]
	for i := range field {
		if field[i][i] != toCompare {
			isWin = false
			break
		}
	}
	if isWin == true && toCompare != "*" {
		return winnerId(toCompare), true
	}

	isWin = true
	lenField := len(field) - 1
	toCompare = field[lenField][lenField]
	for i := range field {
		if field[lenField-i][i] != toCompare {
			isWin = false
			break
		}
	}
	if isWin == true && toCompare != "*" {
		return winnerId(toCompare), true
	}

	return 0, false
}

func (room *Room) isPlayable() bool {
	for _, row := range room.Field {
		for _, el := range row {
			if el == "*" {
				return true
			}
		}
	}
	return false
}

func (r *Room) sendToClients(data interface{}) {
	if r.IsClosed {
		return
	}

	for _, client := range r.Clients {
		if client != nil {
			select {
			case client.Send <- data:
			default:

			}
		}
	}
}

func (r *Room) CloseRoom() {
	if r.IsClosed {
		return
	}
	r.IsClosed = true

	for _, c := range r.Clients {
		if c != nil {
			close(c.Send)
			c.RoomChan = nil
		}
	}
}

func (r *Room) Run() {
	r.sendToClients(GameStart{
		Status: "game",
		Name1:  r.Clients[0].Name,
		Name2:  r.Clients[1].Name,
	})
	for _, client := range r.Clients {
		client.RoomChan = r.Input
	}
	for data := range r.Input {
		switch t := data.(type) {
		case RoomsMovement:
			movement := t
			err := r.isCorrectMove(movement.X, movement.Y, movement.Client)
			if err != nil {
				movement.Client.Send <- RoomsError{
					Status: "error",
					Error:  err.Error(),
				}
			} else {
				marker := movement.Client.Marker
				r.Field[movement.X][movement.Y] = marker
				r.sendToClients(MoveMade{
					Who: marker,
					X:   movement.X,
					Y:   movement.Y,
				})
				r.Turn = 1 - r.Turn

				id, isWin := r.isWin()
				var loserId int
				if id == 0 {
					loserId = 1
				} else {
					loserId = 0
				}

				if isWin {
					winnerName := r.Clients[id].Name
					loserName := r.Clients[loserId].Name
					err := r.Database.Retribution(winnerName, loserName, r.Win, r.Withdraw)
					if err != nil {

					}
					r.sendToClients(AfterGame{
						Status:  "win",
						Name1:   winnerName,
						Name2:   loserName,
						Amount1: r.Win,
						Amount2: r.Withdraw,
					})
					r.DeleteRoomChan <- r.Id
					return
				} else {
					isPl := r.isPlayable()
					if !isPl {
						name1 := r.Clients[0].Name
						name2 := r.Clients[1].Name
						err := r.Database.Retribution(name1, name2, r.Draw, r.Draw)
						if err != nil {

						}
						r.sendToClients(AfterGame{
							Status:  "draw",
							Name1:   name1,
							Name2:   name2,
							Amount1: r.Draw,
							Amount2: r.Draw,
						})
						r.DeleteRoomChan <- r.Id
						return
					}
				}
			}
		case ClientDisconect:
			winId := 0
			if t.Client.Marker == "X" {
				winId = 1
			}
			name1 := t.Client.Name
			name2 := r.Clients[winId].Name
			err := r.Database.Retribution(name2, name1, r.Win, r.Withdraw)
			if err != nil {

			}
			r.sendToClients(AfterGame{
				Status:  "discowin",
				Name1:   name2,
				Name2:   name1,
				Amount1: r.Win,
				Amount2: r.Withdraw,
			})
			r.DeleteRoomChan <- r.Id
			return
		}
	}
}
