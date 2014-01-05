package main

import (
	"./bbbikeng"
	"fmt"
	"github.com/ant0ine/go-json-rest"
	//"log"
	//"net/http"
	//"strings"
)

func main() {

	bbbikeng.ConnectToDatabase()
	defer bbbikeng.Connection.Close()

	bbbikeng.Test()

	/*
	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{"GET", "/search/:name", Search},
	)


	http.ListenAndServe(":8080", &handler) */

}

func Search(w *rest.ResponseWriter, req *rest.Request) {

	query := req.PathParam("name")
	if len(query) > 0 {
		results := bbbikeng.SearchForStreetName(query)
		fmt.Println("Results:", results)
		w.WriteJson(&results)
	}
}

func Route(w *rest.ResponseWriter, req *rest.Request) {

	parameters := req.PathParams
	fmt.Println(parameters)
//	start := parameters["start"]
//	end := parameters["end"]

//	w.WriteJson(&end)

}
