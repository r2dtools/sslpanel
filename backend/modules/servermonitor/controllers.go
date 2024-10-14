package servermonitor

import (
	"backend/controllers"
	"backend/models"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/r2dtools/agentintegration"
)

const DEFAULT_TIME_INTERVAL = 24 // 24 hours
const DISK_CATEGORY = "disk"
const NETWORK_CATEGORY = "network"
const PROCESS_CATEGORY = "process"

// ModuleController controller for Server Monitor module
type ModuleController struct {
	controllers.BaseController
}

func (ctrl *ModuleController) getCategoryData(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	var err error
	var toTimeInt, fromTimeInt int

	category, ok := c.GetQuery("category")
	if !ok {
		ctrl.AbortWithBadRequestError(c, errors.New("category parameter is missed"+category))
		return
	}

	fromTime, ok := c.GetQuery("fromTime")
	// if fromTime is not specified use current time - 12 hours by default
	if !ok {
		interval := -time.Duration(DEFAULT_TIME_INTERVAL) * time.Hour
		fromTimeInt = int(time.Now().Add(interval).Unix())
	} else {
		fromTimeInt, err = strconv.Atoi(fromTime)
		if err != nil {
			ctrl.AbortWithBadRequestError(c, errors.New("fromTime parameter is invalid"))
			return
		}
	}

	toTime, ok := c.GetQuery("toTime")
	// use current time toTime is not specified
	if !ok {
		toTimeInt = int(time.Now().Unix())
	} else {
		toTimeInt, err = strconv.Atoi(toTime)
		if err != nil {
			ctrl.AbortWithBadRequestError(c, errors.New("toTime parameter is invalid"))
			return
		}
	}

	agent, err := ctrl.getAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	requestData := agentintegration.ServerMonitorStatisticsRequestData{Category: category, FromTime: fromTimeInt, ToTime: toTimeInt}
	var data interface{}
	if category == DISK_CATEGORY {
		data, err = agent.loadDiskStatisticsData(&requestData)
	} else if category == NETWORK_CATEGORY {
		data, err = agent.loadNetworkStatisticsData(&requestData)
	} else if category == PROCESS_CATEGORY {
		data, err = agent.loadProcessStatisticsData(&requestData)
	} else {
		data, err = agent.loadStatisticsData(&requestData)
	}

	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"data": data})
}

func (ctrl *ModuleController) getAgent(c *gin.Context, user *models.User) (*agent, error) {
	serverAgent, err := ctrl.GetServerAgent(c, user)
	if err != nil {
		return nil, err
	}

	return &agent{serverAgent}, nil
}
