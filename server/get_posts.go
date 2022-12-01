package server

import (
	"encoding/json"
	d "gritface/database"
	logger "gritface/log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// getPosts function connects to the database and based on the parameters provided or the r.URL.Query values, returns the result of SQLite queries as a JSON string. It also returns an error if there is any.
func getPosts(r *http.Request, uid string, last_post_id int) (string, error) {

	db, err := d.DbConnect()

	if err != nil {
		return "", err
	}

	defer db.Close()
	// base query
	query := "SELECT posts.* from posts"
	limit := true
	if len(r.URL.Query()) > 0 {
		limit = false
		filterValue := r.URL.Query()["filter"]
		useFilter := ""
		if len(filterValue) != 0 {
			useFilter = filterValue[0]
		}
		switch useFilter {
		case "postId":
			// To get one post
			query += " WHERE posts.post_id=" + r.URL.Query()["id"][0]
		case "userPosts":
			// if the user wants to see their own posts
			query += " WHERE posts.user_id=" + uid
		case "liked":
			// if the user wants to see posts liked by them
			query += " INNER JOIN reaction ON posts.post_id=reaction.post_id WHERE reaction.reaction_id='1' AND reaction.user_id=" + uid + " AND reaction.comment_id=0 GROUP BY posts.post_id"
		case "hated":
			// if the user wants to see posts disliked by them
			query += " INNER JOIN reaction ON posts.post_id=reaction.post_id WHERE reaction.comment_id=0 AND reaction.reaction_id='2' AND reaction.user_id=" + uid + " GROUP BY posts.post_id"
		case "category":
			// if the user wants to see posts in a certain category
			cat := r.URL.Query()["cat"][0]
			if err != nil {
				// if there is an error, we return a dummy post that informs the user of an error, but doesn't break the site
				return DummyPost("Invalid category"), nil
			}

			// check if the category_id exists, if not, we return the dummy post
			catCheck := "SELECT * FROM categories WHERE category_id=" + cat
			checkRows, err := db.Query(catCheck)
			if err != nil {
				logger.WTL("Error qith query "+catCheck, false)
			}
			if d.QueryRows(checkRows) < 1 {
				return DummyPost("Invalid category"), nil
			}
			query += " INNER JOIN postsCategory ON posts.post_id=postsCategory.post_id WHERE postsCategory.category_id=" + cat
		default:
			return DummyPost("Invalid filter"), nil
		}
	}
	if last_post_id != 0 {
		// if the user wants to load more posts, we need the last post's post_id
		query += " WHERE post_id < " + strconv.Itoa(last_post_id)
	}
	query += " ORDER BY posts.post_id DESC"
	if limit {
		// if there is no filtering, there is a limit of 20 posts loaded at once
		query += " LIMIT 20"
	}

	structSlice := make(map[int]JSONData)
	rows, err := db.Query(query)
	if err != nil {
		logger.WTL("Error with query "+query, false)
		return "", err
	}

	defer rows.Close()
	nextQuery := ""
	// after the posts' query, the function declares a JSONData object for each of them and assigns the data to the corresponding field
	for rows.Next() {
		rD := &JSONData{
			Comments: make(map[int]JSONComments),
		}
		err = rows.Scan(&rD.Post_id, &rD.User_id, &rD.Heading, &rD.Body, &rD.Closed_user, &rD.Closed_admin, &rD.Closed_date, &rD.Insert_time, &rD.Update_time, &rD.Image)
		if err != nil {
			return "", err
		}

		rD.Body = strings.Replace(rD.Body, "\n", "<br>", -1)

		if rD.User_id < 1 {
			logger.WTL("Post (id = "+strconv.Itoa(rD.Post_id)+") found with user id 0", true)
			continue
		}

		postId := &rD.Post_id

		// getting user's name
		currentUser := make(map[string]string)
		currentUser["user_id"] = strconv.Itoa(rD.User_id)
		users, err := d.GetUsers(db, currentUser)
		if err != nil {
			return "", err
		}
		rD.Username = users[0].Name
		rD.Profile_image = users[0].Profile_image

		// getting post's categories
		currentPost := make(map[string]string)
		currentPost["post_id"] = strconv.Itoa(*postId)
		categories, err := d.GetPostCategories(db, currentPost)
		if err != nil {
			return "", err
		}
		var categoryNames []string
		for _, category := range categories {
			currentCategory := make(map[string]string)
			currentCategory["category_id"] = strconv.Itoa(category.Category_id)
			categoriesName, err := d.GetCategories(db, currentCategory)
			if err != nil {
				return "", err
			}
			categoryNames = append(categoryNames, categoriesName[0].Category_Name)
		}
		rD.Categories = categoryNames

		// getting post's reactions
		currentPost["comment_id"] = "0"
		currentPost["uid"] = uid
		reactions, myReaction, err := d.GetReaction(db, currentPost)
		if err != nil {
			return "", err
		}
		for _, reaction := range reactions {
			if reaction.Comment_id == 0 {
				userReaction := make(map[int]string)
				userReaction[reaction.User_id] = reaction.Reaction_id
				rD.Reactions = append(rD.Reactions, userReaction)
			}
		}
		rD.UserReaction = myReaction

		structSlice[*postId] = *rD

		// thisPostId := &rD.Post_id
		nextQuery += " OR post_id=" + strconv.Itoa(*postId)
	}
	// after the posts' query, we need to query for the comments, if there are any
	if len(nextQuery) > 4 {
		query = "SELECT comment_id, post_id, user_id, body, insert_time FROM comments WHERE " + nextQuery[4:]
		rows, err = db.Query(query)
		if err != nil {
			logger.WTL("Error for query "+query, false)
			return "", err
		}
		defer rows.Close()
		for rows.Next() {
			row := &JSONComments{}
			err = rows.Scan(&row.CommentID, &row.Post_id, &row.User_id, &row.Body, &row.Insert_time)
			if err != nil {
				return "", err
			}
			currentUser := make(map[string]string)
			currentUser["user_id"] = strconv.Itoa(row.User_id)
			users, err := d.GetUsers(db, currentUser)
			if err != nil {
				return "", err
			}
			row.Username = users[0].Name
			row.Profile_image = users[0].Profile_image

			thisPostId := &row.Post_id
			thisCommentId := &row.CommentID
			// getting reactions
			currentComment := make(map[string]string)
			currentComment["comment_id"] = strconv.Itoa(*thisCommentId)
			currentComment["uid"] = uid
			reactions, userReaction, err := d.GetReaction(db, currentComment)
			if err != nil {
				return "", err
			}
			for _, reaction := range reactions {
				userReaction := make(map[int]string)
				userReaction[reaction.User_id] = reaction.Reaction_id
				row.Reactions = append(row.Reactions, userReaction)
			}
			row.UserReaction = userReaction

			structSlice[*thisPostId].Comments[row.CommentID] = *row
		}
	}

	// The output needs to be in a descending order (by post_id), so we save it into a sorted []JSONData, because marshalling maps into JSON data always ends up in ascending order by the keys
	sSlice := make([]JSONData, 0, len(structSlice))
	for _, value := range structSlice {
		sSlice = append(sSlice, value)
	}
	sort.Slice(sSlice, func(i, j int) bool { return sSlice[i].Post_id > sSlice[j].Post_id })

	res, err := json.Marshal(sSlice)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
