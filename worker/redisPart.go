package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"project/user"

	"github.com/redis/go-redis/v9"
)

// можно поменять
func makeKeyForRedis(groupID, channelID int) string {
	return fmt.Sprintf("group:%d;channel:%d", groupID, channelID)
}

func getRedisClient(addr, pass string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})
}

func (w *Worker) getUsersFromRedis(ctx context.Context, groupID, channelID int) (user.PackOfUsers, error) {
	pack, err := w.redisClient.Get(ctx, makeKeyForRedis(groupID, channelID)).Result()
	// какая - то неожиданная ошибка
	if err != nil {
		return user.PackOfUsers{}, err
	}
	var packFromRedis user.PackOfUsers
	err = json.Unmarshal([]byte(pack), &packFromRedis)
	if err != nil {
		return user.PackOfUsers{}, err
	}
	return packFromRedis, nil
}
