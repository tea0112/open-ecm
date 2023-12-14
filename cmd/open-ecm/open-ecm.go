package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"open-ecm/users"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

func main() {
	var appConfig AppConfig

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath("$HOME/.openecm")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	appConfig.Postgres.Username = viper.GetString("POSTGRES_USERNAME")
	appConfig.Postgres.Password = viper.GetString("POSTGRES_PASSWORD")
	appConfig.Postgres.Database = viper.GetString("POSTGRES_DATABASE")
	appConfig.Postgres.Host = viper.GetString("POSTGRES_HOST")
	appConfig.Postgres.Port = viper.GetInt("POSTGRES_PORT")

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", appConfig.Postgres.Host, appConfig.Postgres.Port, appConfig.Postgres.Username, appConfig.Postgres.Password, appConfig.Postgres.Database)
	sugar.Debug(connStr)

	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	appConfig.Port = viper.GetInt("APP_PORT")

	r := chi.NewRouter()

	v1ApiRouter := chi.NewRouter()

	v1ApiRouter.Route("/", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			userController := users.NewController(db)
			r.Post("/", userController.HandleCreateUser)
		})
	})

	r.Mount("/api/v1", v1ApiRouter)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This URL doesn't belong to the OpenECM"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	})

	http.ListenAndServe(fmt.Sprintf(":%d", appConfig.Port), r)
}
