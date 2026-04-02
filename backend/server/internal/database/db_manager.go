package database

import (
	"backend/server/internal/logic"
	"context"
	"fmt"
	"log"
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
	_, err := dbM.Pool.Exec(context.Background(), "INSERT INTO users(name, password, session_token, token_time, rank) VALUES($1, $2, $3, $4, $5)", regData.Name, regData.Password, regData.Token, regData.TokenTime, 0)
	return err
}

func (dbM *DBManager) Login(regData logic.RegisterData) (bool, int, error) {
	var (
		isCorrect bool
		rank      int
	)
	err := dbM.Pool.QueryRow(context.Background(), "SELECT TRUE, rank FROM users WHERE name=$1 AND password=$2", regData.Name, regData.Password).Scan(&isCorrect, &rank)
	return isCorrect, rank, err
}

func (dbM *DBManager) SetToken(name string, token string, tokenTime string) error {
	_, err := dbM.Pool.Exec(context.Background(), "UPDATE users SET session_token=$1, token_time=$2 WHERE name=$3", token, tokenTime, name)
	return err
}

func (dbM *DBManager) GetName(token string) (string, error) {
	var name string
	err := dbM.Pool.QueryRow(context.Background(), "SELECT name FROM users WHERE session_token=$1", token).Scan(&name)
	return name, err
}

func (dbM *DBManager) CheckToken(token string) (bool, string, error) {
	var (
		isCorrect bool
		tokenTime string
	)
	log.Println(token)
	err := dbM.Pool.QueryRow(context.Background(), "SELECT TRUE, token_time FROM users WHERE session_token=$1", token).Scan(&isCorrect, &tokenTime)

	return isCorrect, tokenTime, err
}

func (dbM *DBManager) Retribution(name1, name2 string, amount1, amount2 int) error {
	tx, err := dbM.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "UPDATE users SET rank=rank+$2 WHERE name=$1", name1, amount1)
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), "UPDATE users SET rank=rank+$2 WHERE name=$1", name2, amount2)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
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
