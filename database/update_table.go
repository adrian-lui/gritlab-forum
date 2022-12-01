package database

import (
	"database/sql"
	"errors"
	"time"
)

// UpdateCategoriesData updates the category with the correct category_id with data specified
// in a map[string]string. If category_id doesn't exist or there is something wrong with the statement,
// the function returns an error.
func UpdateCategoriesData(db *sql.DB, data map[string]string, category_id string) error {
	// Checking if category_id exists
	search := "SELECT * FROM categories WHERE category_id=?"
	err := db.QueryRow(search, category_id).Scan()
	if err == sql.ErrNoRows {
		return err
	}

	query := "UPDATE categories SET"
	counter := 0
	for key, value := range data {
		if key == "category_id" {
			return errors.New("ERROR: change of ID not allowed")
		}
		if counter > 0 {
			query += ","
		}
		query += " " + key + "='" + value + "'"
		counter++
	}
	query += " WHERE category_id=" + category_id
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	return err
}

// UpdateCommentsData updates the comment body of the comment specified by the comment_id.
// If the comment_id row doesn't exist or the body cannot be updated for any reason, it returns an error.
func UpdateCommentsData(db *sql.DB, body string, comment_id string) error {
	// Checking if comment_id exists
	search := "SELECT * FROM comments WHERE comment_id=?"
	err := db.QueryRow(search, comment_id).Scan()
	if err == sql.ErrNoRows {
		return err
	}

	query := "UPDATE comments SET body=? WHERE comment_id=?"
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.Exec(body, comment_id)
	return err
}

// UpdatePostsData receives the database, a map[string]string and a post_id and updates post_id's columns
// specified in the map. If the update is not possible for any reason, the function returns an error.
// It also updates the update_time column based on the current time.
func UpdatePostsData(db *sql.DB, data map[string]string, post_id string) error {
	// Checking if post_id exists
	search := "SELECT * FROM posts WHERE post_id=?"
	err := db.QueryRow(search, post_id).Scan()
	if err == sql.ErrNoRows {
		return err
	}

	time := time.Now().Format("2006-01-02 15:04:05")
	query := "UPDATE posts SET"
	count := 0
	// Nobody should be able to directly change IDs and times
	notAllowed := map[string]bool{
		"post_id":     true,
		"user_id":     true,
		"insert_time": true,
		"update_time": true,
	}
	for key, val := range data {
		if notAllowed[key] {
			return errors.New("ERROR: update of IDs and times not allowed")
		}
		if count > 0 {
			query += ","
		}
		query += " " + key + "='" + val + "'"
		count++
	}
	query += ", update_time='" + time + "' WHERE post_id=" + post_id
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	return err
}

// UpdateUserData receives the database, a map[string]string and a user_id and updates user_id's columns
// specified in the map. If the update is not possible for any reason, the function returns an error.
func UpdateUserData(db *sql.DB, data map[string]string, user_id string) error {
	// Checking if post_id exists
	search := "SELECT * FROM users WHERE user_id=?"
	err := db.QueryRow(search, user_id).Scan()
	if err == sql.ErrNoRows {
		return err
	}

	query := "UPDATE users SET"
	count := 0
	for key, val := range data {
		if key == "user_id" {
			return errors.New("ERROR: user_id update not allowed")
		}
		if count > 0 {
			query += ","
		}
		query += " " + key + "='" + val + "'"
		count++
	}
	query += " WHERE user_id=" + user_id
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	return err
}
