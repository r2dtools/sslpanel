package account

import (
	accountService "backend/internal/app/panel/account/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGetAccountByIdHandler(aService accountService.AccountService) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid account ID")) // nolint:errcheck

			return
		}

		account, err := aService.FindAccount(uint(id))

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err) // nolint:errcheck

			return
		}

		if account == nil {
			c.AbortWithError(http.StatusNotFound, errors.New("account not found")) // nolint:errcheck
		}

		c.JSON(http.StatusOK, gin.H{"account": account})
	}
}
