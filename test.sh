#!/usr/bin/env bash

source .env

# Start the DB
docker-compose -f docker-compose_testing.yaml up --build -d >/dev/null 2>&1

# Load env variables and run tests on the given package
ENVIRONMENT=TEST \
GIN_MODE=release \
POSTGRES_HOST=localhost \
POSTGRES_USER=$POSTGRES_USER \
POSTGRES_PASSWORD=$POSTGRES_PASSWORD  \
POSTGRES_PORT=$POSTGRES_PORT \
POSTGRES_DB=$POSTGRES_DB \
go test -failfast \
-cover \
-coverprofile=coverage.txt \
-covermode=atomic \
-v github.com/Brawdunoir/dionysos-server/$"1"

# Stop the DB
docker-compose -f docker-compose_testing.yaml down >/dev/null 2>&1
