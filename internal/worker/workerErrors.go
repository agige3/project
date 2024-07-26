package worker

import "errors"

var (
	ErrNoUsers = errors.New("can`t find users with given parameters")
)
