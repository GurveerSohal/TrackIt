// users table
// has : uuid, created at, hash password

// workout table
// has : workout id (serial number), foriegn key uuid

// set table
// has: exercise name, number of reps, foreign key uuid and workout

package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id uuid.UUID
	Username string
	Hash []byte
	Created time.Time
}

type Database struct {
	db *sql.DB
}

func (d *Database) dropUsersTable() error {
	statement := `DROP TABLE IF EXISTS users;`

	if _, err := d.db.Exec(statement); err != nil {
		return err
	}

	return nil
}

func (d *Database) createUsersTable() error {
	// The SQL standard requires that writing just timestamp be equivalent to timestamp without time zone,
	// and PostgreSQL honors that behavior.

	statement := ` CREATE TABLE IF NOT EXISTS users (
		id uuid PRIMARY KEY,
		username varchar(256) NOT NULL,
		hash char(60) NOT NULL,
		created timestamp NOT NULL DEFAULT now()
	);`

	_, err := d.db.Exec(statement)

	if err != nil {
		return err
	}

	return nil
}

func (d *Database) createDummyUser() error {
	id := uuid.New()
	username := "testAccount"
	password := "testP*ssw0rd!"

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
func (d *Database) createUser(username string, password string) {

}

func (d *Database) getUser(username string) (*User, error) {
	// get the user
	user  := new(User)
	query := `
		SELECT * 
		FROM users
		WHERE 
			username = $1
	`
	row := d.db.QueryRow(query, username)
	if err := row.Scan(&user.Id, &user.Username, &user.Hash, &user.Created); err != nil {
		return nil, err
	}

	return user, nil
}