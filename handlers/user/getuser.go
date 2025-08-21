package user

import (
	"encoding/json"
	"github.com/clementbowe14/sample-uberfx-backend/db"
	"go.uber.org/zap"
	"net/http"
)

type GetUserHandler struct {
	log *zap.SugaredLogger
	db  *db.UserDb
}

type GetUserRequest struct {
	UserId int32 `json:"user_id"`
}

func NewGetUserHandler(logger *zap.SugaredLogger, db *db.UserDb) *GetUserHandler {
	return &GetUserHandler{
		log: logger,
		db:  db,
	}
}

func (h *GetUserHandler) Pattern() string {
	return "GET /users/{user_id}"
}

func (h *GetUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Infof("%s %s", r.Method, r.URL.Path)

	h.log.Infof("%s checking if request contains a valid user_id", r.Method)

	userRequest := &GetUserRequest{}
	err := json.NewDecoder(r.Body).Decode(userRequest)

	if err != nil {
		h.log.Errorw("Error converting user_id to int", "error", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	h.log.Infof("%s checking if request contains userReq", userRequest.UserId)

	if userRequest.UserId == 0 {
		h.log.Infof("%s userReq is empty", userRequest.UserId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.log.Infof("Querying database for the userid %v", userRequest.UserId)
	user, err := h.db.GetUser(userRequest.UserId)
	if err != nil {
		h.log.Errorw("Error getting user from db", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.log.Infof("Successfully retrieved user with the following fields %v, %v, %v, %v", user.UserId, user.FirstName, user.LastName, user.Email)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		h.log.Errorw("Error encoding user to json", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	return
}
