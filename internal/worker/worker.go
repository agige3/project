package worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/agige3/project/internal/user"
)

type Worker struct {
	db          *sql.DB
	redisClient *redis.Client
	ActualTime  time.Duration
}

func NewWorker(r *redis.Client, db *sql.DB, t time.Duration) *Worker {
	return &Worker{db: db, redisClient: r, ActualTime: t}
}

func tryConnectAndGetWorker(DBuser, DBpassword, DBhost, DBdatabaseName string, DBport int, redisAddr, redisPass string, actualTime time.Duration) (worker *Worker, err error) {
	db, err := GetDB(DBuser, DBpassword, DBhost, DBdatabaseName, DBport)
	if err != nil {
		return &Worker{}, fmt.Errorf("cant connect to db")
	}
	red := GetRedisClient(redisAddr, redisPass)
	worker = NewWorker(red, db, actualTime)
	return worker, nil
}

func (worker *Worker) GetUsers(ctx context.Context, tableName string, groupID, channelID int, needMoreActual bool) (user.PackOfUsers, error) {
	if !needMoreActual {
		pack, err := worker.getUsersFromRedis(ctx, groupID, channelID)
		if err == nil {
			return pack, nil
		}
		// redis сломался
		if err != redis.Nil {
			log.Println("error in redis: ", err)
		}
	}
	// если нужны свежие данные или в кеше нет нужной записи
	pack, err := worker.getUsersFromDB(ctx, tableName, groupID, channelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.PackOfUsers{}, ErrNoUsers
		}
		return user.PackOfUsers{}, err
	}
	err = worker.redisClient.Set(ctx, makeKeyForRedis(groupID, channelID), pack, worker.ActualTime).Err()
	if err != nil {
		// добавить в redis не удалось, не критично
		log.Println("can`t write to redis:", err)
	}
	return pack, nil
}

func (worker *Worker) AddUserWithParameters(ctx context.Context, tableName string, age, groupID int, name string, channelIDs []int) (int, error) {
	user := &user.User{Age: age, GroupID: groupID, Name: name, ChannelIDs: channelIDs}
	return worker.addUser(ctx, tableName, user)

}

func (worker *Worker) AddUser(ctx context.Context, tableName string, user *user.User) (int, error) {
	return worker.addUser(ctx, tableName, user)
}
