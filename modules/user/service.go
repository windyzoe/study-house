package userM

import (
	"errors"

	"github.com/windyzoe/study-house/db"
)

func GetUser(name string) (map[string]interface{}, error) {
	mapper := map[string]string{
		"Id":       "u.Id",
		"Name":     "u.Name",
		"Password": "u.Password",
	}
	userMap := db.Query(mapper, "User u where Name='"+name+"'")
	if len(userMap) != 1 {
		return nil, errors.New("No User")
	} else {
		return userMap[0], nil
	}
}
