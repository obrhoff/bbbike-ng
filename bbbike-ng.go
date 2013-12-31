package main

import (
	//"./bbbikeng/helper"
	//"./bbbikeng/model"
	"./misc"
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"
	//"strings"
)

var db *sql.DB

func main() {

	db = util.ConnectToDatabase()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello BBBike :-), %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
