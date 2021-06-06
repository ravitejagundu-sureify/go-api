package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main_() {
	//configure the database and chech errors if any
	db, err := sql.Open("mysql", "root:ravaiteja@sureify@(172.17.0.2:3306)/php-app?parseTime=true")
	if err == nil {
		fmt.Println("Database successfully connected")
		fmt.Println(`For basic CRUD operations refer to this site : https://gowebexamples.com/mysql-database/`)
	} else {
		fmt.Println("Connection Unsuccesfull", err)
	}

	if err != nil {
		db.Ping()
	}
	query := `select * from users`

	rows, err := db.Query(query)
	var (
		id       string
		username string
		email    string
		pw       string
	)

	for rows.Next() {
		rows.Scan(&id, &username, &email, &pw)
		fmt.Println(id, username, email, pw)
	}
}
