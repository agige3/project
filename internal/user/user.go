package user

import (
	"encoding/json"
)

type User struct {
	ID         int    `json:"user_id"`
	Age        int    `json:"user_age"`
	Name       string `json:"user_name"`
	GroupID    int    `json:"group_id"`
	ChannelIDs []int  `json:"channel_ids"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(*u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

func (u *User) CompareWithUser(user *User) bool {
	if u.Age != user.Age || u.Name != user.Name || u.GroupID != user.GroupID || len(u.ChannelIDs) != len(user.ChannelIDs) {
		return false
	}
	for i := 0; i < len(u.ChannelIDs); i++ {
		if u.ChannelIDs[i] != user.ChannelIDs[i] {
			return false
		}
	}
	return true
}

type PackOfUsers struct {
	Users []User `json:"users"`
}

func (p PackOfUsers) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PackOfUsers) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &p)
}
