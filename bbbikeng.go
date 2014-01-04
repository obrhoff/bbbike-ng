package main

import (
	"./bbbikeng"
	"database/sql"
	"fmt"
	"github.com/ant0ine/go-json-rest"
	//"html"
	//"log"
	//"net/http"
	//"strings"
)

func main() {

	bbbikeng.ConnectToDatabase()
	defer bbbike.db.Close()

	testPlacePoint1 := bbbike.MakeNewPoint("52.551080", "13.373370")
	testPlacePoint2 := bbbike.MakeNewPoint("52.492491", "13.428981")
	bbbike.CalculateRoute(testPlacePoint1, testPlacePoint2, db)

	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{"GET", "/search/:name", Search},
	)
	http.ListenAndServe(":8080", &handler)

}

func Search(w *rest.ResponseWriter, req *rest.Request) {

	query := req.PathParam("name")
	if len(query) > 0 {
		results := bbbike.SearchForStreetName(query, db)
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
