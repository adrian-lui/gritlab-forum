package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	d "gritface/database"
	logger "gritface/log"
	"net/http"
	"strconv"
	"strings"
)

func addReaction(w http.ResponseWriter, r *http.Request) (string, error) {
	var recID string
	var count int
	retData := make(map[string]interface{})
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		return "", err
	}
	nUID, err := strconv.Atoi(strings.TrimSpace(uid))
	if err != nil {
		return "", errors.New("Uid could not be converted from string to int : " + err.Error())
	}
	if nUID < 1 {
		return "", errors.New("No active session found for adding reaction")
	}

	db, err := d.DbConnect()
	if err != nil {
		return "", err
	}

	defer db.Close()

	// Gather received data
	sComID := r.URL.Query().Get("comment_id")
	comID, err := strconv.Atoi(sComID)
	if err != nil {
		return "", err
	}
	reactID := r.URL.Query().Get("reaction_id")
	sPostID := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(sPostID)
	if err != nil {
		return "", err
	}
	nreactID, err := strconv.Atoi(reactID)
	if err != nil {
		return "", err
	}
	if postID == 0 || nreactID < 1 || nreactID > 2 {
		return "", errors.New("Invalid reaction request from user " + uid + " trying to add reaction id of " + reactID)
	}

	queryStr := "SELECT count(post_id) as rowCount FROM reaction WHERE user_id=? AND post_id = ? AND comment_id = ? AND reaction_id = ?"
	queryChk, err := db.Prepare(queryStr)
	if err != nil {
		return "", err
	}

	defer queryChk.Close()

	var rowCount string
	err = queryChk.QueryRow(nUID, postID, comID, reactID).Scan(&rowCount)

	// Determin if a reaction is made or should be deleted
	if (err == sql.ErrNoRows || err == nil) && rowCount == "0" {
		// If no rows where found
		_, err = d.InsertReaction(db, nUID, postID, comID, reactID)
		if err != nil {
			logger.WTL(err.Error(), true)
			return "", err
		}
	} else if rowCount == "1" {
		// There is exactly a line like this allready, delete it
		_, err = d.DeleteReaction(db, uid, sPostID, sComID, reactID)
		if err != nil {
			return "", err
		}
		reactID = "0"
	} else { // Something went wrong
		return "", err
	}

	// Gather data to be sent to the frontend

	query := "SELECT reaction_id, COUNT (reaction_id) AS rCount FROM reaction WHERE post_id = " + sPostID + " AND comment_id = " + sComID + " GROUP BY reaction_id"
	res, err := db.Query(query)
	if err != nil {
		logger.WTL("Continue with error : "+err.Error(), true)
	}
	defer res.Close()

	retData["rb1"] = 0
	retData["rb2"] = 0
	for res.Next() {
		err = res.Scan(&recID, &count)
		if err != nil {
			logger.WTL("Continue with error : "+err.Error(), true)
		}
		retData["rb"+recID] = count
	}

	retData["status"] = true
	retData["userReaction"] = reactID
	strLine, err := json.Marshal(retData)
	if err != nil {
		return "", err
	}
	return string(strLine), nil
}
