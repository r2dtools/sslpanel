package auth

import (
	authService "backend/internal/app/panel/auth/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerData struct {
	Email, Password string
}

func CreateMeHandler(appAuth Auth) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := appAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func CreateResetPasswordHandler(appAuthService authService.AuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Query("token")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})

			return
		}

		err := appAuthService.ResetPassword(token)

		if err != nil {
			if errors.Is(err, authService.ErrInvalidConfirmationToken) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else if errors.Is(err, authService.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func CreateRegisterHandler(appAuthService authService.AuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		err := appAuthService.Register(request.Email, request.Password)

		if err != nil {
			if errors.Is(err, authService.ErrAccountAlreadyExists) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func CreateConfirmEmailHandler(appAuthService authService.AuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := struct {
			Token string `json:"token"`
		}{}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		if request.Token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})

			return
		}

		err := appAuthService.ConfirmEmail(request.Token)

		if err != nil {
			if errors.Is(err, authService.ErrInvalidConfirmationToken) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else if errors.Is(err, authService.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func CreateRecoverPasswordHandler(appAuthService authService.AuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		recoverData := struct {
			Email string `json:"email"`
		}{}

		if err := c.ShouldBindJSON(&recoverData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		err := appAuthService.RecoverPassword(recoverData.Email)

		if err != nil {
			if errors.Is(err, authService.ErrUserNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else if errors.Is(err, authService.ErrUserNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}
