package main

import (
	"./bbbikeng"
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"
)

var db *sql.DB

func main() {

	db = bbbike.ConnectToDatabase()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello BBBike :-), %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		getAllRoutes()
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func getAllRoutes() {

	fmt.Println("Getting all Routes...")

}
