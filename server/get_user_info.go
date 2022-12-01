package server

import (
	"encoding/json"
	d "gritface/database"
	logger "gritface/log"
	"net/http"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) (string, bool) {
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		return err.Error(), false
	}
	if uid == "0" {
		return "invalid session", false
	}

	db, err := d.DbConnect()
	if err != nil {
		return err.Error(), false
	}

	defer db.Close()

	user := make(map[string]string)
	user["user_id"] = uid
	users, err := d.GetUsers(db, user)
	if err != nil {
		logger.WTL(err.Error(), true)
	}

	var info Info
	info.Username = users[0].Name

	jsonInfo, err := json.Marshal(info)
	if err != nil {
		return err.Error(), false
	}

	return string(jsonInfo), true
}
