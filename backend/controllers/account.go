package controllers

import (
	"backend/models"
	"errors"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AccountController is a controller to manage accounts
type AccountController struct {
	BaseController
}

var accountModel = new(models.Account)

// GetByID returns user model by id
func (u AccountController) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		u.AbortWithBadRequestError(c, errors.New("account ID is missing"))
		return
	}

	numericID, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		u.AbortWithBadRequestError(c, err)
		return
	}

	account, err := accountModel.GetByID(int(numericID))
	if err != nil {
		u.AbortWithInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}
