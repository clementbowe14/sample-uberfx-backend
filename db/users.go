package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strings"
)

type UserDb struct {
	log  *zap.SugaredLogger
	conn *pgxpool.Pool
}

type GetUserQueryResult struct {
	UserId    int32  `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func NewUserDb(conn *pgxpool.Pool, log *zap.SugaredLogger) *UserDb {
	return &UserDb{
		conn: conn,
		log:  log,
	}
}

func (u *UserDb) CreateNewUser(password string, email string, firstName string, lastName string) error {
	u.log.Infof("Creating new user with the following information: %s %s %s %s", password, email, firstName, lastName)
	res, err := u.conn.Exec(context.Background(), `insert into users(password,  first_name, last_name, email) values ($1, $2, $3, $4)`,
		password, firstName, lastName, email)

	if err != nil {
		u.log.Errorf("res %v error %v", res, err.Error())
		return err
	}

	return nil
}

func (u *UserDb) GetUser(userId int32) (*GetUserQueryResult, error) {
	u.log.Infof("Getting user with id %d", userId)
	user := &GetUserQueryResult{}
	err := u.conn.QueryRow(context.Background(), `select id, first_name, last_name, email from users where id = $1`, userId).Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email)
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.TrimSpace(user.Email)
	if err != nil {
		u.log.Errorf("res %v error %v", user, err.Error())
		return nil, err
	}

	return user, nil
}

func (u *UserDb) UpdateUser(id int32, firstName string, lastName string, email string) error {
	u.log.Infof("Updating user with id %d", id)
	res, err := u.conn.Exec(context.Background(), `update users set first_name = $1, last_name = $2, email = $3 where id = $4`, firstName, lastName, email, id)
	if err != nil {
		u.log.Errorf("res %v error %v", res, err.Error())
		return err
	}

	return nil
}

func (u *UserDb) DeleteUser(id int32) error {
	u.log.Infof("Deleting user with id %d", id)
	res, err := u.conn.Exec(context.Background(), `delete from users where id = $1`, id)
	if err != nil {
		u.log.Errorf("res %v error %v", res, err.Error())
		return err
	}

	return nil
}
