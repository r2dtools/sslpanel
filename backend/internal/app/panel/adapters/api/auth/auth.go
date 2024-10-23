package auth

import (
	userService "backend/internal/app/panel/user/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	userService userService.UserService
}

func (a *Auth) GetCurrentUser(c *gin.Context) *userService.User {
	data, exists := c.Get(IdentityKey)

	if !exists {
		c.AbortWithError(http.StatusUnauthorized, ErrorUnauthorized)

		return nil
	}

	identity, ok := data.(*User)

	if !ok {
		c.AbortWithError(http.StatusUnauthorized, ErrorInvalidUserData)

		return nil
	}

	user, err := a.userService.FindUserByEmail(identity.Email)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return nil
	}

	if user == nil {
		c.AbortWithError(http.StatusNotFound, errors.New("user not found"))
	}

	return user
}

func (a *Auth) CheckAuth(c *gin.Context) bool {
	if user := a.GetCurrentUser(c); user == nil {
		return false
	}

	return true
}

func NewAuth(userService userService.UserService) Auth {
	return Auth{userService: userService}
}
