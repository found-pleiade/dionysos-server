package variables

import (
	"log"
	"os"
)

type Variable struct {
	Key     string
	Value   *string
	Default string
	Needed  bool
}

// Environment is the environment of the API. e.g. PROD, DEV, TEST. See const.go for the possible values.
var Environment string

// Port is the port of the API.
var Port string

// BasePath is the base path of the API. e.g. http://localhost:8080/api/v1 if set to /api/v1.
var BasePath string

// RedisHost is the host of the Redis server.
var RedisHost string

// PostgresHost is the host of the Postgres server.
var PostgresHost string

// PostgresPort is the port of the Postgres server.
var PostgresPort string

// PostgresUser is the user in the Postgres server.
var PostgresUser string

// PostgresPassword is the password of the Postgres user.
var PostgresPassword string

// PostgresDB is the name of the database in the Postgres server.
var PostgresDB string

// All the possible variables within environment.
var env = []Variable{
	{"ENVIRONMENT", &Environment, ENVIRONMENT_PRODUCTION, false},
	{"PORT", &Port, "8080", false},
	{"BASE_PATH", &BasePath, "", false},
	{"REDIS_HOST", &RedisHost, "", false},
	{"POSTGRES_HOST", &PostgresHost, "", true},
	{"POSTGRES_PORT", &PostgresPort, "", true},
	{"POSTGRES_USER", &PostgresUser, "", true},
	{"POSTGRES_PASSWORD", &PostgresPassword, "", true},
	{"POSTGRES_DB", &PostgresDB, "", true},
}

// LoadVariables loads the environment variables and set the default values if not found.
func LoadVariables() {
	log.Print("Loading variables from environment…")

	for _, v := range env {
		value, found := os.LookupEnv(v.Key)
		if !found {
			if v.Needed {
				log.Fatalf("Environment variable NOT found: %s", v.Key)
			} else if v.Default != "" {
				log.Printf("Environment variable NOT found: %s, using default value: %s", v.Key, v.Default)
				*v.Value = v.Default
			} else {
				log.Printf("Environment variable NOT found: %s, using empty value", v.Key)
			}
		} else {
			log.Printf("Environment variable loaded: %s", v.Key)
			*v.Value = value
		}
	}

	log.Print("Variables loaded.")
}
