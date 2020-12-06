package users

import (
	"github.com/ericha1981/bookstore_users-api/datasources/mysql/users_db"
	"github.com/ericha1981/bookstore_users-api/logger"
	"github.com/ericha1981/bookstore_users-api/utils/errors"
)

//====================================================================
// DAO - Data Access Object. CRUD operations (insert, update, delete)
// Access layer to our database. mysql, cassandra, etc.
//====================================================================
const (
	queryInsertUser = "insert into users(first_name, last_name, email, date_created, status, password) values (?, ?, ?, ?, ?, ?);"
	queryGetUser = "select id, first_name, last_name, email, date_created, status from users where id=?;"
	queryUpdateUser = "update users set first_name=?, last_name=?, email=? where id=?;"
	queryDeleteUser = "delete from users where id=?;"
	queryFindUserByStatus = "select id, first_name, last_name, email, date_created, status from users where status=?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Db.Prepare(queryGetUser)
	if err != nil  {
		logger.Error("error trying to prepare the get user statement", err)
		return errors.NewInternalServerError("database error")
	}

	result := stmt.QueryRow(user.Id) // pass WHERE id filter value
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error trying to get user by id", getErr)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Db.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}
	// Really important when you work with MySQL
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error trying to save user", saveErr)
		return errors.NewInternalServerError("database error")
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error trying to get last insert id after creating a new user", err)
		return errors.NewInternalServerError("database error")
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr   {
	stmt, err := users_db.Db.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id) // It's an update so don't need a result back. Just error.
	if err != nil {
		logger.Error("error trying to update user", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Delete() *errors.RestErr  {
	stmt, err := users_db.Db.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error trying to delete user", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Db.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error trying to prepare find users by status statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error trying to find users by status", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close() // VERY IMPORTANT!

	results := make([]User, 0) // initialize
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error scanning user row into user struct", err)
			return nil, errors.NewInternalServerError("database error")
		}
		results = append(results, user) // add user to array
	}

	// no user by the status found.
	if len(results) == 0 {
		// No logger here because this API is driven by user input. They can use my API in a wrong way
		// causing lots of unnecessary logging entries.
		return nil, errors.NewInternalServerError("database error")
	}
	return results, nil
}