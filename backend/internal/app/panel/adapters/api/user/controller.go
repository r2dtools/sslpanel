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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid user ID"))

			return
		}

		user, err := uService.FindUserByID(id)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		if user == nil {
			c.AbortWithError(http.StatusNotFound, errors.New("user not found"))

			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func CreateGetUserByEmailHandler(uService userService.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Param("email")

		if email == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid user email"))

			return
		}

		user, err := uService.FindUserByEmail(email)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		if user == nil {
			c.AbortWithError(http.StatusNotFound, errors.New("user not found"))

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
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		err := uService.ChangePassword(user.ID, data.Password, data.NewPassword)

		if err != nil {
			if errors.Is(err, userService.ErrUserNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else if errors.Is(err, userService.ErrInvalidPassword) {
				c.AbortWithError(http.StatusBadRequest, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}
