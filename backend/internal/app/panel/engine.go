package panel

import (
	"backend/config"
	accountService "backend/internal/app/panel/account/service"
	accountStorage "backend/internal/app/panel/account/storage"
	accountApi "backend/internal/app/panel/adapters/api/account"
	authApi "backend/internal/app/panel/adapters/api/auth"
	domainApi "backend/internal/app/panel/adapters/api/domain"
	serverApi "backend/internal/app/panel/adapters/api/server"
	userApi "backend/internal/app/panel/adapters/api/user"
	authAccount "backend/internal/app/panel/auth/account"
	authService "backend/internal/app/panel/auth/service"
	domainProvider "backend/internal/app/panel/domain/provider"
	domainService "backend/internal/app/panel/domain/service"
	domainStorage "backend/internal/app/panel/domain/storage"
	serverService "backend/internal/app/panel/server/service"
	serverStorage "backend/internal/app/panel/server/storage"
	userService "backend/internal/app/panel/user/service"
	userStorage "backend/internal/app/panel/user/storage"
	"backend/internal/modules"
	"backend/internal/pkg/logger"
	"backend/internal/pkg/notification"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newEngine(
	config *config.Config,
	logger logger.Logger,
	database *gorm.DB,
) (*gin.Engine, error) {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowAllOrigins = true
	engine.Use(cors.New(corsConfig))

	emailNotifocation := notification.EmailNotificationService{
		Config: config,
	}

	appAccountStorage := accountStorage.NewAccountSqlStorage(database)
	appAccountService := accountService.NewAccountService(appAccountStorage)

	appAccountCreator := authAccount.NewAccountCreator(database)

	appUserStorage := userStorage.NewUserSqlStorage(database)
	appUserService := userService.NewUserService(appUserStorage)

	appServerStorage := serverStorage.NewServerSqlStorage(database)
	appServerSevice := serverService.NewServerService(config, appServerStorage, logger)

	appDomainProvider := domainProvider.CreateDomainProvider(appServerStorage, logger)

	appDomainSettingStorage := domainStorage.NewDomainSettingSqlStorage(database)
	appDomainSevice := domainService.NewDomainService(config, appDomainSettingStorage, appServerStorage, appDomainProvider, logger)

	authMiddleware := authApi.AuthMiddleware(config, appUserStorage)

	appAuthService := authService.NewAuthService(
		config,
		appUserStorage,
		appAccountStorage,
		appAccountCreator,
		emailNotifocation,
		logger,
	)

	appAuth := authApi.NewAuth(appUserService)

	v1 := engine.Group("v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.POST("/register", authApi.CreateRegisterHandler(appAuthService))
		v1.POST("/confirm", authApi.CreateConfirmEmailHandler(appAuthService))
		v1.POST("/recover", authApi.CreateRecoverPasswordHandler(appAuthService))
		v1.POST("/reset", authApi.CreateResetPasswordHandler(appAuthService))

		userGroup := v1.Group("users")
		{
			userGroup.Use(authMiddleware.MiddlewareFunc())
			userGroup.GET("/:id", userApi.CreateGetUserByIdHandler(appUserService))
		}

		accountGroup := v1.Group("accounts")
		{
			accountGroup.Use(authMiddleware.MiddlewareFunc())
			accountGroup.GET("/:id", accountApi.CreateGetAccountByIdHandler(appAccountService))
		}

		authGroup := v1.Group("auth")
		{
			authGroup.Use(authMiddleware.MiddlewareFunc())
			authGroup.POST("/me", authApi.CreateMeHandler(appAuth))
		}

		serverGroup := v1.Group("servers")
		{
			serverGroup.Use(authMiddleware.MiddlewareFunc())
			serverGroup.GET("", serverApi.CreateFindAccounServersHandler(appAuth, appServerSevice))
			serverGroup.POST("", serverApi.CreateAddServerHandler(appAuth, appServerSevice))
			serverGroup.POST("/:serverId", serverApi.CreateUpdateServerHandler(appServerSevice))
			serverGroup.DELETE("/:serverId", serverApi.CreateRemoveServerHandler(appAuth, appServerSevice))
			serverGroup.GET("/:serverId", serverApi.CreateGetServerByGuidHandler(appAuth, appServerSevice))
			serverGroup.GET("/:serverId/details", serverApi.CreateGetServerDetailsByGuidHandler(appAuth, appServerSevice))

			serverSettingGroup := serverGroup.Group("/:serverId/settings")
			{
				serverSettingGroup.POST("/certbot-status", serverApi.CreateChangeCertbotStatusHandler(appAuth, appServerSevice))
			}

			domainGroup := serverGroup.Group("/:serverId/domain/:domainName")
			{
				domainGroup.GET("", domainApi.CreateGetDomainHandler(appAuth, appDomainSevice))
				domainGroup.GET("/config", domainApi.CreateGetDomainConfigHandler(appAuth, appDomainSevice))
				domainGroup.GET("/settings", domainApi.CreateFindDomainSettingsHandler(appAuth, appDomainSevice))
				domainGroup.POST("/settings", domainApi.CreateChangeDomainSettingHandler(appAuth, appDomainSevice))
			}
		}

		settingGroup := v1.Group("settings")
		{
			settingGroup.Use(authMiddleware.MiddlewareFunc())
			settingGroup.POST("/change-password", userApi.CreateChangePasswordHandler(appAuth, appUserService))
		}

		modulesGroup := v1.Group("modules")
		{
			modules.InitModulesRouter(modulesGroup, database, authMiddleware, appAuth, appServerStorage, appDomainSettingStorage, logger)
		}
	}

	return engine, nil
}
