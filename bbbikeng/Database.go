/**
 * User: DocterD
 * Date: 28/12/13
 * Time: 11:19
 */

package bbbikeng

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"encoding/json"
	"os"
	"io/ioutil"
)

var Connection *sql.DB

type Config struct {
	User string
	Password string
	Host string
	Port string
	Database string
	Debug bool
}

func InsertPlaceToDatabase(place *WayAttribute) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(place.Geometry()))

	var query string
	if place.Id() != 0 {
		query = fmt.Sprintf("INSERT INTO place (placeid, name, type, geometry) VALUES (%s, '%s', '%s', %s)", place.Id(), place.Name(), place.Type(), points)
	} else {
		query = fmt.Sprintf("INSERT INTO place (name, type, geometry) VALUES ('%s', '%s', %s)", place.Name(), place.Type(), points)
	}
	log.Println("DB:", query)

	_, err = Connection.Exec(query)
	if err != nil {
		log.Fatal("Error inserting Place Into Database: %s", err.Error())
	}
}

func InsertStreetToDatabase(street *WayAttribute) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(street.Geometry()))

	var query string
	if street.Id() != 0 {
		query = fmt.Sprintf("INSERT INTO path(id ,name, type, geometry) VALUES (%s, '%s', '%s', %s)", strconv.Itoa(street.Id()), street.Name(), street.Type(), points)
	} else {
		query = fmt.Sprintf("INSERT INTO path (name, type, geometry) VALUES ('%s', '%s', %s)", street.Name(), street.Type(), points)
	}

	log.Println("DB:", query)
	_, err = Connection.Exec(query)
	if err != nil {
		log.Fatal("Error inserting Street Into Database: %s", err.Error())
	}
}

func InsertCyclePathToDatabase(cyclepath *CyclepathAttribute) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(cyclepath.Geometry()))
	query := fmt.Sprintf("INSERT INTO cyclepath(type, geometry) VALUES ('%s', %s)",cyclepath.Type(), points)

	log.Println("DB:", query)

	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Cyclepath Into Database: %s", err.Error())
	}

}

func InsertGreenToDatabase(green *GreenwayAttribute) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(green.Geometry()))
	query := fmt.Sprintf("INSERT INTO greenpath(type, geometry) VALUES ('%s', %s)",green.Type(), points)

	log.Println("DB:", query)

	_, err = Connection.Exec(query)
	if err != nil {
		log.Fatal("Error inserting Greenway Into Database: %s", err.Error())
	}

}

func InsertQualityToDatabase(quality *QualityAttribute) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(quality.Geometry()))
	query := fmt.Sprintf("INSERT INTO quality(type, geometry) VALUES ('%s', %s)", quality.Type(), points)

	log.Println("DB:", query)
	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Qualitys Into Database: %s", err.Error())
	}

}


func InsertStreetLightToDatabase(trafflight *TrafficLightAttribute) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(trafflight.Geometry()))
	query := fmt.Sprintf("INSERT INTO trafficlight(type, geometry) VALUES ('%s', %s)", trafflight.Type(), points)
	log.Println("DB:", query)
	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Light Into Database: %s", err.Error())
	}
}

func InsertUnlitToDatabase(unlit *UnlitAttribute) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(unlit.Geometry()))
	query := fmt.Sprintf("INSERT INTO unlitpath(geometry) VALUES (%s)", points)
	log.Println("DB:", query)
	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Unlit Into Database: %s", err.Error())
	}
}

func InsertHandicapToDatabase(handicap *HandicapAttribute) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(handicap.Geometry()))
	query := fmt.Sprintf("INSERT INTO handicap(description, type, geometry) VALUES ('%s', '%s', %s)", handicap.Name(), handicap.Type(), points)
	log.Println("DB:", query)
	_, err = Connection.Exec(query)
	if err != nil {
		log.Fatal("Error inserting Handicap Into Database: %s", err.Error())
	}
}


func GetNodeFromId(id int) (node Node) {

	var nodes string
	var ways string
	var geometry string

	var nodesList []int
	var wayList []int

	var newGeometry GeoJSON

	err := Connection.QueryRow("select id, ST_AsGeoJSON(geometry), array_to_json(networks), array_to_json(neighbors), walkable, trafficlight from node where id = $1", id).Scan(&node.NodeID, &geometry, &ways, &nodes, &node.Walkable, &node.TrafficLight)

	if err != nil {
		log.Fatal("Error on getting Node from ID %s", err.Error())
	}

	json.Unmarshal([]byte(nodes), &nodesList)
	json.Unmarshal([]byte(ways), &wayList)
	json.Unmarshal([]byte(geometry), &newGeometry)

	fmt.Println("GeoJSON:", newGeometry)

	node.NodeGeometry = ConvertGeoJSONtoPoint(geometry)

	return node

}

