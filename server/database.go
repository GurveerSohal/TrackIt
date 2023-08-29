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
	Id       uuid.UUID
	Username string
	Hash     []byte
	Created  time.Time
}

type Database struct {
	db *sql.DB
}

func (d *Database) dropTables() error {
	statement := `DROP TABLE IF EXISTS users, workouts, sets;`

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
		printError(err, "error when creating users table")
		return err
	}

	return nil
}

func (d *Database) createWorkoutTable() error {
	// The SQL standard requires that writing just timestamp be equivalent to timestamp without time zone,
	// and PostgreSQL honors that behavior.

	statement := `CREATE TABLE IF NOT EXISTS workouts (
		user_id uuid NOT NULL,
		workout_number smallint,
		created timestamp NOT NULL DEFAULT now(),
		PRIMARY KEY(user_id, workout_number)
	);
	
	-- Create a function to calculate workout_number
	CREATE OR REPLACE FUNCTION calculate_workout_number()
	RETURNS TRIGGER AS $$
	BEGIN
    	NEW.workout_number := COALESCE((SELECT MAX(workout_number) FROM workouts WHERE user_id = NEW.user_id), 0) + 1;
    	RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	
	-- Create a trigger that calls the calculate_workout_number function on INSERT
	DROP TRIGGER IF EXISTS update_workout_number ON workouts;

	CREATE TRIGGER update_workout_number
	BEFORE INSERT ON workouts
	FOR EACH ROW
	EXECUTE FUNCTION calculate_workout_number();
	`

	_, err := d.db.Exec(statement)

	if err != nil {
		printError(err, "error when creating workout table")
		return err
	}

	return nil
}

func (d *Database) createSetTable() error {
	// The SQL standard requires that writing just timestamp be equivalent to timestamp without time zone,
	// and PostgreSQL honors that behavior.

	statement := `CREATE TABLE IF NOT EXISTS sets (
		user_id uuid NOT NULL,
		workout_number smallint NOT NULL,
		set_number smallint NOT NULL,
		reps smallint NOT NULL,
		name VARCHAR(256) NOT NULL,
		PRIMARY KEY(user_id, workout_number, set_number)
	);`

	_, err := d.db.Exec(statement)

	if err != nil {
		printError(err, "error when creating set table")
		return err
	}

	return nil
}

func (d *Database) createUser(username string, password string) error {
	id := uuid.New()
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

func (d *Database) createWorkout(id uuid.UUID) (int, error) {

	statement, err := d.db.Prepare(`
		INSERT INTO workouts VALUES (
			$1
		)
		RETURNING workout_number;
	`); 
	if err != nil {
		printError(err, "error error when creating workout in database (while preparing statment)")
		return -1, err
	}

	defer statement.Close()

	var workout_number int
	err = statement.QueryRow(id, /* other_column_values... */).Scan(&workout_number)
	if err != nil {
		printError(err, "error error when creating workout in database")
		return -1, err
	}
	return workout_number, nil
}

func (d *Database) getUser(username string) (*User, error) {
	// get the user
	user := new(User)
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
