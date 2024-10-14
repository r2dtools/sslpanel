package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	prodMode = false
	logFile  = "var/log/ams.log"
	logLevel = 4
)

var config *Config

// Config is main config of the application
type Config struct {
	DatabaseURI,
	ServerAddress,
	ExecutablePath,
	LogFile,
	SecretKey string
	AgentPort        int
	LogLevel         int
	AmsEmailAddress  string
	AmsEmailPassword string
	SMTPHost         string
	SMTPPort         int
	AgentVersion     string
	InstallerVersion string
}

// GetLoggerFileAbsPath returns absolute path to logger file
func (c *Config) GetLoggerFileAbsPath() string {
	return filepath.Join(c.ExecutablePath, c.LogFile)
}

// GetVarDirAbsPath returns absolute path to var directory
func (c *Config) GetVarDirAbsPath() string {
	return filepath.Join(c.ExecutablePath, "var")
}

// GetConfig returns application config
func GetConfig() *Config {
	if config != nil {
		return config
	}

	var executablePath string
	var err error

	if prodMode {
		executable, err := os.Executable()

		if err != nil {
			panic(err)
		}

		executablePath = filepath.Dir(executable)
	} else {
		executablePath, err = os.Getwd()

		if err != nil {
			panic(err)
		}
	}

	configPath := filepath.Join(executablePath, "config")
	vConfig := viper.New()
	vConfig.SetDefault("LogFile", logFile)
	vConfig.SetDefault("LogLevel", logLevel)
	vConfig.SetConfigType("yaml")
	vConfig.SetConfigName("params")
	vConfig.AddConfigPath(configPath)
	viper.AutomaticEnv()

	if err := vConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	config = &Config{
		DatabaseURI:      vConfig.GetString("DatabaseURI"),
		SecretKey:        vConfig.GetString("SecretKey"),
		ServerAddress:    vConfig.GetString("ServerAddress"),
		AgentPort:        vConfig.GetInt("AgentPort"),
		LogLevel:         vConfig.GetInt("LogLevel"),
		LogFile:          vConfig.GetString("LogFile"),
		AmsEmailAddress:  vConfig.GetString("AmsEmailAddress"),
		AmsEmailPassword: vConfig.GetString("AmsEmailPassword"),
		SMTPHost:         vConfig.GetString("SMTPHost"),
		SMTPPort:         vConfig.GetInt("SMTPPort"),
		AgentVersion:     vConfig.GetString("AgentVersion"),
		InstallerVersion: vConfig.GetString("InstallerVersion"),
		ExecutablePath:   executablePath,
	}

	return config
}
