package user

import (
	"encoding/json"
	"github.com/clementbowe14/sample-uberfx-backend/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

type CreateUserHandler struct {
	log *zap.SugaredLogger
	db  *db.UserDb
}

type CreateUserRequest struct {
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func NewCreateUserHandler(log *zap.SugaredLogger, userDb *db.UserDb) *CreateUserHandler {
	return &CreateUserHandler{
		log: log,
		db:  userDb,
	}
}

func (*CreateUserHandler) Pattern() string {
	return "POST /users"
}

func (handler *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	handler.log.Infof("Starting request to create user %s %s", r.Method, r.URL.Path)
	if err != nil {
		handler.log.Errorw("Error reading body", "error", err)
	}

	var req CreateUserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		handler.log.Errorw("Error unmarshalling body", "error", err)
	}

	if len(req.Email) == 0 {
		handler.log.Errorw("Email address is required", "email", req.Email)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Email address is required"))
		if err != nil {
			handler.log.Errorw("Error writing response", "error", err)
		}
		return
	}

	if len(req.Password) < 8 {
		handler.log.Errorw("Password is too short", "password", req.Password)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Password is too short"))
		if err != nil {
			handler.log.Errorw("Error writing response", "error", err)
		}
		return
	}

	if len(req.FirstName) == 0 {
		handler.log.Errorw("You must add a first name", "first_name", req.FirstName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(req.LastName) == 0 {
		handler.log.Errorw("You must add a last name", "last_name", req.LastName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		handler.log.Errorw("Error generating password", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.log.Infof("Verified all fields have valid information %s %s %s %s", req.FirstName, req.LastName, password, req.Email)
	err = handler.db.CreateNewUser(string(password), req.Email, req.FirstName, req.LastName)

	if err != nil {
		handler.log.Errorw("Error creating new user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("Error creating new user"))
		if err != nil {
			handler.log.Errorw("An error occurred while writing response", "error", err)
			return
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("User created successfully"))
	if err != nil {
		handler.log.Errorw("An error occurred while writing response", "error", err)
		return
	}

	return
}
