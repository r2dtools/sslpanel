package auth

import (
	authService "backend/internal/app/panel/auth/service"
	"backend/internal/pkg/notification"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errSendEmailMessage = "Failed to send email. Check your SMTP server settings"

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
		request := struct {
			Token string `json:"token"`
		}{}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		if request.Token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})

			return
		}

		err := appAuthService.ResetPassword(request.Token)
		var errSendEmail notification.ErrSendEmail

		if err != nil {
			if errors.Is(err, authService.ErrInvalidConfirmationToken) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else if errors.Is(err, authService.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else if errors.As(err, &errSendEmail) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("%s: %v", errSendEmailMessage, err)})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err) // nolint:errcheck
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
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		err := appAuthService.Register(request.Email, request.Password)
		var errSendEmail notification.ErrSendEmail

		if err != nil {
			if errors.Is(err, authService.ErrAccountAlreadyExists) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else if errors.As(err, &errSendEmail) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("%s: %v", errSendEmailMessage, err)})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err) // nolint:errcheck
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
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

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
				c.AbortWithError(http.StatusInternalServerError, err) // nolint:errcheck
			}
		}
	}
}

func CreateRecoverPasswordHandler(appAuthService authService.AuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := struct {
			Email string `json:"email"`
		}{}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		if request.Email == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid email"})

			return
		}

		err := appAuthService.RecoverPassword(request.Email)
		var errSendEmail notification.ErrSendEmail

		if err != nil {
			if errors.Is(err, authService.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else if errors.As(err, &errSendEmail) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("%s: %v", errSendEmailMessage, err)})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err) // nolint:errcheck
			}
		}
	}
}
