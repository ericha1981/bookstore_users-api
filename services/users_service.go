package services

import (
	"github.com/ericha1981/bookstore_users-api/domain/users"
	"github.com/ericha1981/bookstore_users-api/utils/crypto_utils"
	"github.com/ericha1981/bookstore_users-api/utils/date_utils"
	"github.com/ericha1981/bookstore_users-api/utils/errors"
)

/*
	Entire business layer logic to create a user comes here.
	Interacts with external APIS, database etc.
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

	// Set the date created.
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Status = users.StatusActive
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	currUser, err := GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	// Validate the user
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if isPartial { // PATCH method
		if user.FirstName != "" {
			currUser.FirstName = user.FirstName
		}
		if user.LastName != "" {
			currUser.LastName = user.LastName
		}
		if user.Email != "" {
			currUser.Email = user.Email
		}
	} else { // PUT method
		currUser.FirstName = user.FirstName
		currUser.LastName = user.LastName
		currUser.Email = user.Email
	}

	if err := currUser.Update(); err != nil {
		return nil, err
	}
	return currUser, nil
}

func DeleteUser(userId int64) *errors.RestErr {
	currUser := &users.User{ Id: userId }
	return currUser.Delete()
}

func Search(status string) (users.Users, *errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}