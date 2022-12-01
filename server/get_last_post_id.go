package server

import (
	"encoding/json"
	logger "gritface/log"
	"io"
	"net/http"
)

func GetLastPostID(w http.ResponseWriter, r *http.Request) int {
	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.WTL(err.Error(), true)
		return 0
	}

	// Unmarshal
	var lastPost LastPost
	err = json.Unmarshal(req, &lastPost)
	if err != nil {
		logger.WTL(err.Error(), false)
		return 0
	}

	return lastPost.ID
}
