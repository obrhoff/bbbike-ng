/**
 * User: DocterD
 * Date: 28/12/13
 * Time: 11:19
 */

package util

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const user = "root"
const password = "root"
const host = "127.0.0.1"
const port = "5433"
const database = "bbbikeng"

func ConnectToDatabase() (db *sql.DB) {

	connectionParameter := fmt.Sprint("user=", user, " password=", password, " host=", host, " port=", port, " dbname=", database)
	fmt.Println("Connecting to Database:", host)

	database, err := sql.Open("postgres", connectionParameter)
	err = database.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible

	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}

	return database
}
