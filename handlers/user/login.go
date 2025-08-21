package user

import (
	"github.com/clementbowe14/sample-uberfx-backend/db"
	"go.uber.org/zap"
	"net/http"
)

type LoginHandler struct {
	db  *db.UserDb
	log *zap.SugaredLogger
}

func NewLoginHandler(db *db.UserDb, log *zap.SugaredLogger) *LoginHandler {
	return &LoginHandler{
		db:  db,
		log: log,
	}
}

func Pattern() string {
	return "POST /user/login"
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Infof("starting login %s %s", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":true}`))
}
