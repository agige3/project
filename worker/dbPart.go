package worker

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

func getDB(user, password, host, databaseName string, port int) (*sql.DB, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("user=%s password=%s host=%s port=%d database=%s", user, password, host, port, databaseName))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getWorker(r *redis.Client, db *sql.DB, t time.Duration) Worker {
	return Worker{db: db, redisClient: r, actualTime: t}
}

func (w *Worker) addUser(ctx context.Context, age, groupID int, name string, channelIDs []int) error {
	//usr := user.User{Age: age, Name: name, GroupID: groupID, ChannelIDs: channelIDs}
	// делаем запрос к базе данных
	panic("implement me")
}
