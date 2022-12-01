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

func addPostText(w http.ResponseWriter, r *http.Request) (string, bool) {
	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.WTL(err.Error(), true)
		return "Error: reading json log in request from user", false
	}

	// Unmarshal
	var post newPosts
	err = json.Unmarshal(req, &post)
	if err != nil {
		logger.WTL(err.Error(), true)
		return "Error: unsuccessful in unmarshaling log in data from user", false
	}

	// Get session id
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		// No session found,
		logger.WTL(err.Error(), true)
		return err.Error(), false
	}

	uID, err := strconv.Atoi(uid)
	if err != nil {
		logger.WTL(err.Error(), true)
		return err.Error(), false
	}
	// Check for active session
	if uID < 1 {
		logger.WTL("User without active session tried to add post", false)
		return "Not logged in", false
	}

	// Now timetamp (ts)
	nowT := time.Now()
	nowTS := nowT.Unix()

	// Check for session last post insert ts
	lastInsertTS, err := sessionManager.GetSessionVariable(w, r, "last_post_insert")
	if err != nil {
		if err.Error() != "Value not set" {
			// Something went extremely wrong
			logger.WTL(err.Error(), true)
			return err.Error(), false
		}
		// Else value was not set which is okay
	}

	// Check when last post was created (bot-spam prevention)
	if lastInsertTS != nil {
		// This is okay
		if lastInsertTS.(int64)+5 > nowTS {
			logger.WTL("User "+uid+" tried to create a new post during cooldown", false)
			return "Add new post cooldown!", false
		}
	}

	// Get database connection
	db, err := d.DbConnect()
	if err != nil {
		logger.WTL(err.Error(), true)
		return err.Error(), false
	}

	defer db.Close()

	// Insert post to database
	postID, err := d.InsertPost(db, uID, post.Heading, post.Body, time.Now().String()[0:19], "new post image")
	if err != nil {
		logger.WTL(err.Error(), true)
		return err.Error(), false
	}

	// Add category connection to post
	for _, category := range post.Categories {
		catMap := make(map[string]string)
		catMap["category_id"] = category
		categoryArr, err := d.GetCategories(db, catMap)
		if err != nil {
			return err.Error(), false
		}
		if len(categoryArr) < 1 {
			// Category does not exist
			continue
		}
		_, err = d.InsertPostCategory(db, postID, categoryArr[0].Category_id)
		if err != nil {
			return err.Error(), false
		}
	}

	// Prepare return of post id
	postID_str := strconv.Itoa(postID)

	// Store last post insert ts
	err = sessionManager.StoreSessionVariable(w, r, "last_post_insert", nowTS)

	return postID_str, true
}
