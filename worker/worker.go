package worker

import (
	"context"
	"database/sql"
	"fmt"
	"project/user"
	"time"

	"github.com/redis/go-redis/v9"
)

type Worker struct {
	db          *sql.DB
	redisClient *redis.Client
	actualTime  time.Duration
}

func getWorker(r *redis.Client, db *sql.DB, t time.Duration) Worker {
	return Worker{db: db, redisClient: r, actualTime: t}
}

func tryConnectAndGetWorker(DBuser, DBpassword, DBhost, DBdatabaseName string, DBport int, redisAddr, redisPass string, actualTime time.Duration) (w Worker, err error) {
	db, err := getDB(DBuser, DBpassword, DBhost, DBdatabaseName, DBport)
	if err != nil {
		return Worker{}, fmt.Errorf("cant connect to db")
	}
	red := getRedisClient(redisAddr, redisPass)
	w = getWorker(red, db, actualTime)
	return w, nil
}

func (w *Worker) GetUsers(ctx context.Context, groupID, channelID int, needMoreActual bool) (user.PackOfUsers, error) {
	if !needMoreActual {
		pack, err := w.getUsersFromRedis(ctx, groupID, channelID)
		if err != nil {
			return user.PackOfUsers{}, err
		}
		return pack, nil
	}
	// если нужны свежие данные или в кеше нет нужной записи
	pack, err := w.getUsersFromDB(ctx, groupID, channelID)
	if err != nil {
		return user.PackOfUsers{}, err
	}
	err = w.redisClient.Set(ctx, makeKeyForRedis(groupID, channelID), pack, w.actualTime).Err()
	if err != nil {
		// добавить в redis не удалось, не критично
	}
	return pack, nil
}
