package worker

import (
	"context"
	"database/sql"
	"fmt"
	"project/user"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

const (
	selectStatement = `SELECT user_id, user_name, user_age, group_id, channel_ids FROM users 
						WHERE  group_id = %d AND %d = ANY (channel_ids)`
	connectStr = "user=%s password=%s host=%s port=%d database=%s"

	insertStatement = `INSERT INTO users(user_name, user_age, group_id, channel_ids)
	VALUES ($1, $2, $3, $4)`
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

func (w *Worker) addUser(ctx context.Context, age, groupID int, name string, channelIDs []int) error {
	_, err := w.db.ExecContext(
		ctx,
		insertStatement,
		name, age, groupID, channelIDs)
	return err
}

func (w *Worker) getUsersFromDB(ctx context.Context, groupID, channelID int) (user.PackOfUsers, error) {
	rows, err := w.db.QueryContext(ctx,
		fmt.Sprintf(selectStatement, groupID, channelID))
	if err != nil {
		return user.PackOfUsers{}, err
	}
	defer rows.Close()
	pack := user.PackOfUsers{}
	for rows.Next() {
		var id, age, group_id int
		var channel_ids []int
		var name string
		if err := rows.Scan(&id, &name, &age, &group_id, pq.Array(&channel_ids)); err != nil {
			return pack, fmt.Errorf("not all rows readed")
		}
		u := user.User{ID: id, Age: age, Name: name, GroupID: group_id, ChannelIDs: channel_ids}
		pack.Users = append(pack.Users, u)
	}
	return pack, nil
}
