package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"open-ecm/apps"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Controller struct {
	db  *sql.DB
	srv service
}

func NewController(db *sql.DB) Controller {
	return Controller{db: db, srv: newService(db)}
}

func (c Controller) HandleSaveUser(w http.ResponseWriter, r *http.Request) {
	logger := apps.LoggerFromCtx(r.Context())

	var newUser User
	userJSONDecoder := json.NewDecoder(r.Body)
	userJSONDecoder.Decode(&newUser)
	logger.Debug("new user", zap.String("new_user", fmt.Sprintf("%#v", newUser)))

	savedUser, err := c.srv.saveUser(r.Context(), newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}
	logger.Debug("saved user", zap.String("save_user", fmt.Sprintf("%#v", savedUser)))
	w.WriteHeader(http.StatusCreated)
	returnedUser, err := json.Marshal(savedUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}
	w.Write(returnedUser)
}

func (c Controller) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	logger := apps.LoggerFromCtx(r.Context())
	userId := chi.URLParam(r, "userId")
	logger.Debug("User Id Param", zap.String("id", userId))

	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	user, err := c.srv.getUser(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (c Controller) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	logger := apps.LoggerFromCtx(r.Context())
	users, err := c.srv.getUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(usersJSON)
}

func (c Controller) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	logger := apps.LoggerFromCtx(r.Context())

	idParam := chi.URLParam(r, "userId")
	userId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	deletedUser, err := c.srv.deleteUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	deletedUserResponse, err := json.Marshal(deletedUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(deletedUserResponse)
}
