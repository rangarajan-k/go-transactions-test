# Go Transactions Test
Simple Golang application to have transactional APIs

Go Version used  : 1.22.4 \
PostgresDB Server Version : 16.3 \
Main Go Dependencies : Gin, Go-Pg, Viper \

## Directory Structure : 

├── README.md\
├── app\
│└── app.go\
├── config\
│└── config.go\
├── config.json\
├── controller\
│└── transaction.go\
├── datastore\
│├── postgress_client.go\
│└── postgress_store.go\
├── dicontainer\
│└── container.go\
├── go.mod\
├── go.sum\
├── main.go\
├── docs\
│├── docs.go\
│├── swagger.json\
│└── swagger.yaml\
└── router\
  └── router.go

App : Main entry point and initialization of the Application\
Config : Read config file from disk and apply to config map\
Controller : Main API Logic\
DataStore : Handle DB connection & queries\
DiContainer : Dependency injection container when starting application\
Router : Initialize HTTP server endpoints and associate with controllers
Docs : Contains swagger documentation

## Configuration
Default config that will be used when using PostgressDB pulled from docker.

```json
{
  "port": 9000,
  "gin_mode": "debug",
  "db_config": {
    "host": "localhost",
    "port": 5432,
    "username": "postgres",
    "password": "postgres123",
    "database_name": "postgres"
  },
  "swagger_config": {
    "version": 1,
    "host": "http://localhost:9000",
    "base_path": "/",
    "title": "Transaction Service APIs",
    "description": "Transactions Service for creating accounts and posting transactions"
  }
}
```

## Setup Guide
Pre-requisites : Go 1.22.4 \
Supported OS: Windows/ Mac/ Linux \
PostgressDB server : If you already have a server running you can give the address and credentials in `config.json` file

Setup PostgressDB using Docker :\
Make sure docker is installed and running : https://docs.docker.com/engine/install/

Pull the PostgressDB docker image : `docker pull postgres` \
Start the PostgressDB server as a docker container : `docker run --name my-postgres -e POSTGRES_PASSWORD=postgres123 -p 5432:5432 -d postgres`

Next we can build and start out go application

- Clone the repo and navigate into the cloned directory
- Install the go modules : `go get`
- Build & Run application : `go run main.go`

View Swagger Documentation for APIs : http://localhost:9000/swagger/index.html \
Swagger file location : `./docs/swagger.yaml`

### Assumptions: 
- Balance of account can not be in negative ( < 0 )


Postman collection for testing can be found in this file in root dir : `./Transaction Test.postman_collection.json`