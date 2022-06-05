package routes

import (
	"fmt"
	"log"
	"os"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

// client is the client instance for the arango database
var client driver.Client

// db is the database used by the api
var db driver.Database

// cols is an array of collections
var cols = []string{"user"}

// getClient returns a new driver instance for the arango database
func getClient() driver.Client {

	uri, found := os.LookupEnv("ARANGO_URI")
	if !found {
		panic("ARANGO_URI environment variable not found")
	}
	username, found := os.LookupEnv("ARANGO_USERNAME")
	if !found {
		panic("ARANGO_USERNAME environment variable not found")
	}
	password, found := os.LookupEnv("ARANGO_PASSWORD")
	if !found {
		panic("ARANGO_PASSWORD environment variable not found")
	}

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{uri},
	})
	if err != nil {
		log.Fatalf("Failed to create HTTP connection: %v", err)

	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(username, password),
	})

	return client
}

// getDatabase returns a database instance
func getDatabase(name string) (db driver.Database) {
	db_exists, err := client.DatabaseExists(nil, name)
	if err != nil {
		panic(err)
	}

	if db_exists {
		fmt.Println("That db exists already")

		db, err = client.Database(nil, name)

		if err != nil {
			log.Fatalf("Failed to open existing database: %v", err)
		}

	} else {
		db, err = client.CreateDatabase(nil, name, nil)

		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
	}

	return db
}

// SetupCollection takes an array of collections name and creates them if they do not exist
func SetupCollection() {
	for _, collection := range cols {
		coll_exists, err := db.CollectionExists(nil, collection)
		if err != nil {
			panic(err)
		}

		if coll_exists {
			fmt.Printf("%s collection exists already\n", collection)
		} else {

			var col driver.Collection
			col, err = db.CreateCollection(nil, collection, nil)

			if err != nil {
				log.Fatalf("Failed to create %s collection: %v", collection, err)
			}

			fmt.Printf("Created collection '%s' in database '%s'\n", col.Name(), db.Name())
		}
	}
}
