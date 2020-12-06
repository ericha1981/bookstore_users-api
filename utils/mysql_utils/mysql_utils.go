package mysql_utils

import (
	"github.com/ericha1981/bookstore_users-api/utils/errors"
	"github.com/go-sql-driver/mysql"
	"strings"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok { // not able to convert
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError("no record matching for a given id.")
		}
		return errors.NewInternalServerError("error parsing database response")
	}

	switch sqlErr.Number {
	case 1062:
		return errors.NewNotFoundError("invalid data")
	}
	return errors.NewInternalServerError("error processing request")
}
