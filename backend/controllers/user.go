package controllers

import (
	"backend/models"
	"errors"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController controller
type UserController struct {
	BaseController
}

// GetByID returns user model by id
func (u UserController) GetByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		u.AbortWithBadRequestError(c, errors.New("user ID is missing"))

		return
	}

	numericID, err := strconv.ParseUint(id, 10, 0)

	if err != nil {
		u.AbortWithBadRequestError(c, err)

		return
	}

	user, err := models.GetUserByID(uint(numericID))

	if err != nil {
		u.AbortWithInternalServerError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetByEmail returns user model by id
func (u UserController) GetByEmail(c *gin.Context) {
	email, ok := c.GetQuery("email")

	if !ok || email == "" {
		u.AbortWithBadRequestError(c, errors.New("user email is missing"))
		return
	}

	user, err := models.GetUserByEmail(email)

	if err != nil {
		u.AbortWithInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
