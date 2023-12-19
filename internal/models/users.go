package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	//Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP)`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		//Check if email already exists as an entry in the user table.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		
		return err
	}

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