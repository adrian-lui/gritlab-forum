package database

import (
	"database/sql"
	logger "gritface/log"
	"math/rand"
	"strconv"
)

// returns the row affected and error
// insertUsers function inserts a record in the users table
func InsertUsers(db *sql.DB, name string, email string, password string, user_level int) (int, error) {
	insertUsers := `INSERT INTO users(name, email, Password, user_level, profile_image) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertUsers) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		return 0, err
	}
	min := 1
	max := 7
	randomNum := rand.Intn(max-min) + min
	pic := "static/images/raccoon_thumbnail" + strconv.Itoa(randomNum) + ".jpg"
	val, err := statement.Exec(name, email, password, user_level, pic) // Execute statement with parameters
	if err != nil {
		return 0, err
	}
	insertId, _ := val.LastInsertId()
	return int(insertId), nil
}

// returns the row affected and error
// function adds posts to the database
func InsertPost(db *sql.DB, user_id int, heading string, body string, insert_time string, image string) (int, error) {
	insertPost := `INSERT INTO posts(user_id, heading, body, insert_time, image) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertPost) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		return 0, err
	}
	val, err := statement.Exec(user_id, heading, body, insert_time, image) // Execute statement with parameters
	if err != nil {
		return 0, err
	}
	insertId, _ := val.LastInsertId()
	return int(insertId), nil
}

// returns the row affected and error
// function inserts categories into the database
func InsertCategory(db *sql.DB, category_name string) (int, error) {
	insertCategory := `INSERT INTO categories(category_name) VALUES (?)`
	statement, err := db.Prepare(insertCategory) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		return 0, err
	}
	val, err := statement.Exec(category_name) // Execute statement with parameters
	if err != nil {
		return 0, err
	}
	insertId, _ := val.LastInsertId()
	return int(insertId), nil
}

// returns the row affected and error
// function inserts comments into the database
func InsertComment(db *sql.DB, post_id int, user_id int, body string, insert_time string) (int, error) {
	insertComment := `INSERT INTO comments(post_id, user_id, body, insert_time) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertComment) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		return 0, err
	}
	val, err := statement.Exec(post_id, user_id, body, insert_time) // Execute statement with parameters
	if err != nil {
		return 0, err
	}
	insertId, _ := val.LastInsertId()
	return int(insertId), nil
}

// function inserts reaction into the database
// returns the row affected and error
func InsertReaction(db *sql.DB, user_id int, post_id int, comment_id int, reaction_id string) (int, error) {
	insertReaction := `INSERT INTO reaction(user_id, post_id, comment_id, reaction_id) VALUES (?, ?, ?, ?)
	ON CONFLICT(user_id, post_id, comment_id) DO UPDATE SET reaction_id=` + reaction_id

	statement, err := db.Prepare(insertReaction) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		return 0, err
	}
	val, err := statement.Exec(user_id, post_id, comment_id, reaction_id) // Execute statement with parameters
	if err != nil {
		return 0, err
	}

	insertId, _ := val.LastInsertId()
	return int(insertId), nil
}

// function inserts post category into the database
func InsertPostCategory(db *sql.DB, post_id int, category_id int) (int, error) {
	insertPostCategory := `INSERT INTO postscategory(post_id, category_id) VALUES (?, ?)`
	statement, err := db.Prepare(insertPostCategory) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		return 0, err
	}
	val, err := statement.Exec(post_id, category_id) // Execute statement with parameters
	if err != nil {
		return 0, err
	}
	insertId, _ := val.LastInsertId()
	return int(insertId), nil
}

func DeleteReaction(db *sql.DB, user_id string, post_id string, comment_id string, reaction_id string) (bool, error) {
	deleteReaction := "DELETE FROM reaction WHERE user_id =" + user_id + " AND post_id=" + post_id + " AND comment_id=" + comment_id
	stmt, err := db.Prepare(deleteReaction) // Execute statement with parameters
	if err != nil {
		return false, err
	}
	_, err = stmt.Exec()
	if err != nil {
		logger.WTL(err.Error(), true)
		return false, err
	}
	return true, nil
}
