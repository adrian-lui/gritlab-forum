package database

import "database/sql"

func createTable(db *sql.DB) error {
	// users table
	createUsersTable := `CREATE TABLE users (
		"user_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,		
		"name" TEXT NOT NULL UNIQUE,
		"email" TEXT NOT NULL UNIQUE,
		"password" TEXT NOT NULL,
		"profile_image" TEXT DEFAULT '',
		"deactive" INTEGER DEFAULT 0,
		"user_level" TEXT DEFAULT "user"
	  );`

	usersStatement, err := db.Prepare(createUsersTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	usersStatement.Exec() // Execute SQL Statements

	// posts tables
	createPostsTable := `CREATE TABLE posts (
		"post_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
		"user_id" INTEGER NOT NULL,		
		"heading" TEXT NOT NULL,
		"body" TEXT NOT NULL,
		"closed_user" INTEGER default '0',
		"closed_admin" INTEGER default '0',
		"closed_date" TEXT DEFAULT '',
		"insert_time" TEXT NOT NULL,
		"update_time" TEXT NOT NULL DEFAULT '', 
		"image" TEXT
	  );`
	postsStatement, err := db.Prepare(createPostsTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	postsStatement.Exec() // Execute SQL Statements

	// comments table
	createcommentsTable := `CREATE TABLE comments (
		"comment_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
		"post_id" INTEGER NOT NULL,
		"user_id" INTEGER NOT NULL,
		"body" TEXT NOT NULL,
		"insert_time" TEXT NOT NULL,
		FOREIGN KEY("post_id") REFERENCES posts("post_id")
	  );`

	commentsStatement, err := db.Prepare(createcommentsTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	commentsStatement.Exec() // Execute SQL Statements

	// categories table
	createcategoriesTable := `CREATE TABLE categories (
		"category_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
		"category_name" TEXT NOT NULL UNIQUE,
		"closed" INTEGER default 0
	  );`
	categoriesStatement, err := db.Prepare(createcategoriesTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	categoriesStatement.Exec() // Execute SQL Statements

	// reaction table
	createreactionTable := `CREATE TABLE reaction (
		"user_id" INTEGER NOT NULL,
		"post_id" INTEGER NOT NULL,
		"comment_id" INTEGER NOT NULL,
		"reaction_id" TEXT NOT NULL,
		PRIMARY KEY (user_id, post_id, comment_id),
		FOREIGN KEY("post_id") REFERENCES posts("post_id")
	  );`
	reactionStatement, err := db.Prepare(createreactionTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	reactionStatement.Exec() // Execute SQL Statements

	//post category table
	createpostscategoryTable := `CREATE TABLE postsCategory (
		"category_id" INTEGER NOT NULL,
		"post_id" INTEGER NOT NULL,
		FOREIGN KEY("post_id") REFERENCES posts("post_id"),
		FOREIGN KEY("category_id") REFERENCES categories("category_id")
	  );`
	postcategoryStatement, err := db.Prepare(createpostscategoryTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	postcategoryStatement.Exec() // Execute SQL Statements

	createUserLevelTable := `CREATE TABLE userLevel (
		"user_level" TEXT NOT NULL UNIQUE,
		"value" INTEGER NOT NULL,
		PRIMARY KEY (user_level, value)
	)`
	userlevelStatement, err := db.Prepare(createUserLevelTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	userlevelStatement.Exec() // Execute SQL Statements
	return nil
}
