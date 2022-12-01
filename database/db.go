package database

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	logger "gritface/log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// hash password returned the password string as a hash to be stored in the database
// this is doen for security reasons
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// check passwords checks if the password provided by the user matches the one in the database
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Connect to database
func DbConnect() (*sql.DB, error) {
	databaseFile := "forum-db.db"
	forumdb, err := sql.Open("sqlite3", "./"+databaseFile+"?_auth&_auth_user=forum&_auth_pass=forum&_auth_crypt=sha1")

	if err != nil {
		return nil, err
	}

	// Enable foreign key contraints
	enableContraints := `PRAGMA foreign_keys = ON;`

	enableContraintsQuery, err := forumdb.Prepare(enableContraints) // Prepare SQL Statement
	if err != nil {
		logger.WTL(err.Error(), true)
	}
	enableContraintsQuery.Exec()

	return forumdb, nil
}

// function check if database exists, if not it creates it, if it does it opens it
func DatabaseExist() (*sql.DB, error) {
	newDb := false
	databaseFile := "forum-db.db"
	_, err := os.Stat(databaseFile)
	if os.IsNotExist(err) {
		logger.WTL("Creating the forum database ...", true)
		file, err := os.Create(databaseFile) // Create Sqlite file
		if err != nil {
			return nil, err
		}
		file.Close()

		logger.WTL("Database created", true)
		newDb = true
	} else if err != nil {
		return nil, err
	}
	forumdb, err := sql.Open("sqlite3", "./"+databaseFile)
	// Open the created Sqlite3 File
	if err != nil {
		logger.WTL("Database could not be opened", false)
		return nil, err
	}
	conn, err := forumdb.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if newDb {
		err = createTable(forumdb) // Create Database Tables
		if err != nil {
			logger.WTL(err.Error(), true)
		} else {
			// INSERT RECORDS
			exampleDbData(forumdb)
		}
	}

	var requiredTables = map[string]bool{"users": false, "posts": false, "comments": false, "categories": false}
	tables, table_check := forumdb.Query("select name from sqlite_master where type='table' and name not like 'sqlite_%'")

	if table_check == nil {
		for tables.Next() { // Iterate and fetch the records
			var name string
			tables.Scan(&name)
			requiredTables[name] = true // Fetch the record
		}
		for _, value := range requiredTables {
			if !value {
				forumdb.Close() // Close connection to current database
				reader := bufio.NewReader(os.Stdin)
			handleInvalidDatabase:
				logger.WTL("Existing database is not working with this version of forum.\nWould you like to delete current database (y) or rename current database (r), quit (q)?:", true)
				uinput, _ := reader.ReadString('\n')
				uinput = strings.Trim(uinput, "\n")

				if uinput == "y" {
					os.Remove(databaseFile)
				} else if uinput == "r" {
				renameCurrentDB:
					fmt.Print("Rename to: ")
					uinput, _ = reader.ReadString('\n')
					uinput = strings.Trim(uinput, "\n")
					if len(uinput) < 1 {
						fmt.Println("Name can not be empty")
						goto renameCurrentDB
					} else if len(uinput) < 4 { // Add .db to the end as a string 3 char long can not hold name + .db
						uinput += ".db"
					} else if uinput[len(uinput)-4:] != ".db" {
						uinput += ".db"
					}
					os.Rename(databaseFile, uinput)
					fmt.Println("Database renamed to " + uinput)
				} else if uinput == "q" {
					fmt.Println("Exiting....")
					os.Exit(0)
				} else {
					goto handleInvalidDatabase
				}
				return DatabaseExist()
			}
		}
	} else {
		logger.WTL("table not there", true)
	}
	return forumdb, nil
}

// remove this when cleaning up
func exampleDbData(forumdb *sql.DB) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte("123"), 12)
	InsertUsers(forumdb, "admin", "admin@gritface.ax", string(bytes), 1)
	logger.WTL("Admin created as 'admin@gritface.ax', password 123", true)
	InsertPost(forumdb, 1, "Welcome to gritface", "Feel free to discuss anything regarding gritlab or anything else you can come up with. Go wild!", "2022-12-01 00:00:00", "")
	InsertCategory(forumdb, "sport")
	InsertCategory(forumdb, "food")
	InsertCategory(forumdb, "hiking")
	InsertCategory(forumdb, "data science")
	InsertCategory(forumdb, "programming")
	InsertCategory(forumdb, "music")
	InsertCategory(forumdb, "movies")
	InsertCategory(forumdb, "books")
}
