# go-transactions-test
Simple Golang application to have transactional APIs

Go Version used  : 1.22.4 \
PostgresDB Server Version : 16.3 \
Main Go Dependencies : Gin, Go-Pg, Viper \

Directory Structure : 

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
└── router\
  └── router.go