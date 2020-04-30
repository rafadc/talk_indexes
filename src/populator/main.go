package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jaswdr/faker"
)

// Person represents a record to be generated and inserted
type Person struct {
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Company     string
}

var workers = 20

func main() {
	log.Println("Giving time mySQL to start...")
	time.Sleep(15 * time.Second)

	db, err := sql.Open("mysql", "indexes:indexes@tcp(mysql:3306)/indexes?multiStatements=true")
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxOpenConns(50)
	defer db.Close()

	if tableExists(db, "people_without_indexes") {
		log.Println("Already populated")
		os.Exit(1)
	}

	log.Println("mySQL seems not to be populated yet. Starting...")
	loadDBSchema(db)

	populatePeopleTable(db, "people_small", 1_000)
	populatePeopleTable(db, "people_without_indexes", 1_000_000)
	copyTable(db, "people_without_indexes", "people_single_index")
	copyTable(db, "people_without_indexes", "people_multi_column_index")
	copyTable(db, "people_without_indexes", "people_range_query")
}

func tableExists(db *sql.DB, tableName string) bool {
	log.Println("Importing schema")

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

	_, err = db.Exec(string(content))
	if err != nil {
		panic(err.Error())
	}
}

func populatePeopleTable(db *sql.DB, tableName string, numberOfPeople int) {
	log.Println("Populating table ", tableName)

	insertQueryText := fmt.Sprintf("INSERT INTO %s(name,surname,date_of_birth,company) VALUES(?,?,?,?)", tableName)
	insertQuery, err := db.Prepare(insertQueryText)
	defer insertQuery.Close()

	var wg sync.WaitGroup

	if numberOfPeople%workers != 0 {
		panic("The number of people must be divisible by the number of workers")
	}

	wg.Add(numberOfPeople)

	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < numberOfPeople/workers; j++ {
				personToInsert := randomPerson()
				insertQuery.Exec(
					personToInsert.FirstName,
					personToInsert.LastName,
					personToInsert.DateOfBirth,
					personToInsert.Company,
				)
				wg.Done()
			}
		}()
	}
	wg.Wait()

	if err != nil {
		panic(err.Error())
	}
}

func randomPerson() Person {
	var fakeGenerator = faker.New()

	return Person{
		FirstName:   fakeGenerator.Person().FirstName(),
		LastName:    fakeGenerator.Person().LastName(),
		DateOfBirth: fakeGenerator.Time().Time(time.Now()),
		Company:     fakeGenerator.Company().Name(),
	}
}

func copyTable(db *sql.DB, sourceTable string, targetTable string) {
	log.Println("Copying table %s into %s", sourceTable, targetTable)

	createTableQuery := fmt.Sprintf("CREATE TABLE %s LIKE %s", targetTable, sourceTable)
	query, err := db.Query(createTableQuery)

	if err != nil {
		panic(err.Error())
	}
	defer query.Close()

	copyDataQuery := fmt.Sprintf("INSERT %s SELECT * FROM %s", targetTable, sourceTable)
	query, err = db.Query(copyDataQuery)

	if err != nil {
		panic(err.Error())
	}
	defer query.Close()
}
