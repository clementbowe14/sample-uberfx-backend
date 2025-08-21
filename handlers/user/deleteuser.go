package user

import (
	"encoding/json"
	"github.com/clementbowe14/sample-uberfx-backend/db"
	"go.uber.org/zap"
	"net/http"
)

type DeleteUserHandler struct {
	log *zap.SugaredLogger
	db  *db.UserDb
}
type DeleteUserRequest struct {
	UserId int32 `json:"user_id"`
}

func NewDeleteUserHandler(log *zap.SugaredLogger, db *db.UserDb) *DeleteUserHandler {
	return &DeleteUserHandler{
		log: log,
		db:  db,
	}
}

func (h *DeleteUserHandler) Pattern() string { return "DELETE /users/{user_id}" }

func (h *DeleteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.log.Infof("%s %s", r.Method, r.URL.Path)

	h.log.Info("checking if request contains a valid user")

	deleteUserRequest := &DeleteUserRequest{}
	err := json.NewDecoder(r.Body).Decode(deleteUserRequest)
	if err != nil {
		h.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if deleteUserRequest.UserId <= 0 {
		h.log.Warn("user id is zero")
		w.WriteHeader(http.StatusBadRequest)
	}

	h.log.Infof("deleting user %d", deleteUserRequest.UserId)

	err = h.db.DeleteUser(deleteUserRequest.UserId)

	if err != nil {
		h.log.Errorw("Error removing user", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.log.Infof("user deleted %d", deleteUserRequest.UserId)
	w.WriteHeader(http.StatusNoContent)
}
