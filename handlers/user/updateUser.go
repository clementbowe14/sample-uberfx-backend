package user

import (
	"encoding/json"
	"github.com/clementbowe14/sample-uberfx-backend/db"
	"go.uber.org/zap"
	"net/http"
)

type UpdateUserHandler struct {
	log *zap.SugaredLogger
	db  *db.UserDb
}

type UpdateUserRequest struct {
	UserId    int32  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func NewUpdateUserHandler(log *zap.SugaredLogger, db *db.UserDb) *UpdateUserHandler {
	return &UpdateUserHandler{
		log: log,
		db:  db,
	}
}

func (h *UpdateUserHandler) Pattern() string { return "PUT /users/{user_id}" }

func (h *UpdateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Infof("%s %s", r.Method, r.URL.Path)

	h.log.Infof("%s checking if request contains a valid user_id", r.Method)

	userUpdateRequest := &UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(userUpdateRequest)

	if err != nil {
		h.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if userUpdateRequest.UserId == 0 || userUpdateRequest.FirstName == "" || userUpdateRequest.LastName == "" || userUpdateRequest.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.db.UpdateUser(userUpdateRequest.UserId, userUpdateRequest.FirstName, userUpdateRequest.LastName, userUpdateRequest.Email)

	if err != nil {
		h.log.Errorw("Error updating user", "error", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userUpdateRequest)
	if err != nil {
		h.log.Errorw("Error encoding user", "error", err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
