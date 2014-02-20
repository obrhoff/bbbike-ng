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
	"strings"
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

func InsertCityToDatabase(city City) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(city.Geometry))
	fixedName := strings.Replace(city.Name, "'", "''", -1)

	var query string
	if city.ID != 0 {
		query = fmt.Sprintf("INSERT INTO path(id ,name, geometry) VALUES (%s, '%s', '%s', %s)", strconv.Itoa(city.ID), fixedName, points)
	} else {
		query = fmt.Sprintf("INSERT INTO city (name, geometry) VALUES ('%s', %s)", fixedName, points)
	}

	log.Println("DB:", query)
	_, err = Connection.Exec(query)
	if err != nil {
		log.Fatal("Error inserting Street Into Database: %s", err.Error())
	}
}

func InsertPlaceToDatabase(place Street) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(place.Path))
	fixedName := strings.Replace(place.Name, "'", "''", -1)

	var query string
	if place.ID != 0 {
		query = fmt.Sprintf("INSERT INTO place (placeid, name, type, geometry) VALUES (%s, '%s', '%s', %s)", strconv.Itoa(place.ID), fixedName, place.Type, points)
	} else {
		query = fmt.Sprintf("INSERT INTO place (name, type, geometry) VALUES ('%s', '%s', %s)", fixedName, place.Type, points)
	}
	log.Println("DB:", query)

	_, err = Connection.Exec(query)
	if err != nil {
		log.Fatal("Error inserting Place Into Database: %s", err.Error())
	}
}

func InsertStreetToDatabase(street Street) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(street.Path))
	fixedName := strings.Replace(street.Name, "'", "''", -1)

	var query string
	if street.ID != 0 {
		query = fmt.Sprintf("INSERT INTO path(id ,name, type, geometry) VALUES (%s, '%s', '%s', %s)", strconv.Itoa(street.ID), fixedName, street.Type, points)
	} else {
		query = fmt.Sprintf("INSERT INTO path (name, type, geometry) VALUES ('%s', '%s', %s)", fixedName, street.Type, points)
	}

	log.Println("DB:", query)
	_, err = Connection.Exec(query)
	if err != nil {
		log.Fatal("Error inserting Street Into Database: %s", err.Error())
	}
}

func InsertCyclePathToDatabase(cyclepath Street) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(cyclepath.Path))
	query := fmt.Sprintf("INSERT INTO cyclepath(type, geometry) VALUES ('%s', %s)",cyclepath.Type, points)

	log.Println("DB:", query)

	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Cyclepath Into Database: %s", err.Error())
	}

}

func InsertGreenToDatabase(green Street) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(green.Path))
	query := fmt.Sprintf("INSERT INTO greenpath(type, geometry) VALUES ('%s', %s)",green.Type, points)

	log.Println("DB:", query)

	_, err = Connection.Exec(query)
	if err != nil {
		log.Fatal("Error inserting Greenway Into Database: %s", err.Error())
	}

}

func InsertQualityToDatabase(quality Street) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(quality.Path))
	query := fmt.Sprintf("INSERT INTO quality(type, geometry) VALUES ('%s', %s)", quality.Type, points)

	log.Println("DB:", query)
	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Qualitys Into Database: %s", err.Error())
	}

}


func InsertStreetLightToDatabase(light Street) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(light.Path))
	query := fmt.Sprintf("INSERT INTO trafficlight(type, geometry) VALUES ('%s', %s)", light.Type, points)
	log.Println("DB:", query)
	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Light Into Database: %s", err.Error())
	}
}

func InsertUnlitToDatabase(light Street) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(light.Path))
	query := fmt.Sprintf("INSERT INTO unlitpath(geometry) VALUES (%s)", points)
	log.Println("DB:", query)
	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Unlit Into Database: %s", err.Error())
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

func GetStreetFromId(id int) (street Street) {

	street.ID = id
	var geometrys string
	var nodes string

	err := Connection.QueryRow("select wayid, name, type, ST_AsGeoJSON(geometry), nodes from way where wayid = $1", id).Scan(&street.ID, &street.Name, &street.Type, &geometrys, &nodes)
	if err != nil {
		log.Fatal("Error on getting Street from ID %s", err.Error())
	}

	var nodesList []int
	json.Unmarshal([]byte(nodes), &nodesList)
	street.SetPathFromGeoJSON(geometrys)


	return street

}

func GetCyclepathFromId(id int) (cyclepath Street) {

	var geometrys string
	err := Connection.QueryRow("select cycleid, type, ST_AsGeoJSON(geometry) from cyclepaths where cycleid = $1", id).Scan(&cyclepath.ID, &cyclepath.Type, &geometrys)
	if err != nil {
		log.Fatal("Error on getting Cyclepath from ID: %s", err.Error())
	}
	cyclepath.SetPathFromGeoJSON(geometrys)
	return cyclepath

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

	rows, err := Connection.Query("SELECT neighbor.id, networkid, wayid, type, attributesToJson(attributes), name, st_asgeojson(neighbor.node_geo) as nodecoord, st_asgeojson(geometry) as path, neighbor.walkable, neighbor.trafficlight FROM network, ( SELECT parent, id, geometry as node_geo, walkable, trafficlight FROM node JOIN (select id as parent, unnest(neighbors) as id from node where id = $1) x USING (id)) as neighbor WHERE network.nodes @> ARRAY[neighbor.parent,neighbor.id]", node.NodeID)
	if err != nil {
		log.Fatal("Error fetching Neighbor Nodes:", err)
	}

	defer rows.Close()
	for rows.Next() {

		var newNode Node
		var nodeGeometry string
		var pathGeometry string
		var attributes string

		err := rows.Scan(&newNode.NodeID, &newNode.StreetFromParentNode.ID, &newNode.StreetFromParentNode.WayID, &newNode.StreetFromParentNode.Type, &attributes,&newNode.StreetFromParentNode.Name, &nodeGeometry, &pathGeometry, &newNode.Walkable, &newNode.TrafficLight)
		if err != nil {
			log.Fatal("Error Neighbor Nodes:", err)
		}
		newNode.NodeGeometry = ConvertGeoJSONtoPoint(nodeGeometry)
		newNode.StreetFromParentNode.Path = ConvertGeoJSONtoPath(pathGeometry)
		newNode.StreetFromParentNode.ParseAttributes(attributes)
		nodes = append(nodes, newNode)
	}

	return nodes

}

