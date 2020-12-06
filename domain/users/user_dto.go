package users

import (
	"github.com/ericha1981/bookstore_users-api/utils/errors"
	"strings"
)

const (
	StatusActive = "active"
)

//====================================================================
// DTO - Data Transfer Object. Just an object that holds data.
// 		DTO will be passed as value object to DAO layer and DAO layer
//		will use this object to persist data using its CRUD operation
//		methods.
//====================================================================
type User struct {
	Id int64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	DateCreated string `json:"date_created"`
	Status string `json:"status"`
	Password string `json:"password"`
}

type Users []User

func (user *User) Validate() *errors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return errors.NewBadRequestError("invalid password")
	}
	return nil
}