package worker

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"project/internal/user"
	"time"

	"github.com/redis/go-redis/v9"
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
	db, err := getDB(DBuser, DBpassword, DBhost, DBdatabaseName, DBport)
	if err != nil {
		return &Worker{}, fmt.Errorf("cant connect to db")
	}
	red := getRedisClient(redisAddr, redisPass)
	worker = NewWorker(red, db, actualTime)
	return worker, nil
}

func (worker *Worker) GetUsers(ctx context.Context, groupID, channelID int, needMoreActual bool) (user.PackOfUsers, error) {
	if !needMoreActual {
		pack, err := worker.getUsersFromRedis(ctx, groupID, channelID)
		if err == nil {
			return pack, nil
		}
		// redis сломался
		if err != redis.Nil {
			log.Println("error in redis:", err)
		}
	}
	// если нужны свежие данные или в кеше нет нужной записи
	pack, err := worker.getUsersFromDB(ctx, groupID, channelID)
	if err != nil {
		return user.PackOfUsers{}, err
	}
	err = worker.redisClient.Set(ctx, makeKeyForRedis(groupID, channelID), pack, worker.ActualTime).Err()
	if err != nil {
		// добавить в redis не удалось, не критично
	}
	return pack, nil
}
