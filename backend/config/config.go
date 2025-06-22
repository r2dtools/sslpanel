package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	developmentEnv = "development"
	productionEnv  = "production"
)

var config *Config

type Config struct {
	DatabaseURI      string
	PanelHost        string
	ServerAddress    string
	ServerHost       string
	LogFile          string
	SecretKey        string
	AgentPort        int
	LogLevel         int
	AmsEmailAddress  string
	AmsEmailPassword string
	SMTPHost         string
	SMTPPort         int
	DbName           string
	DbHost           string
	DbPort           string
	DbUser           string
	DbPassword       string
	DbType           string
	BaseDir          string
	DbDsn            string
	AllowedHosts     []string
	Environment      string
	IsDevMode        bool
}

func (c *Config) GetVarDirAbsPath() string {
	return filepath.Join(c.BaseDir, "var")
}

func GetConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	viper.AutomaticEnv()
	environment := viper.GetString("CP_ENVIRONMENT")

	if environment == "" {
		environment = developmentEnv
	}

	conf := Config{
		DbName:           viper.GetString("CP_DB_NAME"),
		PanelHost:        viper.GetString("CP_HOST"),
		DbHost:           viper.GetString("CP_DB_HOST"),
		DbPort:           viper.GetString("CP_DB_PORT"),
		DbUser:           viper.GetString("CP_DB_USER"),
		DbPassword:       viper.GetString("CP_DB_PASSWORD"),
		DbType:           viper.GetString("CP_DB_TYPE"),
		SecretKey:        viper.GetString("CP_SERVER_KEY"),
		ServerAddress:    viper.GetString("CP_API_HOST"),
		AgentPort:        viper.GetInt("CP_AGENT_PORT"),
		LogLevel:         viper.GetInt("CP_LOG_LEVEL"),
		LogFile:          viper.GetString("CP_LOG_FILE"),
		AmsEmailAddress:  viper.GetString("CP_EMAIL_ADDRESS"),
		AmsEmailPassword: viper.GetString("CP_EMAIL_PASSWORD"),
		SMTPHost:         viper.GetString("CP_SMTP_HOST"),
		SMTPPort:         viper.GetInt("CP_SMTP_PORT"),
		AllowedHosts:     getAllowedHosts(),
		Environment:      environment,
		IsDevMode:        environment == developmentEnv,
	}

	path, err := getBasePath(environment)

	if err != nil {
		return nil, fmt.Errorf("failed to get base path: %v", err)
	}

	conf.BaseDir = path

	dbDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		conf.DbUser,
		conf.DbPassword,
		conf.DbHost,
		conf.DbPort,
		conf.DbName,
	)
	conf.DbDsn = dbDsn

	dbUri := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		conf.DbHost,
		conf.DbPort,
		conf.DbUser,
		conf.DbName,
		conf.DbPassword,
	)
	conf.DatabaseURI = dbUri

	pUrl, err := url.Parse(conf.ServerAddress)

	if err != nil {
		return nil, err
	}

	conf.ServerHost = pUrl.Host

	config = &conf

	return config, nil
}

func getBasePath(environment string) (string, error) {
	var (
		path string
		err  error
	)

	if environment == productionEnv {
		path, err = os.Executable()

		if err != nil {
			return "", err
		}

		return filepath.Clean(filepath.Dir(path)), nil
	}

	path, err = os.Getwd()

	if err != nil {
		return "", err
	}

	return filepath.Clean(fmt.Sprintf("%s/../../", path)), nil
}

func getAllowedHosts() []string {
	hostsStr := viper.GetString("CP_API_ALLOWED_HOSTS")
	hosts := strings.Split(hostsStr, ",")

	result := []string{}

	for _, host := range hosts {
		result = append(result, strings.TrimSpace(host))
	}

	return result
}
