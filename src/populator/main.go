package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Populating mySQL")

	db, err := sql.Open("mysql", "indexes:indexes@tcp(mysql:3306)/indexes")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	if tableExists(db, "people") {
		fmt.Println("Already populated")
		os.Exit(1)
	}
}

func tableExists(db *sql.DB, tableName string) bool {
	query := `SELECT *
	FROM information_schema.tables
	WHERE table_schema = 'indexes'
		AND table_name = ?
	LIMIT 1;`
	var result string
	err := db.QueryRow(query, tableName).Scan(&result)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("The result is: %s", result)
}
