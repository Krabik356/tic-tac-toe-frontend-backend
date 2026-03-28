package database

import (
	"backend/server/internal/logic"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

type DBManager struct {
	Pool *pgxpool.Pool
}

func NewDBManager() *DBManager {

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", dbUser, dbPassword, dbName))
	if err != nil {

	}
	return &DBManager{
		Pool: pool,
	}
}

func (dbM *DBManager) Register(regData logic.RegisterData) error {
	_, err := dbM.Pool.Exec(context.Background(), "INSERT INTO users(name, password, rank) VALUES($1, $2, $3)", regData.Name, regData.Password, 0)
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
