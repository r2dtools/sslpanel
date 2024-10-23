package auth

import (
	authService "backend/internal/app/panel/auth/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type confirmEmailData struct {
	Code   string
	UserID int
}

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
		var data confirmEmailData

		if err := c.ShouldBindJSON(&data); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		code, err := strconv.Atoi(data.Code)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		err = appAuthService.ResetPassword(data.UserID, code)

		if err != nil {
			if errors.Is(err, authService.ErrInvalidConfirmationCode) {
				c.AbortWithError(http.StatusBadRequest, err)
			} else if errors.Is(err, authService.ErrUserNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

	}
}

func CreateRegisterHandler(appAuthService authService.AuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		rData := struct {
			Email, Password string
		}{}

		if err := c.ShouldBindJSON(&rData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		err := appAuthService.Register(rData.Email, rData.Password)

		if err != nil {
			if errors.Is(err, authService.ErrAccountAlreadyExists) {
				c.AbortWithError(http.StatusBadRequest, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func CreateConfirmEmailHandler(appAuthService authService.AuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data confirmEmailData

		if err := c.ShouldBindJSON(&data); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		code, err := strconv.Atoi(data.Code)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		err = appAuthService.ConfirmEmail(data.UserID, code)

		if err != nil {
			if errors.Is(err, authService.ErrInvalidConfirmationCode) {
				c.AbortWithError(http.StatusBadRequest, err)
			} else if errors.Is(err, authService.ErrUserNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func CreateRecoverPAsswordHandler(appAuthService authService.AuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		recoverData := struct{ Email string }{}

		if err := c.ShouldBindJSON(&recoverData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		user, err := appAuthService.RecoverPassword(recoverData.Email)

		if err != nil {
			if errors.Is(err, authService.ErrUserNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else if errors.Is(err, authService.ErrUserNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
