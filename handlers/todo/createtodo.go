package todo

import (
	"encoding/json"
	"github.com/clementbowe14/sample-uberfx-backend/db"
	"go.uber.org/zap"
	"net/http"
)

type CreateTodoHandler struct {
	log *zap.SugaredLogger
	db  *db.UserDb
}

type CreateTodoRequest struct {
	UserId      int32  `json:"user_id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func NewCreateTodoHandler(logger *zap.SugaredLogger, db *db.UserDb) *CreateTodoHandler {
	return &CreateTodoHandler{
		log: logger,
		db:  db,
	}
}

func (h *CreateTodoHandler) Pattern() string { return "POST /users/{user_id}/todo" }

func (h *CreateTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Infof("%s %s", r.Method, r.URL.Path)

	h.log.Infof("%s %s", r.Method, r.URL.Path)

	todo := CreateTodoRequest{}
	err := json.NewDecoder(r.Body).Decode(todo)

	if err != nil {
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
