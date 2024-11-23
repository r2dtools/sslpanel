package panel

import (
	"backend/config"
	"backend/internal/pkg/db"
	"backend/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	config *config.Config
	logger logger.Logger
	engine *gin.Engine
	db     *gorm.DB
}

func (app *App) Run() error {
	return app.engine.Run(app.config.ServerHost)
}

func GetApp(config *config.Config, logger logger.Logger) (*App, error) {
	database, err := db.GetDB(config)

	if err != nil {
		return nil, err
	}

	engine, err := newEngine(config, logger)

	if err != nil {
		return nil, err
	}

	return &App{
		config: config,
		logger: logger,
		engine: engine,
		db:     database,
	}, nil
}
