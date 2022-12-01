package server

import (
	"encoding/json"
	d "gritface/database"
	logger "gritface/log"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// to store the password in the database
var HashedPassword string

// function to check if email is valid
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)
	return emailRegex.MatchString(e)
}

// function to check if inputs are valid
func IsAscii(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c <= 33 || c > 127 {
			return false
		}
	}
	return true
}

// hash password returned the password string as a hash to be stored in the database
// this is done for security reasons
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// prevents sql injection
func EscapeString(value string) string {
	var sb strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch c {
		case '\\', 0, '\n', '\r', '\'', '"':
			sb.WriteByte('\\')
			sb.WriteByte(c)
		case '\032':
			sb.WriteByte('\\')
			sb.WriteByte('Z')
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

// function to sign up a user
func SignUp(w http.ResponseWriter, r *http.Request) (string, bool) {
	if r.Method != "POST" {
		return "Wrong method", false
	}

	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.WTL(err.Error(), true)
		return "Error: reading json sign up request from user", false
	}

	// Unmarshal
	var user NewUser
	err = json.Unmarshal(req, &user)
	if err != nil {
		return "Error: unsuccessful in unmarshaling data from user", false
	}

	// Check for session
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		// No session found
		return err.Error(), false
	}
	// Check if user is logged in
	if uid != "0" {
		// User is logged in
		errMsg := "Error: user is logged in"
		return errMsg, false
	} // User is not logged in

	// user trying to sign up
	// get form values
	name := EscapeString(user.Name)
	email := strings.ToLower(EscapeString(user.Email))
	// no need to escape password because its hashed before being stored
	password := user.Password
	confirmPwd := user.ConfirmPwd

	if name == "" || email == "" || password == "" || confirmPwd == "" {
		errMsg := "Fill in all the required fields"
		return errMsg, false
	}

	// check if email is valid
	if !isEmailValid(email) {
		errMsg := "Email is not valid"
		return errMsg, false
	}
	if password != confirmPwd {
		errMsg := "Passwords do not match"
		return errMsg, false
	}

	// make sure that no fields are empty or non ascii
	if !IsAscii(name) || !IsAscii(email) || !IsAscii(password) || !IsAscii(confirmPwd) {
		errMsg := "Non-ascii characters found"
		return errMsg, false
	}

	// user level will by 1 by default i.e registered user
	userLevel := 1

	// hash password
	HashedPassword, err := HashPassword(password)
	if err != nil {
		errMsg := "Internal server error in hashing password"
		return errMsg, false
	}
	// connect to database
	db, err := d.DbConnect()
	if err != nil {
		return err.Error(), false
	}

	defer db.Close()

	// add user to db
	UID, err := d.InsertUsers(db, name, email, HashedPassword, userLevel)

	if err != nil {
		// check if pass and confirm pass match
		if strings.Contains(err.Error(), "name") {
			return "Error: user name is already used.", false
		}
		if strings.Contains(err.Error(), "email") {
			return "Error: email is already used.", false
		}

	}

	// Log new signup
	logger.WTL("User signing up. Name:"+name+", email:"+email+" with user id "+strconv.Itoa(UID), false)

	successMsg := "User has been successfully registered"
	return successMsg, true

}
