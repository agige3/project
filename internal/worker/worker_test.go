package worker

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"project/internal/user"
)

const (
	maxIdStatement  = "SELECT user_id FROM testing"
	dropStatement   = "DROP TABLE IF EXISTS testing"
	createStatement = `CREATE TABLE testing (
		user_id SERIAL PRIMARY KEY,
		user_name varchar(50),
		user_age int,
		group_id int,
		channel_ids int[]
	)`
)

func readMaxUserId(worker *Worker) (int, error) {
	row := worker.db.QueryRow(maxIdStatement)
	var userID int
	err := row.Scan(&userID)
	if err != nil {
		return -1, err
	}
	return userID, nil
}

func TestPrepareDB(t *testing.T) {
	db, err := GetDB("gopher", "pass", "localhost", "test", 5432)
	//client := getRedisClient("localhost:6379", "pass")
	//worker := NewWorker(client, db, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
		//fmt.Println(err)
	}
	t.Run("check", func(tt *testing.T) {
		err := db.Ping()
		if err != nil {
			tt.Fatalf("error pinging db: %s", err)
		}
	})

	t.Run("clean", func(tt *testing.T) {
		st, err := db.Prepare(dropStatement)
		if err != nil {
			tt.Fatalf("error preparing statement: %s", err)
		}
		_, err = st.Exec()
		if err != nil {
			tt.Fatalf("error cleaning db: %s", err)
		}
	})

	t.Run("create", func(tt *testing.T) {
		st, err := db.Prepare(createStatement)
		if err != nil {
			tt.Fatalf("error preparing statement: %s", err)
		}
		_, err = st.Exec()
		if err != nil {
			tt.Fatalf("error creating db: %s", err)
		}
	})

	t.Run("check empty", func(tt *testing.T) {
		st, err := db.Prepare("SELECT * FROM testing LIMIT 1")
		if err != nil {
			tt.Fatalf("error preparing statement: %s", err)
		}
		row := st.QueryRow()
		err = row.Scan()
		if err != sql.ErrNoRows {
			tt.Fatalf("db testing must be empty: %s", err)
		}
	})
}

func TestJust(t *testing.T) {
	db, err := GetDB("gopher", "pass", "localhost", "test", 5432)
	client := GetRedisClient("localhost:6379", "pass")
	worker := NewWorker(client, db, 5*time.Minute)
	ctx := context.Background()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("check", func(tt *testing.T) {
		err := db.Ping()
		if err != nil {
			tt.Fatalf("error pinging db: %s", err)
		}
		err = client.Ping(ctx).Err()
		if err != nil {
			tt.Fatalf("error pinging redis: %s", err)
		}
	})

	t.Run("checkWorker", func(tt *testing.T) {
		user := &user.User{Age: 23, Name: "Ivan", GroupID: 1, ChannelIDs: []int{1, 2, 3}}
		worker.addUser(ctx, "testing", user)
		id, err := readMaxUserId(worker)
		fmt.Println(id, err)
		pack, err := worker.GetUsers(ctx, "testing", 1, 1, false)
		//  pupupu от Полины
		if err != nil {
			tt.Fatalf("error getting user: %s", err)
		}
		if pack.Users == nil || len(pack.Users) != 1 {
			tt.Fatalf("incorrent amount of users: %d", len(pack.Users))
		}
		if !user.CompareWithUser(&pack.Users[0]) { // сделать с testify
			tt.Fatalf("users must be equial")
		}
		// rururu от Ралины
	})

}
