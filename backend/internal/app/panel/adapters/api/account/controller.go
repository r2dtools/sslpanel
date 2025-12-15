package account

import (
	accountService "backend/internal/app/panel/account/service"
	"backend/internal/app/panel/adapters/api/auth"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGetAccountByIdHandler(cAuth auth.Auth, aService accountService.AccountService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid account ID")) // nolint:errcheck

			return
		}

		if user.AccountID != id {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Access denied"}) // nolint:errcheck

			return
		}

		account, err := aService.FindAccount(uint(id))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

			return
		}

		if account == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "account not found"})
		}

		c.JSON(http.StatusOK, gin.H{"account": account})
	}
}
