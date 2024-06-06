package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"project/internal/user"
	"project/internal/worker"
)

type Handler struct {
	Worker *worker.Worker
}

func NewHandler(worker *worker.Worker) *Handler {
	return &Handler{Worker: worker}
}

func JSONError(statusCode int, msg string, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	type Error struct {
		Message *string `json:"message,omitempty"`
	}
	_ = json.NewEncoder(w).Encode(
		Error{
			Message: &msg,
		},
	)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query()
	if !u.Has(channelName) {
		JSONError(400, fmt.Sprintf("%s required, but not found", channelName), w)
		return
	}
	channel := u.Get(channelName)
	channelID, err := strconv.Atoi(channel)
	if err != nil {
		JSONError(400, fmt.Sprintf("cannot convert %s: %s to int", channelName, channel), w)
		return
	}
	if !u.Has(groupName) {
		JSONError(400, fmt.Sprintf("%s required, but not foind", groupName), w)
		return
	}
	group := u.Get(groupName)
	groupID, err := strconv.Atoi(group)
	if err != nil {
		JSONError(400, fmt.Sprintf("cannot convert %s: %s to int", groupName, group), w)
		return
	}
	useLastVersion := useLastVersionDefault
	if u.Has(useLastVersionName) {
		useLastVersionString := u.Get(useLastVersionName)
		lastVersionParsing, err := strconv.ParseBool(u.Get(useLastVersionString))
		if err != nil {
			JSONError(400, fmt.Sprintf("cannot convert %s: %s to bool", useLastVersionName, useLastVersionString), w)
			return
		}
		useLastVersion = lastVersionParsing
	}

	pack, err := h.Worker.GetUsers(r.Context(), groupID, channelID, useLastVersion)
	if err != nil {
		// обработка ошибок
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(pack)
}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		// ошибка
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// ошибка
		/*
			response := struct {
				Error error
			}{Error: err}
			encoder := json.NewEncoder(w)
			encoder.Encode(response)
			return
		*/
	}
	user := user.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		// ошибка
	}
	id, err := h.Worker.AddUser(r.Context(), &user)
	if err != nil {
		// ошибка
	}
	w.WriteHeader(201)
	response := struct {
		UserID int
	}{UserID: id}
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}
