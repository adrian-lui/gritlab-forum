package database

import "database/sql"

func QueryRows(rows *sql.Rows) int {
	counter := 0
	for rows.Next() {
		counter++
	}
	return counter
}
