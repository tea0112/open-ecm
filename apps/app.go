package apps

import (
	"fmt"

	"github.com/spf13/viper"
)

type PortgresConfig struct {
	Username string
	Password string
	Database string
	Host     string
	Port     int
}

type AppConfig struct {
	Port     int
	Mode     string
	Postgres PortgresConfig
}

func NewAppConfig() AppConfig {
	var appConfig AppConfig

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath("$HOME/.openecm")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	appConfig.Postgres.Username = viper.GetString("POSTGRES_USERNAME")
	appConfig.Postgres.Password = viper.GetString("POSTGRES_PASSWORD")
	appConfig.Postgres.Database = viper.GetString("POSTGRES_DATABASE")
	appConfig.Postgres.Host = viper.GetString("POSTGRES_HOST")
	appConfig.Postgres.Port = viper.GetInt("POSTGRES_PORT")

	appConfig.Port = viper.GetInt("APP_PORT")

	return appConfig
}
