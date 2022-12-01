package server

import (
	"encoding/json"
	d "gritface/database"
	logger "gritface/log"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var SuccessfulLogin []string

// HashPassword returns the password string as a hash to be stored in the database
func passwordMatch(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(w http.ResponseWriter, r *http.Request) (string, bool) {
	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.WTL(err.Error(), true)
		return "Error: reading json log in request from user", false
	}

	// Unmarshal
	var user d.Users
	err = json.Unmarshal(req, &user)
	if err != nil {
		return "Error: unsuccessful in unmarshaling log in data from user", false
	}

	// Get session
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		// No session found
		logger.WTL(err.Error(), false)
		return err.Error(), false
	}

	// Check if user is logged in
	if uid != "0" {
		// User is logged in, redirect to front page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return "User is logged in", true
	}

	// get the email and password
	email := EscapeString(user.Email)
	password := user.Password

	// check if the email and password are ascii
	if !IsAscii(email) || !IsAscii(password) {
		return "Error: invalid email or password", false
	}

	// retrieve user password from database
	db, err := d.DbConnect()
	if err != nil {
		return err.Error(), false
	}

	defer db.Close()

	LoginUser := make(map[string]string)
	LoginUser["email"] = email
	users, err := d.GetUsers(db, LoginUser)
	if err != nil {
		logger.WTL((err.Error()), true)
	}

	// if no or more than 1 record found, return error
	if len(users) != 1 {
		return "Error: email or password is not found!", false
	}

	// check if the password match the one with database
	if !passwordMatch(password, users[0].Password) {
		return "Error: email or password is not found", false
	}

	num := users[0].User_id

	// Check that no other session is logged in with this user
	for thisSessionId, thisSessionData := range sessionManager.sessions {
		if thisSessionData.UId == strconv.Itoa(num) {
			delete(sessionManager.sessions, thisSessionId)
		}
	}

	err = sessionManager.setSessionUID(num, w, r)
	if err != nil {
		logger.WTL("Failed to set session uid, "+err.Error(), true)
		return err.Error(), false
	}
	return strconv.Itoa(num), true
}
