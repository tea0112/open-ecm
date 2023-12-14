package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Controller struct {
	db  *sql.DB
	srv service
}

func NewController(db *sql.DB) Controller {
	return Controller{db: db, srv: newService(db)}
}

func (c Controller) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	userJSONDecoder := json.NewDecoder(r.Body)
	userJSONDecoder.Decode(&newUser)

	savedUser, err := c.srv.createUser(r.Context(), newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusCreated)
		returnedUser, err := json.Marshal(savedUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
		w.Write(returnedUser)
	}
}
