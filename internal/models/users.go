package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID 			   int
	Name 	  	   string
	Email 		   string
	HashedPassword []byte
	Created 	   time.Time
}

//Define a new UserModel struct which wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

//Adds a new record to the "users" table.
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

//Verifies wether a user exists with the provided email address and password. Returns the relevant user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

//Checks if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}