package server

import (
	"encoding/json"
	"math/rand"
	"strconv"
)

// DummyPost returns a JSONData object with the provided body
func DummyPost(body string) string {
	min := 1
	max := 7
	randomNum := rand.Intn(max-min) + min
	pic := "static/images/raccoon_thumbnail" + strconv.Itoa(randomNum) + ".jpg"
	data := &JSONData{
		Post_id:       0,
		User_id:       0,
		Heading:       "ERROR",
		Body:          body,
		Comments:      make(map[int]JSONComments),
		Profile_image: pic,
	}
	dataSlice := []JSONData{*data}
	dummy, _ := json.Marshal(dataSlice)
	return string(dummy)
}
