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

var (
	UserService userServiceInterface = &userService{}
)

type userService struct {}

type userServiceInterface interface {
	GetUser(int64) (*users.User, *errors.RestErr)
	CreateUser(users.User) (*users.User, *errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	SearchUser(string) (users.Users, *errors.RestErr)
	LoginUser(users.LoginRequest) (*users.User, *errors.RestErr)
}

func (s *userService) GetUser(userId int64) (*users.User, *errors.RestErr) {
	// TODO: Get the user from the database
	result := &users.User{
		Id: userId,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *userService) CreateUser(user users.User) (*users.User, *errors.RestErr) {
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

func (s *userService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	currUser, err := s.GetUser(user.Id)
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

func (s *userService) DeleteUser(userId int64) *errors.RestErr {
	currUser := &users.User{ Id: userId }
	return currUser.Delete()
}

func (s *userService) SearchUser(status string) (users.Users, *errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func (s *userService) LoginUser(request users.LoginRequest) (*users.User, *errors.RestErr) {
	dao := &users.User{
		Email: request.Email,
		Password: crypto_utils.GetMd5(request.Password),
	}
	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}