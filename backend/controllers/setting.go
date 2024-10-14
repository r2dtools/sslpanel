package controllers

import (
	"backend/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// SettingController is acontroller to manage account settings
type SettingController struct {
	BaseController
}

// ChangePassword changes user`s password
func (ctrl SettingController) ChangePassword(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	data := struct{ Password, NewPassword string }{}
	var err error

	if err = c.ShouldBindJSON(&data); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	if !user.CheckPassword(data.Password) {
		ctrl.AbortWithPermissionDeniedError(c, errors.New("the current password is invalid"))
		return
	}

	if err = user.SetPassword(data.NewPassword); err != nil {
		ctrl.AbortWithInternalServerError(c, fmt.Errorf("could not change password: %v", err))
		return
	}

	if err = models.SaveUser(user); err != nil {
		ctrl.AbortWithBadRequestError(c, fmt.Errorf("could not change password: %v", err))
		return
	}
}
