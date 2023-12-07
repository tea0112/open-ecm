package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.openecm")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	postgresUser := viper.GetString("postgres.username")
	postgresPasswd := viper.GetString("postgres.password")
	postgresDb := viper.GetString("postgres.database")
	postgresHost := viper.GetString("postgres.host")
	postgresPort := viper.GetInt32("postgres.port")
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresPort, postgresUser, postgresPasswd, postgresDb)

	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	appPort := viper.GetString("app.port")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServe(fmt.Sprintf(":%s", appPort), r)
}
