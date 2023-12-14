package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"open-ecm/api"
	"open-ecm/apps"

	_ "github.com/lib/pq"
)

func main() {
	appConfig := apps.NewAppConfig()

	logger := apps.NewLogger(appConfig.Mode)
	defer logger.Sync()

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", appConfig.Postgres.Host, appConfig.Postgres.Port, appConfig.Postgres.Username, appConfig.Postgres.Password, appConfig.Postgres.Database)
	logger.Debug(connStr)

	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	r := api.NewChiMainRouter(db, logger)

	http.ListenAndServe(fmt.Sprintf(":%d", appConfig.Port), r)
}
