package user

import (
	"backend/internal/app/panel/adapters/api/auth"
	userService "backend/internal/app/panel/user/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGetUserByIdHandler(uService userService.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid user ID")) // nolint:errcheck

			return
		}

		user, err := uService.FindUserByID(id)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user not found"})

			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func CreateGetUserByEmailHandler(uService userService.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Param("email")

		if email == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid user email")) // nolint:errcheck

			return
		}

		user, err := uService.FindUserByEmail(email)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user not found"})

			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func CreateChangePasswordHandler(appAuth auth.Auth, uService userService.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := appAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		data := struct{ Password, NewPassword string }{}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		err := uService.ChangePassword(user.ID, data.Password, data.NewPassword)

		if err != nil {
			if errors.Is(err, userService.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else if errors.Is(err, userService.ErrInvalidPassword) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
		}
	}
}
