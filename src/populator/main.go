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
	Address     string
}

var workers = 20

func main() {
	log.Println("Giving time mySQL to start...")
	time.Sleep(15 * time.Second)

	db, err := sql.Open("mysql", "indexes:indexes@tcp(mysql:3306)/indexes?multiStatements=true")
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxOpenConns(40)
	defer db.Close()

	if tableExists(db, "people_without_indexes") {
		log.Println("Already populated")
		os.Exit(1)
	}

	log.Println("mySQL seems not to be populated yet. Starting...")
	loadDBSchema(db)

	populatePeopleTable(db, "people_small", 1_000)
	populatePeopleTable(db, "people_without_indexes", 10_000_000)
	copyTable(db, "people_without_indexes", "people_single_index")
	copyTable(db, "people_without_indexes", "people_multi_column_index")
	copyTable(db, "people_without_indexes", "people_range_query")

	createIndexes(db)
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

func executeSQLFile(db *sql.DB, filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		panic(err.Error())
	}
}

func loadDBSchema(db *sql.DB) {
	executeSQLFile(db, "schema.sql")
}

func createIndexes(db *sql.DB) {
	executeSQLFile(db, "indexes.sql")
}

func populatePeopleTable(db *sql.DB, tableName string, numberOfPeople int) {
	log.Println("Populating table ", tableName)

	var wg sync.WaitGroup

	if numberOfPeople%workers != 0 {
		panic("The number of people must be divisible by the number of workers")
	}

	wg.Add(numberOfPeople)

	for i := 0; i < workers; i++ {
		go func() {
			var fakeGenerator = faker.New()
			insertQueryText := fmt.Sprintf("INSERT INTO %s(name,surname,date_of_birth,company, address) VALUES(?,?,?,?,?)", tableName)
			insertQuery, err := db.Prepare(insertQueryText)

			if err != nil {
				panic(err.Error())
			}
			defer insertQuery.Close()

			for j := 0; j < numberOfPeople/workers; j++ {
				personToInsert := randomPerson(fakeGenerator)
				_, err := insertQuery.Exec(
					personToInsert.FirstName,
					personToInsert.LastName,
					personToInsert.DateOfBirth,
					personToInsert.Company,
					personToInsert.Address,
				)
				if err != nil {
					panic(err.Error())
				}
				wg.Done()
			}
		}()
	}
	wg.Wait()
}

func randomPerson(fakeGenerator faker.Faker) Person {
	return Person{
		FirstName:   fakeGenerator.Person().FirstName(),
		LastName:    fakeGenerator.Person().LastName(),
		DateOfBirth: fakeGenerator.Time().Time(time.Now()),
		Company:     fakeGenerator.Company().Name(),
		Address:     fakeGenerator.Address().Address(),
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
