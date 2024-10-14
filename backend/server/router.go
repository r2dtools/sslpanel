package server

import (
	"backend/controllers"
	"backend/middlware"
	"backend/modules"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// NewRouter creates new router for the application
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	authMiddleware := middlware.AuthMiddleware()

	v1 := router.Group("v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)
		auth := new(controllers.AuthController)
		v1.POST("/register", auth.Register)
		v1.POST("/confirm-email", auth.ConfirmEmail)
		v1.POST("/recover-password", auth.RecoverPassword)
		v1.POST("/reset-password", auth.ResetPassword)

		agentGroup := v1.Group("agent")
		{
			agent := new(controllers.AgentController)
			agentGroup.GET("/latest-version", agent.LatestVersion)
		}

		appGroup := v1.Group("app")
		{
			app := new(controllers.AppController)
			appGroup.Use(authMiddleware.MiddlewareFunc())
			appGroup.GET("/get-data", app.GetData)
		}

		userGroup := v1.Group("users")
		{
			user := new(controllers.UserController)
			userGroup.Use(authMiddleware.MiddlewareFunc())
			userGroup.GET("/:id", user.GetByID)
		}

		accountGroup := v1.Group("accounts")
		{
			account := new(controllers.AccountController)
			accountGroup.Use(authMiddleware.MiddlewareFunc())
			accountGroup.GET("/:id", account.GetByID)
		}

		authGroup := v1.Group("auth")
		{
			authGroup.Use(authMiddleware.MiddlewareFunc())
			authGroup.POST("/me", auth.Me)
		}

		serverGroup := v1.Group("servers")
		{
			server := new(controllers.ServerController)
			serverGroup.Use(authMiddleware.MiddlewareFunc())
			serverGroup.GET("", server.GetServers)
			serverGroup.GET("/:serverId", server.GetServer)
			serverGroup.POST("", server.AddServer)
			serverGroup.POST("/:serverId", server.SaveServer)
			serverGroup.DELETE("/:serverId", server.RemoveServer)
			serverGroup.POST("/:serverId/refresh", server.RefreshServer)
			serverGroup.GET("/:serverId/vhosts", server.GetServerVhosts)
			serverGroup.GET("/:serverId/vhost-certificate", server.GetVhostCertificate)
		}

		settingGroup := v1.Group("settings")
		{
			server := new(controllers.SettingController)
			settingGroup.Use(authMiddleware.MiddlewareFunc())
			settingGroup.POST("/change-password", server.ChangePassword)
		}

		modulesGroup := v1.Group("modules")
		{
			modules.InitModulesRouter(modulesGroup, authMiddleware)
		}
	}

	return router
}
