package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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
			r := httptest.NewRequest("GET", u, nil)
			h.GetUsers(w, r)

			//require.Equal(t, http.StatusBadRequest, w.Code)
			if http.StatusBadRequest != w.Code {
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
			r := httptest.NewRequest("GET", u, nil)
			h.GetUsers(w, r)

			//require.Equal(t, http.StatusBadRequest, w.Code)
			if http.StatusBadRequest != w.Code {
				t.Errorf("unexpected StatusCode:%d", w.Code)
			}
		})
	}
}
