package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/leesper/couchdb-golang"
	"github.com/sirupsen/logrus"
)

var (
	dfMeta      *couchdb.Database
	url         string
	couchDbPort int
	serverPort  int
	couchDbName string
	couchDbHost string
	avail       bool
	err         error
)

func init() {
	serverPortStr, avail := os.LookupEnv("LISTEN_PORT")
	if !avail {
		logrus.Fatal("Please set the server listening port to LISTEN_PORT env variable")
	} else {
		serverPort, err = strconv.Atoi(serverPortStr)
		if err != nil {
			logrus.Fatal("Error parsing server listening port")
		}
	}

	couchDbPortStr, avail := os.LookupEnv("COUCHDB_PORT")
	if avail {
		couchDbPort, err = strconv.Atoi(couchDbPortStr)
		if err != nil {
			logrus.Fatal("Error parsing couchdb port")
		} else {
			logrus.Debugf("Connectingto couchdb port  %d ", couchDbPort)
		}
	}

	couchDbName, avail = os.LookupEnv("SERVE_DATABASE")
	if !avail {
		logrus.Fatalf("Please set one CouchDB database to serve over API in env variable SERVE_DATABASE")
	}

	couchDbHost, avail = os.LookupEnv("COUCH_HOST")
	if !avail {
		logrus.Fatal("Please set the CoucchDB host to env variable COUCH_HOST")
	} else {
		logrus.Debugf("The couchdb connecting to %s ", couchDbHost)
	}
}

func main() {
	logrus.Info("Establishing the CouchDB connection.")
	server := setUpCouchDbServer()
	databases, err := server.DBs()
	if err != nil {
		logrus.Fatalf("Error occured while connecting picking the databases in CuchDB error %s", err)
		os.Exit(1)
	}

	for i := 0; i < len(databases); i++ {
		fmt.Printf(" %s\n", databases[i])
	}

	if !contains(databases, couchDbName) {
		logrus.Fatalf("No database available on the %s ", couchDbName)
	} else {
		dfMeta = connectDatabase(server, couchDbName)
	}

	http.HandleFunc("/data", readdata)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil))
}

func setUpCouchDbServer() *couchdb.Server {
	if couchDbPort != 0 {
		url = fmt.Sprintf("http://%s:%d", couchDbHost, couchDbPort)
	} else {
		url = fmt.Sprintf("http://%s", couchDbHost)
	}
	server, err := couchdb.NewServer(url)
	if err != nil {
		logrus.Fatalf("Error occured while connecting to db host %s with error %s", couchDbHost, err)
		os.Exit(1)
	}
	return server
}

func contains(list []string, searchedValue string) bool {
	for _, value := range list {
		if value == searchedValue {
			return true
		}
	}
	return false
}

func connectDatabase(server *couchdb.Server, dbName string) *couchdb.Database {

	logrus.Infof("Connecting to database %s...", dbName)
	db, err := server.Get(couchDbName)
	if err != nil {
		logrus.Fatalf("Connection failed to %s with error %s", dbName, err)
		os.Exit(1)
	}
	logrus.Infof("Connected to %s", dbName)
	return db
}

func readdata(res http.ResponseWriter, req *http.Request) {

	ids, err := dfMeta.DocIDs()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	var data []map[string]interface{}

	for _, id := range ids {
		doc, _ := dfMeta.Get(id, nil)
		data = append(data, doc)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(data)
}
