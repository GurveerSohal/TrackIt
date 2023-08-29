package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (d *Database) createDummyUser(id uuid.UUID, username string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		printError(err, "error when hasing password for dummy user")
		return err
	}

	statement := `
		INSERT INTO users VALUES (
			$1, $2, $3
		);
	`

	if _, err := d.db.Exec(statement, id, username, hash); err != nil {
		printError(err, "error error when creating dummy user in database")
		return err
	}

	return nil
}

func (d *Database) createDummyWorkout(id uuid.UUID) error {
	statement := `
		INSERT INTO workouts VALUES (
			$1
		);
	`

	if _, err := d.db.Exec(statement, id); err != nil {
		printError(err, "error error when creating dummy workout in database")
		return err
	}

	return nil
}

func (d *Database) createDummySet(id uuid.UUID, workout_number int, set_number int, reps int, name string) error {
	statement := `
	INSERT INTO sets VALUES (
		$1, $2, $3, $4, $5
	);
`

if _, err := d.db.Exec(statement, id, workout_number, reps, name, set_number); err != nil {
	printError(err, "error error when creating dummy set in database")
	return err
}

return nil
}


func writeJson(w http.ResponseWriter, status int, v any) error {
	enc := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v != nil {
		if err := enc.Encode(v); err != nil {
			fmt.Println(err.Error())
			fmt.Println("error when writing json in writeJson()")
			return err
		}
	}

	return nil
}

func printError(err error, message string) {
	fmt.Println(err.Error())
	fmt.Println(message)
}
