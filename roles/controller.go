package roles

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"open-ecm/apps"
)

type Controller struct {
	db  *sql.DB
	srv service
}

func NewController(db *sql.DB) Controller {
	return Controller{db: db, srv: NewService(db)}
}

func (c Controller) saveRole(w http.ResponseWriter, r http.Request) {
	logger := apps.LoggerFromCtx(r.Context())

	var newRole Role
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newRole)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	savedRole, err := c.srv.saveRole(r.Context(), newRole)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var responseData bytes.Buffer
	encoder := json.NewEncoder(&responseData)
	err = encoder.Encode(savedRole)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseData.Bytes())
}

func (c Controller) getRole(w http.ResponseWriter, r *http.Request) {

}

func (c Controller) getRoles(w http.ResponseWriter, r *http.Request) {

}

func (c Controller) deleteRole(w http.ResponseWriter, r *http.Request) {

}
