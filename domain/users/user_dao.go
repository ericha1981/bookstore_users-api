package users

import (
	"fmt"
	"github.com/ericha1981/bookstore_users-api/datasources/mysql/users_db"
	"github.com/ericha1981/bookstore_users-api/utils/errors"
	"github.com/ericha1981/bookstore_users-api/utils/mysql_utils"
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
		return errors.NewInternalServerError(err.Error())
	}

	result := stmt.QueryRow(user.Id) // pass WHERE id filter value
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		return mysql_utils.ParseError(getErr)
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Db.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	// Really important when you work with MySQL
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error trying to save user: %s", err.Error()))
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr   {
	stmt, err := users_db.Db.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id) // It's an update so don't need a result back. Just error.
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Delete() *errors.RestErr  {
	stmt, err := users_db.Db.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Db.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close() // VERY IMPORTANT!

	results := make([]User, 0) // initialize
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user) // add user to array
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}