package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jaswdr/faker"
)

func main() {
	fmt.Println("Giving time mySQL to start...")
	time.Sleep(10 * time.Second)

	db, err := sql.Open("mysql", "indexes:indexes@tcp(mysql:3306)/indexes")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	if tableExists(db, "people_without_indexes") {
		fmt.Println("Already populated")
		os.Exit(1)
	}

	fmt.Println("mySQL seems not to be populated yet. Starting...")
	loadDBSchema(db)

	populatePeopleTable(db, "people_without_indexes")
}

func tableExists(db *sql.DB, tableName string) bool {
	fmt.Println("Importing schema")

	query := `SELECT count(*)
	FROM information_schema.tables
	WHERE table_schema = 'indexes'
		AND table_name = ?
	LIMIT 1;`
	var result int8
	err := db.QueryRow(query, tableName).Scan(&result)
	if err != nil {
		panic(err.Error())
	}
	return result != 0
}

func loadDBSchema(db *sql.DB) {
	content, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Query(string(content))
	if err != nil {
		panic(err.Error())
	}
}

func populatePeopleTable(db *sql.DB, tableName string) {
	fmt.Println("Populating table %s", tableName)

	faker := faker.New()

	insertQuery := fmt.Sprintf("INSERT INTO %s(name,surname,date_of_birth,company) VALUES(?,?,?,?)", tableName)
	_, err := db.Query(insertQuery,
		faker.Person().Name(),
		faker.Person().LastName(),
		faker.Time(),
		faker.Company().Name())

	if err != nil {
		panic(err.Error())
	}
}
