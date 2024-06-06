package worker

import (
	"context"
	"database/sql"
	"fmt"
	"project/internal/user"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

const (
	selectStatement = `SELECT user_id, user_name, user_age, group_id, channel_ids FROM users 
						WHERE  group_id = %d AND %d = ANY (channel_ids)`
	connectStr = "user=%s password=%s host=%s port=%d database=%s"

	insertStatement = `INSERT INTO users(user_name, user_age, group_id, channel_ids)
	VALUES ($1, $2, $3, $4) RETURNING user_id`
)

func getDB(user, password, host, databaseName string, port int) (*sql.DB, error) {
	db, err := sql.Open("pgx", fmt.Sprintf(connectStr, user, password, host, port, databaseName))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (worker *Worker) AddUserWithParameters(ctx context.Context, age, groupID int, name string, channelIDs []int) (int, error) {
	row := worker.db.QueryRowContext(
		ctx,
		insertStatement,
		name, age, groupID, channelIDs)
	var user_id int
	err := row.Scan(&user_id)
	return user_id, err
}

func (worker *Worker) AddUser(ctx context.Context, user *user.User) (int, error) {
	row := worker.db.QueryRowContext(
		ctx,
		insertStatement,
		user.Name, user.Age, user.GroupID, user.ChannelIDs)
	var user_id int
	err := row.Scan(&user_id)
	return user_id, err
}

func (worker *Worker) getUsersFromDB(ctx context.Context, groupID, channelID int) (user.PackOfUsers, error) {
	rows, err := worker.db.QueryContext(ctx,
		fmt.Sprintf(selectStatement, groupID, channelID))
	if err != nil {
		return user.PackOfUsers{}, err
	}
	defer rows.Close()
	pack := user.PackOfUsers{}
	for rows.Next() {
		var id, age, group_id int
		var channel_idsINT32 []int32
		var name string
		if err := rows.Scan(&id, &name, &age, &group_id, (*pq.Int32Array)(&channel_idsINT32)); err != nil {
			return pack, fmt.Errorf("not all rows readed")
		}
		channel_ids := make([]int, len(channel_idsINT32))
		for i := 0; i < len(channel_idsINT32); i++ {
			channel_ids[i] = int(channel_idsINT32[i])
		}
		u := user.User{ID: id, Age: age, Name: name, GroupID: group_id, ChannelIDs: channel_ids}
		pack.Users = append(pack.Users, u)
	}
	return pack, nil
}
