package panel

import (
	"backend/config"
	"backend/internal/app/panel/domain/provider"
	domainStorage "backend/internal/app/panel/domain/storage"
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/modules/sslmanager/autorenewal"
	"backend/internal/modules/sslmanager/autorenewal/logwriter"
	"backend/internal/pkg/db"
	"backend/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	config               *config.Config
	logger               logger.Logger
	engine               *gin.Engine
	db                   *gorm.DB
	certRenewalScheduler autorenewal.Scheduler
}

func (app *App) Run() error {
	go app.certRenewalScheduler.Run()

	return app.engine.Run(app.config.ServerHost)
}

func GetApp(config *config.Config, logger logger.Logger) (*App, error) {
	database, err := db.GetDB(config)

	if err != nil {
		return nil, err
	}

	engine, err := newEngine(config, logger, database)

	if err != nil {
		return nil, err
	}

	appServerStorage := serverStorage.NewServerSqlStorage(database)
	appDomainSettingStorage := domainStorage.NewDomainSettingSqlStorage(database)
	domainProvider := provider.CreateDomainProvider(appServerStorage, logger)
	certRenewalManager := autorenewal.CreateAutoRenewalManager(
		appServerStorage,
		appDomainSettingStorage,
		domainProvider,
		config,
		logger,
		logwriter.CreateDefaultLogWriter(logger),
	)

	return &App{
		config:               config,
		logger:               logger,
		engine:               engine,
		db:                   database,
		certRenewalScheduler: autorenewal.CreateScheduler(config, logger, certRenewalManager),
	}, nil
}
