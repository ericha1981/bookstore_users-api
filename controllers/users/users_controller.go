package users

import (
	"github.com/ericha1981/bookstore_users-api/domain/users"
	"github.com/ericha1981/bookstore_users-api/services"
	"github.com/ericha1981/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/*
	Handle all requests here (entry point)
 */

func CreateUser(c *gin.Context) {
	var user users.User

	/*
		// Below logic can be replaced by c.ShouldBindJSON function!!
		bytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return
		}

		if err := json.Unmarshal(bytes, &user); err != nil {
			fmt.Println(err.Error())
			return
		}
	*/

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.CreateUser(user)
	if saveErr != nil { // non pointer type throws an error here: Cannot convert 'nil' to type 'errors.RestErr'.
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Get user from the DB
	user, getErr := services.GetUser(userId)
	if getErr != nil { // non pointer type throws an error here: Cannot convert 'nil' to type 'errors.RestErr'.
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusCreated, user)
}
