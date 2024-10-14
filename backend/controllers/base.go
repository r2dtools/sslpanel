package controllers

import (
	"backend/middlware"
	"backend/models"
	"backend/remote"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BaseController is a controller with base functionality
type BaseController struct{}

func (base *BaseController) abortWithError(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, gin.H{"message": err.Error()})
}

// AbortWithBadRequestError interupt request and return response with error code 400
func (base *BaseController) AbortWithBadRequestError(c *gin.Context, err error) {
	base.abortWithError(c, http.StatusBadRequest, fmt.Errorf("bad request: %w", err))
}

// AbortWithPermissionDeniedError interupt request and return response with error code 403
func (base *BaseController) AbortWithPermissionDeniedError(c *gin.Context, err error) {
	base.abortWithError(c, http.StatusForbidden, fmt.Errorf("permission denied: %w", err))
}

// AbortWithInternalServerError interupt request and return response with error code 500
func (base *BaseController) AbortWithInternalServerError(c *gin.Context, err error) {
	base.abortWithError(c, http.StatusInternalServerError, err)
}

func (base *BaseController) abortWithUnauthorizedError(c *gin.Context, err error) {
	base.abortWithError(c, http.StatusUnauthorized, err)
}

// SuccessJSON returns success response as json
func (base *BaseController) SuccessJSON(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, data)
}

// GetCurrentUser return curren authenticated user
func (base *BaseController) GetCurrentUser(c *gin.Context) *models.User {
	data, exists := c.Get(middlware.IdentityKey)
	if !exists {
		base.abortWithUnauthorizedError(c, middlware.ErrorUnauthorized)
		return nil
	}

	identity, ok := data.(*middlware.User)
	if !ok {
		base.abortWithUnauthorizedError(c, middlware.ErrorInvalidUserData)
		return nil
	}

	user, err := models.GetUserByEmail(identity.Email)
	if err != nil {
		base.AbortWithInternalServerError(c, err)
		return nil
	}

	return user
}

// CheckAuth checks if user is authenticated
func (base *BaseController) CheckAuth(c *gin.Context) bool {
	if user := base.GetCurrentUser(c); user == nil {
		return false
	}

	return true
}

// GetServerByID returns server model by ID in request
func (base *BaseController) GetServerByID(c *gin.Context, user *models.User) (*models.Server, error) {
	id := c.Param("serverId")
	if id == "" {
		return nil, errors.New("server ID is missed")
	}

	nID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	server, err := models.GetAccountServerByID(nID, user.AccountID)
	if err != nil {
		return nil, err
	}
	if server == nil {
		return nil, errors.New("server not found")
	}

	return server, nil
}

// GetServerAgent returns server agent by server ID in request
func (base *BaseController) GetServerAgent(c *gin.Context, user *models.User) (*remote.Agent, error) {
	server, err := base.GetServerByID(c, user)
	if err != nil {
		return nil, err
	}

	agent := &remote.Agent{
		Server: server,
	}

	return agent, nil
}
