package certificatemonitor

import (
	"backend/certificate"
	"backend/controllers"
	"backend/modules/certificatemonitor/models"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/r2dtools/agentintegration"
)

// ModuleController controller for Certificate Monitor module
type ModuleController struct {
	controllers.BaseController
}

type siteData struct {
	Cert *agentintegration.Certificate
}

func (ctrl *ModuleController) getSites(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	sites, err := models.GetUserSites(user.ID)

	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)

		return
	}

	ctrl.SuccessJSON(c, gin.H{"sites": sites})
}

// createSite creates new site
func (ctrl *ModuleController) createSite(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	var site models.Site
	var err error

	if err = c.ShouldBindJSON(&site); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	site.UserID = user.ID
	err = models.SaveSite(&site)

	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}
}

func (ctrl *ModuleController) removeSite(c *gin.Context) {
	if !ctrl.CheckAuth(c) {
		return
	}

	site, err := ctrl.GetSiteByID(c)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	if err = site.Remove(); err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}
}

func (ctrl *ModuleController) refreshSite(c *gin.Context) {
	if !ctrl.CheckAuth(c) {
		return
	}

	site, err := ctrl.GetSiteByID(c)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	cert, err := certificate.GetCertificateForDomainFromRequest(site.URL)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, fmt.Errorf("could not get site certificate: %v", err))
		return
	}

	if err = site.UpdateCertData(cert); err != nil {
		ctrl.AbortWithInternalServerError(c, fmt.Errorf("could not update site data: %v", err))
		return
	}

	ctrl.SuccessJSON(c, gin.H{"site": site})
}

// GetSiteByID returns site model by ID in request
func (ctrl *ModuleController) GetSiteByID(c *gin.Context) (*models.Site, error) {
	id := c.Param("siteId")
	if id == "" {
		return nil, errors.New("site ID is missing")
	}

	nID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	// TODO: Check permissions
	site, err := models.GetSiteByID(nID)
	if err != nil {
		return nil, err
	}

	return site, nil
}
