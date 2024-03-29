package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
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
	database.createUsersTable()
	database.createWorkoutTable()
	database.createSetTable()

	// TO DO remove these later
	database.dropTables()
	database.createUsersTable()
	database.createWorkoutTable()
	database.createSetTable()
	id1 := uuid.MustParse("fd1117b6-f2d4-48c9-b334-1676d95cfc0a")
	id2 := uuid.MustParse("54fb8829-a3f4-4bd8-8f63-b3e532365667")
	database.createDummyUser(id1, "user1", "pwd1")
	database.createDummyUser(id2, "user2", "pwd2")
	database.createDummyWorkout(id1)
	database.createDummyWorkout(id1)
	database.createDummyWorkout(id2)
	database.createDummySet(id1, 1, 1, 4, "bench")
	database.createDummySet(id1, 1, 2, 4, "bench")

	router := chi.NewRouter()
	println(os.Getenv("CLIENT_URL"))
	println(os.Getenv("JWT_SECRET"))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("CLIENT_URL")}, // Use this to allow specific origin hosts
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	  }))

	server := Server{
		database: &database,
		router: router,
	}

	server.init()
}
