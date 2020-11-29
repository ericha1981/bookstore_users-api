package services

import (
	"github.com/ericha1981/bookstore_users-api/domain/users"
	"github.com/ericha1981/bookstore_users-api/utils/errors"
)

/*
	Entire business logic to create a user comes here
*/

func GetUser(userId int64) (*users.User, *errors.RestErr) {
	// TODO: Get the user from the database
	result := &users.User{
		Id: userId,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	// TODO: Save the new user to the database
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}