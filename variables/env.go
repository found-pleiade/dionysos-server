package variables

import "os"

// All the possible values for the env variables.

// BasePath is the base path of the API. e.g. http://localhost:8080/api/v1 if set to /api/v1.
var BasePath string = os.Getenv("BASE_PATH")

// Environment is the environment of the API. e.g. PROD, DEV, TEST. See const.go for the possible values.
var Environment string = os.Getenv("ENVIRONMENT")

// RedisHost is the host of the Redis server.
var RedisHost string = os.Getenv("REDIS_HOST")
