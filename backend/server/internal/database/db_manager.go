package database

import (
	"backend/server/internal/logic"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBManager struct {
	Pool *pgxpool.Pool
}

func NewDBManager() *DBManager {

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@db:5432/%s", dbUser, dbPassword, dbName))
	if err != nil {

	}
	return &DBManager{
		Pool: pool,
	}
}

func (dbM *DBManager) Register(regData logic.RegisterData) error {
	_, err := dbM.Pool.Exec(context.Background(), "INSERT INTO users(name, password, rank, session_token) VALUES($1, $2, $3, $4)", regData.Name, regData.Password, regData.Token, 0)
	return err
}

func (dbM *DBManager) Login(regData logic.RegisterData) (bool, int, error) {
	var (
		isCorrect bool
		rank      int
	)
	err := dbM.Pool.QueryRow(context.Background(), "SELECT 1, rank FROM users WHERE name=$1 AND password=$2", regData.Name, regData.Password).Scan(&isCorrect, &rank)
	return isCorrect, rank, err
}

func (dbM *DBManager) SetInfo(name string, info []string, what []string) error {
	if len(info) != len(what) {
		return errors.New("Info != what")
	}

	tx, err := dbM.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	for id, which := range info {
		command := fmt.Sprintf("UPDATE users SET %s=$2 WHERE name=$3", which)
		_, err := tx.Exec(context.Background(), command, what[id], name)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (dbM *DBManager) CheckToken(token string) (bool, error) {
	var isCorrect bool

	err := dbM.Pool.QueryRow(context.Background(), "SELECT 1 FROM users WHERE session_token=$1", token).Scan(&isCorrect)

	return isCorrect, err
}

func (dbM *DBManager) GetLeaderBoard() (logic.LeaderBoard, error) {
	data, err := dbM.Pool.Query(context.Background(), "SELECT name, rank FROM users ORDER BY rank DESC LIMIT 10")
	if err != nil {
		return logic.LeaderBoard{}, err
	}
	defer data.Close()
	var result logic.LeaderBoard
	for data.Next() {
		var (
			name string
			rank int
		)
		data.Scan(&name, &rank)
		result.Data = append(result.Data, logic.Leaders{
			Name: name,
			Rank: rank,
		})
	}
	return result, nil
}
