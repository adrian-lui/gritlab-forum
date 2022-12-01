package server

import (
	"encoding/json"
	d "gritface/database"
	logger "gritface/log"
	"io"
	"net/http"
	"strconv"
	"time"
)

func addComment(w http.ResponseWriter, r *http.Request) (string, bool) {
	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.WTL(err.Error(), true)
		return "Error: reading json log in request from user", false
	}
	// Unmarshal
	var comment Comments
	err = json.Unmarshal(req, &comment)
	if err != nil {
		logger.WTL(err.Error(), true)
		return "Error: unsuccessful in unmarshaling log in data from user", false
	}

	// Get session
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		logger.WTL(err.Error(), true)
		return err.Error(), false
	}

	// Check if session is alive
	uID, err := strconv.Atoi(uid)
	if err != nil {
		// For some reason the uid is not numeric
		logger.WTL(err.Error(), true)
		return err.Error(), false
	}
	if uID < 1 {
		// No active session
		logger.WTL("Non-logged in user tried to make a comment", true)
		return "User is not logged in", false
	}

	// Connect to db
	db, err := d.DbConnect()
	if err != nil {
		logger.WTL(err.Error(), true)
		return err.Error(), false
	}

	defer db.Close()

	// Insert comment to database
	commentID, err := d.InsertComment(db, comment.Post_id, uID, comment.Body, time.Now().String())
	if err != nil {
		logger.WTL(err.Error(), true)
		return err.Error(), false
	}

	// Prepare return of new comment id
	commentID_str := strconv.Itoa(commentID)
	return commentID_str, true
}