func GetStreetFromId(id int) (street WayAttribute) {

	street.SetId(id)
	var geometrys string
	var nodes string
	var streetName string
	var streetType string

	err := Connection.QueryRow("select name, type, ST_AsGeoJSON(geometry), nodes from way where wayid = $1", id).Scan(&streetName, &streetType, &geometrys, &nodes)
	if err != nil {
		log.Fatal("Error on getting Street from ID %s", err.Error())
	}

	var nodesList []int
	json.Unmarshal([]byte(nodes), &nodesList)

	street.SetId(id)
	street.SetName(streetName)
	street.SetType(streetType)
	street.SetPathFormGeoJSONString(geometrys)

	return street

}

func FindNearestNode(point Point) (closestNode Node){

	lat, lng := point.LatitudeLongitudeAsString()
	makePoint := ("ST_Distance(geometry, ST_GeomFromText('POINT(" + lng + " " + lat + ")', 4326))")
	query := fmt.Sprintf("SELECT id FROM node ORDER BY %s LIMIT 1", makePoint)

	log.Println("DB:", query)

	var nodeid int
	err := Connection.QueryRow(query).Scan(&nodeid)
	if err != nil {
		log.Fatal("Error on getting Closest Node from ID: %s", err.Error())
	}

	return GetNodeFromId(nodeid)
}


func GetNeighborNodesFromNode(node Node) (nodes []Node) {

	rows, err := Connection.Query("SELECT neighbor.id, networkid, wayid, type, attributesToJson(defaults) as defaults, attributesToJson(normal) as normal, attributesToJson(reversed) as reversed , name, st_asgeojson(neighbor.node_geo) as nodecoord, st_asgeojson(geometry) as path, neighbor.walkable, neighbor.trafficlight FROM network, ( SELECT parent, id, geometry as node_geo, walkable, trafficlight FROM node JOIN (select id as parent, unnest(neighbors) as id from node where id = $1) x USING (id)) as neighbor WHERE network.nodes @> ARRAY[neighbor.parent,neighbor.id]", node.NodeID)
	if err != nil {
		log.Fatal("Error fetching Neighbor Nodes:", err)
	}

	defer rows.Close()
	for rows.Next() {

		var newNode Node
		var nodeGeometry string
		var pathGeometry string
		var globalAttributes string
		var normalAttributes string
		var flippedAttributes string

		err := rows.Scan(&newNode.NodeID, &newNode.StreetFromParentNode.ID, &newNode.StreetFromParentNode.WayID, &newNode.StreetFromParentNode.Type, &globalAttributes, &normalAttributes, &flippedAttributes,&newNode.StreetFromParentNode.Name, &nodeGeometry, &pathGeometry, &newNode.Walkable, &newNode.TrafficLight)
		if err != nil {
			log.Fatal("Error Neighbor Nodes:", err)
		}
		newNode.NodeGeometry = ConvertGeoJSONtoPoint(nodeGeometry)
		newNode.StreetFromParentNode.Path = ConvertGeoJSONtoPath(pathGeometry)
	//	newNode.StreetFromParentNode.GlobalAttribute = ParseAttributes(globalAttributes)
		newNode.StreetFromParentNode.NormalAttribute = ParseAttributes(normalAttributes)
	//–	newNode.StreetFromParentNode.FlippedAttribute = ParseAttributes(flippedAttributes)

		nodes = append(nodes, newNode)
	}

	return nodes

}

func SearchForNearestStreetFromPoint(point Point) (street WayAttribute) {

	//SELECT * FROM streetpath ORDER BY ST_Distance(path, ST_GeomFromText('POINT(13.373370 52.551080)', 4326)) LIMIT 1;
	var geometrys string
	var streetName string
	var streetType string
	var streetId string

	latPath, lngPath := point.LatitudeLongitudeAsString()

	makePoint := ("ST_Distance(path, ST_GeomFromText('POINT(" + lngPath + " " + latPath + ")', 4326))")
	query := fmt.Sprintf("SELECT networkid, name, type, ST_AsGeoJSON(path)  FROM streetpaths ORDER BY %s LIMIT 1", makePoint)
	err := Connection.QueryRow(query).Scan(&streetId, &streetName, &streetType, &geometrys)

	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.Atoi(streetId)
	street.SetName(streetName)
	street.SetType(streetType)
	street.SetId(id)
	street.SetPathFormGeoJSONString(geometrys)
	return street

}

func ConnectToDatabase() {

	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var newConfig Config
	json.Unmarshal(file, &newConfig)
	connectionParameter := fmt.Sprint("user=", newConfig.User, " password=", newConfig.Password, " host=", newConfig.Host, " port=", newConfig.Port, " dbname=", newConfig.Database)
	log.Println("Connecting to Database:", connectionParameter)

	var err error
	Connection, err = sql.Open("postgres", connectionParameter)
	err = Connection.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible

	if err != nil {
		log.Panic("Error on opening database connection: %s", err.Error())
	}


}
