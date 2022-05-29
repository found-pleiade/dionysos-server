package routes

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Neo4jConfiguration struct {
	Url      string
	Username string
	Password string
	Database string
}

func newDriverFromConfig() neo4j.Driver {
	// Configuring database
	configuration := parseConfiguration()
	//Client Database instance
	driver, err := configuration.newDriver()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to Neo4j driver !")
	}
	defer unsafeClose(driver)
	return driver
}

func (nc Neo4jConfiguration) newDriver() (neo4j.Driver, error) {
	return neo4j.NewDriver(nc.Url, neo4j.BasicAuth(nc.Username, nc.Password, ""))
}

func parseConfiguration() *Neo4jConfiguration {
	return &Neo4jConfiguration{
		Url:      lookupEnvOrGetDefault("NEO4J_URI", "http://localhost:7474/"),
		Username: lookupEnvOrGetDefault("NEO4J_USER", "dionysos"),
		Password: lookupEnvOrGetDefault("NEO4J_PASSWORD", "dionysos"),
		Database: lookupEnvOrGetDefault("NEO4J_DATABASE", "dionysos"),
	}
}

func lookupEnvOrGetDefault(key string, defaultValue string) string {
	if env, found := os.LookupEnv(key); !found {
		return defaultValue
	} else {
		return env
	}
}

func unsafeClose(closeable io.Closer) {
	if err := closeable.Close(); err != nil {
		log.Fatal(fmt.Errorf("could not close resource: %w", err))
	}
}
