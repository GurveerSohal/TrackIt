package main

import (
	"database/sql"
	"fmt"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

// open a connection to the database, pass it to server and start

func main() {
	connStr := "user=postgres dbname=trackitdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if (err != nil) {
		fmt.Println("error when opening postgres")
		return
	}

	if err := db.Ping(); (err != nil) {
		fmt.Println(err.Error())
		fmt.Println("error when pinging db")
		return
	}

	database := Database{db: db}
	database.dropUsersTable()
	database.createUsersTable()

	router := chi.NewRouter()

	server := Server{
		database: &database,
		router: router,
	}

	server.init()
}
