package controllers

import (
	"backend/logger"
	"backend/models"
	"backend/notification"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type registerData struct {
	Email, Password string
}

type confirmEmailData struct {
	Code   string
	UserID uint
}

// AuthController controller
type AuthController struct {
	BaseController
}

// Me returns authorized user
func (ctrl AuthController) Me(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Register a new user
func (ctrl AuthController) Register(c *gin.Context) {
	var account *models.Account
	var rData registerData
	var err error

	if err = c.ShouldBindJSON(&rData); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	user, err := models.GetUserByEmail(rData.Email)

	if err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	if user == nil {
		// create a new user with account
		account = new(models.Account)

		if err := models.SaveAccount(account); err != nil {
			ctrl.AbortWithBadRequestError(c, fmt.Errorf("could not create an account: %v", err))
			return
		}

		user = new(models.User)
		user.AccountID = account.ID
		user.AccountOwner = 1
		user.Email = rData.Email
		user.SetPassword(rData.Password)
		user.Account = *account

		if err = models.SaveUser(user); err != nil {
			ctrl.AbortWithBadRequestError(c, fmt.Errorf("could not create user: %v", err))
			return
		}
	} else {
		account, err = user.GetAccount()

		if err != nil {
			ctrl.AbortWithBadRequestError(c, err)
			return
		}

		// user with such email already exists
		if !user.IsAccountOwner() || account.IsConfirmed() {
			ctrl.AbortWithBadRequestError(c, fmt.Errorf("account with the email '%s' already exists", rData.Email))
			return
		}

		// update confirmation code
		account.GenerateConfirmationCode()

		if err = models.SaveAccount(account); err != nil {
			ctrl.AbortWithBadRequestError(c, fmt.Errorf("could not save an account: %v", err))
			return
		}
	}

	// send confirmation code to provided email
	eNotification := &notification.EmailNotification{}
	data := struct{ Code uint }{account.ConfirmationCode}

	if err = eNotification.CreateAndSendPlainNotification("confirmEmail", "signup-confirm-email-template", user.Email, "Email confirmation", data); err != nil {
		logger.Error(err.Error())
	}

	ctrl.SuccessJSON(c, gin.H{"user": user})
}

// ConfirmEmail a new user
func (ctrl AuthController) ConfirmEmail(c *gin.Context) {
	var data confirmEmailData
	var err error

	if err = c.ShouldBindJSON(&data); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	code, err := strconv.Atoi(data.Code)

	if err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	user, err := models.GetUserByID(data.UserID)

	if err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	account, err := user.GetAccount()

	if err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	if !account.VerifyCode(uint(code)) {
		ctrl.AbortWithBadRequestError(c, errors.New("confirmation code is invalid"))
		return
	}

	user.Active = 1

	if err = models.SaveUser(user); err != nil {
		ctrl.AbortWithBadRequestError(c, fmt.Errorf("could not save user: %v", err))
		return
	}

	account.Confirmed = 1

	if err = models.SaveAccount(account); err != nil {
		ctrl.AbortWithBadRequestError(c, fmt.Errorf("could not save an account: %v", err))
		return
	}
}

func (ctrl AuthController) RecoverPassword(c *gin.Context) {
	recoverData := struct{ Email string }{}
	var err error

	if err = c.ShouldBindJSON(&recoverData); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	user, err := models.GetUserByEmail(recoverData.Email)

	if err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	if user == nil {
		ctrl.AbortWithBadRequestError(c, fmt.Errorf("user with email %s does not exist", recoverData.Email))
		return
	}

	// update confirmation code
	user.GenerateConfirmationCode()

	if err = models.SaveUser(user); err != nil {
		ctrl.AbortWithBadRequestError(c, fmt.Errorf("could not save user: %v", err))
		return
	}

	// send confirmation code to provided email
	eNotification := &notification.EmailNotification{}
	data := struct{ Code uint }{user.ConfirmationCode}

	if err = eNotification.CreateAndSendPlainNotification("confirmEmail", "recover-confirm-email-template", user.Email, "Email confirmation", data); err != nil {
		logger.Error(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// ConfirmEmail a new user
func (ctrl AuthController) ResetPassword(c *gin.Context) {
	var data confirmEmailData
	var err error

	if err = c.ShouldBindJSON(&data); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	code, err := strconv.Atoi(data.Code)

	if err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	user, err := models.GetUserByID(data.UserID)

	if err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	if !user.VerifyCode(uint(code)) {
		ctrl.AbortWithBadRequestError(c, errors.New("confirmation code is invalid"))
		return
	}

	password := user.GeneratePassword()

	if err = models.SaveUser(user); err != nil {
		ctrl.AbortWithBadRequestError(c, fmt.Errorf("could not reset user password: %v", err))
		return
	}

	// send confirmation code to provided email
	eNotification := &notification.EmailNotification{}
	tplData := struct{ Password string }{password}

	if err = eNotification.CreateAndSendPlainNotification("passwordReset", "reset-password-email-template", user.Email, "Password reset", tplData); err != nil {
		logger.Error(err.Error())
	}
}
