package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"project/internal/worker"
)

func TestGetUsers2_invalid(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h := NewHandler(nil)
	h.GetUsers(w, r)
	if w.Code != 400 {
		t.Error("expected error")
	}
	//u := url.Values{}
	//u.Add("")
}

func TestGetUsers_ok(t *testing.T) {

	db, err := worker.GetDB("gopher", "pass", "localhost", "test", 5432)
	if err != nil {
		t.Fatal(err)
	}
	client := worker.GetRedisClient("localhost:6379", "pass")
	worker := worker.NewWorker(client, db, 5*time.Minute)
	h := NewHandler(worker)

	type TestCase struct {
		ChannelID      string `json:"channel_id"`
		GroupID        string `json:"group_id"`
		UseLastVersion string `json:"use_last_version"`
	}

	testName := func(tc *TestCase) string {
		data, err := json.Marshal(tc)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		//require.NoError(t, err)
		return string(data)
	}

	for _, tc := range []*TestCase{
		{ChannelID: "1", GroupID: "2"},
		{ChannelID: "1", GroupID: "2", UseLastVersion: "true"},
		{ChannelID: "1", GroupID: "2", UseLastVersion: "false"},
	} {
		t.Run(testName(tc), func(t *testing.T) {
			v := url.Values{}
			if tc.ChannelID != "" {
				v.Add(channelName, tc.ChannelID)
			}
			if tc.GroupID != "" {
				v.Add(groupName, tc.GroupID)
			}
			if tc.UseLastVersion != "" {
				v.Add(useLastVersionName, tc.UseLastVersion)
			}
			u := fmt.Sprintf("/?%s", v.Encode())
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, u, nil)
			h.GetUsers(w, r)

			//require.Equal(t, http.StatusBadRequest, w.Code)
			if http.StatusOK != w.Code {
				t.Errorf("unexpected StatusCode:%d", w.Code)
			}
		})
	}
}

func TestGetUsers_invalid(t *testing.T) {

	type TestCase struct {
		ChannelID      string `json:"channel_id"`
		GroupID        string `json:"group_id"`
		UseLastVersion string `json:"use_last_version"`
	}

	testName := func(tc *TestCase) string {
		data, err := json.Marshal(tc)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		//require.NoError(t, err)
		return string(data)
	}

	for _, tc := range []*TestCase{
		{ChannelID: "1"},
		{ChannelID: "abc"},
		{GroupID: "2"},
		{GroupID: "abc"},
		{ChannelID: "1", GroupID: "abc"},
		{ChannelID: "abc", GroupID: "1"},
		{ChannelID: "1", UseLastVersion: "true"},
		{GroupID: "2", UseLastVersion: "false"},
		{UseLastVersion: "true"},
	} {
		t.Run(testName(tc), func(t *testing.T) {
			h := NewHandler(nil)
			v := url.Values{}
			if tc.ChannelID != "" {
				v.Add(channelName, tc.ChannelID)
			}
			if tc.GroupID != "" {
				v.Add(groupName, tc.GroupID)
			}
			if tc.UseLastVersion != "" {
				v.Add(useLastVersionName, tc.UseLastVersion)
			}
			u := fmt.Sprintf("/?%s", v.Encode())
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, u, nil)
			h.GetUsers(w, r)

			//require.Equal(t, http.StatusBadRequest, w.Code)
			if http.StatusBadRequest != w.Code {
				t.Errorf("unexpected StatusCode:%d", w.Code)
			}
		})
	}
}

func TestAddUser_ok(t *testing.T) {

	db, err := worker.GetDB("gopher", "pass", "localhost", "test", 5432)
	if err != nil {
		t.Fatal(err)
	}
	client := worker.GetRedisClient("localhost:6379", "pass")
	worker := worker.NewWorker(client, db, 5*time.Minute)
	h := NewHandler(worker)

	type TestCase struct {
		ChannelID      string `json:"channel_id"`
		GroupID        string `json:"group_id"`
		UseLastVersion string `json:"use_last_version"`
	}

	testName := func(tc *TestCase) string {
		data, err := json.Marshal(tc)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		//require.NoError(t, err)
		return string(data)
	}

	for _, tc := range []*TestCase{
		{ChannelID: "1", GroupID: "2"},
		{ChannelID: "1", GroupID: "2", UseLastVersion: "true"},
		{ChannelID: "1", GroupID: "2", UseLastVersion: "false"},
	} {
		t.Run(testName(tc), func(t *testing.T) {

			message := map[string]interface{}{
				"user_age":    12,
				"user_name":   "Ivan",
				"group_id":    5,
				"channel_ids": []int{1, 2, 3},
			}

			bytesRepresentation, err := json.Marshal(message)
			if err != nil {
				t.Fatalf("error during preparing body: %s", err)
			}

			body := bytes.NewBuffer(bytesRepresentation)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:6029", body)
			h.AddUser(w, r)

			//require.Equal(t, http.StatusBadRequest, w.Code)
			if http.StatusCreated != w.Code {
				t.Errorf("unexpected StatusCode:%d", w.Code)
			}
		})
	}
}

func TestAddUser_invalid(t *testing.T) {

	h := NewHandler(nil)

	type TestCase struct {
		Name string
		Body io.Reader
	}

	for _, tc := range []*TestCase{
		{Name: "noBody", Body: nil},
		{Name: "incorrectBody", Body: bytes.NewBufferString(`\lastname:\`)},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:6029", tc.Body)
			h.AddUser(w, r)

			if http.StatusBadRequest != w.Code {
				t.Errorf("unexpected StatusCode:%d", w.Code)
			}
		})
	}
}
