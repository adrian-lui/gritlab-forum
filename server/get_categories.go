package server

import (
	"encoding/json"
	d "gritface/database"
	"net/http"
)

func GetCategories(w http.ResponseWriter, r *http.Request) (string, bool) {
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

	categories := make(map[int]string)
	rows, err := db.Query("select * from categories")
	if err != nil {
		return err.Error(), false
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var category d.Categories
		if err := rows.Scan(&category.Category_id, &category.Category_Name, &category.Closed); // Fetch the record
		err != nil {
			return err.Error(), false
		}

		categories[category.Category_id] = category.Category_Name
	}

	res, err := json.Marshal(categories)
	if err != nil {
		return err.Error(), false
	}

	return string(res), true
}
