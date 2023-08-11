// users table
// has : uuid, created at, hash password

// workout table
// has : workout id (serial number), foriegn key uuid

// set table
// has: exercise name, number of reps, foreign key uuid and workout

package main

import (
	"database/sql"
)

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
		created timestamp NOT NULL
	);`

	_, err := d.db.Exec(statement)

	if err != nil {
		return err
	}

	return nil
}

func (database *Database) createUser(username string, password string) {

}
