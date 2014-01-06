package main

import (
	"./bbbikeng"
	"./Import"
	"fmt"
	"flag"
	"os"
	"github.com/ant0ine/go-json-rest"
	//"log"
	"net/http"
	//"strings"
)

var startFlag = flag.Bool("run",true, "start bbbikeng")
var dataImportFlag = flag.Bool("import-data", false, "import bbbike path")
var dataImportPathFlag = flag.String("import-path", "./bbbike/data", "bbbike data path")

func main() {

	flag.Parse()

	if *startFlag && !*dataImportFlag {
		StartBBBikeServer()
	} else if *dataImportFlag{
		StartParsingBBBikeData(*dataImportPathFlag)
	} else {
		fmt.Printf("--import-data --import-path=/bbbike/data\n")
		os.Exit(1)
	}

}

func StartBBBikeServer() {

	bbbikeng.ConnectToDatabase()
	defer bbbikeng.Connection.Close()

	handler := rest.ResourceHandler{}
		handler.SetRoutes(
		rest.Route{"GET", "/search/:name", Search},
	)

	http.ListenAndServe(":8080", &handler)

}

func StartParsingBBBikeData(path string) {

	bbbikeng.ConnectToDatabase()
	defer bbbikeng.Connection.Close()

	Import.ParseData(path)


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
