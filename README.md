# Dionysos Server
[![GitHub Super-Linter](https://github.com/Brawdunoir/dionysos-server/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)
[![Go Report Card](https://goreportcard.com/badge/github.com/Brawdunoir/dionysos-server)](https://goreportcard.com/report/github.com/Brawdunoir/dionysos-server)

Dionysos Server is a golang RESTful API for the [dionysos-client](https://github.com/Brawdunoir/dionysos-client) project, enabling users to **share cinematic experiences**.

## Developping another client
This API is hosted at https://ipa.dionysos.live and is publicly available. The default client using this API is located at https://dionysos.live.

You are free to use the API and developping another client. See [documentation](https://github.com/found-pleiade/dionysos-server#documentation).

Users using your client will be able to interact with users using the default client and vice-versa.

## Self-hosting
If you have your own server setup and you don't want to use the default API to connect with your friends, feel free to host your own dionysos server instance.

Note that your server and its ports need to be publicly accessible from your friends. We recommend using a reverse proxy such as traefik or nginx.

### Using docker 🐳
The preferred method is to use the official dionysos API image along with its postgresSQL database using a `docker-compose.yaml` file:

```yaml
---
version: '3.5'

networks:
  dionysos:

services:
  api:
    container_name: api
    image: brawdunoir/dionysos-server:0.2
      your-api-local-port:8080
    networks:
      dionysos:
    environment:
      - ENVIRONMENT=PROD
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
      - POSTGRES_HOST=${POSTGRES_HOST}
      - POSTGRES_PORT=${POSTGRES_PORT}
      - BASE_PATH=${BASE_PATH}
    depends_on:
      - postgres
  postgres:
    container_name: ${POSTGRES_HOST}
    image: postgres:14
    command: -p ${POSTGRES_PORT}
    volumes:
      - ./path-to-your-db-folder:/var/lib/postgresql/data
    networks:
      dionysos:
    environment:
      - POSTGRES_DB=${POSTGRES_DB}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
```

Add a `.env` (check the one in the repo) along with this `docker-compose.yaml`. We recommend to change default passwords.

## I want to participate 🍵
First of all you can:
- Fill issues for enhancements or bugs, we will try to fix them asap
- Fork this repository, make code changes and create a pull request (see below for tools)

We ship several docker-based tools that help the development workflow. You need to have `docker` and `docker-compose` installed to leverage these tools.
### Run the API
Simply run `docker-compose up` and you should be good to go !

Hot-reloading is supported out of the box.

The API uses port 8080. A pgAdmin instance runs on port 8081. Check `.env` to see credentials.

### Run tests
We have a script that runs a postgresSQL database and run `go test` automatically so you can run:

`./test.sh <package>`

For example: `./test.sh routes`

## Documentation
WIP #59
