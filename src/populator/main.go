package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jaswdr/faker"
)

const maxRetries = 50

// Person represents a record to be generated and inserted
type Person struct {
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Company     string
	Address     string
	Happy       bool
}

func main() {
	db := connectToDB()
	defer db.Close()

	if tableExists(db, "people_without_indexes") {
		log.Println("Already populated")
		os.Exit(1)
	}

	log.Println("mySQL seems not to be populated yet. Starting...")
	loadDBSchema(db)

	numberOfRecords, err := strconv.Atoi(os.Getenv("NUMBER_OF_RECORDS"))
	if err != nil {
		panic(err.Error())
	}

	populatePeopleTable(db, "people_small", 1_000)
	populatePeopleTable(db, "people_without_indexes", numberOfRecords)
	copyTable(db, "people_without_indexes", "people_single_index")
	copyTable(db, "people_without_indexes", "people_multi_column_index")
	copyTable(db, "people_without_indexes", "people_range_query")
	copyTable(db, "people_without_indexes", "people_full_text_search")

	createIndexes(db)

	log.Println("Done!")
}

func connectToDB() *sql.DB {
	log.Println("Trying to connect to mySQL...")

	db, err := sql.Open("mysql", "indexes:indexes@tcp(mysql:3306)/indexes?multiStatements=true")
	if err != nil {
		panic("Incorrect URL in mySQL")
	}

	for retry := 0; retry < maxRetries; retry++ {
		err = db.Ping()

		if err != nil {
			time.Sleep(1 * time.Second)
			log.Println("retrying...")
		} else {
			db.SetMaxOpenConns(40)
			db.SetMaxIdleConns(workers())
			return db
		}
	}

	panic("Couldn't get a connection to database")
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
	log.Printf("Populating table %s with %d records", tableName, numberOfPeople)

	var wg sync.WaitGroup

	if numberOfPeople%workers() != 0 {
		panic("The number of people must be divisible by the number of workers")
	}

	wg.Add(numberOfPeople)

	for i := 0; i < workers(); i++ {
		go func() {
			var fakeGenerator = faker.New()
			insertQueryText := fmt.Sprintf("INSERT INTO %s(name,surname,date_of_birth,company, address, happy) VALUES(?,?,?,?,?,?)", tableName)
			insertQuery, err := db.Prepare(insertQueryText)

			if err != nil {
				panic(err.Error())
			}
			defer insertQuery.Close()

			for j := 0; j < numberOfPeople/workers(); j++ {
				personToInsert := randomPerson(fakeGenerator)
				_, err := insertQuery.Exec(
					personToInsert.FirstName,
					personToInsert.LastName,
					personToInsert.DateOfBirth,
					personToInsert.Company,
					personToInsert.Address,
					personToInsert.Happy,
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
		Happy:       rand.Intn(10_000) > 0,
	}
}

func copyTable(db *sql.DB, sourceTable string, targetTable string) {
	log.Printf("Copying table %s into %s\n", sourceTable, targetTable)

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

func workers() int {
	workers, err := strconv.Atoi(os.Getenv("WORKERS"))
	if err != nil {
		panic(err.Error())
	}

	return workers
}