func GetCyclepathFromStreet(street Street) (cyclepaths []Cyclepath) {

	rows, err := Connection.Query("select s2.type, ST_AsGeoJSON(s2.path) from streetpaths s1, cyclepaths s2 where s1.streetpathid = $1 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path))", street.ID)
	if err != nil {
		log.Fatal("Error on getting Cyclepaths: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var newCyclepath Cyclepath
		var geometrys string
		err := rows.Scan(&newCyclepath.Type, &geometrys)
		if err != nil {
			log.Fatal(err)
		}
		newCyclepath.Path = ConvertGeoJSONtoPath(geometrys)

		if len(newCyclepath.Path) > 0 {
			cyclepaths = append(cyclepaths, newCyclepath)
		}

	}

	fmt.Println("Cyclepaths:", cyclepaths)
	return cyclepaths
}

func GetQualityFromStreet(street Street) (qualitys []Quality) {

	rows, err := Connection.Query("select s2.type, ST_AsGeoJSON(s2.path) from streetpaths s1, qualitys s2 where s1.streetpathid = $1 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path))", street.ID)
	if err != nil {
		log.Fatal("Error on getting Qualitys: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var newQuality Quality
		var geometrys string
		err := rows.Scan(&newQuality.Type, &geometrys)
		if err != nil {
			log.Fatal(err)
		}
		newQuality.Path = ConvertGeoJSONtoPath(geometrys)
		qualitys = append(qualitys, newQuality)
	}

	fmt.Println("Quality:", qualitys)

	return qualitys
}

func GetGreenwaysFromStreet(street Street) (greenways []Greenway) {

	rows, err := Connection.Query("select s2.type, ST_AsGeoJSON(s2.path) from streetpaths s1, greenways s2 where s1.streetpathid = $1 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path))", street.ID)
	if err != nil {
		log.Fatal("Error on getting Greenways: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var newGreenway Greenway
		var geometrys string
		err := rows.Scan(&newGreenway.Type, &geometrys)
		if err != nil {
			log.Fatal(err)
		}
		newGreenway.Path = ConvertGeoJSONtoPath(geometrys)
		greenways = append(greenways, newGreenway)
	}

	fmt.Println("Greenway:", greenways)

	return greenways
}

func SearchForStreetName(name string) (streets []Street) {

	log.Println("Searching for Streetname:", name)

	rows, err := Connection.Query("select wayid, name, type, ST_AsGeoJSON(geometry) from way where name ilike $1 LIMIT 10", ("%" + name + "%"))
	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {

		var newStreet Street
		var geometrys string
		err := rows.Scan(&newStreet.ID, &newStreet.Name, &newStreet.Type, &geometrys)
		if err != nil {
			log.Fatal(err)
		}

		newStreet.SetPathFromGeoJSON(geometrys)
		streets = append(streets, newStreet)
	}

	return streets
}

func SearchForNearestStreetFromPoint(point Point) (street Street) {

	//SELECT * FROM streetpath ORDER BY ST_Distance(path, ST_GeomFromText('POINT(13.373370 52.551080)', 4326)) LIMIT 1;
	var geometrys string
	latPath, lngPath := point.LatitudeLongitudeAsString()

	makePoint := ("ST_Distance(path, ST_GeomFromText('POINT(" + lngPath + " " + latPath + ")', 4326))")
	query := fmt.Sprintf("SELECT streetpathid, name, type, ST_AsGeoJSON(path)  FROM streetpaths ORDER BY %s LIMIT 1", makePoint)
	err := Connection.QueryRow(query).Scan(&street.ID, &street.Name, &street.Type, &geometrys)

	if err != nil {
		log.Fatal(err)
	}


	street.SetPathFromGeoJSON(geometrys)
	return street

}

func SearchForNearestCyclepathFromPoint(point Point) (cyclepath Street) {

	var geometrys string
	latPath, lngPath := point.LatitudeLongitudeAsString()

	makePoint := ("ST_Distance(path, ST_GeomFromText('POINT(" + lngPath + " " + latPath + ")', 4326))")
	query := fmt.Sprintf("SELECT cycleID, type, ST_AsGeoJSON(path) FROM cyclepaths ORDER BY %s LIMIT 1", makePoint)

	err := Connection.QueryRow(query).Scan(&cyclepath.ID, &cyclepath.Type, &geometrys)

	if err != nil {
		log.Fatal(err)
	}

	cyclepath.SetPathFromGeoJSON(geometrys)

	return cyclepath

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
