package database

import (
	"database/sql"
	"errors"
	logger "gritface/log"
	"strconv"
)

// Get users from database
func GetUsers(db *sql.DB, userData map[string]string) ([]Users, error) {
	query := "select * from users WHERE"
	count := 0
	for k, v := range userData {
		if k == "password" {
			return nil, errors.New("password is not a valid search parameter")
		}
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'" + " COLLATE NOCASE"
		}
		count++
	}
	var users []Users
	rows, err := db.Query(query)
	if err != nil {
		logger.WTL("Error while trying to query "+query+"\n"+err.Error(), false)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var user Users
		if err := rows.Scan(&user.User_id, &user.Name, &user.Email, &user.Password, &user.Profile_image, &user.Deactive, &user.User_level); // Fetch the record
		err != nil {
			logger.WTL(err.Error(), true)
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Get posts from db
func GetPosts(db *sql.DB, postData map[string]string) ([]Posts, error) {
	query := "select * from posts WHERE"
	count := 0
	for k, v := range postData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var posts []Posts
	rows, err := db.Query(query)
	if err != nil {
		logger.WTL(err.Error(), false)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var post Posts
		if err := rows.Scan(&post.Post_id, &post.User_id, &post.Heading, &post.Body, &post.Closed_user, &post.Closed_admin, &post.Closed_date, &post.Insert_time, &post.Update_time, &post.Image); // Fetch the record
		err != nil {
			logger.WTL(err.Error(), false)
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// Get comments from db
func GetComments(db *sql.DB, commentData map[string]string) ([]Comments, error) {
	query := "select * from comments WHERE"
	count := 0
	for k, v := range commentData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var comments []Comments
	rows, err := db.Query(query)
	if err != nil {
		logger.WTL(err.Error(), false)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var comment Comments
		if err := rows.Scan(&comment.Comment_id, &comment.User_id, &comment.Post_id, &comment.Body, &comment.Insert_time); // Fetch the record
		err != nil {
			logger.WTL(err.Error(), false)
			return comments, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// get categories from db
func GetCategories(db *sql.DB, categoryData map[string]string) ([]Categories, error) {
	query := "select * from categories WHERE"
	count := 0
	for k, v := range categoryData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var categories []Categories

	rows, err := db.Query(query)
	if err != nil {
		logger.WTL(err.Error(), false)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var category Categories
		if err := rows.Scan(&category.Category_id, &category.Category_Name, &category.Closed); // Fetch the record
		err != nil {
			logger.WTL(err.Error(), false)
			return categories, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// get reaction from db
func GetReaction(db *sql.DB, reactionData map[string]string) ([]Reaction, string, error) {
	uid := reactionData["uid"]
	delete(reactionData, "uid")

	userReaction := "0"
	query := "select * from reaction WHERE"
	count := 0
	for k, v := range reactionData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var reactions []Reaction
	rows, err := db.Query(query)
	if err != nil {
		logger.WTL(err.Error(), false)
		return nil, userReaction, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var reaction Reaction
		if err := rows.Scan(&reaction.User_id, &reaction.Post_id, &reaction.Comment_id, &reaction.Reaction_id); // Fetch the record
		err != nil {
			logger.WTL(err.Error(), false)
			return nil, userReaction, err
		}

		if strconv.Itoa(reaction.User_id) == uid {
			userReaction = reaction.Reaction_id
		}

		reactions = append(reactions, reaction)
	}
	return reactions, userReaction, nil
}

// get post categories from db
func GetPostCategories(db *sql.DB, postCategoriesData map[string]string) ([]PostCategory, error) {
	query := "select * from postsCategory WHERE"
	count := 0
	for k, v := range postCategoriesData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var postCategories []PostCategory
	rows, err := db.Query(query)
	if err != nil {
		logger.WTL(err.Error(), false)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var postCategory PostCategory
		if err := rows.Scan(&postCategory.Category_id, &postCategory.Post_id); // Fetch the record
		err != nil {
			logger.WTL(err.Error(), false)
			return postCategories, err
		}
		postCategories = append(postCategories, postCategory)
	}
	return postCategories, nil
}

// get user_level from db

func GetUserLevel(db *sql.DB, userLevelData map[string]string) ([]UserLevel, error) {
	query := "select * from userLevel WHERE"
	count := 0
	for k, v := range userLevelData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var userLevels []UserLevel
	rows, err := db.Query(query)
	if err != nil {
		logger.WTL(err.Error(), false)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var userLevel UserLevel
		if err := rows.Scan(&userLevel.User_level, &userLevel.value); // Fetch the record
		err != nil {
			logger.WTL(err.Error(), false)
			return userLevels, err
		}
		userLevels = append(userLevels, userLevel)
	}
	return userLevels, nil
}
